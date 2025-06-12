package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
	"players/app/http/controllers"
	inertiaHelper "players/app/http/inertia"
	"players/app/http/middleware"
)

func Web() {
	// Register the Inertia middleware globally
	facades.Route().GlobalMiddleware(inertiaMiddleware)

	authController := controllers.NewAuthController()
	utilController := controllers.NewUtilController()
	dashboardController := controllers.NewDashboardController()

	facades.Route().Post("/login", authController.Login)
	//register una
	facades.Route().Get("/una", utilController.ShowUnaPage)

	// Authenticated routes
	facades.Route().Middleware(middleware.JwtAuth()).Group(func(router route.Router) {
		router.Post("/logout", authController.Logout)

		router.Get("/", func(ctx http.Context) http.Response {
			return inertiaHelper.Render(ctx, "auth/Login", map[string]interface{}{
				"version": support.Version,
			})
		})

		router.Get("/settings", func(ctx http.Context) http.Response {
			return inertiaHelper.Render(ctx, "settings/Index", map[string]interface{}{
				"version": support.Version,
			})
		})

		router.Get("/account", func(ctx http.Context) http.Response {
			return inertiaHelper.Render(ctx, "auth/Profile", map[string]interface{}{
				"version": support.Version,
			})
		})

		// Admin Dashboard - requires auth only
		router.Get("/dashboard", dashboardController.Show)
	})

	// Add more routes as needed
}

// inertiaMiddleware wraps the Inertia middleware
func inertiaMiddleware(ctx http.Context) {
	inertiaHelper.Middleware(func(c http.Context) http.Response {
		return nil
	})(ctx)
}
