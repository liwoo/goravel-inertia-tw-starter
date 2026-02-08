package middleware

import (
	"books-database/app/tenant"

	contractshttp "github.com/goravel/framework/contracts/http"
)

// TenantFromSlug resolves the current tenant from the {tenant} route parameter
func TenantFromSlug() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		slug := ctx.Request().Route("tenant")

		if slug == "" {
			// No tenant slug in URL - resolve main tenant
			mainTenant, err := tenant.GetMainTenant()
			if err != nil {
				ctx.Request().AbortWithStatusJson(500, contractshttp.Json{
					"error": "Failed to resolve main tenant",
				})
				return
			}
			tenant.SetInContext(ctx, mainTenant)
			ctx.Request().Next()
			return
		}

		// Resolve tenant from slug
		t, err := tenant.ResolveBySlug(slug)
		if err != nil {
			ctx.Request().AbortWithStatusJson(404, contractshttp.Json{
				"error": "Tenant not found",
			})
			return
		}

		tenant.SetInContext(ctx, t)
		ctx.Request().Next()
	}
}
