package middleware

import (
	"strings"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/models"
)

// Require2FA returns a middleware function that enforces 2FA when required globally.
// This should be used after JwtAuth middleware on routes that need 2FA enforcement.
// It will redirect web users to /2fa-required and return 403 for API users.
func Require2FA() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		// Check if 2FA is required globally
		require2FA := facades.Config().GetBool("auth.require_2fa", false)
		if !require2FA {
			// 2FA not required, continue
			ctx.Request().Next()
			return
		}

		// Get the current URL to check if we're on an allowed path
		currentPath := ctx.Request().Path()

		// Allow certain paths without 2FA (the setup page itself, logout, 2FA API endpoints)
		allowedPaths := []string{
			"/2fa-required",
			"/logout",
			"/api/2fa/",    // All 2FA setup endpoints
			"/api/auth/",   // Auth endpoints
			"/api/account", // Account profile (needed for 2FA status)
		}

		for _, allowed := range allowedPaths {
			if strings.HasPrefix(currentPath, allowed) {
				ctx.Request().Next()
				return
			}
		}

		// Get authenticated user
		var user models.User
		if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
			// Not authenticated - let JwtAuth handle this
			ctx.Request().Next()
			return
		}

		// Check if user has 2FA enabled
		if user.TOTPEnabled {
			// User has 2FA enabled, continue
			ctx.Request().Next()
			return
		}

		// User needs to set up 2FA - handle based on request type
		xInertiaHeader := ctx.Request().Header("X-Inertia", "")
		isAPIRequest := strings.HasPrefix(currentPath, "/api/")

		if isAPIRequest && xInertiaHeader != "true" {
			// API request - return JSON error
			ctx.Request().AbortWithStatusJson(contractshttp.StatusForbidden, contractshttp.Json{
				"success": false,
				"message": "Two-factor authentication is required to access this resource",
				"error":   "2fa_required",
				"data": contractshttp.Json{
					"redirect": "/2fa-required",
				},
			})
			return
		}

		// Web/Inertia request - redirect to 2FA setup page
		if xInertiaHeader == "true" {
			ctx.Response().Header("X-Inertia-Location", "/2fa-required")
			ctx.Request().AbortWithStatus(contractshttp.StatusConflict)
		} else {
			ctx.Response().Header("Location", "/2fa-required")
			ctx.Request().AbortWithStatus(contractshttp.StatusFound)
		}
	}
}
