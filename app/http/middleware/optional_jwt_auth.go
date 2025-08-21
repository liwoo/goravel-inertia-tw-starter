package middleware

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/models"
	"strings"
)

// OptionalJwtAuth is a middleware that attempts to parse JWT token if present
// but doesn't require authentication. This allows endpoints to be publicly accessible
// while still identifying authenticated users when they provide a token.
func OptionalJwtAuth() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		// Check for token in Authorization header
		authHeader := ctx.Request().Header("Authorization", "")
		tokenString := ""

		if authHeader != "" {
			// Split "Bearer <token>"
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
				tokenString = headerParts[1]
			}
		}

		// If token not found in header, try cookie
		if tokenString == "" {
			cookieToken := ctx.Request().Cookie("token")
			if cookieToken != "" {
				tokenString = cookieToken
			}
		}

		// If we have a token, try to parse it
		if tokenString != "" {
			// Attempt to parse the token, but don't fail if it's invalid
			_, err := facades.Auth(ctx).Parse(tokenString)
			if err == nil {
				// If token is valid, make sure the user is loaded in context
				// This ensures GetAuthenticatedUser will work in permission checks
				var user models.User
				facades.Auth(ctx).User(&user)
			}
		}

		// Continue to the next middleware/handler regardless of auth status
		ctx.Request().Next()
	}
}
