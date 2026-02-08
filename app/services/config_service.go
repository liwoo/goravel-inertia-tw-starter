package services

import (
	"books-database/app/contracts"
	"books-database/app/http/requests"
	"books-database/app/models"

	"github.com/goravel/framework/facades"
)

// ConfigService implements business logic for configs using the builder pattern
type ConfigService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// NewConfigService creates a new Config service using the builder pattern
func NewConfigService() *ConfigService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Config]("configs", "id").
		WithSearchFields("name", "config_type", "description").                                 // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "name", "config_type", "description"). // Fields that can be used for sorting results
		WithFilterFields("name", "config_type").                                                // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                             // Validation rules for create/update operations
			"name":        "required|string|max:255",
			"config_type": "required|string|max:255",
			"description": "string|max:255",
		}).
		WithDefaultSort("created_at", "DESC").       // Default sorting when none specified
		WithScopeFiltering("configs", "created_by"). // Enable permission-based filtering
		WithTenantAwareness().
		Build() // Returns a fully configured CrudServiceContract

	configServiceInstance := &ConfigService{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, configServiceInstance, "ConfigService")

	return configServiceInstance
}

// Override GetColumnMapping to include config-specific mappings
func (s *ConfigService) GetColumnMapping() map[string]string {
	mapping := s.CrudServiceContract.GetColumnMapping()
	// Add config-specific mappings
	mapping["configType"] = "config_type"
	mapping["createdAt"] = "created_at"
	mapping["updatedAt"] = "updated_at"
	mapping["name"] = "name"
	mapping["description"] = "description"
	return mapping
}

// GetConfigStatistics returns statistics about configs grouped by type
func (s *ConfigService) GetConfigStatistics() (map[string]interface{}, error) {
	var stats struct {
		FinancingCount           int64
		ImprovementAspectsCount  int64
		BusinessCategoriesCount  int64
		IndustriesCount          int64
		SectorsCount             int64
		RegistrationStatusCount  int64
		DevelopmentPartnersCount int64
		TotalConfigs             int64
	}

	// Get total configs
	stats.TotalConfigs, _ = facades.Orm().Query().Model(&models.Config{}).Count()

	// Get counts for each config type
	stats.FinancingCount, _ = facades.Orm().Query().Model(&models.Config{}).
		Where("config_type = ?", requests.ConfigTypeFinancing).Count()

	stats.ImprovementAspectsCount, _ = facades.Orm().Query().Model(&models.Config{}).
		Where("config_type = ?", requests.ConfigTypeImprovementAspects).Count()

	stats.BusinessCategoriesCount, _ = facades.Orm().Query().Model(&models.Config{}).
		Where("config_type = ?", requests.ConfigTypeBusinessCategories).Count()

	stats.IndustriesCount, _ = facades.Orm().Query().Model(&models.Config{}).
		Where("config_type = ?", requests.ConfigTypeIndustries).Count()

	stats.SectorsCount, _ = facades.Orm().Query().Model(&models.Config{}).
		Where("config_type = ?", requests.ConfigTypeSectors).Count()

	stats.RegistrationStatusCount, _ = facades.Orm().Query().Model(&models.Config{}).
		Where("config_type = ?", requests.ConfigTypeRegistrationStatus).Count()

	stats.DevelopmentPartnersCount, _ = facades.Orm().Query().Model(&models.Config{}).
		Where("config_type = ?", requests.ConfigTypeDevelopmentPartners).Count()

	return map[string]interface{}{
		"totalConfigs":             stats.TotalConfigs,
		"financingCount":           stats.FinancingCount,
		"improvementAspectsCount":  stats.ImprovementAspectsCount,
		"businessCategoriesCount":  stats.BusinessCategoriesCount,
		"industriesCount":          stats.IndustriesCount,
		"sectorsCount":             stats.SectorsCount,
		"registrationStatusCount":  stats.RegistrationStatusCount,
		"developmentPartnersCount": stats.DevelopmentPartnersCount,
	}, nil
}
