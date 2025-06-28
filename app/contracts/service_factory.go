package contracts

import (
	"fmt"
	"reflect"
)

// ServiceFactory ensures all services implement required contracts
type ServiceFactory struct {
	registeredServices map[string]CompleteCrudService
	validationEnabled  bool
}

// NewServiceFactory creates a new service factory
func NewServiceFactory() *ServiceFactory {
	return &ServiceFactory{
		registeredServices: make(map[string]CompleteCrudService),
		validationEnabled:  true,
	}
}

// RegisterService registers a service and validates its implementation
func (sf *ServiceFactory) RegisterService(name string, service CompleteCrudService) error {
	// Validate that service implements all required contracts
	if sf.validationEnabled {
		if err := sf.validateServiceContracts(service); err != nil {
			return fmt.Errorf("service '%s' validation failed: %w", name, err)
		}
	}
	
	sf.registeredServices[name] = service
	return nil
}

// MustRegisterService registers a service and panics if validation fails
func (sf *ServiceFactory) MustRegisterService(name string, service CompleteCrudService) {
	if err := sf.RegisterService(name, service); err != nil {
		panic(fmt.Sprintf("Failed to register service '%s': %v", name, err))
	}
}

// GetService retrieves a registered service
func (sf *ServiceFactory) GetService(name string) (CompleteCrudService, error) {
	service, exists := sf.registeredServices[name]
	if !exists {
		return nil, fmt.Errorf("service '%s' not found", name)
	}
	return service, nil
}

// ListServices returns all registered service names
func (sf *ServiceFactory) ListServices() []string {
	names := make([]string, 0, len(sf.registeredServices))
	for name := range sf.registeredServices {
		names = append(names, name)
	}
	return names
}

// ValidateAllServices validates all registered services
func (sf *ServiceFactory) ValidateAllServices() map[string]ServiceValidationResult {
	results := make(map[string]ServiceValidationResult)
	
	for name, service := range sf.registeredServices {
		results[name] = ValidateServiceImplementation(service)
	}
	
	return results
}

// validateServiceContracts validates that a service implements all required contracts
func (sf *ServiceFactory) validateServiceContracts(service interface{}) error {
	serviceType := reflect.TypeOf(service)
	
	// Check if service implements CompleteCrudService
	completeCrudType := reflect.TypeOf((*CompleteCrudService)(nil)).Elem()
	if !serviceType.Implements(completeCrudType) {
		return fmt.Errorf("service must implement CompleteCrudService interface")
	}
	
	// Validate specific method implementations
	requiredMethods := []string{
		"GetList", "GetListAdvanced", "GetByID", "Create", "Update", "Delete",
		"GetPaginatedList", "ValidatePaginationParams", "GetMaxPageSize", "GetDefaultPageSize",
		"GetSortableFields", "ValidateSortField", "ValidateSortDirection", "GetDefaultSort", "MapSortField",
		"GetFilterableFields", "ValidateFilterField", "ValidateFilterValue", "GetSearchableFields", "BuildFilterQuery",
		"Search", "ValidateSearchQuery",
		"BulkCreate", "BulkUpdate", "BulkDelete", "ValidateBulkOperation",
		"GetTableName", "GetPrimaryKey", "GetModel", "GetValidationRules", "GetColumnMapping",
	}
	
	missing := []string{}
	for _, methodName := range requiredMethods {
		if _, hasMethod := serviceType.MethodByName(methodName); !hasMethod {
			missing = append(missing, methodName)
		}
	}
	
	if len(missing) > 0 {
		return fmt.Errorf("service is missing required methods: %v", missing)
	}
	
	return nil
}

// EnableValidation enables or disables service validation
func (sf *ServiceFactory) EnableValidation(enabled bool) {
	sf.validationEnabled = enabled
}

// ServiceRegistry provides a global service registry
var GlobalServiceRegistry = NewServiceFactory()

// RegisterCrudService registers a service in the global registry
func RegisterCrudService(name string, service CompleteCrudService) error {
	return GlobalServiceRegistry.RegisterService(name, service)
}

// MustRegisterCrudService registers a service in the global registry and panics on failure
func MustRegisterCrudService(name string, service CompleteCrudService) {
	GlobalServiceRegistry.MustRegisterService(name, service)
}

// GetCrudService retrieves a service from the global registry
func GetCrudService(name string) (CompleteCrudService, error) {
	return GlobalServiceRegistry.GetService(name)
}

// ServiceInitializer is a function type for service initialization
type ServiceInitializer func() CompleteCrudService

// AutoRegisterServices automatically registers services with validation
func AutoRegisterServices(services map[string]ServiceInitializer) error {
	for name, initializer := range services {
		service := initializer()
		if err := RegisterCrudService(name, service); err != nil {
			return fmt.Errorf("failed to register service '%s': %w", name, err)
		}
	}
	return nil
}