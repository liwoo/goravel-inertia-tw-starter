package inertia

import (
	"encoding/json"
	"log"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/petaki/inertia-go"

	"players/app/models" // Import the User model
	"players/app/auth"   // Import auth for permission helper
)

// Version represents the current asset version
const Version = "1.0.0"

var Manager *inertia.Inertia

// GetManager returns the Inertia manager instance
func GetManager() *inertia.Inertia {
	manager, _ := facades.App().Make("inertia")
	return manager.(*inertia.Inertia)
}

// Render renders an Inertia page
func Render(ctx http.Context, component string, props map[string]interface{}) http.Response {
	// Prepare shared props, including auth user
	sharedProps := make(map[string]interface{})

	// Add authenticated user information
	var authUser *models.User

	// Default to no authenticated user in shared props
	sharedProps["auth"] = map[string]interface{}{
		"user": nil,
		"permissions": map[string]interface{}{}, // Empty permissions object
	}

	err := facades.Auth(ctx).User(&authUser)

	if err != nil {
		// Only log if the error is unexpected (i.e., not the "token not parsed" error, which is normal on public pages)
		// or other standard unauthenticated errors if we knew them (e.g., auth.ErrUnauthenticated).
		// For now, we target the specific string from logs.
		if err.Error() != "authentication token must be parsed first" {
			log.Printf("Unexpected error fetching authenticated user for Inertia props: %v", err)
		}
		// auth.user remains nil as per default
	} else if authUser != nil && authUser.ID != 0 { // User successfully fetched and seems valid
		// Load user with roles for proper RBAC functionality
		var userWithRoles models.User
		err = facades.Orm().Query().
			Where("id = ?", authUser.ID).
			With("Roles").
			First(&userWithRoles)
		
		if err != nil {
			// Fallback to basic user info if roles loading fails
			log.Printf("DEBUG: Error loading user roles: %v", err)
			
			// Get permission helper to build permissions map even without roles
			permHelper := auth.GetPermissionHelper()
			
			// Build empty permissions for all services
			allPermissions := make(map[string]map[string]bool)
			allServices := auth.GetAllServiceRegistries()
			
			for _, service := range allServices {
				// This will return false for all permissions since no roles are loaded
				servicePerms := permHelper.BuildPermissionsMap(ctx, string(service))
				allPermissions[string(service)] = servicePerms
			}
			
			sharedProps["auth"] = map[string]interface{}{
				"user": map[string]interface{}{
					"id":    authUser.ID,
					"name":  authUser.Name,
					"email": authUser.Email,
					"role":  authUser.Role,
					"roles": []map[string]interface{}{}, // Empty roles array
					"permissions": []string{}, // Empty permissions array
					"isSuperAdmin": authUser.Role == "ADMIN", // Check legacy role
					"isAdmin": authUser.Role == "ADMIN",
				},
				"permissions": allPermissions,
			}
		} else {
			log.Printf("DEBUG: User %d loaded with %d roles", userWithRoles.ID, len(userWithRoles.Roles))
			for _, role := range userWithRoles.Roles {
				log.Printf("DEBUG: User has role: %s (active: %t)", role.Slug, role.IsActive)
			}
			
			// Get permission helper to build permissions map
			permHelper := auth.GetPermissionHelper()
			
			// Build permissions for all services the user might access
			allPermissions := make(map[string]map[string]bool)
			allServices := auth.GetAllServiceRegistries()
			
			for _, service := range allServices {
				servicePerms := permHelper.BuildPermissionsMap(ctx, string(service))
				allPermissions[string(service)] = servicePerms
			}
			
			// Include roles data for frontend RBAC checks
			rolesList := make([]map[string]interface{}, 0, len(userWithRoles.Roles))
			for _, role := range userWithRoles.Roles {
				if role.IsActive {
					rolesList = append(rolesList, map[string]interface{}{
						"id":          role.ID,
						"name":        role.Name,
						"slug":        role.Slug,
						"description": role.Description,
						"is_active":   role.IsActive,
					})
				}
			}
			
			// Get user's actual permissions from their roles
			userPermissions := permHelper.GetUserPermissions(ctx)
			log.Printf("DEBUG: User permissions loaded: %v", userPermissions)
			
			sharedProps["auth"] = map[string]interface{}{
				"user": map[string]interface{}{
					"id":          userWithRoles.ID,
					"name":        userWithRoles.Name,
					"email":       userWithRoles.Email,
					"role":        userWithRoles.Role,
					"roles":       rolesList,
					"permissions": userPermissions,
					"isSuperAdmin": userWithRoles.IsSuperAdmin(),
					"isAdmin":     userWithRoles.IsAdmin(),
				},
				"permissions": allPermissions,
			}
		}
	}
	// If err == nil but authUser is nil or authUser.ID == 0, auth.user remains nil (covered by default and the else if condition)

	// Merge controller-specific props with shared props
	// Controller props take precedence if keys overlap, though 'auth' should be unique to shared
	finalProps := make(map[string]interface{})
	for k, v := range sharedProps {
		finalProps[k] = v
	}
	for k, v := range props {
		finalProps[k] = v
	}

	// Get the URL safely
	requestURL := ctx.Request().FullUrl()
	if requestURL == "" {
		requestURL = ctx.Request().Url()
	}
	
	// Create the page data
	pageMap := map[string]interface{}{
		"component": component,
		"props":     finalProps, // Use the merged props
		"url":       requestURL,
		"version":   Version,
	}
	
	// Debug logging
	log.Printf("DEBUG: Inertia page data - component: %s, url: %s, version: %s", component, requestURL, Version)
	log.Printf("DEBUG: Props keys: %v", getMapKeys(finalProps))

	// Check if this is an Inertia request
	if ctx.Request().Header("X-Inertia", "") == "true" {
		// For Inertia requests, return JSON
		return ctx.Response().Header("X-Inertia", "true").
			Header("Vary", "Accept").
			Status(200).
			Json(pageMap)
	}

	// For regular requests, marshal the page data to JSON and render the template
	pageJSON, err := json.Marshal(pageMap)
	if err != nil {
		// Log the error and return an appropriate error response
		log.Printf("Error marshalling Inertia page data: %v", err)
		// Depending on your error handling strategy, you might return a 500 error page
		// For simplicity, returning a basic error response here
		return ctx.Response().String(500, "Error preparing page data")
	}

	return ctx.Response().
		Header("X-Inertia", "true").
		Header("Vary", "Accept").
		View().Make("app.tmpl", map[string]interface{}{
		"page":    string(pageJSON),
		"appName": facades.Config().GetString("app.name", "Goravel"),
		"isDev":   facades.Config().GetString("app.env", "production") != "production",
	})
}

// getMapKeys returns the keys of a map for debugging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Middleware wraps the Inertia middleware
func Middleware(next http.HandlerFunc) http.HandlerFunc {

	return func(ctx http.Context) http.Response {
		// Check if this is an Inertia request
		if ctx.Request().Header("X-Inertia", "") == "true" {
			// Set Inertia headers
			ctx.Response().Header("X-Inertia", "true")
			ctx.Response().Header("Vary", "Accept")

			// Check for version mismatch
			if ctx.Request().Header("X-Inertia-Version", "") != Version {
				// If there's a version mismatch and this is a GET request, force a full page reload
				if ctx.Request().Method() == "GET" {
					return ctx.Response().
						Header("X-Inertia-Location", ctx.Request().FullUrl()).
						Status(409).
						Json(http.Json{})
				}
			}
		}

		// Continue with the next handler
		return next(ctx)
	}
}
