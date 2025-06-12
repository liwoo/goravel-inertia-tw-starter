package middleware

import (
	"fmt"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	. "net/http"

	"players/app/models" // Module 'players' from go.mod
)

// Authenticated returns a middleware handler function.
// This handler checks if a user is authenticated and has a valid ID.
func Authenticated() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		var appUser models.User
		// Attempt to retrieve the authenticated user.
		err := facades.Auth(ctx).User(&appUser)
		xInertiaHeader := ctx.Request().Header("X-Inertia", "")
		facades.Log().Infof("[AuthMiddleware] Checking auth. Error: %v, UserID: %d. X-Inertia: '%s'", err, appUser.ID, xInertiaHeader)

		isAuthenticated := err == nil && appUser.ID != 0

		if !isAuthenticated {
			facades.Log().Infof("[AuthMiddleware] Unauthenticated. Error: %v, UserID: %d. X-Inertia: '%s'", err, appUser.ID, xInertiaHeader)
			if xInertiaHeader == "true" {
				facades.Log().Info("[AuthMiddleware] Inertia request: sending 409 Conflict with X-Inertia-Location to /una")
				ctx.Response().Header("X-Inertia-Location", "/una")
				ctx.Response().Status(StatusConflict) // 409 Conflict
			} else {
				facades.Log().Info("[AuthMiddleware] Non-Inertia request: sending 302 Found to /una")
				ctx.Response().Redirect(StatusFound, "/una") // 302 Found
			}
			ctx.Request().Abort()
			return
		}

		//check if the route is / and redirect to /dashboard
		if ctx.Request().Url() == "/" {
			fmt.Println("[AuthMiddleware] Redirecting to /dashboard")
			facades.Log().Info("[AuthMiddleware] Redirecting to /dashboard")
			ctx.Response().Redirect(StatusFound, "/dashboard") // 302 Found
			ctx.Request().Abort()
			return
		}

		facades.Log().Info("[AuthMiddleware] Authenticated. Proceeding to next handler.")
		// User is authenticated and has a valid ID. Proceed to the next handler in the chain.
		ctx.Request().Next()
	}
}
