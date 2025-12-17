package services

import (
	"fmt"

	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
)

// BusinessFormalisationService implements business_formalisation-specific business logic using the builder pattern
type BusinessFormalisationService struct {
	contracts.CrudServiceContract // Embedded - automatically exposes all methods!
}

// NewBusinessFormalisationService creates a new business_formalisation service using the builder pattern
func NewBusinessFormalisationService() *BusinessFormalisationService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.BusinessFormalisation]("business_formalisation", "id").
		WithSearchFields("sme_id").                                                                                                          // REQUIRED
		WithSortFields("id", "sme_id", "formalisation_score", "annual_turnover", "estimated_value_of_assets", "created_at", "updated_at").   // REQUIRED
		WithFilterFields("sme_id", "has_bank_account", "has_tax_clarification", "is_registered_for_vat",
			"is_member_of_association", "is_affiliated", "has_export_license", "has_accessed_bds"). // REQUIRED
		WithValidationRules(map[string]interface{}{ // REQUIRED
			"sme_id":                    "required|numeric|min:1",
			"has_bank_account":          "boolean",
			"has_tax_clarification":     "boolean",
			"is_registered_for_vat":     "boolean",
			"is_member_of_association":  "boolean",
			"is_affiliated":             "boolean",
			"has_export_license":        "boolean",
			"has_accessed_bds":          "boolean",
			"annual_turnover":           "numeric|min:0",
			"estimated_value_of_assets": "numeric|min:0",
			// Note: formalisation_score is intentionally not included in validation rules
			// as it should be read-only and calculated internally
		}).
		WithDefaultSort("created_at", "DESC").                      // Optional
		WithSoftDeletes().                                          // Optional
		WithScopeFiltering("business_formalisation", "created_by"). // Optional
		WithBeforeCreate(func(data map[string]interface{}) error {  // Optional
			// Check if a BusinessFormalisation already exists for this SME
			// to prevent duplicate records
			if smeID, exists := data["sme_id"]; exists && smeID != nil {
				var smeIDInt int
				switch v := smeID.(type) {
				case int:
					smeIDInt = v
				case int64:
					smeIDInt = int(v)
				case float64:
					smeIDInt = int(v)
				case uint:
					smeIDInt = int(v)
				}

				if smeIDInt > 0 {
					var existing models.BusinessFormalisation
					err := facades.Orm().Query().Where("sme_id = ?", smeIDInt).First(&existing)
					if err == nil && existing.ID > 0 {
						// Record already exists - return error with existing ID for upsert handling
						return fmt.Errorf("EXISTING_RECORD:%d", existing.ID)
					}
				}
			}

			// Ensure boolean fields have default values
			boolFields := []string{
				"has_bank_account", "has_tax_clarification", "is_registered_for_vat",
				"is_member_of_association", "is_affiliated", "has_export_license", "has_accessed_bds",
			}

			for _, field := range boolFields {
				if _, exists := data[field]; !exists {
					data[field] = false
				}
			}

			// Ensure numeric fields have default values
			if _, exists := data["annual_turnover"]; !exists {
				data["annual_turnover"] = 0.0
			}
			if _, exists := data["estimated_value_of_assets"]; !exists {
				data["estimated_value_of_assets"] = 0.0
			}

			// Remove formalisation_score if it was accidentally included
			delete(data, "formalisation_score")

			return nil
		}).
		WithBeforeUpdate(func(id uint, data map[string]interface{}) error { // Optional
			// Remove formalisation_score if it was accidentally included
			delete(data, "formalisation_score")

			return nil
		}).
		WithAfterCreate(func(model *models.BusinessFormalisation) error {
			// Recalculate formalisation score and classification after creating a business formalisation record
			if model != nil && model.SmeID > 0 {
				smeService := NewSmeService()
				_, err := smeService.CalculateFormalisationScore(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate formalisation score after create: %v", err)
				}
				// Recalculate classification (uses turnover and assets from BusinessFormalisation)
				_, err = smeService.CalculateClassification(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate classification after create: %v", err)
				}
			}
			return nil
		}).
		WithAfterUpdate(func(model *models.BusinessFormalisation) error {
			// Recalculate formalisation score and classification after updating a business formalisation record
			if model != nil && model.SmeID > 0 {
				smeService := NewSmeService()
				_, err := smeService.CalculateFormalisationScore(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate formalisation score after update: %v", err)
				}
				// Recalculate classification (uses turnover and assets from BusinessFormalisation)
				_, err = smeService.CalculateClassification(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate classification after update: %v", err)
				}
			}
			return nil
		}).
		Build()

	businessFormalisationServiceInstance := &BusinessFormalisationService{
		CrudServiceContract: service,
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, businessFormalisationServiceInstance, "BusinessFormalisationService")

	return businessFormalisationServiceInstance
}

// CreateOrUpdate creates a new BusinessFormalisation or updates an existing one for the given SME
// This is the recommended method to use from frontend to avoid duplicate records
func (s *BusinessFormalisationService) CreateOrUpdate(data map[string]interface{}) (interface{}, error) {
	// Extract sme_id
	smeID, exists := data["sme_id"]
	if !exists || smeID == nil {
		return nil, fmt.Errorf("sme_id is required")
	}

	var smeIDInt int
	switch v := smeID.(type) {
	case int:
		smeIDInt = v
	case int64:
		smeIDInt = int(v)
	case float64:
		smeIDInt = int(v)
	case uint:
		smeIDInt = int(v)
	}

	if smeIDInt <= 0 {
		return nil, fmt.Errorf("invalid sme_id")
	}

	// Check if a record already exists for this SME
	var existing models.BusinessFormalisation
	err := facades.Orm().Query().Where("sme_id = ?", smeIDInt).First(&existing)

	if err == nil && existing.ID > 0 {
		// Record exists - update it
		facades.Log().Info("BusinessFormalisation already exists for SME, updating instead of creating", map[string]interface{}{
			"sme_id":      smeIDInt,
			"existing_id": existing.ID,
		})
		return s.Update(uint(existing.ID), data)
	}

	// No existing record - create new one
	// Remove the sme_id check since we've already verified
	return s.Create(data)
}
