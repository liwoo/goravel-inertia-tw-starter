package services

import (
	"smedi-sme-db/app/contracts"
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
		WithFilterFields("nationality", "is_intern", "is_part_time", "sme_id").
		WithValidationRules(map[string]interface{}{
			"first_name":         "required|string|max:100",
			"last_name":          "required|string|max:100",
			"other_names":        "string|max:100",
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
