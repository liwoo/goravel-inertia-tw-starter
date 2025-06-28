package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"players/app/http/controllers/auth"
	"players/app/http/controllers/books"

	"players/app/http/middleware"
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

	userController := auth.NewUserController()

	bookController := books.NewBookController()
	authController := auth.NewAuthController()
	rolesController := &auth.RolesController{}
	permissionsController := &auth.PermissionsController{}
	jwtAuth := middleware.JwtAuth()

	// Book resource routes
	router.Get("/books", bookController.Index)
	router.Get("/books/{id}", bookController.Show)
	router.Get("/books/isbn/{isbn}", bookController.GetByISBN)
	router.Get("/books/author/{author}", bookController.GetByAuthor)
	router.Get("/books/available", bookController.GetAvailable)
	router.Get("/books/advanced", bookController.Advanced)

	// Protected routes (require authentication)
	router.Middleware(jwtAuth).Group(func(protectedRouter route.Router) {
		// Book routes
		protectedRouter.Post("/books", bookController.Store)
		protectedRouter.Put("/books/{id}", bookController.Update)
		protectedRouter.Delete("/books/{id}", bookController.Delete)
		protectedRouter.Post("/books/{id}/borrow", bookController.Borrow)
		protectedRouter.Post("/books/{id}/return", bookController.Return)

		// Role management routes
		protectedRouter.Get("/roles", rolesController.Index)
		protectedRouter.Post("/roles", rolesController.Store)
		protectedRouter.Get("/roles/{id}", rolesController.Show)
		protectedRouter.Put("/roles/{id}", rolesController.Update)
		protectedRouter.Delete("/roles/{id}", rolesController.Destroy)

		// Permission assignment routes
		protectedRouter.Post("/permissions/assign", permissionsController.Assign)
		protectedRouter.Delete("/permissions/revoke", permissionsController.Revoke)

		// User management routes (super admin only)
		protectedRouter.Get("/users", userController.Index)
		protectedRouter.Get("/users/{id}", userController.Show)
		protectedRouter.Post("/users", userController.Store)
		protectedRouter.Put("/users/{id}", userController.Update)
		protectedRouter.Delete("/users/{id}", userController.Delete)
		protectedRouter.Get("/users/roles", userController.GetRoles)
	})

	// This Prefix("auth") group will also be relative to the router passed in.
	// If called from RouteServiceProvider's /api group, this becomes /api/auth
	router.Prefix("auth").Group(func(authRouter route.Router) {
		authRouter.Post("/login", authController.Login)
		authRouter.Middleware(jwtAuth).Post("/logout", authController.Logout)
	})
}
