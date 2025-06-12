package middleware

import (
	stdhttp "net/http" // Alias for standard library http

	contractshttp "github.com/goravel/framework/contracts/http" // Alias for Goravel contracts
	"github.com/goravel/framework/facades"

	"players/app/models"
)

// AdminAuth returns a middleware function that conforms to contractshttp.Middleware (func(Context)).
func AdminAuth() contractshttp.Middleware {
	return func(ctx contractshttp.Context) { // Signature now func(Context), returns nothing
		var user models.User
		err := facades.Auth(ctx).User(&user)

		if err != nil {
			ctx.Response().Json(stdhttp.StatusInternalServerError, contractshttp.Json{
				"message": "Authentication system error",
			})
			return // Terminate processing for this request
		}

		if user.ID == 0 {
			ctx.Response().Json(stdhttp.StatusUnauthorized, contractshttp.Json{
				"message": "Unauthorized: Authentication required",
			})
			return // Terminate processing for this request
		}

		if user.Role != "ADMIN" {
			ctx.Response().Json(stdhttp.StatusForbidden, contractshttp.Json{
				"message": "Forbidden: Administrator access required",
			})
			return // Terminate processing for this request
		}

		// User is authenticated and is an ADMIN. Proceed to the next handler.
		ctx.Request().Next()
		// After Next(), control will eventually return here. The response is handled by downstream handlers.
	}
}
