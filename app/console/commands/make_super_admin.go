package commands

import (
	"errors"
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"players/app/models"
)

type MakeSuperAdmin struct {
}

// Signature The name and signature of the console command.
func (receiver *MakeSuperAdmin) Signature() string {
	return "user:make-super-admin {email}"
}

// Description The console command description.
func (receiver *MakeSuperAdmin) Description() string {
	return "Make a user a super admin"
}

// Extend The console command extend.
func (receiver *MakeSuperAdmin) Extend() command.Extend {
	return command.Extend{
		Category: "user",
	}
}

// Handle Execute the console command.
func (receiver *MakeSuperAdmin) Handle(ctx console.Context) error {
	email := ctx.Argument(0)
	if email == "" {
		return errors.New("email is required")
	}

	var user models.User
	err := facades.Orm().Query().Where("email = ?", email).First(&user)
	if err != nil {
		ctx.Error(fmt.Sprintf("User with email '%s' not found", email))
		return err
	}

	// Update user to be super admin
	_, err = facades.Orm().Query().Model(&user).Where("id = ?", user.ID).Update("is_super_admin", true)
	if err != nil {
		ctx.Error(fmt.Sprintf("Failed to update user: %v", err))
		return err
	}

	ctx.Success(fmt.Sprintf("User '%s' is now a super admin!", user.Name))
	return nil
}