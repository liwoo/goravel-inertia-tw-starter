package feature

import (
	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
	"starter-project/app/models"
	"starter-project/tests"
	"testing"
)

func TestJWTWorkaround(t *testing.T) {
	suite := &tests.TestCase{}
	suite.RefreshDatabase()

	// Create multiple users to check sequence
	hashedPassword, err := facades.Hash().Make("password")
	assert.NoError(t, err)

	// Create first user
	user1 := &models.User{
		Name:     "User 1",
		Email:    "user1@example.com",
		Password: hashedPassword,
		IsActive: true,
	}
	assert.NoError(t, facades.Orm().Query().Create(user1))
	t.Logf("First user ID: %d", user1.ID)

	// Create second user
	user2 := &models.User{
		Name:     "User 2",
		Email:    "user2@example.com",
		Password: hashedPassword,
		IsActive: true,
	}
	assert.NoError(t, facades.Orm().Query().Create(user2))
	t.Logf("Second user ID: %d", user2.ID)

	// Try raw SQL to insert user with ID=1
	_, err = facades.Orm().Query().Exec(`
		INSERT INTO users (id, name, email, password, role, is_active, is_super_admin, email_verified, created_at, updated_at) 
		VALUES (1, 'JWT User', 'jwt@example.com', ?, 'USER', 1, 0, 0, datetime('now'), datetime('now'))
	`, hashedPassword)

	if err != nil {
		t.Logf("Failed to insert user with ID=1: %v", err)

		// Check if user with ID=1 already exists
		var existingUser models.User
		err2 := facades.Orm().Query().Where("id = ?", 1).First(&existingUser)
		if err2 == nil {
			t.Logf("User with ID=1 already exists: %+v", existingUser)
		}
	} else {
		t.Log("Successfully created user with ID=1")
	}

	// Count total users
	count, err := facades.Orm().Query().Model(&models.User{}).Count()
	t.Logf("Total users in database: %d", count)
}
