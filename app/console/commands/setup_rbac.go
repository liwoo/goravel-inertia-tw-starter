package commands

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	"players/app/auth"
	"players/app/models"
	"players/database/seeders"
)

type SetupRBAC struct {
}

// Signature The name and signature of the console command.
func (receiver *SetupRBAC) Signature() string {
	return "rbac:setup"
}

// Description The console command description.
func (receiver *SetupRBAC) Description() string {
	return "Set up RBAC system with roles, permissions, and upgrade existing users"
}

// Extend The console command extend.
func (receiver *SetupRBAC) Extend() command.Extend {
	return command.Extend{
		Category: "rbac",
	}
}

// Handle Execute the console command.
func (receiver *SetupRBAC) Handle(ctx console.Context) error {
	ctx.Info("Setting up RBAC system...")

	// Step 1: Run RBAC seeder to create roles and permissions
	ctx.Info("Creating roles and permissions...")
	rbacSeeder := &seeders.RBACSeeder{}
	if err := rbacSeeder.Run(); err != nil {
		ctx.Error(fmt.Sprintf("Failed to create roles and permissions: %v", err))
		return err
	}
	ctx.Success("✓ Roles and permissions created successfully!")

	// Step 2: Upgrade existing users with RBAC roles
	ctx.Info("Upgrading existing users with RBAC roles...")
	if err := receiver.upgradeExistingUsers(ctx); err != nil {
		ctx.Warning(fmt.Sprintf("Some users could not be upgraded: %v", err))
	} else {
		ctx.Success("✓ Existing users upgraded successfully!")
	}

	// Step 3: Display summary
	ctx.Info("RBAC Setup Summary:")
	ctx.Info("==================")
	
	// Count roles and permissions
	var roleCount, permissionCount int64
	facades.Orm().Query().Model(&models.Role{}).Where("is_active = ?", true).Count(&roleCount)
	facades.Orm().Query().Model(&models.Permission{}).Where("is_active = ?", true).Count(&permissionCount)
	
	ctx.Info(fmt.Sprintf("• Roles created: %d", roleCount))
	ctx.Info(fmt.Sprintf("• Permissions created: %d", permissionCount))
	
	// Count users by role
	var adminCount, userCount int64
	facades.Orm().Query().Model(&models.User{}).Where("role = ?", "ADMIN").Count(&adminCount)
	facades.Orm().Query().Model(&models.User{}).Where("role = ?", "USER").Count(&userCount)
	
	ctx.Info(fmt.Sprintf("• Admin users: %d (upgraded to super-admin)", adminCount))
	ctx.Info(fmt.Sprintf("• Regular users: %d (upgraded to member)", userCount))

	ctx.Success("RBAC system setup completed successfully!")
	ctx.Info("You can now use the permission system in your controllers and services.")
	ctx.Info("To create a new admin user: go run . artisan user:create-admin")
	
	return nil
}

// upgradeExistingUsers assigns RBAC roles to existing users based on their legacy roles
func (receiver *SetupRBAC) upgradeExistingUsers(ctx console.Context) error {
	permissionService := auth.GetPermissionService()
	
	// Get all existing users
	var users []models.User
	err := facades.Orm().Query().Find(&users)
	if err != nil {
		return fmt.Errorf("failed to fetch existing users: %w", err)
	}

	upgraded := 0
	for _, user := range users {
		var targetRole string
		
		// Map legacy roles to RBAC roles
		switch user.Role {
		case "ADMIN":
			targetRole = "super-admin"
		case "MODERATOR":
			targetRole = "moderator"
		case "USER":
			targetRole = "member"
		default:
			targetRole = "member" // Default fallback
		}
		
		// Assign RBAC role
		err := permissionService.AssignRole(&user, targetRole, nil)
		if err != nil {
			ctx.Warning(fmt.Sprintf("Failed to assign role '%s' to user '%s': %v", targetRole, user.Email, err))
			continue
		}
		
		upgraded++
		ctx.Info(fmt.Sprintf("• Upgraded %s (%s) → %s role", user.Name, user.Email, targetRole))
	}

	ctx.Success(fmt.Sprintf("Successfully upgraded %d users", upgraded))
	return nil
}