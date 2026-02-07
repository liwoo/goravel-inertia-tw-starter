package http

import (
	"books-database/app/http/middleware"

	"github.com/goravel/framework/contracts/http"
	sessionMiddleware "github.com/goravel/framework/session/middleware"
)

type Kernel struct {
}

// The application's global HTTP middleware stack.
// These middleware are run during every request to your application.
func (kernel Kernel) Middleware() []http.Middleware {
	return []http.Middleware{
		sessionMiddleware.StartSession(), // Add session middleware globally
	}
}

// The application's route middleware groups.
func (kernel Kernel) RouteMiddleware() map[string]http.Middleware {
	return map[string]http.Middleware{
		"admin":         middleware.AdminAuth(),
		"authenticated": middleware.Authenticated(), // Added new middleware
	}
}
