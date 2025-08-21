package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support/collect"
)

// MakeCrudGenericCommand creates CRUD resources using the generic implementation
type MakeCrudGenericCommand struct{}

func (receiver *MakeCrudGenericCommand) Signature() string {
	return "make:crud-generic {name : The name of the resource} {--fields : Comma-separated list of searchable fields} {--sorts : Comma-separated list of sortable fields} {--filters : Comma-separated list of filterable fields}"
}

func (receiver *MakeCrudGenericCommand) Description() string {
	return "Create a minimal CRUD resource using generic implementations (Controller, Service)"
}

func (receiver *MakeCrudGenericCommand) Extend() command.Extend {
	return command.Extend{}
}

func (receiver *MakeCrudGenericCommand) Handle(ctx console.Context) error {
	name := ctx.Argument(0)
	if name == "" {
		return fmt.Errorf("resource name is required")
	}

	// Convert to proper case
	resourceName := strings.Title(name)
	resourceNameLower := strings.ToLower(name)
	resourceNamePlural := resourceNameLower + "s"

	// Parse options
	searchFields := parseFields(ctx.Option("fields"), []string{"name", "title"})
	sortFields := parseFields(ctx.Option("sorts"), []string{"id", "created_at", "updated_at"})
	filterFields := parseFields(ctx.Option("filters"), []string{})

	ctx.Info(fmt.Sprintf("Creating generic CRUD for resource: %s", resourceName))

	// 1. Create Service
	if err := createGenericService(resourceName, searchFields, sortFields, filterFields); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	ctx.Info("✓ Service created")

	// 2. Create Controller
	if err := createGenericController(resourceName); err != nil {
		return fmt.Errorf("failed to create controller: %w", err)
	}
	ctx.Info("✓ Controller created")

	// 3. Show next steps
	ctx.Info("")
	ctx.Info("🎉 Generic CRUD resource created successfully!")
	ctx.Info("")
	ctx.Info("Next steps:")
	ctx.Info(fmt.Sprintf("1. Create the model: go run . artisan make:model %s", resourceName))
	ctx.Info(fmt.Sprintf("2. Create requests: go run . artisan make:request %sCreateRequest", resourceName))
	ctx.Info(fmt.Sprintf("3. Create requests: go run . artisan make:request %sUpdateRequest", resourceName))
	ctx.Info(fmt.Sprintf("4. Add validation rules to the service in app/services/%s_service.go", resourceNameLower))
	ctx.Info(fmt.Sprintf("5. Add routes for /%s", resourceNamePlural))
	ctx.Info("")
	ctx.Info("Example routes:")
	ctx.Info(fmt.Sprintf(`
	%sController := %s.New%sController()
	
	// Public routes
	router.Get("/%s", %sController.Index)
	router.Get("/%s/{id}", %sController.Show) 
	router.Get("/%s/search", %sController.Search)
	
	// Protected routes
	router.Middleware(middleware.JwtAuth()).Group(func(router route.Router) {
		router.Post("/%s", %sController.Store)
		router.Put("/%s/{id}", %sController.Update)
		router.Delete("/%s/{id}", %sController.Delete)
	})
`,
		resourceNameLower, resourceNameLower, resourceName,
		resourceNamePlural, resourceNameLower,
		resourceNamePlural, resourceNameLower,
		resourceNamePlural, resourceNameLower,
		resourceNamePlural, resourceNameLower,
		resourceNamePlural, resourceNameLower,
		resourceNamePlural, resourceNameLower))

	return nil
}

func parseFields(option string, defaults []string) []string {
	if option == "" {
		return defaults
	}
	fields := strings.Split(option, ",")
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}
	return fields
}

func createGenericService(resourceName string, searchFields, sortFields, filterFields []string) error {
	resourceNameLower := strings.ToLower(resourceName)
	servicePath := filepath.Join("app", "services")
	if err := os.MkdirAll(servicePath, 0755); err != nil {
		return err
	}

	filename := filepath.Join(servicePath, resourceNameLower+"_service.go")

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("service already exists")
	}

	searchFieldsStr := formatStringSlice(searchFields)
	sortFieldsStr := formatStringSlice(sortFields)
	filterFieldsStr := formatStringSlice(filterFields)

	content := fmt.Sprintf(`package services

import (
	"fmt"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"players/app/contracts"
	"players/app/models"
)

// %sService handles %s business logic using generic CRUD implementation
type %sService struct {
	*contracts.GenericCrudService[models.%s]
}

// New%sService creates a new %s service with minimal code
func New%sService() *%sService {
	// Create the generic service
	genericService := contracts.NewGenericCrudService[models.%s]("%s", "id")
	
	// Configure the service
	genericService.
		SetSearchFields(%s).
		SetSortFields(%s).
		SetFilterFields(%s).
		SetValidationRules(map[string]interface{}{
			// TODO: Add your validation rules here
			// "name": "required|string|max:255",
			// "email": "required|email|unique:%s,email",
			// "status": "string|in:active,inactive",
		}).
		SetBeforeCreate(func(data map[string]interface{}) error {
			// TODO: Add any custom logic before creating
			// Example: Set defaults, check uniqueness, etc.
			return nil
		}).
		SetBeforeUpdate(func(id uint, data map[string]interface{}) error {
			// TODO: Add any custom logic before updating
			// Example: Check permissions, validate state transitions, etc.
			return nil
		})

	service := &%sService{
		GenericCrudService: genericService,
	}
	
	// Register service
	contracts.MustRegisterCrudService("%s", service)
	
	return service
}

// GetColumnMapping returns database column mappings
func (s *%sService) GetColumnMapping() map[string]string {
	return map[string]string{
		"id":        "id",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		// TODO: Add your column mappings here
		// "firstName": "first_name",
		// "isActive": "is_active",
	}
}

// Add any custom methods specific to %s below this line
// Example:
// func (s *%sService) GetByEmail(email string) (*models.%s, error) {
//     var %s models.%s
//     if err := facades.Orm().Query().Where("email = ?", email).First(&%s); err != nil {
//         return nil, fmt.Errorf("%s not found with email %%s: %%w", email, err)
//     }
//     return &%s, nil
// }
`,
		resourceName, resourceNameLower,
		resourceName, resourceName,
		resourceName, resourceNameLower,
		resourceName, resourceName,
		resourceName, resourceNameLower,
		searchFieldsStr,
		sortFieldsStr,
		filterFieldsStr,
		resourceNameLower+"s",
		resourceName,
		resourceNameLower+"s",
		resourceName,
		resourceNameLower,
		resourceName, resourceName,
		resourceNameLower, resourceName,
		resourceNameLower, resourceNameLower,
		resourceNameLower,
	)

	return os.WriteFile(filename, []byte(content), 0644)
}

func createGenericController(resourceName string) error {
	resourceNameLower := strings.ToLower(resourceName)
	resourceNamePlural := resourceNameLower + "s"
	controllerPath := filepath.Join("app", "http", "controllers", resourceNamePlural)
	if err := os.MkdirAll(controllerPath, 0755); err != nil {
		return err
	}

	filename := filepath.Join(controllerPath, resourceNameLower+"_controller.go")

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("controller already exists")
	}

	content := fmt.Sprintf(`package %s

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
)

// %sController handles HTTP requests for %s resources using generic CRUD
type %sController struct {
	*contracts.GenericCrudController[models.%s, requests.%sCreateRequest, requests.%sUpdateRequest]
	%sService *services.%sService
}

// New%sController creates a new %s controller with minimal code
func New%sController() *%sController {
	%sService := services.New%sService()
	
	// Create the generic controller
	genericController := contracts.NewGenericCrudController[models.%s, requests.%sCreateRequest, requests.%sUpdateRequest](
		"%s",
		%sService,
	)
	
	controller := &%sController{
		GenericCrudController: genericController,
		%sService:            %sService,
	}
	
	// Configure authorization
	genericController.SetAuthCheck(func(ctx http.Context, action string, resource interface{}) error {
		// TODO: Configure your authorization logic
		// Example: Public read, authenticated write
		if action == "viewAny" || action == "view" {
			return nil // Public access
		}
		
		// Check authentication for other actions
		permHelper := auth.GetPermissionHelper()
		user := permHelper.GetAuthenticatedUser(ctx)
		if user == nil {
			return fmt.Errorf("authentication required")
		}
		
		// Map actions to permissions
		permissionMap := map[string]string{
			"create": "%s_create",
			"update": "%s_update",
			"delete": "%s_delete",
		}
		
		if permission, ok := permissionMap[action]; ok {
			_, err := permHelper.RequirePermission(ctx, permission)
			return err
		}
		
		return nil
	})
	
	// Configure request bindings
	genericController.SetRequestBindings(
		// Bind create request
		func(ctx http.Context) (requests.%sCreateRequest, error) {
			var req requests.%sCreateRequest
			err := ctx.Request().Bind(&req)
			return req, err
		},
		// Transform create request
		func(req requests.%sCreateRequest) map[string]interface{} {
			return req.ToCreateData()
		},
		// Bind update request
		func(ctx http.Context, id uint) (requests.%sUpdateRequest, error) {
			var req requests.%sUpdateRequest
			req.ID = id
			err := ctx.Request().Bind(&req)
			return req, err
		},
		// Transform update request
		func(req requests.%sUpdateRequest) map[string]interface{} {
			return req.ToUpdateData()
		},
	)
	
	// Register controller
	contracts.MustRegisterCrudController("%s", controller)
	
	return controller
}

// Contract method implementations
func (c *%sController) GetSearchableFields() []string {
	return c.%sService.GetSearchableFields()
}

func (c *%sController) GetValidationRules() map[string]interface{} {
	return c.%sService.GetValidationRules()
}

func (c *%sController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	return c.GenericCrudController.checkAuth(ctx, permission, resource)
}

func (c *%sController) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	return permHelper.GetAuthenticatedUser(ctx)
}

func (c *%sController) RequireAuthentication(ctx http.Context) error {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func (c *%sController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	permHelper := auth.GetPermissionHelper()
	return permHelper.BuildPermissionsMap(ctx, resourceType)
}

// Add any custom endpoints beyond basic CRUD below this line
// Example:
// func (c *%sController) GetByEmail(ctx http.Context) http.Response {
//     email := ctx.Request().Query("email")
//     if email == "" {
//         return c.BadRequestResponse(ctx, "Email is required", nil)
//     }
//     
//     %s, err := c.%sService.GetByEmail(email)
//     if err != nil {
//         return c.ResourceNotFoundResponse(ctx, "%s", 0)
//     }
//     
//     return c.SuccessResponse(ctx, %s, "%s found")
// }
`,
		resourceNamePlural,
		resourceName, resourceNameLower,
		resourceName,
		resourceName, resourceName, resourceName,
		resourceNameLower, resourceName,
		resourceName, resourceNameLower,
		resourceName, resourceName,
		resourceNameLower, resourceName,
		resourceName, resourceName, resourceName,
		resourceNameLower,
		resourceNameLower,
		resourceName,
		resourceNameLower, resourceNameLower,
		resourceNamePlural,
		resourceNamePlural,
		resourceNamePlural,
		resourceName, resourceName,
		resourceName,
		resourceName, resourceName,
		resourceName,
		resourceNamePlural,
		resourceName,
		resourceNameLower,
		resourceName,
		resourceNameLower,
		resourceName,
		resourceName,
		resourceName,
		resourceName,
		resourceName,
		resourceNameLower, resourceNameLower,
		resourceNameLower,
		resourceNameLower, resourceName,
	)

	return os.WriteFile(filename, []byte(content), 0644)
}

func formatStringSlice(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	quoted := collect.Map(fields, func(field string, index int) string {
		return fmt.Sprintf(`"%s"`, field)
	})
	return strings.Join(quoted, ", ")
}
