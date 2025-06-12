package http

import (
	"github.com/goravel/framework/contracts/http"
	"players/app/http/middleware"
)

type Kernel struct {
}

// The application's global HTTP middleware stack.
// These middleware are run during every request to your application.
func (kernel Kernel) Middleware() []http.Middleware {
	return []http.Middleware{}
}

// The application's route middleware groups.
func (kernel Kernel) RouteMiddleware() map[string]http.Middleware {
	return map[string]http.Middleware{
		"admin":         middleware.AdminAuth(),
		"authenticated": middleware.Authenticated(), // Added new middleware
	}
}
