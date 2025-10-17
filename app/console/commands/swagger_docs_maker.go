package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type SwaggerDocsMaker struct {
}

// Signature The name and signature of the console command.
func (r *SwaggerDocsMaker) Signature() string {
	return "make:swagger-docs"
}

// Description The console command description.
func (r *SwaggerDocsMaker) Description() string {
	return "Generate Swagger documentation wrapper methods for a CRUD controller"
}

// Extend The console command extend.
func (r *SwaggerDocsMaker) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:     "controller",
				Aliases:  []string{"c"},
				Usage:    "Controller name (e.g., Book, Lender, User)",
				Required: true,
			},
			&command.StringFlag{
				Name:    "tag",
				Aliases: []string{"t"},
				Usage:   "Swagger tag (defaults to lowercase plural of controller)",
			},
		},
	}
}

// Handle Execute the console command.
func (r *SwaggerDocsMaker) Handle(ctx console.Context) error {
	controllerName := ctx.Option("controller")
	if controllerName == "" {
		return fmt.Errorf("controller name is required (use --controller=Book)")
	}

	// Capitalize controller name
	controllerName = strings.Title(strings.ToLower(controllerName))

	// Determine tag (plural, lowercase)
	tag := ctx.Option("tag")
	if tag == "" {
		tag = strings.ToLower(controllerName) + "s"
	}

	// Determine resource name (lowercase)
	resourceName := strings.ToLower(controllerName)
	resourceNamePlural := tag

	// Find the controller file
	controllerDir := filepath.Join("app", "http", "controllers", resourceNamePlural)
	controllerFile := filepath.Join(controllerDir, resourceName+"_controller.go")

	// Check if file exists
	if _, err := os.Stat(controllerFile); os.IsNotExist(err) {
		return fmt.Errorf("controller file not found: %s", controllerFile)
	}

	// Read existing file
	content, err := os.ReadFile(controllerFile)
	if err != nil {
		return fmt.Errorf("failed to read controller file: %v", err)
	}

	fileContent := string(content)

	// Check if already has wrapper methods
	if strings.Contains(fileContent, "// CRUD Methods with Swagger annotations") {
		ctx.Warning("Controller already has Swagger wrapper methods")
		return nil
	}

	// Generate wrapper methods
	wrapperMethods := r.generateWrapperMethods(controllerName, resourceName, resourceNamePlural, tag)

	// Find insertion point (after NewController function, before custom methods)
	insertMarker := fmt.Sprintf("func New%sController() *%sController {", controllerName, controllerName)
	insertPos := strings.Index(fileContent, insertMarker)
	if insertPos == -1 {
		return fmt.Errorf("could not find New%sController function", controllerName)
	}

	// Find the end of the NewController function
	closingBracePos := r.findFunctionEnd(fileContent, insertPos)
	if closingBracePos == -1 {
		return fmt.Errorf("could not find end of New%sController function", controllerName)
	}

	// Insert wrapper methods after the NewController function
	newContent := fileContent[:closingBracePos+1] + "\n\n" + wrapperMethods + "\n" + fileContent[closingBracePos+1:]

	// Write back to file
	if err := os.WriteFile(controllerFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write controller file: %v", err)
	}

	ctx.Success(fmt.Sprintf("✓ Swagger wrapper methods added to %s", controllerFile))
	ctx.NewLine()
	ctx.Info("Run 'go run . artisan swagger:generate' to update Swagger docs")

	return nil
}

// generateWrapperMethods generates the wrapper methods with Swagger annotations
func (r *SwaggerDocsMaker) generateWrapperMethods(controllerName, resourceName, resourceNamePlural, tag string) string {
	// Capitalize first letter for type names
	modelType := "models." + strings.Title(resourceName)
	createRequestType := "requests." + controllerName + "CreateRequest"
	updateRequestType := "requests." + controllerName + "UpdateRequest"

	return fmt.Sprintf(`// ============================================================================
// CRUD Methods with Swagger annotations
// ============================================================================

// Index godoc
// @Summary      List all %s
// @Description  Get paginated list of %s with filtering, sorting, and search
// @Tags         %s
// @Accept       json
// @Produce      json
// @Param        page      query  int     false  "Page number" default(1)
// @Param        pageSize  query  int     false  "Items per page" default(20) Enums(5, 10, 20, 30, 50, 100)
// @Param        search    query  string  false  "Search query"
// @Param        sort      query  string  false  "Sort field"
// @Param        direction query  string  false  "Sort direction" Enums(ASC, DESC)
// @Success      200  {object}  contracts.ResponseFormat{data=contracts.PaginatedResponse{data=[]%s}}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /%s [get]
func (c *%sController) Index(ctx http.Context) http.Response {
	return c.CrudController.Index(ctx)
}

// Show godoc
// @Summary      Get %s by ID
// @Description  Retrieve a specific %s
// @Tags         %s
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "%s ID"
// @Success      200  {object}  contracts.ResponseFormat{data=%s}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Router       /%s/{id} [get]
func (c *%sController) Show(ctx http.Context) http.Response {
	return c.CrudController.Show(ctx)
}

// Store godoc
// @Summary      Create a new %s
// @Description  Create a %s with validation
// @Tags         %s
// @Accept       json
// @Produce      json
// @Param        %s  body  %s  true  "%s data"
// @Success      201  {object}  contracts.ResponseFormat{data=%s}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      422  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /%s [post]
func (c *%sController) Store(ctx http.Context) http.Response {
	return c.CrudController.Store(ctx)
}

// Update godoc
// @Summary      Update a %s
// @Description  Update an existing %s
// @Tags         %s
// @Accept       json
// @Produce      json
// @Param        id    path  int  true  "%s ID"
// @Param        %s  body  %s  true  "%s data"
// @Success      200  {object}  contracts.ResponseFormat{data=%s}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Failure      422  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /%s/{id} [put]
func (c *%sController) Update(ctx http.Context) http.Response {
	return c.CrudController.Update(ctx)
}

// Delete godoc
// @Summary      Delete a %s
// @Description  Delete a %s by ID
// @Tags         %s
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "%s ID"
// @Success      204  "%s deleted"
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /%s/{id} [delete]
func (c *%sController) Delete(ctx http.Context) http.Response {
	return c.CrudController.Delete(ctx)
}

// Search godoc
// @Summary      Search %s
// @Description  Full-text search for %s
// @Tags         %s
// @Accept       json
// @Produce      json
// @Param        q         query  string  true   "Search query (min 2 chars)"
// @Param        page      query  int     false  "Page number"
// @Param        pageSize  query  int     false  "Items per page"
// @Success      200  {object}  contracts.ResponseFormat{data=contracts.PaginatedResponse{data=[]%s}}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /%s/search [get]
func (c *%sController) Search(ctx http.Context) http.Response {
	return c.CrudController.Search(ctx)
}

// FilterMetadata godoc
// @Summary      Get %s filter metadata
// @Description  Returns available filters for %s
// @Tags         %s
// @Accept       json
// @Produce      json
// @Success      200  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /%s/filters [get]
func (c *%sController) FilterMetadata(ctx http.Context) http.Response {
	return c.CrudController.FilterMetadata(ctx)
}`,
		// Index args: 6 (added modelType)
		resourceNamePlural, resourceNamePlural, tag, modelType, resourceNamePlural, controllerName,
		// Show args: 8 (added modelType)
		resourceName, resourceName, tag, controllerName, modelType, resourceNamePlural, controllerName,
		// Store args: 10 (added createRequestType and modelType)
		resourceName, resourceName, tag, resourceName, createRequestType, resourceName, modelType, resourceNamePlural, controllerName,
		// Update args: 12 (added updateRequestType and modelType)
		resourceName, resourceName, tag, controllerName, resourceName, updateRequestType, resourceName, modelType, resourceNamePlural, controllerName,
		// Delete args: 7 (unchanged)
		resourceName, resourceName, tag, controllerName, controllerName, resourceNamePlural, controllerName,
		// Search args: 6 (added modelType)
		resourceNamePlural, resourceNamePlural, tag, modelType, resourceNamePlural, controllerName,
		// FilterMetadata args: 5 (unchanged)
		resourceName, resourceNamePlural, tag, resourceNamePlural, controllerName,
	)
}

// findFunctionEnd finds the closing brace of a function
func (r *SwaggerDocsMaker) findFunctionEnd(content string, startPos int) int {
	braceCount := 0
	inFunction := false

	for i := startPos; i < len(content); i++ {
		char := content[i]

		if char == '{' {
			braceCount++
			inFunction = true
		} else if char == '}' {
			braceCount--
			if inFunction && braceCount == 0 {
				return i
			}
		}
	}

	return -1
}
