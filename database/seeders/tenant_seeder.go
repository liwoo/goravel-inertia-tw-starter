package seeders

import (
	"books-database/app/models"

	"github.com/goravel/framework/facades"
)

type TenantSeeder struct{}

func (s *TenantSeeder) Signature() string {
	return "TenantSeeder"
}

func (s *TenantSeeder) Run() error {
	// Check if main tenant already exists
	var existing models.Tenant
	facades.Orm().Query().Where("is_main = ?", true).First(&existing)
	if existing.ID > 0 {
		facades.Log().Info("Main tenant already exists, skipping...")
		return nil
	}

	// Get default slug from config
	slug := facades.Config().GetString("tenancy.default_slug", "main")

	mainTenant := models.Tenant{
		Name:     facades.Config().GetString("app.name", "Main"),
		Slug:     slug,
		IsActive: true,
		IsMain:   true,
	}

	if err := facades.Orm().Query().Create(&mainTenant); err != nil {
		return err
	}
	facades.Log().Info("Main tenant created successfully")

	// Assign all existing super admin users to the main tenant
	var superAdmins []models.User
	facades.Orm().Query().Where("is_super_admin = ?", true).Find(&superAdmins)

	for _, admin := range superAdmins {
		ut := models.UserTenant{
			UserID:   admin.ID,
			TenantID: mainTenant.ID,
			IsActive: true,
		}
		if err := facades.Orm().Query().Create(&ut); err != nil {
			facades.Log().Warning("Failed to assign super admin to main tenant: " + err.Error())
		}
	}

	// Assign all other existing users to the main tenant
	var users []models.User
	facades.Orm().Query().Where("is_super_admin = ? OR is_super_admin IS NULL", false).Find(&users)

	for _, user := range users {
		ut := models.UserTenant{
			UserID:   user.ID,
			TenantID: mainTenant.ID,
			IsActive: true,
		}
		facades.Orm().Query().Create(&ut)
	}

	facades.Log().Info("All existing users assigned to main tenant")
	return nil
}
