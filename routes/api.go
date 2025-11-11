package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"smedi-sme-db/app/http/controllers"
	"smedi-sme-db/app/http/controllers/additional_business_members"
	"smedi-sme-db/app/http/controllers/auth"
	"smedi-sme-db/app/http/controllers/auth/perimissions"
	"smedi-sme-db/app/http/controllers/auth/roles"
	"smedi-sme-db/app/http/controllers/auth/users"
	"smedi-sme-db/app/http/controllers/books"
	"smedi-sme-db/app/http/controllers/configs"
	"smedi-sme-db/app/http/controllers/messages"
	"smedi-sme-db/app/http/controllers/primary_business_owners"
	"smedi-sme-db/app/http/controllers/smes"

	"smedi-sme-db/app/http/middleware"
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
	smeController := smes.NewSmeController()
	primaryBusinessOwnerController := primary_business_owners.NewPrimaryBusinessOwnerController()
	additionalBusinessMemberController := additional_business_members.NewAdditionalBusinessMemberController()

	jwtAuth := middleware.JwtAuth()
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

		//smes
		optionalAuthRouter.Get("/smes", smeController.Index)
		optionalAuthRouter.Get("/smes/search", smeController.Search)
		optionalAuthRouter.Get("/smes/filters", smeController.FilterMetadata)
		optionalAuthRouter.Get("/smes/{id}", smeController.Show)
		optionalAuthRouter.Get("/smes/{id}/primary_business_owner", smeController.FetchPrimaryBusinessOwner)
		optionalAuthRouter.Get("/smes/{id}/additional_business_members", smeController.FetchAdditionalBusinessMembers)
		optionalAuthRouter.Get("/smes/{id}/business_formalisation", smeController.FetchBusinessFormalisation)
		optionalAuthRouter.Get("/smes/{id}/business_employee_summary", smeController.FetchBusinessEmployeeSummary)

		//primary business owners
		optionalAuthRouter.Get("/primary_business_owners", primaryBusinessOwnerController.Index)
		optionalAuthRouter.Get("/primary_business_owners/search", primaryBusinessOwnerController.Search)
		optionalAuthRouter.Get("/primary_business_owners/filters", primaryBusinessOwnerController.FilterMetadata)
		optionalAuthRouter.Get("/primary_business_owners/{id}", primaryBusinessOwnerController.Show)

		//additional business members
		optionalAuthRouter.Get("/additional_business_members", additionalBusinessMemberController.Index)
		optionalAuthRouter.Get("/additional_business_members/search", additionalBusinessMemberController.Search)
		optionalAuthRouter.Get("/additional_business_members/filters", additionalBusinessMemberController.FilterMetadata)
		optionalAuthRouter.Get("/additional_business_members/{id}", additionalBusinessMemberController.Show)

	})

	// Protected routes (require authentication)
	router.Middleware(jwtAuth).Group(func(protectedRouter route.Router) {
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

		// SME routes
		protectedRouter.Post("/smes", smeController.Store)
		protectedRouter.Put("/smes/{id}", smeController.Update)
		protectedRouter.Delete("/smes/{id}", smeController.Delete)

		// Primary Business Owner routes
		protectedRouter.Post("/primary_business_owners", primaryBusinessOwnerController.Store)
		protectedRouter.Put("/primary_business_owners/{id}", primaryBusinessOwnerController.Update)
		protectedRouter.Delete("/primary_business_owners/{id}", primaryBusinessOwnerController.Delete)

		// Additional Business Member routes
		protectedRouter.Post("/additional_business_members", additionalBusinessMemberController.Store)
		protectedRouter.Put("/additional_business_members/{id}", additionalBusinessMemberController.Update)
		protectedRouter.Delete("/additional_business_members/{id}", additionalBusinessMemberController.Delete)

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
			messageRouter.Get("/conversations", messageController.GetConversations) // List all conversations
			messageRouter.Get("/conversation/{userId}", messageController.GetConversation)
			messageRouter.Put("/{id}/read", messageController.MarkAsRead)
			messageRouter.Get("/users", messageController.GetMessagableUsers) // Now implemented
			// messageRouter.Get("/search-users", messageController.SearchUsers) // TODO: implement
			// messageRouter.Get("/unread-count", messageController.GetUnreadCount) // TODO: implement as endpoint

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

	})

	// This Prefix("auth") group will also be relative to the router passed in.
	// If called from RouteServiceProvider's /api group, this becomes /api/auth
	router.Prefix("auth").Group(func(authRouter route.Router) {
		authRouter.Post("/login", apiAuthController.Login)
		authRouter.Middleware(jwtAuth).Post("/logout", apiAuthController.Logout)
		authRouter.Middleware(jwtAuth).Get("/me", apiAuthController.Me)
	})
}
