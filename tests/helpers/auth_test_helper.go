package helpers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	mockhttp "github.com/goravel/framework/mocks/http"
	mockauth "github.com/goravel/framework/mocks/auth"
	"players/app/models"
)

// CreateAuthenticatedContext creates a mock HTTP context with an authenticated user
func CreateAuthenticatedContext(user *models.User) http.Context {
	// Create mock context
	mockContext := &mockhttp.Context{}
	
	// Create mock auth
	mockAuth := &mockauth.Auth{}
	
	// Setup auth to return the user
	mockAuth.On("User", &models.User{}).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.User)
		*arg = *user
	})
	
	// Setup auth to return user ID
	mockAuth.On("ID").Return(user.ID)
	
	// Make facades.Auth return our mock
	facades.Auth = func(ctx ...http.Context) auth.Auth {
		return mockAuth
	}
	
	// Setup request methods
	mockRequest := &mockhttp.Request{}
	mockContext.On("Request").Return(mockRequest)
	
	return mockContext
}

// CreateUnauthenticatedContext creates a mock HTTP context without authentication
func CreateUnauthenticatedContext() http.Context {
	mockContext := &mockhttp.Context{}
	mockAuth := &mockauth.Auth{}
	
	// Setup auth to return error (no user)
	mockAuth.On("User", &models.User{}).Return(errors.New("unauthenticated"))
	mockAuth.On("ID").Return(uint(0))
	
	facades.Auth = func(ctx ...http.Context) auth.Auth {
		return mockAuth
	}
	
	mockRequest := &mockhttp.Request{}
	mockContext.On("Request").Return(mockRequest)
	
	return mockContext
}