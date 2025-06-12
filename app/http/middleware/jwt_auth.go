package middleware

import (
	"fmt"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"strings"
)

// JwtAuth returns a middleware function that handles JWT authentication.
func JwtAuth() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		authHeader := ctx.Request().Header("Authorization", "")
		xInertiaHeader := ctx.Request().Header("X-Inertia", "")
		var tokenString string

		if authHeader != "" {
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) == 2 && strings.ToLower(headerParts[0]) == "bearer" {
				tokenString = headerParts[1]
			} else {
				// Invalid header format, proceed to handleAuthFailure which will be caught by empty tokenString
			}
		}

		// If token not found in header, try cookie
		if tokenString == "" {
			cookieToken := ctx.Request().Cookie("token") // Default Goravel JWT cookie name
			if cookieToken != "" {
				tokenString = cookieToken
			}
		}

		handleAuthFailure := func(logMessage string) {
			// Log the failure reason if needed, perhaps using facades.Log() once configured
			if xInertiaHeader == "true" {
				if ctx.Request().Url() == "/" {
					ctx.Response().Header("X-Inertia-Location", "/login")
				} else {
					ctx.Response().Header("X-Inertia-Location", "/una")
				}
				ctx.Request().AbortWithStatus(contractshttp.StatusConflict)
			} else {
				if ctx.Request().Url() == "/" {
					ctx.Response().Header("Location", "/login")
				} else {
					ctx.Response().Header("Location", "/una")
				}
				ctx.Request().AbortWithStatus(contractshttp.StatusFound)
			}
		}

		if tokenString == "" {
			handleAuthFailure("Authentication token not found in header or cookie")
			return
		}

		if _, err := facades.Auth(ctx).Parse(tokenString); err != nil {
			handleAuthFailure("Invalid or expired token: " + err.Error())
			return
		}

		//check if the route is / and redirect to /dashboard
		if ctx.Request().Url() == "/" {
			fmt.Println("[AuthMiddleware] Redirecting to /dashboard")
			facades.Log().Info("[AuthMiddleware] Redirecting to /dashboard")
			//redirect using inertia
			if xInertiaHeader == "true" {
				ctx.Response().Header("X-Inertia-Location", "/dashboard")
				ctx.Request().AbortWithStatus(contractshttp.StatusConflict)
			} else {
				ctx.Response().Header("Location", "/dashboard")
				ctx.Request().AbortWithStatus(contractshttp.StatusFound)
			}
		}

		ctx.Request().Next()
	}
}
