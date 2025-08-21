package commands

import (
	"fmt"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"players/app/models"
)

type CheckRolesCommand struct{}

func NewCheckRolesCommand() *CheckRolesCommand {
	return &CheckRolesCommand{}
}

func (receiver *CheckRolesCommand) Signature() string {
	return "roles:check"
}

func (receiver *CheckRolesCommand) Description() string {
	return "Check user roles assignments"
}

func (receiver *CheckRolesCommand) Extend() command.Extend {
	return command.Extend{}
}

func (receiver *CheckRolesCommand) Handle(ctx console.Context) error {
	// Check user_roles table
	var userRoles []models.UserRole
	err := facades.Orm().Query().
		With("User", "Role").
		Find(&userRoles)

	if err != nil {
		return fmt.Errorf("error loading user roles: %v", err)
	}

	ctx.Info("=== User Roles Table ===")
	for _, ur := range userRoles {
		ctx.Info(fmt.Sprintf("User: %s (ID:%d), Role: %s (ID:%d), IsActive: %v, AssignedAt: %v",
			ur.User.Email, ur.UserID,
			ur.Role.Name, ur.RoleID,
			ur.IsActive, ur.AssignedAt))
	}

	// Check if Lucky has any roles
	var user models.User
	err = facades.Orm().Query().
		Where("email = ?", "lucky@test.com").
		First(&user)

	if err == nil {
		ctx.Info("\n=== Lucky's Roles (direct check) ===")
		var luckyRoles []models.UserRole
		facades.Orm().Query().
			Where("user_id = ?", user.ID).
			With("Role").
			Find(&luckyRoles)

		for _, ur := range luckyRoles {
			ctx.Info(fmt.Sprintf("Role: %s, IsActive: %v", ur.Role.Name, ur.IsActive))
		}
	}

	// Check all roles
	var roles []models.Role
	facades.Orm().Query().Find(&roles)
	ctx.Info("\n=== All Roles ===")
	for _, role := range roles {
		ctx.Info(fmt.Sprintf("Role: %s (ID:%d), IsActive: %v", role.Name, role.ID, role.IsActive))
	}

	return nil
}
