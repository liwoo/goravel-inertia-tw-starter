package inertia

import (
	"encoding/json"
	"log"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/petaki/inertia-go"

	"players/app/models" // Import the User model
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
		sharedProps["auth"] = map[string]interface{}{
			"user": map[string]interface{}{
				"id":    authUser.ID,
				"name":  authUser.Name,
				"email": authUser.Email,
				"role":  authUser.Role,
			},
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

	// Create the page data
	pageMap := map[string]interface{}{
		"component": component,
		"props":     finalProps, // Use the merged props
		"url":       ctx.Request().FullUrl(),
		"version":   Version,
	}

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
