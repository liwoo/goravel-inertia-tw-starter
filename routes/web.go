package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
	"players/app/http/controllers"
	"players/app/http/controllers/auth"
	"players/app/http/controllers/books"
	inertiaHelper "players/app/http/inertia"
	"players/app/http/middleware"
)

func Web() {
	// Register the Inertia middleware globally
	facades.Route().GlobalMiddleware(inertiaMiddleware)

	authController := auth.NewAuthController()
	utilController := controllers.NewUtilController()
	dashboardController := controllers.NewDashboardController()
	booksPageController := books.NewBooksPageController()
	permissionsPageController := auth.NewPermissionsPageController()
	userPageController := auth.NewUserPageController()

	facades.Route().Post("/login", authController.Login)
	facades.Route().Get("/login", func(ctx http.Context) http.Response {
		return inertiaHelper.Render(ctx, "auth/Login", map[string]interface{}{
			"version": support.Version,
		})
	})
	//register una
	facades.Route().Get("/una", utilController.ShowUnaPage)

	// Public route for home/login, redirect to dashboard if already authenticated
	facades.Route().Middleware(middleware.RedirectIfAuthenticated()).Get("/", func(ctx http.Context) http.Response {
		return inertiaHelper.Render(ctx, "auth/Login", map[string]interface{}{
			"version": support.Version,
		})
	})

	// Authenticated routes
	facades.Route().Middleware(middleware.JwtAuth()).Group(func(router route.Router) {
		router.Post("/logout", authController.Logout)

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

		// Books management page
		router.Get("/admin/books", booksPageController.Index)

		// Permissions/Role management pages
		router.Get("/admin/permissions", permissionsPageController.Index)
		router.Get("/admin/permissions/roles/create", permissionsPageController.RoleCreate)
		router.Get("/admin/permissions/roles/:id", permissionsPageController.RoleShow)
		router.Get("/admin/permissions/roles/:id/edit", permissionsPageController.RoleEdit)

		// User management pages (super admin only)
		router.Get("/admin/users", userPageController.Index)
	})

	// Add more routes as needed
}

// inertiaMiddleware wraps the Inertia middleware
func inertiaMiddleware(ctx http.Context) {
	// The Inertia middleware should only set headers, not handle the full request
	// Just set the headers directly here
	if ctx.Request().Header("X-Inertia", "") == "true" {
		ctx.Response().Header("X-Inertia", "true")
		ctx.Response().Header("Vary", "Accept")
	}
}
