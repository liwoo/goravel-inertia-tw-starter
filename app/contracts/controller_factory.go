package contracts

import (
	"fmt"
	"reflect"
)

// ControllerFactory ensures all controllers implement required contracts
type ControllerFactory struct {
	registeredControllers map[string]ResourceControllerContract
	registeredPages       map[string]PageControllerContract
	validationEnabled     bool
}

// NewControllerFactory creates a new controller factory
func NewControllerFactory() *ControllerFactory {
	return &ControllerFactory{
		registeredControllers: make(map[string]ResourceControllerContract),
		registeredPages:       make(map[string]PageControllerContract),
		validationEnabled:     true,
	}
}

// RegisterController registers a resource controller and validates its implementation
func (cf *ControllerFactory) RegisterController(name string, controller ResourceControllerContract) error {
	// Validate that controller implements all required contracts
	if cf.validationEnabled {
		if err := cf.validateControllerContracts(controller); err != nil {
			return fmt.Errorf("controller '%s' validation failed: %w", name, err)
		}
	}
	
	cf.registeredControllers[name] = controller
	return nil
}

// MustRegisterController registers a controller and panics if validation fails
func (cf *ControllerFactory) MustRegisterController(name string, controller ResourceControllerContract) {
	if err := cf.RegisterController(name, controller); err != nil {
		panic(fmt.Sprintf("Failed to register controller '%s': %v", name, err))
	}
}

// RegisterPageController registers a page controller and validates its implementation
func (cf *ControllerFactory) RegisterPageController(name string, controller PageControllerContract) error {
	// Validate that controller implements all required contracts
	if cf.validationEnabled {
		if err := cf.validatePageControllerContracts(controller); err != nil {
			return fmt.Errorf("page controller '%s' validation failed: %w", name, err)
		}
	}
	
	cf.registeredPages[name] = controller
	return nil
}

// MustRegisterPageController registers a page controller and panics if validation fails
func (cf *ControllerFactory) MustRegisterPageController(name string, controller PageControllerContract) {
	if err := cf.RegisterPageController(name, controller); err != nil {
		panic(fmt.Sprintf("Failed to register page controller '%s': %v", name, err))
	}
}

// GetController retrieves a registered controller
func (cf *ControllerFactory) GetController(name string) (ResourceControllerContract, error) {
	controller, exists := cf.registeredControllers[name]
	if !exists {
		return nil, fmt.Errorf("controller '%s' not found", name)
	}
	return controller, nil
}

// GetPageController retrieves a registered page controller
func (cf *ControllerFactory) GetPageController(name string) (PageControllerContract, error) {
	controller, exists := cf.registeredPages[name]
	if !exists {
		return nil, fmt.Errorf("page controller '%s' not found", name)
	}
	return controller, nil
}

// ListControllers returns all registered controller names
func (cf *ControllerFactory) ListControllers() []string {
	names := make([]string, 0, len(cf.registeredControllers))
	for name := range cf.registeredControllers {
		names = append(names, name)
	}
	return names
}

// ListPageControllers returns all registered page controller names
func (cf *ControllerFactory) ListPageControllers() []string {
	names := make([]string, 0, len(cf.registeredPages))
	for name := range cf.registeredPages {
		names = append(names, name)
	}
	return names
}

// ValidateAllControllers validates all registered controllers
func (cf *ControllerFactory) ValidateAllControllers() map[string]ControllerValidationResult {
	results := make(map[string]ControllerValidationResult)
	
	for name, controller := range cf.registeredControllers {
		results[name] = ValidateControllerImplementation(controller)
	}
	
	for name, controller := range cf.registeredPages {
		results[name+"_page"] = ValidateControllerImplementation(controller)
	}
	
	return results
}

// validateControllerContracts validates that a controller implements all required contracts
func (cf *ControllerFactory) validateControllerContracts(controller interface{}) error {
	controllerType := reflect.TypeOf(controller)
	
	// Check if controller implements ResourceControllerContract
	resourceControllerType := reflect.TypeOf((*ResourceControllerContract)(nil)).Elem()
	if !controllerType.Implements(resourceControllerType) {
		return fmt.Errorf("controller must implement ResourceControllerContract interface")
	}
	
	// Validate specific method implementations for CRUD operations
	requiredMethods := []string{
		"Index", "Show", "Store", "Update", "Delete",
		"ValidatePaginationRequest", "GetPaginationDefaults", "BuildPaginatedResponse",
		"ValidateCreateRequest", "ValidateUpdateRequest", "ValidateID", "GetValidationRules",
		"SuccessResponse", "CreatedResponse", "NoContentResponse", 
		"BadRequestResponse", "NotFoundResponse", "ForbiddenResponse",
		"ValidationErrorResponse", "InternalErrorResponse",
		"ResourceNotFoundResponse", "ResourceCreatedResponse", 
		"ResourceUpdatedResponse", "ResourceDeletedResponse",
		"CheckPermission", "GetCurrentUser", "RequireAuthentication", "BuildPermissionsMap",
	}
	
	missing := []string{}
	for _, methodName := range requiredMethods {
		if _, hasMethod := controllerType.MethodByName(methodName); !hasMethod {
			missing = append(missing, methodName)
		}
	}
	
	if len(missing) > 0 {
		return fmt.Errorf("controller is missing required methods: %v", missing)
	}
	
	return nil
}

// validatePageControllerContracts validates that a page controller implements all required contracts
func (cf *ControllerFactory) validatePageControllerContracts(controller interface{}) error {
	controllerType := reflect.TypeOf(controller)
	
	// Check if controller implements PageControllerContract
	pageControllerType := reflect.TypeOf((*PageControllerContract)(nil)).Elem()
	if !controllerType.Implements(pageControllerType) {
		return fmt.Errorf("page controller must implement PageControllerContract interface")
	}
	
	// Validate specific method implementations for page operations
	requiredMethods := []string{
		"Index",
		"ValidatePaginationRequest", "GetPaginationDefaults", "BuildPaginatedResponse",
		"CheckPermission", "GetCurrentUser", "RequireAuthentication", "BuildPermissionsMap",
		"BuildPageProps", "GetPageMetadata", "ValidatePageRequest",
	}
	
	missing := []string{}
	for _, methodName := range requiredMethods {
		if _, hasMethod := controllerType.MethodByName(methodName); !hasMethod {
			missing = append(missing, methodName)
		}
	}
	
	if len(missing) > 0 {
		return fmt.Errorf("page controller is missing required methods: %v", missing)
	}
	
	return nil
}

// EnableValidation enables or disables controller validation
func (cf *ControllerFactory) EnableValidation(enabled bool) {
	cf.validationEnabled = enabled
}

// ControllerRegistry provides a global controller registry
var GlobalControllerRegistry = NewControllerFactory()

// RegisterCrudController registers a controller in the global registry
func RegisterCrudController(name string, controller ResourceControllerContract) error {
	return GlobalControllerRegistry.RegisterController(name, controller)
}

// MustRegisterCrudController registers a controller in the global registry and panics on failure
func MustRegisterCrudController(name string, controller ResourceControllerContract) {
	GlobalControllerRegistry.MustRegisterController(name, controller)
}

// RegisterPageController registers a page controller in the global registry
func RegisterPageController(name string, controller PageControllerContract) error {
	return GlobalControllerRegistry.RegisterPageController(name, controller)
}

// MustRegisterPageController registers a page controller in the global registry and panics on failure
func MustRegisterPageController(name string, controller PageControllerContract) {
	GlobalControllerRegistry.MustRegisterPageController(name, controller)
}

// GetCrudController retrieves a controller from the global registry
func GetCrudController(name string) (ResourceControllerContract, error) {
	return GlobalControllerRegistry.GetController(name)
}

// GetPageController retrieves a page controller from the global registry
func GetPageController(name string) (PageControllerContract, error) {
	return GlobalControllerRegistry.GetPageController(name)
}

// ControllerInitializer is a function type for controller initialization
type ControllerInitializer func() ResourceControllerContract

// PageControllerInitializer is a function type for page controller initialization
type PageControllerInitializer func() PageControllerContract

// AutoRegisterControllers automatically registers controllers with validation
func AutoRegisterControllers(controllers map[string]ControllerInitializer) error {
	for name, initializer := range controllers {
		controller := initializer()
		if err := RegisterCrudController(name, controller); err != nil {
			return fmt.Errorf("failed to register controller '%s': %w", name, err)
		}
	}
	return nil
}

// AutoRegisterPageControllers automatically registers page controllers with validation
func AutoRegisterPageControllers(controllers map[string]PageControllerInitializer) error {
	for name, initializer := range controllers {
		controller := initializer()
		if err := RegisterPageController(name, controller); err != nil {
			return fmt.Errorf("failed to register page controller '%s': %w", name, err)
		}
	}
	return nil
}