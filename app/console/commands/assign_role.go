package commands

import (
	"errors"
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	"players/app/auth"
	"players/app/models"
)

type AssignRole struct {
}

// Signature The name and signature of the console command.
func (receiver *AssignRole) Signature() string {
	return "rbac:assign {email} {role}"
}

// Description The console command description.
func (receiver *AssignRole) Description() string {
	return "Assign an RBAC role to a user by email"
}

// Extend The console command extend.
func (receiver *AssignRole) Extend() command.Extend {
	return command.Extend{
		Category: "rbac",
	}
}

// Handle Execute the console command.
func (receiver *AssignRole) Handle(ctx console.Context) error {
	email := ctx.Argument(0)
	roleSlug := ctx.Argument(1)

	if email == "" || roleSlug == "" {
		ctx.Error("Usage: go run . artisan rbac:assign <email> <role>")
		ctx.Info("Available roles: super-admin, admin, librarian, moderator, member, guest")
		return errors.New("missing arguments")
	}

	// Find user by email
	var user models.User
	err := facades.Orm().Query().Where("email = ?", email).First(&user)
	if err != nil {
		ctx.Error(fmt.Sprintf("User with email '%s' not found", email))
		return err
	}

	// Check if role exists
	var role models.Role
	err = facades.Orm().Query().Where("slug = ? AND is_active = ?", roleSlug, true).First(&role)
	if err != nil {
		ctx.Error(fmt.Sprintf("Role '%s' not found or not active", roleSlug))
		ctx.Info("Available roles:")
		receiver.listAvailableRoles(ctx)
		return err
	}

	// Assign role using permission service
	permissionService := auth.GetPermissionService()
	err = permissionService.AssignRole(&user, roleSlug, nil)
	if err != nil {
		ctx.Error(fmt.Sprintf("Failed to assign role: %v", err))
		return err
	}

	ctx.Success(fmt.Sprintf("Successfully assigned '%s' role to '%s' (%s)", role.Name, user.Name, user.Email))
	
	// Show user's current roles
	receiver.showUserRoles(ctx, &user)
	
	return nil
}

// listAvailableRoles displays all available roles
func (receiver *AssignRole) listAvailableRoles(ctx console.Context) {
	var roles []models.Role
	facades.Orm().Query().Where("is_active = ?", true).Order("level DESC").Find(&roles)
	
	for _, role := range roles {
		ctx.Info(fmt.Sprintf("• %s (%s) - Level %d - %s", role.Slug, role.Name, role.Level, role.Description))
	}
}

// showUserRoles displays current user roles
func (receiver *AssignRole) showUserRoles(ctx console.Context, user *models.User) {
	// Load user with roles
	var userWithRoles models.User
	facades.Orm().Query().Where("id = ?", user.ID).First(&userWithRoles)
	
	ctx.Info(fmt.Sprintf("Current roles for %s:", user.Name))
	
	// Show legacy role
	if userWithRoles.Role != "" {
		ctx.Info(fmt.Sprintf("• Legacy Role: %s", userWithRoles.Role))
	}
	
	// Show RBAC roles (simplified for now - would need to load relationships)
	ctx.Info("• RBAC Roles: [Role relationships would be displayed here]")
}