package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type ControllerMaker struct {
}

// Signature The name and signature of the console command.
func (r *ControllerMaker) Signature() string {
	return "make:ctrl"
}

// Description The console command description.
func (r *ControllerMaker) Description() string {
	return "Generate a CRUD controller based on a model"
}

// Extend The console command extend.
func (r *ControllerMaker) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "The model to base the controller on",
			},
		},
	}
}

// Handle Execute the console command.
func (r *ControllerMaker) Handle(ctx console.Context) error {
	// Get the controller name from argument
	controllerName := ctx.Argument(0)
	if controllerName == "" {
		return fmt.Errorf("controller name is required")
	}

	// Get the model name from flag or derive from controller name
	modelName := ctx.Option("model")
	if modelName == "" {
		// Derive model name from controller name
		// e.g., "LenderController" -> "Lender", "lender" -> "Lender"
		modelName = strings.TrimSuffix(controllerName, "Controller")
		modelName = strings.Title(strings.ToLower(modelName))
	}

	// Ensure proper capitalization
	controllerBaseName := strings.TrimSuffix(controllerName, "Controller")
	controllerBaseName = strings.Title(strings.ToLower(controllerBaseName))
	controllerName = controllerBaseName + "Controller"

	// Resource name for routes and messages
	resourceName := strings.ToLower(controllerBaseName)
	resourceNamePlural := resourceName + "s"

	// Create controllers directory structure
	controllerDir := filepath.Join("app", "http", "controllers", resourceNamePlural)
	if err := os.MkdirAll(controllerDir, 0755); err != nil {
		return err
	}

	// Generate filename
	filename := filepath.Join(controllerDir, resourceName+"_controller.go")

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("controller %s already exists", filename)
	}

	// Generate controller content
	content := r.generateControllerContent(controllerName, controllerBaseName, modelName, resourceName, resourceNamePlural)

	// Write file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return err
	}

	ctx.Success(fmt.Sprintf("Controller created: %s", filename))
	ctx.NewLine()

	// Generate route registration instructions
	r.printRouteInstructions(ctx, resourceName, resourceNamePlural, controllerName)

	return nil
}

// generateControllerContent generates the controller file content
func (r *ControllerMaker) generateControllerContent(controllerName, baseName, modelName, resourceName, resourceNamePlural string) string {
	return fmt.Sprintf(`package %s

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// %s handles API endpoints for %s management
type %s struct {
	*contracts.CrudController[models.%s, *requests.%sCreateRequest, *requests.%sUpdateRequest]
	%sService *services.%sService
}

// New%s creates a new %s controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func New%s() *%s {
	%sService := services.New%sService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.%s, *requests.%sCreateRequest, *requests.%sUpdateRequest](
		"%s",
		%sService,
	).
		WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
			scopedHelper := auth.GetScopedPermissionHelper()

			// Map generic actions to permission actions
			var permAction auth.CorePermissionAction
			switch action {
			case "viewAny", "view":
				permAction = auth.PermissionRead
			case "create":
				permAction = auth.PermissionCreate
			case "update":
				permAction = auth.PermissionUpdate
			case "delete":
				permAction = auth.PermissionDelete
			default:
				permAction = auth.PermissionManage
			}

			// For specific resource actions, pass the resource
			if resource != nil {
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.Service%s, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.Service%s, permAction, nil)
			return err
		}).
		Build()

	controller := &%s{
		CrudController: crudController,
		%sService:    %sService,
	}

	// Set custom hooks
	controller.SetBeforeStore(func(ctx http.Context, data map[string]interface{}) error {
		// Set created_by from authenticated user
		var user models.User
		if err := facades.Auth(ctx).User(&user); err == nil && user.ID > 0 {
			data["created_by"] = user.ID
		}
		return nil
	})

	controller.SetBeforeUpdate(func(ctx http.Context, id uint, data map[string]interface{}) error {
		// Any custom logic before updating a %s
		return nil
	})

	return controller
}

// Add custom domain-specific methods below this line
// Examples:
//
// // CustomAction GET /api/%s/{id}/custom-action
// func (c *%s) CustomAction(ctx http.Context) http.Response {
//     // Get ID from URL
//     id, err := c.ValidateID(ctx, "id")
//     if err != nil {
//         return c.BadRequestResponse(ctx, "Invalid %s ID", nil)
//     }
//
//     // Check permissions
//     if err := c.CheckAuth(ctx, "view", nil); err != nil {
//         return c.ForbiddenResponse(ctx, "Access denied")
//     }
//
//     // Your custom logic here
//     // result, err := c.%sService.CustomMethod(id)
//
//     return c.SuccessResponse(ctx, nil, "Action completed successfully")
// }
`,
		resourceNamePlural,
		controllerName, resourceName, controllerName,
		modelName, baseName, baseName,
		resourceName, baseName,
		controllerName, resourceName,
		controllerName, controllerName,
		resourceName, baseName,
		modelName, baseName, baseName,
		resourceName,
		resourceName,
		strings.Title(resourceNamePlural), strings.Title(resourceNamePlural),
		controllerName,
		resourceName, resourceName,
		resourceName,
		resourceNamePlural, controllerName,
		resourceName,
		resourceName,
	)
}

// printRouteInstructions outputs the route registration instructions
func (r *ControllerMaker) printRouteInstructions(ctx console.Context, resourceName, resourceNamePlural, controllerName string) {
	ctx.Info("To register the controller routes, add the following to routes/api.go:")
	ctx.NewLine()

	ctx.Comment("1. Import the controller package:")
	ctx.Line(fmt.Sprintf("   \"%s/app/http/controllers/%s\"", "players", resourceNamePlural))
	ctx.NewLine()

	ctx.Comment("2. Initialize the controller in the Api() function:")
	ctx.Line(fmt.Sprintf("   %sController := %s.New%s()", resourceName, resourceNamePlural, controllerName))
	ctx.NewLine()

	ctx.Comment("3. Add routes with optional auth (for scoped permissions):")
	ctx.Line("   router.Middleware(optionalAuth).Group(func(optionalAuthRouter route.Router) {")
	ctx.Line(fmt.Sprintf("       optionalAuthRouter.Get(\"/%s\", %sController.Index)", resourceNamePlural, resourceName))
	ctx.Line(fmt.Sprintf("       optionalAuthRouter.Get(\"/%s/search\", %sController.Search)", resourceNamePlural, resourceName))
	ctx.Line(fmt.Sprintf("       optionalAuthRouter.Get(\"/%s/filters\", %sController.FilterMetadata)", resourceNamePlural, resourceName))
	ctx.Line(fmt.Sprintf("       optionalAuthRouter.Get(\"/%s/{id}\", %sController.Show)", resourceNamePlural, resourceName))
	ctx.Line("   })")
	ctx.NewLine()

	ctx.Comment("4. Add protected routes (require authentication):")
	ctx.Line("   router.Middleware(jwtAuth).Group(func(protectedRouter route.Router) {")
	ctx.Line(fmt.Sprintf("       protectedRouter.Post(\"/%s\", %sController.Store)", resourceNamePlural, resourceName))
	ctx.Line(fmt.Sprintf("       protectedRouter.Put(\"/%s/{id}\", %sController.Update)", resourceNamePlural, resourceName))
	ctx.Line(fmt.Sprintf("       protectedRouter.Delete(\"/%s/{id}\", %sController.Delete)", resourceNamePlural, resourceName))
	ctx.Line("       // Add custom routes here if needed")
	ctx.Line("   })")
	ctx.NewLine()

	ctx.Info("Middleware explanation:")
	ctx.Line("- optionalAuth: Allows public access but identifies authenticated users (for scoped permissions)")
	ctx.Line("- jwtAuth: Requires authentication - blocks request if no valid token")
	ctx.Line("")
	ctx.Comment("NOTE: If your resource contains sensitive data (like users), use jwtAuth for ALL routes.")
}
