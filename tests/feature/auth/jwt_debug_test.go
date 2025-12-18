package feature

import (
	"fmt"
	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
	"starter-project/app/models"
	"starter-project/tests"
	"testing"
)

func TestJWTDebug(t *testing.T) {
	suite := &tests.TestCase{}
	suite.RefreshDatabase()

	// Create a user
	hashedPassword, err := facades.Hash().Make("password")
	assert.NoError(t, err)

	user := &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
		IsActive: true,
	}
	assert.NoError(t, facades.Orm().Query().Create(user))
	fmt.Printf("Created user with ID: %d\n", user.ID)

	// Get a mock context for auth
	// This is just to check how JWT encodes the user ID
	fmt.Printf("User struct before login: %+v\n", user)

	// Check what the ID field contains
	fmt.Printf("User.ID: %d\n", user.ID)
	fmt.Printf("User.Model.ID: %d\n", user.Model.ID)
}
