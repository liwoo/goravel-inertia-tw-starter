package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type MakeCrudCommand struct{}

func (receiver *MakeCrudCommand) Signature() string {
	return "make:crud {name : The name of the resource} {--model : Also create model} {--migration : Also create migration} {--routes : Show route examples}"
}

func (receiver *MakeCrudCommand) Description() string {
	return "Create a complete CRUD resource (Repository, Service, Requests, Controller)"
}

func (receiver *MakeCrudCommand) Extend() command.Extend {
	return command.Extend{}
}

func (receiver *MakeCrudCommand) Handle(ctx console.Context) error {
	name := ctx.Argument(0)
	if name == "" {
		return fmt.Errorf("resource name is required")
	}

	// Convert to proper case
	resourceName := strings.Title(name)

	ctx.Info(fmt.Sprintf("Creating CRUD for resource: %s", resourceName))

	// 1. Create Model (if --model flag is set)
	if ctx.OptionBool("model") {
		if err := createModel(ctx, resourceName); err != nil {
			ctx.Error(fmt.Sprintf("Failed to create model: %v", err))
		} else {
			ctx.Info("✓ Model created")
		}
	}

	// 2. Create Migration (if --migration flag is set)
	if ctx.OptionBool("migration") {
		if err := createMigration(ctx, resourceName); err != nil {
			ctx.Error(fmt.Sprintf("Failed to create migration: %v", err))
		} else {
			ctx.Info("✓ Migration created")
		}
	}

	// 3. Create Repository
	if err := createRepository(ctx, resourceName); err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}
	ctx.Info("✓ Repository created")

	// 4. Create Service
	if err := createService(ctx, resourceName); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	ctx.Info("✓ Service created")

	// 5. Create Create Request
	if err := createCreateRequest(ctx, resourceName); err != nil {
		return fmt.Errorf("failed to create create request: %w", err)
	}
	ctx.Info("✓ Create request created")

	// 6. Create Update Request
	if err := createUpdateRequest(ctx, resourceName); err != nil {
		return fmt.Errorf("failed to create update request: %w", err)
	}
	ctx.Info("✓ Update request created")

	// 7. Create Controller
	if err := createController(ctx, resourceName); err != nil {
		return fmt.Errorf("failed to create controller: %w", err)
	}
	ctx.Info("✓ Controller created")

	// 8. Create Gates (optional)
	if err := createGates(ctx, resourceName); err != nil {
		ctx.Warning(fmt.Sprintf("Failed to create gates: %v", err))
	} else {
		ctx.Info("✓ Authorization gates created")
	}

	// 9. Show route suggestion (if --routes flag is set)
	if ctx.OptionBool("routes") {
		showRouteExample(ctx, resourceName)
	}

	ctx.Info("")
	ctx.Info("🎉 CRUD resource created successfully!")
	ctx.Info("")
	ctx.Info("Next steps:")
	ctx.Info("1. Update the model fields in app/models/" + strings.ToLower(resourceName) + ".go")
	ctx.Info("2. Update validation rules in app/http/requests/*_request.go")
	ctx.Info("3. Add domain-specific methods to the repository")
	ctx.Info("4. Register routes for the new resource")
	ctx.Info("5. Run migration if created: go run . artisan migrate")
	ctx.Info("")
	ctx.Info("To see route examples, run: go run . artisan make:crud " + resourceName + " --routes")

	return nil
}

func createModel(ctx console.Context, name string) error {
	ctx.Info("Creating model...")
	ctx.Info(fmt.Sprintf("Please run: go run . artisan make:model %s", name))
	return nil
}

func createMigration(ctx console.Context, name string) error {
	ctx.Info("Creating migration...")
	tableName := strings.ToLower(name) + "s"
	migrationName := "create_" + tableName + "_table"
	ctx.Info(fmt.Sprintf("Please run: go run . artisan make:migration %s", migrationName))
	return nil
}

func createRepository(ctx console.Context, name string) error {
	ctx.Info("Creating repository...")
	ctx.Info(fmt.Sprintf("Please run: go run . artisan make:repository %sRepository", name))
	return nil
}

func createService(ctx console.Context, name string) error {
	ctx.Info("Creating service...")
	ctx.Info(fmt.Sprintf("Please run: go run . artisan make:service %sService", name))
	return nil
}

func createCreateRequest(ctx console.Context, name string) error {
	ctx.Info("Creating create request...")
	ctx.Info(fmt.Sprintf("Please run: go run . artisan make:request %sCreateRequest", name))
	return nil
}

func createUpdateRequest(ctx console.Context, name string) error {
	ctx.Info("Creating update request...")
	ctx.Info(fmt.Sprintf("Please run: go run . artisan make:request %sUpdateRequest --update", name))
	return nil
}

func createController(ctx console.Context, name string) error {
	ctx.Info("Creating controller...")
	return createControllerFile(name)
}

func createControllerFile(resourceName string) error {
	controllerName := resourceName + "Controller"
	controllerPath := filepath.Join("app", "http", "controllers")
	if err := os.MkdirAll(controllerPath, 0755); err != nil {
		return err
	}

	filename := filepath.Join(controllerPath, strings.ToLower(resourceName)+"_controller.go")

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("controller %s already exists", controllerName)
	}

	// Generate controller content
	content, err := generateControllerContent(controllerName)
	if err != nil {
		return err
	}

	// Write file
	return os.WriteFile(filename, []byte(content), 0644)
}

func generateControllerContent(controllerName string) (string, error) {
	resourceName := strings.TrimSuffix(controllerName, "Controller")
	modelName := strings.Title(resourceName)
	serviceName := resourceName + "Service"
	createRequestName := resourceName + "CreateRequest"
	updateRequestName := resourceName + "UpdateRequest"
	resourceNameLower := strings.ToLower(resourceName)
	resourceNamePlural := resourceNameLower + "s"
	serviceField := resourceNameLower + "Service"

	tmpl := `package controllers

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/http/requests"
	"players/app/services"
)

// {{.ControllerName}} - enhanced with validation and authorization
type {{.ControllerName}} struct {
	{{.ServiceField}} *services.{{.ServiceName}}
	authHelper  contracts.AuthHelper
}

// New{{.ControllerName}} creates a new {{.ResourceNameLower}} controller
func New{{.ControllerName}}() *{{.ControllerName}} {
	return &{{.ControllerName}}{
		{{.ServiceField}}: services.New{{.ServiceName}}(),
		authHelper:  helpers.NewAuthHelper(),
	}
}

// Index GET /{{.ResourceNamePlural}}
func (c *{{.ControllerName}}) Index(ctx http.Context) http.Response {
	// Check authorization
	if response := facades.Gate().Inspect("viewAny.{{.ResourceNamePlural}}", ctx); response.Denied() {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": response.Message(),
		})
	}

	var req helpers.ListRequest
	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid request parameters",
		})
	}

	result, err := c.{{.ServiceField}}.GetList(req)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, result)
}

// Show GET /{{.ResourceNamePlural}}/{id}
func (c *{{.ControllerName}}) Show(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid {{.ResourceNameLower}} ID",
		})
	}

	// Get the {{.ResourceNameLower}} first to pass to authorization
	{{.ResourceNameLower}}, err := c.{{.ServiceField}}.GetByID(uint(id))
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "{{.ModelName}} not found",
		})
	}

	// Check authorization
	if response := facades.Gate().Inspect("view.{{.ResourceNamePlural}}", ctx, {{.ResourceNameLower}}); response.Denied() {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": response.Message(),
		})
	}

	return ctx.Response().Json(http.StatusOK, {{.ResourceNameLower}})
}

// Store POST /{{.ResourceNamePlural}}
func (c *{{.ControllerName}}) Store(ctx http.Context) http.Response {
	// Check authorization
	if response := facades.Gate().Inspect("create.{{.ResourceNamePlural}}", ctx); response.Denied() {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": response.Message(),
		})
	}

	// Create and validate request
	var createRequest requests.{{.CreateRequestName}}
	errors, err := ctx.Request().ValidateRequest(&createRequest)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, map[string]interface{}{
			"error":  "Validation failed",
			"errors": errors.All(),
		})
	}

	// Create the {{.ResourceNameLower}} using validated data
	{{.ResourceNameLower}}, err := c.{{.ServiceField}}.Create(createRequest.ToCreateData())
	if err != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusCreated, {{.ResourceNameLower}})
}

// Update PUT /{{.ResourceNamePlural}}/{id}
func (c *{{.ControllerName}}) Update(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid {{.ResourceNameLower}} ID",
		})
	}

	// Get the {{.ResourceNameLower}} first to pass to authorization
	{{.ResourceNameLower}}, err := c.{{.ServiceField}}.GetByID(uint(id))
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "{{.ModelName}} not found",
		})
	}

	// Check authorization
	if response := facades.Gate().Inspect("update.{{.ResourceNamePlural}}", ctx, {{.ResourceNameLower}}); response.Denied() {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": response.Message(),
		})
	}

	// Create and validate update request
	var updateRequest requests.{{.UpdateRequestName}}
	updateRequest.ID = uint(id) // Set the ID for validation context

	errors, err := ctx.Request().ValidateRequest(&updateRequest)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, map[string]interface{}{
			"error":  "Validation failed",
			"errors": errors.All(),
		})
	}

	// Update the {{.ResourceNameLower}} using validated data
	updated{{.ModelName}}, err := c.{{.ServiceField}}.Update(uint(id), updateRequest.ToUpdateData())
	if err != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, updated{{.ModelName}})
}

// Destroy DELETE /{{.ResourceNamePlural}}/{id}
func (c *{{.ControllerName}}) Destroy(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid {{.ResourceNameLower}} ID",
		})
	}

	// Get the {{.ResourceNameLower}} first to pass to authorization
	{{.ResourceNameLower}}, err := c.{{.ServiceField}}.GetByID(uint(id))
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "{{.ModelName}} not found",
		})
	}

	// Check authorization
	if response := facades.Gate().Inspect("delete.{{.ResourceNamePlural}}", ctx, {{.ResourceNameLower}}); response.Denied() {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": response.Message(),
		})
	}

	err = c.{{.ServiceField}}.Delete(uint(id))
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusNoContent, nil)
}

// Advanced GET /{{.ResourceNamePlural}}/advanced - with filters
func (c *{{.ControllerName}}) Advanced(ctx http.Context) http.Response {
	// Public endpoint - no authorization needed for viewing
	var req helpers.ListRequest
	if err := ctx.Request().Bind(&req); err != nil {
		req = helpers.ListRequest{} // Use defaults
	}

	// Parse filters from query parameters
	filters := make(map[string]interface{})

	if status := ctx.Request().Query("status"); status != "" {
		filters["status"] = status
	}
	if name := ctx.Request().Query("name"); name != "" {
		filters["name"] = name
	}
	// TODO: Add more filters specific to {{.ModelName}}

	result, err := c.{{.ServiceField}}.GetListAdvanced(req, filters)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, result)
}
`

	t, err := template.New("controller").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = t.Execute(&result, map[string]string{
		"ControllerName":      controllerName,
		"ServiceName":         serviceName,
		"ServiceField":        serviceField,
		"CreateRequestName":   createRequestName,
		"UpdateRequestName":   updateRequestName,
		"ModelName":           modelName,
		"ResourceName":        resourceName,
		"ResourceNameLower":   resourceNameLower,
		"ResourceNamePlural":  resourceNamePlural,
	})

	return result.String(), err
}

func createGates(ctx console.Context, resourceName string) error {
	ctx.Info("Creating authorization gates...")
	
	gatePath := filepath.Join("app", "providers")
	filename := filepath.Join(gatePath, strings.ToLower(resourceName)+"_gate_provider.go")

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("gate provider already exists")
	}

	content := generateGateContent(resourceName)
	return os.WriteFile(filename, []byte(content), 0644)
}

func generateGateContent(resourceName string) string {
	resourceNameLower := strings.ToLower(resourceName)
	resourceNamePlural := resourceNameLower + "s"

	return fmt.Sprintf(`package providers

import (
	"players/app/contracts"
	"players/app/helpers"

	"github.com/goravel/framework/contracts/auth/access"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type %sGateServiceProvider struct{}

func (receiver *%sGateServiceProvider) Register(app foundation.Application) {
	// Register any dependencies
}

func (receiver *%sGateServiceProvider) Boot(app foundation.Application) {
	// Register %s resource gates
	gateHelper := helpers.NewGateHelper()

	// Define %s-specific gate configuration
	%sGateConfig := contracts.GateConfig{
		ViewAnyHandler: func(ctx http.Context, user interface{}) access.Response {
			// Everyone can view %s lists
			return access.NewAllowResponse()
		},
		ViewHandler: func(ctx http.Context, user interface{}, model interface{}) access.Response {
			// Everyone can view individual %ss
			return access.NewAllowResponse()
		},
		CreateHandler: func(ctx http.Context, user interface{}) access.Response {
			// Only moderators and admins can create %ss
			return gateHelper.RoleBasedAccess("ADMIN", "MODERATOR")(ctx, user)
		},
		UpdateHandler: func(ctx http.Context, user interface{}, model interface{}) access.Response {
			// Only moderators and admins can update %ss
			return gateHelper.RoleBasedAccess("ADMIN", "MODERATOR")(ctx, user)
		},
		DeleteHandler: func(ctx http.Context, user interface{}, model interface{}) access.Response {
			// Only admins can delete %ss
			return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
		},
	}

	// Register the gates for %ss
	gateHelper.RegisterResourceGates("%s", %sGateConfig)

	// TODO: Add custom gates specific to %s operations
	// Example:
	// facades.Gate().Define("activate.%s", func(ctx http.Context, user interface{}, args ...interface{}) access.Response {
	//     return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
	// })
}
`,
		strings.Title(resourceName),        // %sGateServiceProvider struct name
		strings.Title(resourceName),        // %sGateServiceProvider Register method
		strings.Title(resourceName),        // %sGateServiceProvider Boot method  
		resourceNameLower,                  // Register %s resource gates
		resourceNameLower,                  // Define %s-specific gate configuration
		resourceNameLower,                  // %sGateConfig variable
		resourceNameLower,                  // Everyone can view %s lists
		resourceNameLower,                  // Everyone can view individual %ss
		resourceNameLower,                  // Only moderators and admins can create %ss
		resourceNameLower,                  // Only moderators and admins can update %ss
		resourceNameLower,                  // Only admins can delete %ss
		resourceNamePlural,                 // Register the gates for %ss
		resourceNamePlural,                 // gateHelper.RegisterResourceGates("%s", %sGateConfig)
		resourceNameLower,                  // %sGateConfig
		resourceName,                       // TODO: Add custom gates specific to %s operations
		resourceNamePlural,                 // facades.Gate().Define("activate.%s"
	)
}

func showRouteExample(ctx console.Context, resourceName string) {
	resourceNameLower := strings.ToLower(resourceName)
	resourceNamePlural := resourceNameLower + "s"
	controllerName := resourceName + "Controller"

	example := fmt.Sprintf(`
📋 Route Examples for %s:

Add the following to your routes/api.go or routes/web.go:

package routes

import (
    "github.com/goravel/framework/contracts/route"
    "github.com/goravel/framework/facades"
    "players/app/http/controllers"
    "players/app/http/middleware"
)

func %sRoutes() {
    %sController := controllers.New%s()

    // Public routes - anyone can view
    facades.Route().Group(func(router route.Router) {
        router.Get("/%s", %sController.Index)                    // List %ss
        router.Get("/%s/{id}", %sController.Show)               // Get single %s
        router.Get("/%s/advanced", %sController.Advanced)       // Advanced filtering
    })

    // Protected routes - require authentication
    facades.Route().Middleware(middleware.JwtAuth()).Group(func(router route.Router) {
        router.Post("/%s", %sController.Store)                  // Create %s
        router.Put("/%s/{id}", %sController.Update)             // Update %s
        router.Delete("/%s/{id}", %sController.Destroy)         // Delete %s
    })
}

Don't forget to call %sRoutes() in your main route registration!

API Endpoints:
GET    /%s                   - List %ss with pagination/search
GET    /%s/{id}              - Get single %s
GET    /%s/advanced          - Advanced filtering
POST   /%s                   - Create %s (auth required)
PUT    /%s/{id}              - Update %s (auth required)
DELETE /%s/{id}              - Delete %s (auth required)
`,
		resourceName,
		resourceName,
		strings.ToLower(controllerName[:len(controllerName)-10]),
		controllerName,
		resourceNamePlural, strings.ToLower(controllerName[:len(controllerName)-10]), resourceNamePlural,
		resourceNamePlural, strings.ToLower(controllerName[:len(controllerName)-10]), resourceNameLower,
		resourceNamePlural, strings.ToLower(controllerName[:len(controllerName)-10]),
		resourceNamePlural, strings.ToLower(controllerName[:len(controllerName)-10]), resourceNameLower,
		resourceNamePlural, strings.ToLower(controllerName[:len(controllerName)-10]), resourceNameLower,
		resourceNamePlural, strings.ToLower(controllerName[:len(controllerName)-10]), resourceNameLower,
		resourceName,
		resourceNamePlural, resourceNamePlural,
		resourceNamePlural, resourceNameLower,
		resourceNamePlural,
		resourceNamePlural, resourceNameLower,
		resourceNamePlural, resourceNameLower,
		resourceNamePlural, resourceNameLower,
	)

	ctx.Info(example)
}