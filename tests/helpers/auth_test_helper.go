package helpers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/testing/mock"
	"players/app/models"
)

// CreateAuthenticatedContext creates a mock context with an authenticated user
func CreateAuthenticatedContext(user *models.User) http.Context {
	mockFactory := mock.Factory()
	ctx := mockFactory.Context()
	
	if user != nil {
		auth := mockFactory.Auth(ctx)
		auth.On("User").Return(user)
		auth.On("Check").Return(true)
		auth.On("ID").Return(user.ID)
	}
	
	// Mock the Log facade that services might use
	mockFactory.Log()
	
	// Mock the ORM facade that services might use
	// Note: This returns a mock ORM but actual queries won't work
	// For integration tests that need real DB queries, use real facades instead
	mockFactory.Orm()
	
	return ctx
}

// NewTestContext creates a mock context with an authenticated user
func NewTestContext(user *models.User) http.Context {
	return CreateAuthenticatedContext(user)
}

// NewUnauthenticatedTestContext creates a mock context without authentication
func NewUnauthenticatedTestContext() http.Context {
	mockFactory := mock.Factory()
	ctx := mockFactory.Context()
	
	auth := mockFactory.Auth(ctx)
	auth.On("User").Return(nil)
	auth.On("Check").Return(false)
	auth.On("ID").Return("")
	
	// Mock the Log facade that services might use
	mockFactory.Log()
	
	// Mock the ORM facade that services might use
	// Note: This returns a mock ORM but actual queries won't work
	// For integration tests that need real DB queries, use real facades instead
	mockFactory.Orm()
	
	return ctx
}