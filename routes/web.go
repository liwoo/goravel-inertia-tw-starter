package routes

import (
	"smedi-sme-db/app/http/controllers"
	"smedi-sme-db/app/http/controllers/applications"
	"smedi-sme-db/app/http/controllers/auth"
	"smedi-sme-db/app/http/controllers/auth/perimissions"
	"smedi-sme-db/app/http/controllers/auth/users"
	"smedi-sme-db/app/http/controllers/bdsps"
	"smedi-sme-db/app/http/controllers/books"
	"smedi-sme-db/app/http/controllers/configs"
	"smedi-sme-db/app/http/controllers/directory"
	"smedi-sme-db/app/http/controllers/events"
	"smedi-sme-db/app/http/controllers/members"
	"smedi-sme-db/app/http/controllers/procurementnotices"
	"smedi-sme-db/app/http/controllers/smes"
	inertiaHelper "smedi-sme-db/app/http/inertia"
	"smedi-sme-db/app/http/middleware"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

func Web() {
	// Readiness probe - verifies database connection
	facades.Route().Get("/ready", func(ctx http.Context) http.Response {
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
	facades.Route().Get("/health", func(ctx http.Context) http.Response {
		return ctx.Response().Json(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// Serve static files from the public directory
	facades.Route().Static("/images", "./public/images")
	facades.Route().Static("/css", "./public/css")
	facades.Route().Static("/js", "./public/js")

	// Register the Inertia middleware globally
	facades.Route().GlobalMiddleware(inertiaMiddleware)

	authController := auth.NewAuthController()
	utilController := controllers.NewUtilController()
	dashboardController := controllers.NewDashboardController()
	booksPageController := books.NewBooksPageController()
	permissionsPageController := perimissions.NewPermissionsPageController()
	userPageController := users.NewUserPageController()
	configsPageController := configs.NewConfigPageController()
	smesPageController := smes.NewSmePageController()
	bdspsPageController := bdsps.NewBdspPageController()
	eventsPageController := events.NewEventPageController()
	procurementnoticesPageController := procurementnotices.NewProcurementNoticePageController()
	applicationsPageController := applications.NewApplicationPageController()
	membersPageController := members.NewMemberPageController()
	directoryController := directory.NewDirectoryController()

	facades.Route().Post("/login", authController.Login)
	facades.Route().Post("/verify-2fa", authController.Verify2FA) // 2FA verification during web login
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

	// Public Application Page
	facades.Route().Get("/apply", applicationsPageController.ShowPublicApply)

	// Authenticated routes with 2FA enforcement
	// The Require2FA middleware checks if AUTH_REQUIRE_2FA is enabled and redirects
	// users without 2FA to /2fa-required
	facades.Route().Middleware(middleware.JwtAuth(), middleware.Require2FA()).Group(func(router route.Router) {
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

		// SMEs management page
		router.Get("/admin/smes", smesPageController.Index)

		// BDSPs management page
		router.Get("/admin/bdsps", bdspsPageController.Index)

		// Event management page
		router.Get("/admin/events", eventsPageController.Index)

		// Procurement Notice management page
		router.Get("/admin/procurement-notices", procurementnoticesPageController.Index)

		// Applications management page
		router.Get("/admin/applications", applicationsPageController.Index)

		// Configurations management page
		router.Get("/admin/configs", configsPageController.Index)

		// Permissions/Role management pages
		router.Get("/admin/permissions", permissionsPageController.Index)

		router.Get("/admin/roles/:id/permissions", permissionsPageController.RolePermissions)

		// User management pages (super admin only)
		router.Get("/admin/users", userPageController.Index)

		router.Get("/portal", membersPageController.Index)

		// SME Directory (accessible to all authenticated users)
		router.Get("/directory", directoryController.ShowDirectory)

		// SSE Test page (for development/testing)
		router.Get("/test/sse", func(ctx http.Context) http.Response {
			return inertiaHelper.Render(ctx, "test/SSETest", map[string]interface{}{
				"version": support.Version,
			})
		})
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
