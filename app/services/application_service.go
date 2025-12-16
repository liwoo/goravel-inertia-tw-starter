package services

import (
	"fmt"
	"math/rand"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
)

// ApplicationService implements business logic for applications using the builder pattern
type ApplicationService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
	cacheService *CacheService
}

// NewApplicationService creates a new Application service using the builder pattern
func NewApplicationService() *ApplicationService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Application]("applications", "id").
		WithSearchFields("sme", "registrant_name", "email", "phone", "sme_registration_number", "sme_tax_identification_number", "status").                                 // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "sme", "registrant_name", "email", "phone", "sme_registration_number", "sme_tax_identification_number", "status"). // Fields that can be used for sorting results
		WithFilterFields("sme", "registrant_name", "email", "phone", "sme_registration_number", "sme_tax_identification_number", "status").                                 // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                                                                                                   // Validation rules for create/update operations
			"sme":                           "required|string|max:255",
			"registrant_name":               "required|string|max:255",
			"email":                         "required|string|max:255",
			"phone":                         "required|string|max:255",
			"sme_registration_number":       "string|max:255",
			"sme_tax_identification_number": "string|max:255",
		}).
		WithDefaultSort("created_at", "DESC").            // Default sorting when none specified
		WithScopeFiltering("applications", "created_by"). // Enable permission-based filtering
		WithBeforeCreate(func(data map[string]interface{}) error {
			data["status"] = "Pending"
			return nil
		}).
		Build() // Returns a fully configured CrudServiceContract

	applicationServiceInstance := &ApplicationService{
		CrudServiceContract: service, // Set the embedded interface
		cacheService:        GetCacheService(),
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, applicationServiceInstance, "ApplicationService")

	return applicationServiceInstance
}

// GetFilterDefinitions returns filter definitions for the applications resource
// This method is OPTIONAL - only implement if you need custom filter UI components
// Uncomment and customize the implementation below if needed:
//
// func (s *ApplicationService) GetFilterDefinitions() []contracts.FilterDefinition {
// 	return []contracts.FilterDefinition{
// 		// String field example
// 		contracts.NewFilterDefinition(
// 			"field_name",           // Field name in database
// 			"Display Name",         // Human-readable label
// 			contracts.FilterTypeString,
// 			nil,                    // nil = use all string operators (equals, contains, starts_with, etc.)
// 		),
//
// 		// Enum field example (dropdown)
// 		contracts.NewFilterDefinition(
// 			"status",
// 			"Status",
// 			contracts.FilterTypeEnum,
// 			&[]string{"ACTIVE", "INACTIVE", "PENDING"}, // Available options
// 		),
//
// 		// Number field example
// 		contracts.NewFilterDefinition(
// 			"price",
// 			"Price",
// 			contracts.FilterTypeNumber,
// 			nil,                    // nil = use all number operators (equals, greater_than, less_than, etc.)
// 		),
//
// 		// Date field example
// 		contracts.NewFilterDefinition(
// 			"created_at",
// 			"Created Date",
// 			contracts.FilterTypeDate,
// 			nil,                    // nil = use all date operators (equals, before, after, between, etc.)
// 		),
//
// 		// Boolean field example
// 		contracts.NewFilterDefinition(
// 			"is_active",
// 			"Active Status",
// 			contracts.FilterTypeBoolean,
// 			nil,
// 		),
// 	}
// }

// Add domain-specific methods below this line

// ApproveApplication approves an application, creates a user account, and adds primary business owner to the selected SME
// Returns the created user data including the plain text password
func (s *ApplicationService) ApproveApplication(applicationID uint, smeID uint, approverUserID uint) (map[string]interface{}, error) {
	// 1. Fetch and validate the application
	var application models.Application
	if err := facades.Orm().Query().Where("id = ?", applicationID).First(&application); err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// Check if application is in Pending status
	if application.Status != "Pending" {
		return nil, fmt.Errorf("application is not in Pending status (current status: %s)", application.Status)
	}

	// 2. Validate that the SME exists
	var sme models.Sme
	if err := facades.Orm().Query().Where("id = ?", smeID).First(&sme); err != nil {
		return nil, fmt.Errorf("SME not found: %w", err)
	}

	// 3. Check if user with this email already exists
	var existingUser models.User
	err := facades.Orm().Query().Where("email = ?", application.Email).First(&existingUser)
	if err == nil && existingUser.ID > 0 {
		return nil, fmt.Errorf("user with email %s already exists", application.Email)
	}

	// 4. Create user account from application data
	// Generate a random password
	randomPassword := generateRandomPassword(12)
	hashedPassword, err := facades.Hash().Make(randomPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		Name:     application.RegistrantName,
		Email:    application.Email,
		Password: hashedPassword,
		IsActive: true,
		Role:     "User",
	}

	if err := facades.Orm().Query().Create(&user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	permissionService := auth.GetPermissionService()
	if err := permissionService.AssignRole(&user, "sme-user", nil); err != nil {
		facades.Log().Warning("Failed to assign sme-user role to new user", map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
			"error":   err.Error(),
		})
	} else {
		facades.Log().Info("Assigned sme-user role to new user", map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
		})
	}

	// Log the generated password (in production, this should be emailed to the user)
	facades.Log().Info("User created for approved application", map[string]interface{}{
		"user_id":        user.ID,
		"email":          user.Email,
		"application_id": applicationID,
		"password":       randomPassword, // TODO: Email this to the user instead of logging
	})

	// 5. Check if PrimaryBusinessOwner already exists for this SME (may have been created during SME creation)
	var existingOwner models.PrimaryBusinessOwner
	ownerErr := facades.Orm().Query().Where("sme_id = ?", smeID).First(&existingOwner)
	if ownerErr != nil || existingOwner.ID == 0 {
		// No existing owner, create one from application data
		primaryOwner := models.PrimaryBusinessOwner{
			FirstName:              application.FirstName,
			LastName:               application.LastName,
			OtherNames:             stringToPointer(application.OtherNames),
			Nationality:            application.Nationality,
			NationalIdNumber:       application.NationalIDNumber,
			DateOfBirth:            parseDate(application.DateOfBirth),
			Gender:                 application.Gender,
			EducationLevel:         application.EducationLevel,
			MalawianStatus:         application.MalawianStatus,
			HasSpecialNeeds:        application.HasSpecialNeeds,
			PhoneNumber:            application.Phone,
			LandlineNumber:         stringToPointer(application.LandlineNumber),
			Email:                  stringToPointer(application.Email),
			PhysicalAddress:        stringToPointer(application.PhysicalAddress),
			PostalAddress:          stringToPointer(application.PostalAddress),
			Region:                 stringToPointer(application.Region),
			District:               stringToPointer(application.District),
			TraditionalAuthority:   stringToPointer(application.TraditionalAuthority),
			AltContactName:         stringToPointer(application.AltContactName),
			AltContactRelationship: stringToPointer(application.AltContactRelationship),
			AltContactPhone:        stringToPointer(application.AltContactPhone),
			SmeID:                  int(smeID),
			CreatedBy:              intToPointer(int(approverUserID)),
		}

		if err := facades.Orm().Query().Create(&primaryOwner); err != nil {
			return nil, fmt.Errorf("failed to create primary business owner: %w", err)
		}
		facades.Log().Info("Created new primary business owner for SME", map[string]interface{}{
			"primary_owner_id": primaryOwner.ID,
			"sme_id":           smeID,
		})
	} else {
		facades.Log().Info("Primary business owner already exists for SME, skipping creation", map[string]interface{}{
			"existing_owner_id": existingOwner.ID,
			"sme_id":            smeID,
		})
	}

	updateData := map[string]interface{}{
		"status": "Approved",
	}
	// 6. Update application status to Approved
	if _, err := s.Update(applicationID, updateData); err != nil {
		return nil, fmt.Errorf("failed to update application status: %w", err)
	}

	facades.Log().Info("Application approved successfully", map[string]interface{}{
		"application_id":   applicationID,
		"sme_id":           smeID,
		"user_id":          user.ID,
		"approver_user_id": approverUserID,
	})

	return map[string]interface{}{
		"user": map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"password": randomPassword,
		},
	}, nil
}

// RejectApplication rejects an application by updating its status
func (s *ApplicationService) RejectApplication(applicationID uint, rejectorUserID uint) error {
	// Fetch and validate the application
	var application models.Application
	if err := facades.Orm().Query().Where("id = ?", applicationID).First(&application); err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Check if application is in Pending status
	if application.Status != "Pending" {
		return fmt.Errorf("application is not in Pending status (current status: %s)", application.Status)
	}

	// Update application status to Rejected
	if _, err := facades.Orm().Query().Model(&application).Update("status", "Rejected"); err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}

	facades.Log().Info("Application rejected", map[string]interface{}{
		"application_id":   applicationID,
		"rejector_user_id": rejectorUserID,
	})

	return nil
}

// Helper functions

// generateRandomPassword generates a random password of the specified length
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}

// stringToPointer converts a string to a pointer
func stringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// intToPointer converts an int to a pointer
func intToPointer(i int) *int {
	return &i
}

// parseDate parses a date string and returns a carbon.DateTime
func parseDate(dateStr string) carbon.DateTime {
	// Try to parse the date string
	// The application stores dates as strings, so we need to parse them
	parsed := carbon.Parse(dateStr)
	if parsed.Error != nil {
		// Return zero time if parsing fails
		return carbon.DateTime{}
	}
	return *carbon.NewDateTime(parsed)
}
