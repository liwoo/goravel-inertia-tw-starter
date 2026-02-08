package routes

import (
	"books-database/app/http/controllers"
	"books-database/app/http/controllers/applications"
	"books-database/app/http/controllers/auth"
	"books-database/app/http/controllers/auth/perimissions"
	"books-database/app/http/controllers/auth/users"
	"books-database/app/http/controllers/authors"
	"books-database/app/http/controllers/books"
	"books-database/app/http/controllers/configs"
	"books-database/app/http/controllers/lenders"
	tenants_ctrl "books-database/app/http/controllers/tenants"
	inertiaHelper "books-database/app/http/inertia"
	"books-database/app/http/middleware"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

func Web(router route.Router) {
	// Readiness probe - verifies database connection
	router.Get("/ready", func(ctx http.Context) http.Response {
		var result int
		err := facades.Orm().Query().Raw("SELECT 1").Scan(&result)
		if err != nil {
			return ctx.Response().Json(http.StatusServiceUnavailable, map[string]string{
				"status": "not ready",
				"error":  "database connection failed",
			})
		}
		return ctx.Response().Json(http.StatusOK, map[string]string{
			"status": "ready",
		})
	})

	// Liveness probe - simple health check
	router.Get("/health", func(ctx http.Context) http.Response {
		return ctx.Response().Json(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// Serve static files from the public directory

	authController := auth.NewAuthController()
	utilController := controllers.NewUtilController()
	dashboardController := controllers.NewDashboardController()
	booksPageController := books.NewBooksPageController()
	permissionsPageController := perimissions.NewPermissionsPageController()
	userPageController := users.NewUserPageController()
	configsPageController := configs.NewConfigPageController()
	lendersPageController := lenders.NewLenderPageController()
	authorsPageController := authors.NewAuthorsPageController()
	applicationsPageController := applications.NewApplicationPageController()
	tenantPageController := tenants_ctrl.NewTenantPageController()

	router.Post("/login", authController.Login)
	router.Post("/verify-2fa", authController.Verify2FA) // 2FA verification during web login
	router.Get("/login", func(ctx http.Context) http.Response {
		return inertiaHelper.Render(ctx, "auth/Login", map[string]interface{}{
			"version": support.Version,
		})
	})
	//register una
	router.Get("/una", utilController.ShowUnaPage)

	// Public route for home/login, redirect to dashboard if already authenticated
	router.Middleware(middleware.RedirectIfAuthenticated()).Get("/", func(ctx http.Context) http.Response {
		return inertiaHelper.Render(ctx, "auth/Login", map[string]interface{}{
			"version": support.Version,
		})
	})

	// Authenticated routes with 2FA enforcement
	// The Require2FA middleware checks if AUTH_REQUIRE_2FA is enabled and redirects
	// users without 2FA to /2fa-required
	router.Middleware(middleware.JwtAuth(), middleware.Require2FA()).Group(func(router route.Router) {
		router.Post("/logout", authController.Logout)

		// 2FA required setup page - accessible by authenticated users who need to set up 2FA
		// (Require2FA middleware allows this path even without 2FA)
		router.Get("/2fa-required", func(ctx http.Context) http.Response {
			return inertiaHelper.Render(ctx, "auth/TwoFactorRequired", map[string]interface{}{
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

		// Books management page
		router.Get("/admin/books", booksPageController.Index)

		// Lenders management page
		router.Get("/admin/lenders", lendersPageController.Index)

		// Authors management page
		router.Get("/admin/authors", authorsPageController.Index)

		// Applications management page
		router.Get("/admin/applications", applicationsPageController.Index)

		// Configurations management page
		router.Get("/admin/configs", configsPageController.Index)

		// Permissions/Role management pages
		router.Get("/admin/permissions", permissionsPageController.Index)

		router.Get("/admin/roles/:id/permissions", permissionsPageController.RolePermissions)

		// Tenant management (super admin only)
		router.Get("/admin/tenants", tenantPageController.Index)

		// User management pages (super admin only)
		router.Get("/admin/users", userPageController.Index)
	})

	// Add more routes as needed
}

// inertiaMiddleware wraps the Inertia middleware
func InertiaMiddleware(ctx http.Context) {
	// The Inertia middleware should only set headers, not handle the full request
	// Just set the headers directly here
	if ctx.Request().Header("X-Inertia", "") == "true" {
		ctx.Response().Header("X-Inertia", "true")
		ctx.Response().Header("Vary", "Accept")
	}
}
