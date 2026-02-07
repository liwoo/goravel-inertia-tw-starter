package services

import (
	"fmt"
	"math/rand"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"

	"github.com/goravel/framework/facades"
)

// ApplicationService implements business logic for applications using the builder pattern
type ApplicationService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
	cacheService                  *CacheService
	notificationService           *NotificationService
}

// NewApplicationService creates a new Application service using the builder pattern
func NewApplicationService() *ApplicationService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Application]("applications", "id").
		WithSearchFields("sme", "registrant_name", "email", "phone", "status", "type").                                 // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "sme", "registrant_name", "email", "phone", "status", "type"). // Fields that can be used for sorting results
		WithFilterFields("sme", "registrant_name", "email", "phone", "status", "type").                                 // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                                                     // Validation rules for create/update operations
			"sme":             "required|string|max:255",
			"registrant_name": "required|string|max:255",
			"email":           "required|string|max:255",
			"phone":           "required|string|max:255",
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
		Build() // Returns a fully configured CrudServiceContract

	applicationServiceInstance := &ApplicationService{
		CrudServiceContract: service, // Set the embedded interface
		cacheService:        GetCacheService(),
		notificationService: NewNotificationService(),
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, applicationServiceInstance, "ApplicationService")

	return applicationServiceInstance
}

// Add domain-specific methods below this line

// ApproveApplication approves an application by updating its status and creating a user
func (s *ApplicationService) ApproveApplication(applicationID uint, approverUserID uint) (map[string]interface{}, error) {
	// Fetch and validate the application
	var application models.Application
	if err := facades.Orm().Query().Where("id = ?", applicationID).First(&application); err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// Check if application is in Pending status
	if application.Status != models.ApplicationStatusPending {
		return nil, fmt.Errorf("application is not in Pending status (current status: %s)", application.Status)
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
	if err := permissionService.AssignRole(&user, "member", nil); err != nil {
		facades.Log().Warning("Failed to assign member role to new user", map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
			"error":   err.Error(),
		})
	} else {
		facades.Log().Info("Assigned member role to new user", map[string]interface{}{
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

	// Update application status to Approved
	if _, err := s.Update(applicationID, map[string]interface{}{"status": models.ApplicationStatusApproved}); err != nil {
		return nil, fmt.Errorf("failed to update application status: %w", err)
	}

	facades.Log().Info("Application approved successfully", map[string]interface{}{
		"application_id":   applicationID,
		"user_id":          user.ID,
		"approver_user_id": approverUserID,
	})

	// Send notification to the new user
	s.notificationService.CreateNotification(
		user.ID,
		"Welcome",
		"Your application has been approved. Welcome!",
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

	// For signup applications, we don't have a user yet - use email notification instead
	// For now, just log it
	facades.Log().Info("Application rejected - email notification should be sent", map[string]interface{}{
		"email":  application.Email,
		"reason": reason,
	})

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

// intToPointer converts an int to a pointer
func intToPointer(i int) *int {
	return &i
}
