package helpers

import (
	"books-database/app/models"
	"fmt"
	"github.com/goravel/framework/facades"
	"time"
)

// SetupJWTUser creates or updates a user with the given email, password, and role for testing
func SetupJWTUser(email, password string, role *models.Role) (*models.User, error) {
	fmt.Printf("SetupJWTUser called with email: %s\n", email)

	hashedPassword, err := facades.Hash().Make(password)
	if err != nil {
		return nil, err
	}

	// Check if user with this email already exists
	var user models.User
	err = facades.Orm().Query().Where("email = ?", email).First(&user)
	if err == nil {
		// User exists, update it
		user.Email = email // Ensure email is set
		user.Password = hashedPassword
		user.IsActive = true
		err = facades.Orm().Query().Save(&user)
		if err != nil {
			return nil, err
		}

		// Clear existing roles for this user
		facades.Orm().Query().Where("user_id = ?", user.ID).Delete(&models.UserRole{})
	} else {
		// Create a new user
		user = models.User{
			Name:     "Test User",
			Email:    email,
			Password: hashedPassword,
			IsActive: true,
		}
		err = facades.Orm().Query().Create(&user)
		if err != nil {
			return nil, err
		}
	}

	// Assign new role
	if role != nil {
		userRole := &models.UserRole{
			UserID:     user.ID,
			RoleID:     role.ID,
			AssignedAt: time.Now(),
			IsActive:   true,
		}
		err = facades.Orm().Query().Create(userRole)
		if err != nil {
			return nil, err
		}
	}

	// Reload the user to get all fields populated
	var fullUser models.User
	err = facades.Orm().Query().Where("id = ?", user.ID).First(&fullUser)
	if err != nil {
		return nil, err
	}

	fmt.Printf("SetupJWTUser returning user: ID=%d, Email=%s\n", fullUser.ID, fullUser.Email)
	return &fullUser, nil
}
