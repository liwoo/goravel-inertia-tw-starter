package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"books-database/app/http"
	"books-database/app/http/middleware"
	"books-database/routes"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Register(app foundation.Application) {
}

func (receiver *RouteServiceProvider) Boot(app foundation.Application) {
	// Add HTTP middleware
	facades.Route().GlobalMiddleware(http.Kernel{}.Middleware()...)

	receiver.configureRateLimiting()

	// Static files and inertia middleware (global, run once)
	facades.Route().Static("/images", "./public/images")
	facades.Route().Static("/css", "./public/css")
	facades.Route().Static("/js", "./public/js")
	facades.Route().GlobalMiddleware(routes.InertiaMiddleware)

	// Main routes (no tenant prefix - backward compatible)
	routes.Web(facades.Route())

	// API routes prefixed with /api
	facades.Route().Prefix("api").Group(func(apiRouter route.Router) {
		routes.Api(apiRouter)
	})

	// Tenant-prefixed routes: /t/{tenant}/...
	tenantMw := middleware.TenantFromSlug()
	facades.Route().Prefix("t/{tenant}").Middleware(tenantMw).Group(func(tenantRouter route.Router) {
		routes.Web(tenantRouter)
	})
	facades.Route().Prefix("t/{tenant}/api").Middleware(tenantMw).Group(func(tenantApiRouter route.Router) {
		routes.Api(tenantApiRouter)
	})
}

func (receiver *RouteServiceProvider) configureRateLimiting() {

}
