package middleware

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"strings"
	"fmt" // For logging detailed info
)

// RedirectIfAuthenticated redirects authenticated users from / to /dashboard.
func RedirectIfAuthenticated() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		facades.Log().Info(fmt.Sprintf("[RedirectIfAuthenticated] Middleware hit for URL: %s, X-Inertia: %s", ctx.Request().Url(), ctx.Request().Header("X-Inertia", "")))

		tokenString := ""
		authHeader := ctx.Request().Header("Authorization", "")
		if authHeader != "" {
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) == 2 && strings.ToLower(headerParts[0]) == "bearer" {
				tokenString = headerParts[1]
				facades.Log().Info("[RedirectIfAuthenticated] Token found in Authorization header.")
			}
		}

		if tokenString == "" {
			cookieToken := ctx.Request().Cookie("token")
			if cookieToken != "" {
				tokenString = cookieToken
				facades.Log().Info("[RedirectIfAuthenticated] Token found in 'token' cookie.")
			} else {
				facades.Log().Info("[RedirectIfAuthenticated] No 'token' cookie found.")
			}
		}

		if tokenString == "" {
			facades.Log().Info("[RedirectIfAuthenticated] No token string available. Proceeding to next handler.")
			ctx.Request().Next()
			return
		}

		facades.Log().Info("[RedirectIfAuthenticated] Attempting to parse token.")
		_, err := facades.Auth(ctx).Parse(tokenString)
		if err != nil {
			facades.Log().Error(fmt.Sprintf("[RedirectIfAuthenticated] Token parsing error: %v", err))
			facades.Log().Info("[RedirectIfAuthenticated] Token invalid or expired. Proceeding to next handler.")
			ctx.Request().Next()
			return
		}

		facades.Log().Info("[RedirectIfAuthenticated] Token parsed successfully. User is considered authenticated.")

		if ctx.Request().Url() == "/" {
			facades.Log().Info("[RedirectIfAuthenticated] User is authenticated and on '/'. Attempting to redirect to /dashboard.")
			xInertiaHeader := ctx.Request().Header("X-Inertia", "")
			if xInertiaHeader == "true" {
				facades.Log().Info("[RedirectIfAuthenticated] Inertia request detected. Setting X-Inertia-Location to /dashboard and aborting with 409.")
				ctx.Response().Header("X-Inertia-Location", "/dashboard")
				ctx.Request().AbortWithStatus(contractshttp.StatusConflict)
			} else {
				facades.Log().Info("[RedirectIfAuthenticated] Non-Inertia request. Setting Location to /dashboard and aborting with 302.")
				ctx.Response().Header("Location", "/dashboard")
				ctx.Request().AbortWithStatus(contractshttp.StatusFound)
			}
			return // Ensure no further processing after redirect/abort
		}

		facades.Log().Info(fmt.Sprintf("[RedirectIfAuthenticated] User authenticated but on URL '%s', not '/'. Proceeding to next handler.", ctx.Request().Url()))
		ctx.Request().Next()
	}
}
