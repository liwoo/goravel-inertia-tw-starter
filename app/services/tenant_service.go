package services

import (
	"fmt"

	"books-database/app/contracts"
	"books-database/app/models"

	"github.com/goravel/framework/facades"
)

type TenantService struct {
	contracts.CrudServiceContract
}

func NewTenantService() *TenantService {
	service := contracts.NewServiceBuilder[models.Tenant]("tenants", "id").
		WithSearchFields("name", "slug", "description").
		WithSortFields("id", "name", "slug", "created_at", "updated_at").
		WithFilterFields("is_active", "is_main").
		WithValidationRules(map[string]interface{}{
			"name": "required|string|max:255",
			"slug": "required|string|max:100",
		}).
		WithRelations("Creator", "Updater").
		WithDefaultSort("created_at", "DESC").
		WithSoftDeletes().
		Build()

	svc := &TenantService{CrudServiceContract: service}
	contracts.SetActualServiceHelper(service, svc, "TenantService")
	return svc
}

// AssignUserToTenant adds a user to a tenant
func (s *TenantService) AssignUserToTenant(userID, tenantID uint, roleID *uint) error {
	// Check if already assigned
	var existing models.UserTenant
	facades.Orm().Query().
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		First(&existing)
	if existing.ID > 0 {
		// Reactivate if inactive
		if !existing.IsActive {
			_, err := facades.Orm().Query().Model(&existing).Update("is_active", true)
			return err
		}
		return fmt.Errorf("user is already assigned to this tenant")
	}

	ut := models.UserTenant{
		UserID:   userID,
		TenantID: tenantID,
		RoleID:   roleID,
		IsActive: true,
	}
	return facades.Orm().Query().Create(&ut)
	// Note: Goravel's Create returns error only
}

// RemoveUserFromTenant deactivates a user's tenant membership
func (s *TenantService) RemoveUserFromTenant(userID, tenantID uint) error {
	_, err := facades.Orm().Query().
		Model(&models.UserTenant{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Update("is_active", false)
	return err
}

// GetTenantUsers returns users assigned to a tenant
func (s *TenantService) GetTenantUsers(tenantID uint) ([]models.UserTenant, error) {
	var userTenants []models.UserTenant
	if err := facades.Orm().Query().
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		With("User").
		With("Role").
		Find(&userTenants); err != nil {
		return nil, err
	}
	return userTenants, nil
}
