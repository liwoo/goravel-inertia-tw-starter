package services

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/events"
	"smedi-sme-db/app/models"
)

// AdditionalBusinessMemberService implements business logic for additional business members using the builder pattern
type AdditionalBusinessMemberService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// NewAdditionalBusinessMemberService creates a new AdditionalBusinessMember service using the builder pattern
func NewAdditionalBusinessMemberService() *AdditionalBusinessMemberService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.AdditionalBusinessMember]("additional_business_members", "id").
		WithSearchFields("first_name", "last_name", "other_names", "national_id_number", "phone_number", "email").
		WithSortFields("id", "created_at", "updated_at", "first_name", "last_name", "date_of_birth").
		WithFilterFields("nationality", "gender", "is_intern", "is_part_time", "sme_id").
		WithValidationRules(map[string]interface{}{
			"first_name":         "required|string|max:100",
			"last_name":          "required|string|max:100",
			"other_names":        "string|max:100",
			"gender":             "string",
			"nationality":        "required|string|max:100",
			"national_id_number": "required|string|max:50",
			"date_of_birth":      "date",
			"email":              "email|max:100",
			"phone_number":       "required|string|max:20",
			"is_intern":          "boolean",
			"is_part_time":       "boolean",
			"sme_id":             "required|numeric",
		}).
		WithDefaultSort("created_at", "DESC").
		WithScopeFiltering("additional_business_members", "created_by").
		WithSoftDeletes().
		Build()

	additionalBusinessMemberServiceInstance := &AdditionalBusinessMemberService{
		CrudServiceContract: service,
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, additionalBusinessMemberServiceInstance, "AdditionalBusinessMemberService")

	return additionalBusinessMemberServiceInstance
}

// Override GetColumnMapping to include AdditionalBusinessMember-specific mappings
func (s *AdditionalBusinessMemberService) GetColumnMapping() map[string]string {
	mapping := s.CrudServiceContract.GetColumnMapping()
	// Add AdditionalBusinessMember-specific camelCase to snake_case mappings
	mapping["firstName"] = "first_name"
	mapping["lastName"] = "last_name"
	mapping["otherNames"] = "other_names"
	mapping["gender"] = "gender"
	mapping["nationalIdNumber"] = "national_id_number"
	mapping["dateOfBirth"] = "date_of_birth"
	mapping["phoneNumber"] = "phone_number"
	mapping["isIntern"] = "is_intern"
	mapping["isPartTime"] = "is_part_time"
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

// GetFilterDefinitions returns filter definitions for the AdditionalBusinessMember resource
func (s *AdditionalBusinessMemberService) GetFilterDefinitions() []contracts.FilterDefinition {
	return []contracts.FilterDefinition{
		contracts.NewFilterDefinition(
			"first_name",
			"First Name",
			contracts.FilterTypeString,
			nil,
		),
		contracts.NewFilterDefinition(
			"last_name",
			"Last Name",
			contracts.FilterTypeString,
			nil,
		),
		contracts.NewFilterDefinition(
			"nationality",
			"Nationality",
			contracts.FilterTypeString,
			nil,
		),
		contracts.NewFilterDefinition(
			"is_intern",
			"Is Intern",
			contracts.FilterTypeBoolean,
			nil,
		),
		contracts.NewFilterDefinition(
			"is_part_time",
			"Is Part Time",
			contracts.FilterTypeBoolean,
			nil,
		),
		contracts.NewFilterDefinition(
			"date_of_birth",
			"Date of Birth",
			contracts.FilterTypeDate,
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

// Add domain-specific methods below this line

// Create overrides the base Create method to fire AdditionalBusinessMemberCreated event
// and recalculate formalisation score
func (s *AdditionalBusinessMemberService) Create(data map[string]interface{}) (interface{}, error) {
	// Call the base Create method
	result, err := s.CrudServiceContract.Create(data)
	if err != nil {
		return nil, err
	}

	// Fire the AdditionalBusinessMemberCreated event and recalculate formalisation score
	if member, ok := result.(*models.AdditionalBusinessMember); ok {
		if err := facades.Event().Job(&events.AdditionalBusinessMemberCreated{}, []event.Arg{
			{Type: "int", Value: member.SmeId},
		}).Dispatch(); err != nil {
			facades.Log().Warningf("Failed to dispatch AdditionalBusinessMemberCreated event: %v", err)
		}

		// Recalculate formalisation score after creating an additional business member
		if member.SmeId > 0 {
			smeService := NewSmeService()
			_, err := smeService.CalculateFormalisationScore(uint(member.SmeId))
			if err != nil {
				facades.Log().Warningf("Failed to recalculate formalisation score after additional member create: %v", err)
			}
		}
	}

	return result, nil
}

// Update overrides the base Update method to fire AdditionalBusinessMemberUpdated event
// and recalculate formalisation score
func (s *AdditionalBusinessMemberService) Update(id uint, data map[string]interface{}) (interface{}, error) {
	// Call the base Update method
	result, err := s.CrudServiceContract.Update(id, data)
	if err != nil {
		return nil, err
	}

	// Fire the AdditionalBusinessMemberUpdated event and recalculate formalisation score
	if member, ok := result.(*models.AdditionalBusinessMember); ok {
		if err := facades.Event().Job(&events.AdditionalBusinessMemberUpdated{}, []event.Arg{
			{Type: "int", Value: member.SmeId},
		}).Dispatch(); err != nil {
			facades.Log().Warningf("Failed to dispatch AdditionalBusinessMemberUpdated event: %v", err)
		}

		// Recalculate formalisation score after updating an additional business member
		if member.SmeId > 0 {
			smeService := NewSmeService()
			_, err := smeService.CalculateFormalisationScore(uint(member.SmeId))
			if err != nil {
				facades.Log().Warningf("Failed to recalculate formalisation score after additional member update: %v", err)
			}
		}
	}

	return result, nil
}

// Delete overrides the base Delete method to fire AdditionalBusinessMemberDeleted event
// and recalculate formalisation score
func (s *AdditionalBusinessMemberService) Delete(id uint) error {
	// Get the member before delete to capture sme_id
	memberInterface, err := s.CrudServiceContract.GetByID(id)
	if err != nil {
		return err
	}
	member, _ := memberInterface.(*models.AdditionalBusinessMember)

	// Call the base Delete method
	err = s.CrudServiceContract.Delete(id)
	if err != nil {
		return err
	}

	// Fire the AdditionalBusinessMemberDeleted event and recalculate formalisation score
	if member != nil {
		if err := facades.Event().Job(&events.AdditionalBusinessMemberDeleted{}, []event.Arg{
			{Type: "int", Value: member.SmeId},
		}).Dispatch(); err != nil {
			facades.Log().Warningf("Failed to dispatch AdditionalBusinessMemberDeleted event: %v", err)
		}

		// Recalculate formalisation score after deleting an additional business member
		if member.SmeId > 0 {
			smeService := NewSmeService()
			_, err := smeService.CalculateFormalisationScore(uint(member.SmeId))
			if err != nil {
				facades.Log().Warningf("Failed to recalculate formalisation score after additional member delete: %v", err)
			}
		}
	}

	return nil
}
