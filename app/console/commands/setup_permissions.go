package commands

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/models"
)

// SetupPermissionsCommand sets up the standard permission system
type SetupPermissionsCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *SetupPermissionsCommand) Signature() string {
	return "permissions:setup"
}

// Description The console command description.
func (receiver *SetupPermissionsCommand) Description() string {
	return "Set up the standard permission system with all service-action combinations"
}

// Extend The console command extend.
func (receiver *SetupPermissionsCommand) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *SetupPermissionsCommand) Handle(ctx console.Context) error {
	ctx.Info("Setting up standard permission system...")

	// Get all services and actions
	services := auth.GetAllServiceRegistries()
	
	permissionsCreated := 0
	permissionsSkipped := 0

	for _, service := range services {
		actions := auth.GetServiceActions(service)
		
		for _, action := range actions {
			permissionSlug := auth.BuildPermissionSlug(service, action)
			
			// Check if permission already exists
			var existingPermission models.Permission
			err := facades.Orm().Query().
				Where("slug = ?", permissionSlug).
				First(&existingPermission)
			
			if err == nil {
				// Permission already exists
				permissionsSkipped++
				continue
			}

			// Create new permission
			permission := models.Permission{
				Name:        fmt.Sprintf("%s %s", auth.GetServiceDisplayName(service), auth.GetActionDisplayName(action)),
				Slug:        permissionSlug,
				Description: fmt.Sprintf("Allows %s on %s", auth.GetActionDisplayName(action), auth.GetServiceDisplayName(service)),
				Category:    string(service),
				Resource:    string(service),
				Action:      string(action),
				IsActive:    true,
			}

			err = facades.Orm().Query().Create(&permission)
			if err != nil {
				ctx.Error(fmt.Sprintf("Failed to create permission %s: %v", permissionSlug, err))
				continue
			}

			permissionsCreated++
			ctx.Line(fmt.Sprintf("Created permission: %s", permissionSlug))
		}
	}

	ctx.Success(fmt.Sprintf("Permission setup complete! Created: %d, Skipped: %d", permissionsCreated, permissionsSkipped))
	
	// Show next steps
	ctx.Info("Next steps:")
	ctx.Line("1. Run 'go run . artisan permissions:list' to see all permissions")
	ctx.Line("2. Run 'go run . artisan permissions:assign-role <role> <permission>' to assign permissions")
	ctx.Line("3. Visit /admin/permissions to manage permissions via the web interface")
	
	return nil
}