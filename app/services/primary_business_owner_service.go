package services

import (
	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
)

// PrimaryBusinessOwnerService implements business logic for primary business owners using the builder pattern
type PrimaryBusinessOwnerService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
	// Store sme_id for use in after delete hook
	pendingDeleteSmeID int
}

// NewPrimaryBusinessOwnerService creates a new PrimaryBusinessOwner service using the builder pattern
func NewPrimaryBusinessOwnerService() *PrimaryBusinessOwnerService {
	// Create service instance first to allow referencing in hooks
	primaryBusinessOwnerServiceInstance := &PrimaryBusinessOwnerService{}

	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.PrimaryBusinessOwner]("primary_business_owner", "id").
		WithSearchFields("first_name", "last_name", "other_names", "national_id_number", "phone_number", "email").
		WithSortFields("id", "created_at", "updated_at", "first_name", "last_name", "date_of_birth").
		WithFilterFields("gender", "education_level", "malawian_status", "nationality", "region", "district", "has_special_needs", "sme_id").
		WithValidationRules(map[string]interface{}{
			"first_name":               "required|string|max:100",
			"last_name":                "required|string|max:100",
			"other_names":              "string|max:100",
			"nationality":              "required|string|max:100",
			"national_id_number":       "required|string|max:50",
			"date_of_birth":            "required|date",
			"gender":                   "required|string|max:20",
			"education_level":          "required|string|max:50",
			"malawian_status":          "required|string|max:50",
			"has_special_needs":        "boolean",
			"phone_number":             "required|string|max:20",
			"landline_number":          "string|max:20",
			"email":                    "email|max:100",
			"physical_address":         "string|max:255",
			"postal_address":           "string|max:255",
			"region":                   "string|max:100",
			"district":                 "string|max:100",
			"traditional_authority":    "string|max:100",
			"alt_contact_name":         "string|max:100",
			"alt_contact_relationship": "string|max:50",
			"alt_contact_phone":        "string|max:20",
			"sme_id":                   "required|numeric",
		}).
		WithDefaultSort("created_at", "DESC").
		WithScopeFiltering("primary_business_owner", "created_by").
		WithSoftDeletes().
		WithAfterCreate(func(model *models.PrimaryBusinessOwner) error {
			// Recalculate formalisation score and classification after creating a primary business owner
			if model != nil && model.SmeID > 0 {
				smeService := NewSmeService()
				_, err := smeService.CalculateFormalisationScore(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate formalisation score after primary owner create: %v", err)
				}
				// Recalculate classification (primary owner counts as employee)
				_, err = smeService.CalculateClassification(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate classification after primary owner create: %v", err)
				}
			}
			return nil
		}).
		WithAfterUpdate(func(model *models.PrimaryBusinessOwner) error {
			// Recalculate formalisation score and classification after updating a primary business owner
			if model != nil && model.SmeID > 0 {
				smeService := NewSmeService()
				_, err := smeService.CalculateFormalisationScore(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate formalisation score after primary owner update: %v", err)
				}
				// Recalculate classification (primary owner counts as employee)
				_, err = smeService.CalculateClassification(uint(model.SmeID))
				if err != nil {
					facades.Log().Warningf("Failed to recalculate classification after primary owner update: %v", err)
				}
			}
			return nil
		}).
		WithBeforeDelete(func(id uint) error {
			// Fetch the primary business owner to get the sme_id before deletion
			var owner models.PrimaryBusinessOwner
			err := facades.Orm().Query().Where("id = ?", id).First(&owner)
			if err == nil && owner.SmeID > 0 {
				primaryBusinessOwnerServiceInstance.pendingDeleteSmeID = owner.SmeID
			}
			return nil
		}).
		WithAfterDelete(func(id uint) error {
			// Recalculate formalisation score and classification after deleting a primary business owner
			if primaryBusinessOwnerServiceInstance.pendingDeleteSmeID > 0 {
				smeService := NewSmeService()
				smeID := uint(primaryBusinessOwnerServiceInstance.pendingDeleteSmeID)
				_, err := smeService.CalculateFormalisationScore(smeID)
				if err != nil {
					facades.Log().Warningf("Failed to recalculate formalisation score after primary owner delete: %v", err)
				}
				// Recalculate classification (primary owner counts as employee)
				_, err = smeService.CalculateClassification(smeID)
				if err != nil {
					facades.Log().Warningf("Failed to recalculate classification after primary owner delete: %v", err)
				}
				primaryBusinessOwnerServiceInstance.pendingDeleteSmeID = 0 // Reset after use
			}
			return nil
		}).
		Build()

	primaryBusinessOwnerServiceInstance.CrudServiceContract = service

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, primaryBusinessOwnerServiceInstance, "PrimaryBusinessOwnerService")

	return primaryBusinessOwnerServiceInstance
}

// Override GetColumnMapping to include PrimaryBusinessOwner-specific mappings
func (s *PrimaryBusinessOwnerService) GetColumnMapping() map[string]string {
	mapping := s.CrudServiceContract.GetColumnMapping()
	// Add PrimaryBusinessOwner-specific camelCase to snake_case mappings
	mapping["firstName"] = "first_name"
	mapping["lastName"] = "last_name"
	mapping["otherNames"] = "other_names"
	mapping["nationalIdNumber"] = "national_id_number"
	mapping["dateOfBirth"] = "date_of_birth"
	mapping["educationLevel"] = "education_level"
	mapping["malawianStatus"] = "malawian_status"
	mapping["hasSpecialNeeds"] = "has_special_needs"
	mapping["phoneNumber"] = "phone_number"
	mapping["landlineNumber"] = "landline_number"
	mapping["physicalAddress"] = "physical_address"
	mapping["postalAddress"] = "postal_address"
	mapping["traditionalAuthority"] = "traditional_authority"
	mapping["altContactName"] = "alt_contact_name"
	mapping["altContactRelationship"] = "alt_contact_relationship"
	mapping["altContactPhone"] = "alt_contact_phone"
	mapping["smeId"] = "sme_id"
	mapping["createdAt"] = "created_at"
	mapping["updatedAt"] = "updated_at"
	mapping["createdBy"] = "created_by"
	mapping["updatedBy"] = "updated_by"
	mapping["deletedBy"] = "deleted_by"
	mapping["ipAddress"] = "ip_address"
	mapping["userAgent"] = "user_agent"
	return mapping
}
