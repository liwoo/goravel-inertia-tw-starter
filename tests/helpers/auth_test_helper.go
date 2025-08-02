package helpers

import (
	"errors"
	
	"github.com/goravel/framework/contracts/http"
	mockauth "github.com/goravel/framework/mocks/auth"
	mockhttp "github.com/goravel/framework/mocks/http"
	"github.com/stretchr/testify/mock"
	
	"players/app/models"
)

// CreateAuthenticatedContext creates a mock HTTP context with an authenticated user
func CreateAuthenticatedContext(user *models.User) http.Context {
	// Create mock context
	mockContext := &mockhttp.Context{}
	
	// Create mock auth
	mockAuth := &mockauth.Auth{}
	
	// Setup auth to return the user
	// Use mock.Anything instead of &models.User{} to match any argument
	mockAuth.On("User", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.User)
		*arg = *user
	})
	
	// Setup auth to return user ID
	mockAuth.On("ID").Return(user.ID)
	
	// Setup context to return our mock auth
	// The framework calls ctx.Value("GoravelAuth") to get the auth instance
	mockContext.On("Value", "GoravelAuth").Return(mockAuth)
	
	// Setup request methods
	mockRequest := &mockhttp.ContextRequest{}
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
	
	// Setup context to return our mock auth
	mockContext.On("Auth").Return(mockAuth)
	
	mockRequest := &mockhttp.ContextRequest{}
	mockContext.On("Request").Return(mockRequest)
	
	return mockContext
}