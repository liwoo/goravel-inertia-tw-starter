package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
)

// FormalisationAmendmentData represents the data structure for formalisation amendment requests
type FormalisationAmendmentData struct {
	// SME-level formalisation fields
	RegistrationNumber      *string `json:"registration_number"`
	TaxIdentificationNumber *string `json:"tax_identification_number"`

	// BusinessFormalisation fields
	HasBankAccount         *bool    `json:"has_bank_account"`
	HasTaxClarification    *bool    `json:"has_tax_clarification"`
	IsRegisteredForVat     *bool    `json:"is_registered_for_vat"`
	IsMemberOfAssociation  *bool    `json:"is_member_of_association"`
	IsAffiliated           *bool    `json:"is_affiliated"`
	HasExportLicense       *bool    `json:"has_export_license"`
	HasAccessedBds         *bool    `json:"has_accessed_bds"`
	AnnualTurnover         *float64 `json:"annual_turnover"`
	EstimatedValueOfAssets *float64 `json:"estimated_value_of_assets"`

	// Metadata
	ChangeReason string `json:"change_reason"`
}

// AmendmentApplicationData wraps current and proposed values for comparison
type AmendmentApplicationData struct {
	Current  FormalisationAmendmentData `json:"current"`
	Proposed FormalisationAmendmentData `json:"proposed"`
	SmeName  string                     `json:"sme_name"`
}

// ApplicationService implements business logic for applications using the builder pattern
type ApplicationService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
	cacheService        *CacheService
	notificationService *NotificationService
}

// NewApplicationService creates a new Application service using the builder pattern
func NewApplicationService() *ApplicationService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Application]("applications", "id").
		WithSearchFields("sme", "registrant_name", "email", "phone", "sme_registration_number", "sme_tax_identification_number", "status", "type").                                 // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "sme", "registrant_name", "email", "phone", "sme_registration_number", "sme_tax_identification_number", "status", "type"). // Fields that can be used for sorting results
		WithFilterFields("sme", "registrant_name", "email", "phone", "sme_registration_number", "sme_tax_identification_number", "status", "type", "sme_id").                       // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                                                                                                                 // Validation rules for create/update operations
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
			// Set default type if not provided
			if _, exists := data["type"]; !exists {
				data["type"] = models.ApplicationTypeSignup
			}
			return nil
		}).
		WithRelations("Sme"). // Load SME relation for amendment applications
		Build()               // Returns a fully configured CrudServiceContract

	applicationServiceInstance := &ApplicationService{
		CrudServiceContract: service, // Set the embedded interface
		cacheService:        GetCacheService(),
		notificationService: NewNotificationService(),
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

// HasPendingAmendment checks if an SME has a pending formalisation amendment application
func (s *ApplicationService) HasPendingAmendment(smeID uint) (bool, error) {
	count, err := facades.Orm().Query().Model(&models.Application{}).
		Where("sme_id = ?", smeID).
		Where("type = ?", models.ApplicationTypeAmendFormalisation).
		Where("status = ?", models.ApplicationStatusPending).
		Count()
	if err != nil {
		return false, fmt.Errorf("failed to check pending amendments: %w", err)
	}
	return count > 0, nil
}

// CreateAmendmentApplication creates a new formalisation amendment application
func (s *ApplicationService) CreateAmendmentApplication(smeID uint, submitterUserID uint, proposed FormalisationAmendmentData) (*models.Application, error) {
	// Check for existing pending amendment
	hasPending, err := s.HasPendingAmendment(smeID)
	if err != nil {
		return nil, err
	}
	if hasPending {
		return nil, fmt.Errorf("you already have a pending formalisation amendment request")
	}
	// 1. Fetch the SME and its formalisation data
	var sme models.Sme
	if err := facades.Orm().Query().With("BusinessFormalisation").Where("id = ?", smeID).First(&sme); err != nil {
		return nil, fmt.Errorf("SME not found: %w", err)
	}

	// 2. Build current data from SME
	current := FormalisationAmendmentData{
		RegistrationNumber:      sme.RegistrationNumber,
		TaxIdentificationNumber: sme.TaxIdentificationNumber,
	}

	// Add formalisation fields if they exist
	if sme.BusinessFormalisation != nil {
		bf := sme.BusinessFormalisation
		current.HasBankAccount = &bf.HasBankAccount
		current.HasTaxClarification = &bf.HasTaxClarification
		current.IsRegisteredForVat = &bf.IsRegisteredForVat
		current.IsMemberOfAssociation = &bf.IsMemberOfAssociation
		current.IsAffiliated = &bf.IsAffiliated
		current.HasExportLicense = &bf.HasExportLicense
		current.HasAccessedBds = &bf.HasAccessedBds
		current.AnnualTurnover = &bf.AnnualTurnover
		current.EstimatedValueOfAssets = &bf.EstimatedValueOfAssets
	}

	// 3. Build amendment data JSON
	amendmentData := AmendmentApplicationData{
		Current:  current,
		Proposed: proposed,
		SmeName:  sme.Name,
	}

	dataJSON, err := json.Marshal(amendmentData)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize amendment data: %w", err)
	}
	dataStr := string(dataJSON)

	// 4. Create the application
	application := &models.Application{
		Type:           models.ApplicationTypeAmendFormalisation,
		SmeID:          &smeID,
		Data:           &dataStr,
		Status:         models.ApplicationStatusPending,
		SME:            sme.Name,
		RegistrantName: "", // Not applicable for amendment
		Email:          sme.ContactEmail,
		Phone:          sme.ContactPhone,
	}

	// Set audit fields
	submitterID := submitterUserID
	application.CreatedBy = &submitterID
	application.UpdatedBy = &submitterID

	if err := facades.Orm().Query().Create(application); err != nil {
		return nil, fmt.Errorf("failed to create amendment application: %w", err)
	}

	facades.Log().Info("Amendment application created", map[string]interface{}{
		"application_id": application.ID,
		"sme_id":         smeID,
		"submitter_id":   submitterUserID,
	})

	// 5. Send notification to submitter confirming submission
	s.notificationService.CreateNotification(
		submitterUserID,
		"Amendment Request Submitted",
		fmt.Sprintf("Your formalisation amendment request for %s has been submitted and is pending review.", sme.Name),
		"application",
		nil,
		strPtr("application"),
		&application.ID,
		"normal",
		nil,
		"{}",
	)

	// Notification count update for admins will happen automatically via pending_applications count

	return application, nil
}

// ApproveApplication approves an application based on its type
func (s *ApplicationService) ApproveApplication(applicationID uint, smeID uint, approverUserID uint) (map[string]interface{}, error) {
	// Fetch and validate the application
	var application models.Application
	if err := facades.Orm().Query().Where("id = ?", applicationID).First(&application); err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// Check if application is in Pending status
	if application.Status != models.ApplicationStatusPending {
		return nil, fmt.Errorf("application is not in Pending status (current status: %s)", application.Status)
	}

	// Route to appropriate handler based on type
	switch application.Type {
	case models.ApplicationTypeAmendFormalisation:
		return s.approveAmendmentApplication(applicationID, application, approverUserID)
	default: // signup
		return s.approveSignupApplication(applicationID, smeID, approverUserID)
	}
}

// approveSignupApplication handles approval of signup applications
func (s *ApplicationService) approveSignupApplication(applicationID uint, smeID uint, approverUserID uint) (map[string]interface{}, error) {
	// Fetch the application
	var application models.Application
	if err := facades.Orm().Query().Where("id = ?", applicationID).First(&application); err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// Validate that the SME exists
	var sme models.Sme
	if err := facades.Orm().Query().Where("id = ?", smeID).First(&sme); err != nil {
		return nil, fmt.Errorf("SME not found: %w", err)
	}

	// Check if user with this email already exists
	var existingUser models.User
	err := facades.Orm().Query().Where("email = ?", application.Email).First(&existingUser)
	if err == nil && existingUser.ID > 0 {
		return nil, fmt.Errorf("user with email %s already exists", application.Email)
	}

	// Create user account from application data
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

	facades.Log().Info("User created for approved application", map[string]interface{}{
		"user_id":        user.ID,
		"email":          user.Email,
		"application_id": applicationID,
		"password":       randomPassword,
	})

	// Check if PrimaryBusinessOwner already exists for this SME
	var existingOwner models.PrimaryBusinessOwner
	ownerErr := facades.Orm().Query().Where("sme_id = ?", smeID).First(&existingOwner)
	if ownerErr != nil || existingOwner.ID == 0 {
		primaryOwner := models.PrimaryBusinessOwner{
			FirstName:              derefString(application.FirstName),
			LastName:               derefString(application.LastName),
			OtherNames:             application.OtherNames,
			Nationality:            derefString(application.Nationality),
			NationalIdNumber:       derefString(application.NationalIDNumber),
			DateOfBirth:            parseDateFromPointer(application.DateOfBirth),
			Gender:                 derefString(application.Gender),
			EducationLevel:         derefString(application.EducationLevel),
			MalawianStatus:         derefString(application.MalawianStatus),
			HasSpecialNeeds:        application.HasSpecialNeeds,
			PhoneNumber:            application.Phone,
			LandlineNumber:         application.LandlineNumber,
			Email:                  stringToPointer(application.Email),
			PhysicalAddress:        application.PhysicalAddress,
			PostalAddress:          application.PostalAddress,
			Region:                 application.Region,
			District:               application.District,
			TraditionalAuthority:   application.TraditionalAuthority,
			AltContactName:         application.AltContactName,
			AltContactRelationship: application.AltContactRelationship,
			AltContactPhone:        application.AltContactPhone,
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

	// Update application status to Approved
	if _, err := s.Update(applicationID, map[string]interface{}{"status": models.ApplicationStatusApproved}); err != nil {
		return nil, fmt.Errorf("failed to update application status: %w", err)
	}

	facades.Log().Info("Signup application approved successfully", map[string]interface{}{
		"application_id":   applicationID,
		"sme_id":           smeID,
		"user_id":          user.ID,
		"approver_user_id": approverUserID,
	})

	// Send notification to the new user
	s.notificationService.CreateNotification(
		user.ID,
		"Welcome to SMEDI",
		"Your application has been approved. Welcome to the SME Database!",
		"application_approved",
		&approverUserID,
		strPtr("application"),
		&applicationID,
		"normal",
		nil,
		"{}",
	)

	return map[string]interface{}{
		"user": map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"password": randomPassword,
		},
	}, nil
}

// approveAmendmentApplication handles approval of formalisation amendment applications
func (s *ApplicationService) approveAmendmentApplication(applicationID uint, application models.Application, approverUserID uint) (map[string]interface{}, error) {
	if application.SmeID == nil || application.Data == nil {
		return nil, fmt.Errorf("invalid amendment application: missing SME ID or data")
	}

	// Parse the amendment data
	var amendmentData AmendmentApplicationData
	if err := json.Unmarshal([]byte(*application.Data), &amendmentData); err != nil {
		return nil, fmt.Errorf("failed to parse amendment data: %w", err)
	}

	proposed := amendmentData.Proposed
	smeID := *application.SmeID

	// Fetch the SME with formalisation
	var sme models.Sme
	if err := facades.Orm().Query().With("BusinessFormalisation").Where("id = ?", smeID).First(&sme); err != nil {
		return nil, fmt.Errorf("SME not found: %w", err)
	}

	// Update SME-level fields if provided
	smeUpdates := make(map[string]interface{})
	if proposed.RegistrationNumber != nil {
		smeUpdates["registration_number"] = *proposed.RegistrationNumber
	}
	if proposed.TaxIdentificationNumber != nil {
		smeUpdates["tax_identification_number"] = *proposed.TaxIdentificationNumber
	}

	if len(smeUpdates) > 0 {
		smeUpdates["updated_by"] = approverUserID
		if _, err := facades.Orm().Query().Model(&models.Sme{}).Where("id = ?", smeID).Update(smeUpdates); err != nil {
			return nil, fmt.Errorf("failed to update SME: %w", err)
		}
	}

	// Update or create BusinessFormalisation
	if sme.BusinessFormalisation != nil {
		// Update existing formalisation
		bfUpdates := make(map[string]interface{})
		if proposed.HasBankAccount != nil {
			bfUpdates["has_bank_account"] = *proposed.HasBankAccount
		}
		if proposed.HasTaxClarification != nil {
			bfUpdates["has_tax_clarification"] = *proposed.HasTaxClarification
		}
		if proposed.IsRegisteredForVat != nil {
			bfUpdates["is_registered_for_vat"] = *proposed.IsRegisteredForVat
		}
		if proposed.IsMemberOfAssociation != nil {
			bfUpdates["is_member_of_association"] = *proposed.IsMemberOfAssociation
		}
		if proposed.IsAffiliated != nil {
			bfUpdates["is_affiliated"] = *proposed.IsAffiliated
		}
		if proposed.HasExportLicense != nil {
			bfUpdates["has_export_license"] = *proposed.HasExportLicense
		}
		if proposed.HasAccessedBds != nil {
			bfUpdates["has_accessed_bds"] = *proposed.HasAccessedBds
		}
		if proposed.AnnualTurnover != nil {
			bfUpdates["annual_turnover"] = *proposed.AnnualTurnover
		}
		if proposed.EstimatedValueOfAssets != nil {
			bfUpdates["estimated_value_of_assets"] = *proposed.EstimatedValueOfAssets
		}

		if len(bfUpdates) > 0 {
			bfUpdates["updated_by"] = approverUserID
			if _, err := facades.Orm().Query().Model(&models.BusinessFormalisation{}).Where("sme_id = ?", smeID).Update(bfUpdates); err != nil {
				return nil, fmt.Errorf("failed to update business formalisation: %w", err)
			}

			// Recalculate formalisation score
			smeService := NewSmeService()
			if _, err := smeService.CalculateFormalisationScore(smeID); err != nil {
				facades.Log().Warning("Failed to recalculate formalisation score", map[string]interface{}{
					"sme_id": smeID,
					"error":  err.Error(),
				})
			}
		}
	}

	// Update application status to Approved
	if _, err := s.Update(applicationID, map[string]interface{}{"status": models.ApplicationStatusApproved}); err != nil {
		return nil, fmt.Errorf("failed to update application status: %w", err)
	}

	facades.Log().Info("Amendment application approved successfully", map[string]interface{}{
		"application_id":   applicationID,
		"sme_id":           smeID,
		"approver_user_id": approverUserID,
	})

	// Send notification to the SME user who created the application
	if application.CreatedBy != nil {
		s.notificationService.CreateNotification(
			uint(*application.CreatedBy),
			"Amendment Approved",
			fmt.Sprintf("Your formalisation amendment for %s has been approved. Changes have been applied.", amendmentData.SmeName),
			"application_approved",
			&approverUserID,
			strPtr("application"),
			&applicationID,
			"normal",
			nil,
			"{}",
		)
	}

	return map[string]interface{}{
		"sme_id":  smeID,
		"message": "Formalisation amendment approved and applied",
	}, nil
}

// RejectApplication rejects an application by updating its status with an optional reason
func (s *ApplicationService) RejectApplication(applicationID uint, rejectorUserID uint, reason string) error {
	// Fetch and validate the application
	var application models.Application
	if err := facades.Orm().Query().Where("id = ?", applicationID).First(&application); err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Check if application is in Pending status
	if application.Status != models.ApplicationStatusPending {
		return fmt.Errorf("application is not in Pending status (current status: %s)", application.Status)
	}

	// Build update data
	updateData := map[string]interface{}{
		"status": models.ApplicationStatusRejected,
	}
	if reason != "" {
		updateData["rejection_reason"] = reason
	}

	// Update application status to Rejected
	if _, err := facades.Orm().Query().Model(&application).Update(updateData); err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}

	facades.Log().Info("Application rejected", map[string]interface{}{
		"application_id":   applicationID,
		"rejector_user_id": rejectorUserID,
		"type":             application.Type,
		"reason":           reason,
	})

	// Send notification to the applicant
	var notifyUserID uint
	var smeName string

	if application.Type == models.ApplicationTypeAmendFormalisation {
		// For amendment applications, notify the creator
		if application.CreatedBy != nil {
			notifyUserID = uint(*application.CreatedBy)
		}
		// Try to get SME name from data
		if application.Data != nil {
			var amendmentData AmendmentApplicationData
			if err := json.Unmarshal([]byte(*application.Data), &amendmentData); err == nil {
				smeName = amendmentData.SmeName
			}
		}
	} else {
		// For signup applications, we don't have a user yet - use email notification instead
		// For now, just log it
		facades.Log().Info("Signup application rejected - email notification should be sent", map[string]interface{}{
			"email":  application.Email,
			"reason": reason,
		})
	}

	// Send notification if we have a user ID
	if notifyUserID > 0 {
		message := "Your application has been rejected."
		if reason != "" {
			message = fmt.Sprintf("Your application has been rejected. Reason: %s", reason)
		}
		title := "Application Rejected"
		if application.Type == models.ApplicationTypeAmendFormalisation && smeName != "" {
			title = fmt.Sprintf("Amendment Rejected - %s", smeName)
		}

		s.notificationService.CreateNotification(
			notifyUserID,
			title,
			message,
			"application_rejected",
			&rejectorUserID,
			strPtr("application"),
			&applicationID,
			"high", // High priority for rejections
			nil,
			"{}",
		)
	}

	return nil
}

// strPtr is a helper function to convert string to pointer
func strPtr(s string) *string {
	return &s
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

// derefString safely dereferences a string pointer, returning empty string if nil
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// parseDateFromPointer parses a date string pointer and returns a carbon.DateTime
func parseDateFromPointer(dateStr *string) carbon.DateTime {
	if dateStr == nil || *dateStr == "" {
		return carbon.DateTime{}
	}
	parsed := carbon.Parse(*dateStr)
	if parsed.Error != nil {
		return carbon.DateTime{}
	}
	return *carbon.NewDateTime(parsed)
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
