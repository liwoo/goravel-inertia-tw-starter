package services

import (
	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
)

// BusinessEmployeeSummaryService implements business logic for business employee summaries using the builder pattern
type BusinessEmployeeSummaryService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// NewBusinessEmployeeSummaryService creates a new BusinessEmployeeSummary service using the builder pattern
func NewBusinessEmployeeSummaryService() *BusinessEmployeeSummaryService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.BusinessEmployeeSummary]("business_employee_summary", "id").
		WithSearchFields("sme_id"). // REQUIRED
		WithSortFields("id", "sme_id", "full_time_males", "full_time_females", "part_time_males", "part_time_females",
			"intern_males", "intern_females", "created_at", "updated_at"). // REQUIRED
		WithFilterFields("sme_id"). // REQUIRED
		WithValidationRules(map[string]interface{}{ // REQUIRED
			"sme_id":           "required|numeric|min:1",
			"full_time_males":  "numeric|min:0",
			"full_time_females": "numeric|min:0",
			"part_time_males":  "numeric|min:0",
			"part_time_females": "numeric|min:0",
			"intern_males":     "numeric|min:0",
			"intern_females":   "numeric|min:0",
		}).
		WithDefaultSort("created_at", "DESC").                        // Optional
		WithSoftDeletes().                                            // Optional
		WithScopeFiltering("business_employee_summary", "created_by"). // Optional
		WithBeforeCreate(func(data map[string]interface{}) error {
			// Ensure numeric fields have default values
			numericFields := []string{
				"full_time_males", "full_time_females", "part_time_males",
				"part_time_females", "intern_males", "intern_females",
			}

			for _, field := range numericFields {
				if _, exists := data[field]; !exists {
					data[field] = 0
				}
			}

			return nil
		}).
		WithAfterCreate(func(model *models.BusinessEmployeeSummary) error {
			// Recalculate formalisation score and classification after creating an employee summary record
			if model != nil && model.SmeID > 0 {
				smeService := NewSmeService()
				_, err := smeService.CalculateFormalisationScore(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate formalisation score after employee summary create: %v", err)
				}
				// Recalculate classification (uses employee count from BusinessEmployeeSummary)
				_, err = smeService.CalculateClassification(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate classification after employee summary create: %v", err)
				}
			}
			return nil
		}).
		WithAfterUpdate(func(model *models.BusinessEmployeeSummary) error {
			// Recalculate formalisation score and classification after updating an employee summary record
			if model != nil && model.SmeID > 0 {
				smeService := NewSmeService()
				_, err := smeService.CalculateFormalisationScore(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate formalisation score after employee summary update: %v", err)
				}
				// Recalculate classification (uses employee count from BusinessEmployeeSummary)
				_, err = smeService.CalculateClassification(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate classification after employee summary update: %v", err)
				}
			}
			return nil
		}).
		Build()

	businessEmployeeSummaryServiceInstance := &BusinessEmployeeSummaryService{
		CrudServiceContract: service,
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, businessEmployeeSummaryServiceInstance, "BusinessEmployeeSummaryService")

	return businessEmployeeSummaryServiceInstance
}

// Override GetColumnMapping to include BusinessEmployeeSummary-specific mappings
func (s *BusinessEmployeeSummaryService) GetColumnMapping() map[string]string {
	mapping := s.CrudServiceContract.GetColumnMapping()
	// Add BusinessEmployeeSummary-specific camelCase to snake_case mappings
	mapping["smeId"] = "sme_id"
	mapping["fullTimeMales"] = "full_time_males"
	mapping["fullTimeFemales"] = "full_time_females"
	mapping["partTimeMales"] = "part_time_males"
	mapping["partTimeFemales"] = "part_time_females"
	mapping["internMales"] = "intern_males"
	mapping["internFemales"] = "intern_females"
	mapping["createdAt"] = "created_at"
	mapping["updatedAt"] = "updated_at"
	mapping["createdBy"] = "created_by"
	mapping["updatedBy"] = "updated_by"
	mapping["deletedBy"] = "deleted_by"
	mapping["ipAddress"] = "ip_address"
	mapping["userAgent"] = "user_agent"
	return mapping
}

// GetFilterDefinitions returns filter definitions for the BusinessEmployeeSummary resource
func (s *BusinessEmployeeSummaryService) GetFilterDefinitions() []contracts.FilterDefinition {
	return []contracts.FilterDefinition{
		contracts.NewFilterDefinition(
			"sme_id",
			"MSME ID",
			contracts.FilterTypeNumber,
			nil,
		),
		contracts.NewFilterDefinition(
			"full_time_males",
			"Full Time Males",
			contracts.FilterTypeNumber,
			nil,
		),
		contracts.NewFilterDefinition(
			"full_time_females",
			"Full Time Females",
			contracts.FilterTypeNumber,
			nil,
		),
		contracts.NewFilterDefinition(
			"created_at",
			"Date Created",
			contracts.FilterTypeDateTime,
			nil,
		),
	}
}
