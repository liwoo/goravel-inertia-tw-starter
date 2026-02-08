package routes

import (
	"books-database/app/http/controllers"
	"books-database/app/http/controllers/account"
	"books-database/app/http/controllers/applications"
	"books-database/app/http/controllers/auth"
	"books-database/app/http/controllers/auth/perimissions"
	"books-database/app/http/controllers/auth/roles"
	"books-database/app/http/controllers/auth/users"
	"books-database/app/http/controllers/authors"
	"books-database/app/http/controllers/books"
	"books-database/app/http/controllers/configs"
	"books-database/app/http/controllers/lenders"
	"books-database/app/http/controllers/messages"
	tenants_api "books-database/app/http/controllers/tenants"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"

	"books-database/app/http/middleware"
)

// Api defines the routes for the API.
// It now accepts a router argument to scope routes correctly within a group.
func Api(router route.Router) {
	// This GET "/" will be relative to the router passed in.
	// If called from RouteServiceProvider's /api group, this becomes /api/
	router.Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().Success().Json(http.Json{
			"Hello": "Goravel API is live!",
		})
	})

	userController := users.NewUserController()

	bookController := books.NewBookController()
	apiAuthController := auth.NewAPIAuthController()
	rolesController := roles.NewRolesController()
	permissionsController := perimissions.NewPermissionsController()
	searchController := controllers.NewSearchController()
	messageController := messages.NewMessageController()
	notificationController := messages.NewNotificationController()
	swaggerController := controllers.NewSwaggerController()
	configController := configs.NewConfigController()
	lenderController := lenders.NewLenderController()
	tenantController := tenants_api.NewTenantController()
	authorController := authors.NewAuthorController()
	applicationController := applications.NewApplicationController()
	accountController := account.NewAccountController()
	totpController := auth.NewTOTPController()
	sseController := controllers.NewSSEController()
	presenceController := controllers.NewPresenceController()

	jwtAuth := middleware.JwtAuth()
	require2FA := middleware.Require2FA()
	optionalAuth := middleware.OptionalJwtAuth()

	// Swagger API documentation routes (public access)
	router.Prefix("swagger").Group(func(swaggerRouter route.Router) {
		swaggerRouter.Get("/", swaggerController.ServeSwagger)
		swaggerRouter.Get("/{any}", swaggerController.ServeSwagger)
	})
	router.Prefix("docs").Group(func(docsRouter route.Router) {
		docsRouter.Get("/swagger.json", swaggerController.ServeSwaggerJSON)
		docsRouter.Get("/swagger.yaml", swaggerController.ServeSwaggerYAML)
	})

	// Book resource routes (with optional auth for scoped permissions)
	router.Middleware(optionalAuth).Group(func(optionalAuthRouter route.Router) {
		optionalAuthRouter.Get("/books", bookController.Index)
		optionalAuthRouter.Get("/books/search", bookController.Search)          // Search endpoint (must be before {id})
		optionalAuthRouter.Get("/books/filters", bookController.FilterMetadata) // Filter metadata endpoint
		optionalAuthRouter.Get("/books/available", bookController.Available)
		optionalAuthRouter.Get("/books/isbn/{isbn}", bookController.GetByISBN)
		optionalAuthRouter.Get("/books/author/{author}", bookController.GetByAuthor)
		optionalAuthRouter.Get("/books/{id}", bookController.Show) // Must be last to avoid conflicts

		//configurations
		optionalAuthRouter.Get("/configs", configController.Index)
		optionalAuthRouter.Get("/configs/search", configController.Search)
		optionalAuthRouter.Get("/configs/filters", configController.FilterMetadata)
		optionalAuthRouter.Get("/configs/{id}", configController.Show)

		//lenders
		optionalAuthRouter.Get("/lenders", lenderController.Index)
		optionalAuthRouter.Get("/lenders/search", lenderController.Search)
		optionalAuthRouter.Get("/lenders/filters", lenderController.FilterMetadata)
		optionalAuthRouter.Get("/lenders/{id}", lenderController.Show)

		// Authors
		optionalAuthRouter.Get("/authors", authorController.Index)
		optionalAuthRouter.Get("/authors/search", authorController.Search)          // Search endpoint (must be before {id})
		optionalAuthRouter.Get("/authors/filters", authorController.FilterMetadata) // Filter metadata endpoint
		optionalAuthRouter.Get("/authors/{id}", authorController.Show)              // Must be last to avoid conflicts

		// Public application submission (no auth required)
		optionalAuthRouter.Post("/applications", applicationController.PublicStore)
	})

	// Protected routes (require authentication and 2FA when enabled)
	router.Middleware(jwtAuth, require2FA).Group(func(protectedRouter route.Router) {
		// Global search
		protectedRouter.Get("/search", searchController.GlobalSearch)

		// Book routes
		protectedRouter.Post("/books", bookController.Store)
		protectedRouter.Put("/books/{id}", bookController.Update)
		protectedRouter.Delete("/books/{id}", bookController.Delete)

		// Custom endpoints
		protectedRouter.Post("/books/{id}/borrow", bookController.Borrow)
		protectedRouter.Post("/books/{id}/return", bookController.Return)
		protectedRouter.Get("/books/statistics", bookController.Statistics)

		// Configuration routes
		protectedRouter.Post("/configs", configController.Store)
		protectedRouter.Put("/configs/{id}", configController.Update)
		protectedRouter.Delete("/configs/{id}", configController.Delete)

		//lenders
		protectedRouter.Post("/lenders", lenderController.Store)
		protectedRouter.Put("/lenders/{id}", lenderController.Update)
		protectedRouter.Delete("/lenders/{id}", lenderController.Delete)

		// Author routes
		protectedRouter.Get("/authors/statistics", authorController.Statistics) // Must be before {id}
		protectedRouter.Post("/authors", authorController.Store)
		protectedRouter.Put("/authors/{id}", authorController.Update)
		protectedRouter.Delete("/authors/{id}", authorController.Delete)

		protectedRouter.Get("/applications", applicationController.Index)
		protectedRouter.Get("/applications/search", applicationController.Search)
		protectedRouter.Get("/applications/filters", applicationController.FilterMetadata)
		protectedRouter.Get("/applications/{id}", applicationController.Show)
		protectedRouter.Put("/applications/{id}", applicationController.Update)
		protectedRouter.Delete("/applications/{id}", applicationController.Delete)
		protectedRouter.Post("/applications/{id}/approve", applicationController.ApproveApplication)
		protectedRouter.Post("/applications/{id}/reject", applicationController.RejectApplication)

		// Role management routes
		protectedRouter.Get("/roles", rolesController.Index)
		protectedRouter.Post("/roles", rolesController.Store)
		protectedRouter.Get("/roles/{id}", rolesController.Show)
		protectedRouter.Put("/roles/{id}", rolesController.Update)
		protectedRouter.Delete("/roles/{id}", rolesController.Delete)
		protectedRouter.Put("/roles/{id}/permissions", rolesController.UpdatePermissions)

		// Permission assignment routes
		protectedRouter.Post("/permissions/assign", permissionsController.Assign)
		protectedRouter.Delete("/permissions/revoke", permissionsController.Revoke)

		// User management routes (super admin only)
		protectedRouter.Get("/users", userController.Index)
		protectedRouter.Get("/users/filters", userController.FilterMetadata)
		protectedRouter.Get("/users/{id}", userController.Show)
		protectedRouter.Post("/users", userController.Store)
		protectedRouter.Put("/users/{id}", userController.Update)
		protectedRouter.Delete("/users/{id}", userController.Delete)
		// protectedRouter.Get("/users/roles", userController.GetRoles) // TODO: implement

		// Messaging routes
		protectedRouter.Prefix("messages").Group(func(messageRouter route.Router) {
			// Send and manage messages
			messageRouter.Post("/", messageController.SendMessage)
			messageRouter.Post("/broadcast", messageController.SendBroadcast)              // Broadcast to multiple users
			messageRouter.Post("/broadcast-to-role", messageController.BroadcastToRole)    // Super admin: broadcast to all users in a role
			messageRouter.Get("/broadcast-history", messageController.GetBroadcastHistory) // Super admin: get broadcast history
			messageRouter.Get("/conversations", messageController.GetConversations)        // List all conversations
			messageRouter.Get("/conversation/{userId}", messageController.GetConversation)
			messageRouter.Put("/conversation/{userId}/read", messageController.MarkConversationAsRead)
			messageRouter.Get("/inbox", messageController.GetInbox)           // User's inbox
			messageRouter.Get("/sent", messageController.GetSentMessages)     // User's sent messages
			messageRouter.Get("/unread", messageController.GetUnreadMessages) // User's unread messages
			messageRouter.Put("/{id}/read", messageController.MarkAsRead)
			messageRouter.Get("/users", messageController.GetMessagableUsers)
			messageRouter.Get("/unread-count", messageController.GetUnreadCount)
			// messageRouter.Get("/search-users", messageController.SearchUsers) // TODO: implement

			// Individual message management
			// messageRouter.Get("/{id}", messageController.GetMessage) // TODO: use Show method
			// messageRouter.Put("/{id}", messageController.EditMessage) // TODO: implement
			messageRouter.Delete("/{id}", messageController.Delete)
			// messageRouter.Post("/{id}/reply", messageController.ReplyToMessage) // TODO: implement
		})

		// Notification routes
		protectedRouter.Prefix("notifications").Group(func(notificationRouter route.Router) {
			// User notifications
			notificationRouter.Get("/", notificationController.GetNotifications)
			notificationRouter.Get("/counts", notificationController.GetCounts)
			notificationRouter.Get("/type/{type}", notificationController.GetByType)

			// Mark as read/dismiss
			notificationRouter.Put("/{id}/read", notificationController.MarkAsRead)
			notificationRouter.Put("/read-all", notificationController.MarkAllAsRead)
			notificationRouter.Delete("/{id}", notificationController.DismissNotification)
			notificationRouter.Delete("/", notificationController.DismissAllNotifications)

			// Batch operations
			notificationRouter.Put("/batch/read", notificationController.BatchMarkAsRead)
			notificationRouter.Delete("/batch", notificationController.BatchDismiss)

			// Admin operations
			notificationRouter.Post("/", notificationController.CreateNotification)
			notificationRouter.Post("/system", notificationController.CreateSystemNotification)
			notificationRouter.Post("/cleanup", notificationController.CleanupExpired)
		})

		// Account routes (user profile and settings)
		protectedRouter.Prefix("account").Group(func(accountRouter route.Router) {
			accountRouter.Get("/profile", accountController.GetProfile)
			accountRouter.Put("/profile", accountController.UpdateProfile)
			accountRouter.Put("/password", accountController.ChangePassword)
			accountRouter.Get("/activities", accountController.GetActivities)
			accountRouter.Get("/activities/types", accountController.GetActivityTypes)
			accountRouter.Get("/activities/summary", accountController.GetActivitySummary)
			accountRouter.Get("/recent-activities", accountController.GetRecentActivities)
		})

		// Tenant management routes (super admin only)
		protectedRouter.Prefix("tenants").Group(func(tenantRouter route.Router) {
			tenantRouter.Get("/", tenantController.Index)
			tenantRouter.Get("/search", tenantController.Search)
			tenantRouter.Get("/filters", tenantController.FilterMetadata)
			tenantRouter.Get("/{id}", tenantController.Show)
			tenantRouter.Post("/", tenantController.Store)
			tenantRouter.Put("/{id}", tenantController.Update)
			tenantRouter.Delete("/{id}", tenantController.Delete)
			tenantRouter.Post("/{id}/users", tenantController.AssignUser)
			tenantRouter.Delete("/{id}/users/{userId}", tenantController.RemoveUser)
			tenantRouter.Get("/{id}/users", tenantController.GetUsers)
		})

		// Two-factor authentication routes
		protectedRouter.Prefix("2fa").Group(func(twoFaRouter route.Router) {
			twoFaRouter.Get("/status", totpController.Status)
			twoFaRouter.Post("/setup", totpController.Setup)
			twoFaRouter.Post("/verify", totpController.Verify)
			twoFaRouter.Post("/disable", totpController.Disable)
			twoFaRouter.Post("/backup-codes", totpController.RegenerateBackupCodes)
		})

	})

	// This Prefix("auth") group will also be relative to the router passed in.
	// If called from RouteServiceProvider's /api group, this becomes /api/auth
	router.Prefix("auth").Group(func(authRouter route.Router) {
		authRouter.Post("/login", apiAuthController.Login)
		authRouter.Post("/verify-2fa", apiAuthController.Verify2FA) // 2FA verification during login (no auth required)
		authRouter.Middleware(jwtAuth, require2FA).Post("/logout", apiAuthController.Logout)
		authRouter.Middleware(jwtAuth, require2FA).Get("/me", apiAuthController.Me)
	})

	// SSE (Server-Sent Events) for real-time updates
	router.Middleware(jwtAuth, require2FA).Get("/sse/stream", sseController.Stream)

	// Presence (online status) routes
	router.Middleware(jwtAuth, require2FA).Prefix("presence").Group(func(presenceRouter route.Router) {
		presenceRouter.Get("/online", presenceController.GetOnlineUsers)
		presenceRouter.Get("/stats", presenceController.GetStats)
		presenceRouter.Get("/check", presenceController.CheckOnline)
		presenceRouter.Get("/user/{id}", presenceController.GetUserStatus)
		presenceRouter.Post("/bulk", presenceController.GetBulkStatus)
	})
}
