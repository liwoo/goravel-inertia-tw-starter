package services

import (
	"encoding/json"
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
		WithSearchFields("name", "postal_address", "physical_address", "registration_status", "partners", "product_types", "service_list"). // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "name", "postal_address", "physical_address", "registration_status").              // Fields that can be used for sorting results
		WithFilterFields("name", "postal_address", "physical_address", "registration_status", "partners", "product_types", "service_list"). // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                                                                         // Validation rules for create/update operations
			"registration_status": "string|max:255",
			"name":                "required|string|max:255",
			"postal_address":      "required|string|max:255",
			"physical_address":    "string|max:255",
			"partners":            "array",
			"product_types":       "array",
			"service_list":        "array",
		}).
		WithDefaultSort("created_at", "DESC").     // Default sorting when none specified
		WithScopeFiltering("bdsps", "created_by"). // Enable permission-based filtering
		WithBeforeUpdate(func(id uint, data map[string]interface{}) error {
			if partners, ok := data["partners"]; ok {
				if partnersSlice, ok := partners.([]interface{}); ok {
					bytes, err := json.Marshal(partnersSlice)
					if err == nil {
						data["partners_json"] = string(bytes)
					}
				}
				delete(data, "partners")
			}

			if productTypes, ok := data["product_types"]; ok {
				if productTypesSlice, ok := productTypes.([]interface{}); ok {
					bytes, err := json.Marshal(productTypesSlice)
					if err == nil {
						data["product_types_json"] = string(bytes)
					}
				}
				delete(data, "product_types")
			}

			if serviceList, ok := data["service_list"]; ok {
				if serviceListSlice, ok := serviceList.([]interface{}); ok {
					bytes, err := json.Marshal(serviceListSlice)
					if err == nil {
						data["service_list_json"] = string(bytes)
					}
				}
				delete(data, "service_list")
			}

			return nil
		}).
		Build() // Returns a fully configured CrudServiceContract

	bdspServiceInstance := &BdspService{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, bdspServiceInstance, "BdspService")

	return bdspServiceInstance
}

// GetFilterDefinitions returns filter definitions for the bdsps resource
func (s *BdspService) GetFilterDefinitions() []contracts.FilterDefinition {
	registrationStatuses := []string{
		models.RegistrationPending,
		models.RegistrationConfirmed,
		models.RegistrationRejected,
		models.RegistrationSuspended,
	}

	return []contracts.FilterDefinition{
		// Name - string search
		contracts.NewFilterDefinition(
			"name",
			"Name",
			contracts.FilterTypeString,
			nil,
		),
		// Registration Status - enum
		contracts.NewFilterDefinition(
			"registration_status",
			"Registration Status",
			contracts.FilterTypeEnum,
			&registrationStatuses,
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
		// Partners - string search
		contracts.NewFilterDefinition(
			"partners",
			"Partners",
			contracts.FilterTypeString,
			nil,
		),
		// Product Types - string search
		contracts.NewFilterDefinition(
			"product_types",
			"Product Types",
			contracts.FilterTypeString,
			nil,
		),
		// Service List - string search
		contracts.NewFilterDefinition(
			"service_list",
			"Service List",
			contracts.FilterTypeString,
			nil,
		),
		// Created Date - datetime filter
		contracts.NewFilterDefinition(
			"created_at",
			"Date Created",
			contracts.FilterTypeDateTime,
			nil,
		),
	}
}

// GetBdspStatistics returns statistics about bdsps
func (s *BdspService) GetBdspStatistics() (map[string]interface{}, error) {
	var stats struct {
		TotalBdsps     int64
		PendingBdsps   int64
		ActiveBdsps    int64
		RejectedBdsps  int64
		SuspendedBdsps int64
	}

	// Get total bdsps (excluding soft deleted)
	stats.TotalBdsps, _ = facades.Orm().Query().Model(&models.Bdsp{}).Where("deleted_at IS NULL").Count()

	// Get pending bdsps
	stats.PendingBdsps, _ = facades.Orm().Query().Model(&models.Bdsp{}).Where("registration_status = ? AND deleted_at IS NULL", models.RegistrationPending).Count()

	// Get active bdsps
	stats.ActiveBdsps, _ = facades.Orm().Query().Model(&models.Bdsp{}).Where("registration_status = ? AND deleted_at IS NULL", models.RegistrationConfirmed).Count()

	// Get rejected bdsps
	stats.RejectedBdsps, _ = facades.Orm().Query().Model(&models.Bdsp{}).Where("registration_status = ? AND deleted_at IS NULL", models.RegistrationRejected).Count()

	// Get suspended bdsps
	stats.SuspendedBdsps, _ = facades.Orm().Query().Model(&models.Bdsp{}).Where("registration_status = ? AND deleted_at IS NULL", models.RegistrationSuspended).Count()

	return map[string]interface{}{
		"totalBdsps":     stats.TotalBdsps,
		"pendingBdsps":   stats.PendingBdsps,
		"activeBdsps":    stats.ActiveBdsps,
		"rejectedBdsps":  stats.RejectedBdsps,
		"suspendedBdsps": stats.SuspendedBdsps,
	}, nil
}
