package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	"starter-project/app/auth"
	"starter-project/app/models"
)

type CreateAdminUser struct {
}

// Signature The name and signature of the console command.
func (receiver *CreateAdminUser) Signature() string {
	return "user:create-admin"
}

// Description The console command description.
func (receiver *CreateAdminUser) Description() string {
	return "Create a new admin user"
}

// Extend The console command extend.
func (receiver *CreateAdminUser) Extend() command.Extend {
	return command.Extend{
		Category: "user",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Admin user name (for non-interactive mode)",
			},
			&command.StringFlag{
				Name:    "email",
				Aliases: []string{"e"},
				Usage:   "Admin user email (for non-interactive mode)",
			},
			&command.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "Admin user password (for non-interactive mode)",
			},
			&command.BoolFlag{
				Name:    "skip-existing",
				Aliases: []string{"s"},
				Usage:   "Skip if admin user already exists (don't error)",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *CreateAdminUser) Handle(ctx console.Context) error {
	// Check for non-interactive mode (CLI flags)
	name := ctx.Option("name")
	email := ctx.Option("email")
	password := ctx.Option("password")
	skipExisting := ctx.OptionBool("skip-existing")

	// If all required flags are provided, run in non-interactive mode
	if name != "" && email != "" && password != "" {
		return receiver.createAdminNonInteractive(ctx, name, email, password, skipExisting)
	}

	// Interactive mode
	var err error
	name, err = ctx.Ask("Enter admin name:", console.AskOption{
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return errors.New("name cannot be empty")
			}
			return nil
		},
	})
	if err != nil {
		return err
	}

	email, err = ctx.Ask("Enter admin email:", console.AskOption{
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return errors.New("email cannot be empty")
			}
			if !strings.Contains(input, "@") {
				return errors.New("invalid email format")
			}
			return nil
		},
	})
	if err != nil {
		return err
	}

	password, err = ctx.Secret("Enter admin password (min 8 characters):", console.SecretOption{
		Validate: func(input string) error {
			if len(input) < 8 {
				return errors.New("password must be at least 8 characters long")
			}
			return nil
		},
	})
	if err != nil {
		return err
	}

	// Check if user already exists
	userCount, queryErr := facades.Orm().Query().Model(&models.User{}).Where("email = ?", email).Count()
	if queryErr != nil {
		ctx.Error(fmt.Sprintf("Error checking for existing user (count query): %v", queryErr))
		return queryErr
	}

	if userCount > 0 {
		ctx.Error(fmt.Sprintf("User with email '%s' already exists.", email))
		return errors.New("user already exists")
	}
	// If we reach here, userCount is 0, meaning the user does not exist.

	hashedPassword, err := facades.Hash().Make(password)
	if err != nil {
		ctx.Error(fmt.Sprintf("Error hashing password: %v", err))
		return err
	}

	adminUser := models.User{
		Name:          name,
		Email:         email,
		Password:      hashedPassword,
		Role:          "ADMIN", // Keep legacy role for backward compatibility
		IsActive:      true,
		IsSuperAdmin:  true, // Set super admin flag
		EmailVerified: true, // Admin users are pre-verified
	}

	createErr := facades.Orm().Query().Create(&adminUser)
	if createErr != nil {
		ctx.Error(fmt.Sprintf("Error creating admin user: %v", createErr))
		return createErr
	}

	// Assign super-admin role using RBAC system
	permissionService := auth.GetPermissionService()
	err = permissionService.AssignRole(&adminUser, "super-admin", nil)
	if err != nil {
		// Log warning but don't fail - user was created successfully
		ctx.Warning(fmt.Sprintf("User created but failed to assign super-admin role: %v", err))
		ctx.Warning("You may need to run 'go run . artisan seed --seeder=rbac' first, then manually assign the role.")
	} else {
		ctx.Success(fmt.Sprintf("Super-admin role assigned to '%s' successfully!", adminUser.Name))
	}

	ctx.Success(fmt.Sprintf("Admin user '%s' (%s) created successfully with super-admin privileges!", adminUser.Name, adminUser.Email))
	ctx.Info("User has both legacy ADMIN role and super-admin RBAC role for maximum compatibility.")
	return nil
}

// createAdminNonInteractive creates admin user in non-interactive mode (for CI/CD)
func (receiver *CreateAdminUser) createAdminNonInteractive(ctx console.Context, name, email, password string, skipExisting bool) error {
	ctx.Info("Running in non-interactive mode...")

	// Validate inputs
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check if user already exists
	userCount, queryErr := facades.Orm().Query().Model(&models.User{}).Where("email = ?", email).Count()
	if queryErr != nil {
		ctx.Error(fmt.Sprintf("Error checking for existing user: %v", queryErr))
		return queryErr
	}

	if userCount > 0 {
		if skipExisting {
			ctx.Info(fmt.Sprintf("Admin user with email '%s' already exists. Skipping creation.", email))
			return nil
		}
		ctx.Error(fmt.Sprintf("User with email '%s' already exists.", email))
		return errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := facades.Hash().Make(password)
	if err != nil {
		ctx.Error(fmt.Sprintf("Error hashing password: %v", err))
		return err
	}

	// Create admin user
	adminUser := models.User{
		Name:          name,
		Email:         email,
		Password:      hashedPassword,
		Role:          "ADMIN",
		IsActive:      true,
		IsSuperAdmin:  true,
		EmailVerified: true,
	}

	createErr := facades.Orm().Query().Create(&adminUser)
	if createErr != nil {
		ctx.Error(fmt.Sprintf("Error creating admin user: %v", createErr))
		return createErr
	}

	// Assign super-admin role using RBAC system
	permissionService := auth.GetPermissionService()
	err = permissionService.AssignRole(&adminUser, "super-admin", nil)
	if err != nil {
		ctx.Warning(fmt.Sprintf("User created but failed to assign super-admin role: %v", err))
	} else {
		ctx.Success(fmt.Sprintf("Super-admin role assigned to '%s' successfully!", adminUser.Name))
	}

	ctx.Success(fmt.Sprintf("Admin user '%s' (%s) created successfully!", adminUser.Name, adminUser.Email))
	return nil
}
