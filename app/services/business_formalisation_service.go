package services

import (
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
		WithSearchFields("sme_id").                                                                             // REQUIRED
		WithSortFields("id", "sme_id", "formalisation_score", "annual_turnover", "estimated_value_of_assets", "created_at", "updated_at"). // REQUIRED
		WithFilterFields("sme_id", "has_bank_account", "has_tax_clarification", "is_registered_for_vat",
			"is_member_of_association", "is_affiliated", "has_export_license", "has_accessed_bds").          // REQUIRED
		WithValidationRules(map[string]interface{}{                                                             // REQUIRED
			"sme_id":                     "required|numeric|min:1",
			"has_bank_account":           "boolean",
			"has_tax_clarification":      "boolean",
			"is_registered_for_vat":      "boolean",
			"is_member_of_association":   "boolean",
			"is_affiliated":              "boolean",
			"has_export_license":         "boolean",
			"has_accessed_bds":           "boolean",
			"annual_turnover":            "numeric|min:0",
			"estimated_value_of_assets":  "numeric|min:0",
			// Note: formalisation_score is intentionally not included in validation rules
			// as it should be read-only and calculated internally
		}).
		WithDefaultSort("created_at", "DESC").                     // Optional
		WithSoftDeletes().                                         // Optional
		WithScopeFiltering("business_formalisation", "created_by"). // Optional
		WithBeforeCreate(func(data map[string]interface{}) error { // Optional
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
		Build()

	return &BusinessFormalisationService{
		CrudServiceContract: service,
	}
}

// Additional custom methods can be added here if needed
// For example, a method to calculate the formalisation score could be added:
/*
func (s *BusinessFormalisationService) CalculateFormalisationScore(id uint) error {
	// Implementation for calculating and updating the formalisation score
	// This would be the only way to update the formalisation_score field
	return nil
}
*/