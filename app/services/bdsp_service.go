package services

import (
	"encoding/json"
	"fmt"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
	"strconv"
	"time"

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
		WithSearchFields("ubdsp_number", "name", "postal_address", "physical_address", "registration_status", "partners_json", "product_types_json", "service_list_json"). // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "ubdsp_number", "name", "postal_address", "physical_address", "registration_status").                             // Fields that can be used for sorting results
		WithFilterFields("ubdsp_number", "name", "postal_address", "physical_address", "registration_status", "partners_json", "product_types_json", "service_list_json"). // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                                                                                                        // Validation rules for create/update operations
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
		WithBeforeCreate(func(data map[string]interface{}) error {
			// Generate UBDSP number if not provided
			if _, exists := data["ubdsp_number"]; !exists || data["ubdsp_number"] == "" {
				ubi, err := generateBdspUBI()
				if err != nil {
					return fmt.Errorf("failed to generate UBDSP number: %w", err)
				}
				data["ubdsp_number"] = ubi
			}

			// Handle partners array to JSON conversion
			if partners, ok := data["partners"]; ok {
				if partnersSlice, ok := partners.([]interface{}); ok {
					bytes, err := json.Marshal(partnersSlice)
					if err == nil {
						data["partners_json"] = string(bytes)
					}
				}
				delete(data, "partners")
			}

			// Handle product_types array to JSON conversion
			if productTypes, ok := data["product_types"]; ok {
				if productTypesSlice, ok := productTypes.([]interface{}); ok {
					bytes, err := json.Marshal(productTypesSlice)
					if err == nil {
						data["product_types_json"] = string(bytes)
					}
				}
				delete(data, "product_types")
			}

			// Handle service_list array to JSON conversion
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
		WithBeforeUpdate(func(id uint, data map[string]interface{}) error {
			if _, exists := data["ubdsp_number"]; !exists || data["ubdsp_number"] == "" {
				ubi, err := generateBdspUBI()
				if err != nil {
					return fmt.Errorf("failed to generate UBI: %w", err)
				}
				data["ubdsp_number"] = ubi
			}

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

func generateBdspUBI() (string, error) {
	// Get current year
	year := time.Now().Year()

	// Get next sequential number
	sequentialNumber, err := getNextSequentialBdspNumber(year)
	if err != nil {
		return "", err
	}

	// Build UBI without check digit
	ubiWithoutCheck := fmt.Sprintf("MW-%04d-%06d", year, sequentialNumber)

	// Calculate check digit
	checkDigit := calculateLuhnCheckDigit(ubiWithoutCheck)

	// Return complete UBI
	return fmt.Sprintf("%s-%d", ubiWithoutCheck, checkDigit), nil
}

// getNextSequentialBdspNumber gets the next sequential number for the given year, district, and category
func getNextSequentialBdspNumber(year int) (int, error) {
	// Query the database to find the highest sequential number for this year/district/category combination
	var maxSequence int

	// Pattern to match: MW-YYYY-
	pattern := fmt.Sprintf("MW-%04d-%%", year)

	var bdsp models.Bdsp
	err := facades.Orm().Query().
		Where("ubdsp_number LIKE ?", pattern).
		Order("ubdsp_number DESC").
		First(&bdsp)

	if err != nil || bdsp.ID == 0 {
		// No existing records, start from 1
		return 1, nil
	}

	// Extract the sequential number from the ubdsp_number
	// Format: MW-YYYY-NNNNNN-C (4 parts: [0]=MW, [1]=YYYY, [2]=NNNNNN, [3]=C)
	parts := splitUBI(bdsp.UbdspNumber)
	if len(parts) >= 3 {
		var parseErr error
		maxSequence, parseErr = strconv.Atoi(parts[2])
		if parseErr != nil {
			// If parsing fails, start from 1
			return 1, nil
		}
	}

	// Return next number
	return maxSequence + 1, nil
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
		// UBDSP Number - string search
		contracts.NewFilterDefinition(
			"ubdsp_number",
			"UBDSP Number",
			contracts.FilterTypeString,
			nil,
		),
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
		// Partners - array search
		contracts.NewFilterDefinition(
			"partners_json",
			"Partners",
			contracts.FilterTypeArray,
			nil,
		),
		// Product Types - array search
		contracts.NewFilterDefinition(
			"product_types_json",
			"Product Types",
			contracts.FilterTypeArray,
			nil,
		),
		// Service List - array search
		contracts.NewFilterDefinition(
			"service_list_json",
			"Service List",
			contracts.FilterTypeArray,
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

	// Get total bdsps (GORM automatically handles deleted_at IS NULL for soft delete models)
	stats.TotalBdsps, _ = facades.Orm().Query().Model(&models.Bdsp{}).Count()

	// Get pending bdsps
	stats.PendingBdsps, _ = facades.Orm().Query().Model(&models.Bdsp{}).Where("registration_status = ?", models.RegistrationPending).Count()

	// Get active bdsps
	stats.ActiveBdsps, _ = facades.Orm().Query().Model(&models.Bdsp{}).Where("registration_status = ?", models.RegistrationConfirmed).Count()

	// Get rejected bdsps
	stats.RejectedBdsps, _ = facades.Orm().Query().Model(&models.Bdsp{}).Where("registration_status = ?", models.RegistrationRejected).Count()

	// Get suspended bdsps
	stats.SuspendedBdsps, _ = facades.Orm().Query().Model(&models.Bdsp{}).Where("registration_status = ?", models.RegistrationSuspended).Count()

	return map[string]interface{}{
		"totalBdsps":     stats.TotalBdsps,
		"pendingBdsps":   stats.PendingBdsps,
		"activeBdsps":    stats.ActiveBdsps,
		"rejectedBdsps":  stats.RejectedBdsps,
		"suspendedBdsps": stats.SuspendedBdsps,
	}, nil
}
