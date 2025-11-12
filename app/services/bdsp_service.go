package services

import (
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"

	"github.com/goravel/framework/facades"
)

// BdspService implements business logic for bdsps using the builder pattern
type BdspService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// NewBdspService creates a new Bdsp service using the builder pattern
func NewBdspService() *BdspService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Bdsp]("bdsps", "id").
		WithSearchFields("name", "postal_address", "physical_address", "registration_status", "product_types_json", "service_list_json", "associated_partners_json", "ip_address", "user_agent").                                                                           // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "name", "postal_address", "physical_address", "registration_status", "product_types_json", "service_list_json", "associated_partners_json", "created_by", "updated_by", "deleted_by", "ip_address", "user_agent"). // Fields that can be used for sorting results
		WithFilterFields("name", "postal_address", "physical_address", "registration_status", "product_types_json", "service_list_json", "associated_partners_json", "created_by", "updated_by", "deleted_by", "ip_address", "user_agent").                                 // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                                                                                                                                                                                                         // Validation rules for create/update operations
			"deleted_by":          "numeric",
			"ip_address":          "string|max:255",
			"user_agent":          "string|max:255",
			"name":                "required|string|max:255",
			"physical_address":    "string|max:255",
			"product_types":       "required|array",
			"service_list":        "required",
			"created_by":          "numeric",
			"updated_by":          "numeric",
			"postal_address":      "string|max:255",
			"registration_status": "required|string|max:255",
			"associated_partners": "required|array",
		}).
		WithDefaultSort("created_at", "DESC").     // Default sorting when none specified
		WithScopeFiltering("bdsps", "created_by"). // Enable permission-based filtering
		WithSoftDeletes().                         // Enable soft delete support
		Build()                                    // Returns a fully configured CrudServiceContract

	bdspServiceInstance := &BdspService{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, bdspServiceInstance, "BdspService")

	return bdspServiceInstance
}

// GetFilterDefinitions returns filter definitions for the bdsps resource
// This method is OPTIONAL - only implement if you need custom filter UI components
func (s *BdspService) GetFilterDefinitions() []contracts.FilterDefinition {
	// Get registration status options from config
	registrationStatuses := s.getConfigValues("Registration Status")

	return []contracts.FilterDefinition{
		// Name - string search
		contracts.NewFilterDefinition(
			"name",
			"BDSP Name",
			contracts.FilterTypeString,
			nil, // Use all string operators (equals, contains, starts_with, etc.)
		),

		// Postal Address - string search
		contracts.NewFilterDefinition(
			"postal_address",
			"Postal Address",
			contracts.FilterTypeString,
			nil,
		),

		// Physical Address - string search
		contracts.NewFilterDefinition(
			"physical_address",
			"Physical Address",
			contracts.FilterTypeString,
			nil,
		),

		// Registration Status - enum from config
		contracts.NewFilterDefinition(
			"registration_status",
			"Registration Status",
			contracts.FilterTypeEnum,
			&registrationStatuses, // Available options from config
		),

		// Product Types - string search (searches within JSON)
		contracts.NewFilterDefinition(
			"product_types_json",
			"Product Types",
			contracts.FilterTypeString,
			nil, // Use all string operators for JSON search
		),

		// Associated Partners - string search (searches within JSON)
		contracts.NewFilterDefinition(
			"associated_partners_json",
			"Associated Partners",
			contracts.FilterTypeString,
			nil,
		),

		// Created Date - datetime filter
		contracts.NewFilterDefinition(
			"created_at",
			"Date Created",
			contracts.FilterTypeDateTime,
			nil, // Use all datetime operators
		),

		// Updated Date - datetime filter
		contracts.NewFilterDefinition(
			"updated_at",
			"Date Updated",
			contracts.FilterTypeDateTime,
			nil,
		),
	}
}

// getConfigValues retrieves config values for a given config type
func (s *BdspService) getConfigValues(configType string) []string {
	var configs []models.Config
	facades.Orm().Query().
		Where("config_type = ?", configType).
		Order("name ASC").
		Get(&configs)

	values := make([]string, len(configs))
	for i, config := range configs {
		values[i] = config.Name
	}

	return values
}

// Add domain-specific methods below this line
// Examples:
// - GetByStatus(status string) ([]*models.Bdsp, error)
// - GetActive() ([]*models.Bdsp, error)
// - Custom business logic methods specific to bdsps
