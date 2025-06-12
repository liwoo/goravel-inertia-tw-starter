package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"players/app/http"
	"players/routes"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Register(app foundation.Application) {
}

func (receiver *RouteServiceProvider) Boot(app foundation.Application) {
	// Add HTTP middleware
	facades.Route().GlobalMiddleware(http.Kernel{}.Middleware()...)

	receiver.configureRateLimiting()

	// Add routes
	routes.Web()

	// API routes will be prefixed with /api
	facades.Route().Prefix("api").Group(func(apiRouter route.Router) {
		routes.Api(apiRouter)
	})
}

func (receiver *RouteServiceProvider) configureRateLimiting() {

}
