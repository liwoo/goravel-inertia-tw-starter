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

type MakeRequestCommand struct{}

func (receiver *MakeRequestCommand) Signature() string {
	return "make:request {name : The name of the request} {--update : Create update request variant}"
}

func (receiver *MakeRequestCommand) Description() string {
	return "Create a new validation request class with comprehensive validation"
}

func (receiver *MakeRequestCommand) Extend() command.Extend {
	return command.Extend{}
}

func (receiver *MakeRequestCommand) Handle(ctx console.Context) error {
	name := ctx.Argument(0)
	if name == "" {
		return fmt.Errorf("request name is required")
	}

	// Ensure name ends with "Request"
	if !strings.HasSuffix(name, "Request") {
		name += "Request"
	}

	requestPath := filepath.Join("app", "http", "requests")
	if err := os.MkdirAll(requestPath, 0755); err != nil {
		return err
	}

	isUpdate := ctx.OptionBool("update")
	
	var filename string
	var content string
	var err error

	if isUpdate {
		// Create update request
		updateName := strings.Replace(name, "Request", "UpdateRequest", 1)
		filename = filepath.Join(requestPath, strings.ToLower(strings.ReplaceAll(updateName, "Request", ""))+"_request.go")
		content, err = generateUpdateRequestContent(updateName)
	} else {
		// Create create request (default)
		createName := strings.Replace(name, "Request", "CreateRequest", 1)
		filename = filepath.Join(requestPath, strings.ToLower(strings.ReplaceAll(createName, "Request", ""))+"_request.go")
		content, err = generateCreateRequestContent(createName)
	}

	if err != nil {
		return err
	}

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("request %s already exists", filename)
	}

	// Write file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return err
	}

	ctx.Info(fmt.Sprintf("Request [%s] created successfully", filename))
	
	if !isUpdate {
		ctx.Info("To create the corresponding update request, run:")
		baseName := strings.TrimSuffix(strings.TrimSuffix(name, "Request"), "Create")
		ctx.Info(fmt.Sprintf("  go run . artisan make:request %s --update", baseName))
	}
	
	return nil
}

func generateCreateRequestContent(requestName string) (string, error) {
	resourceName := strings.TrimSuffix(strings.TrimSuffix(requestName, "Request"), "Create")
	resourceNameLower := strings.ToLower(resourceName)

	tmpl := `package requests

import (
	"fmt"
	"players/app/contracts"
	"strings"

	"github.com/goravel/framework/contracts/http"
)

// {{.RequestName}} handles {{.ResourceNameLower}} creation validation
type {{.RequestName}} struct {
	// TODO: Add your fields here - update according to your model
	Name        string   {{.BackTick}}form:"name" json:"name"{{.BackTick}}
	Description string   {{.BackTick}}form:"description" json:"description"{{.BackTick}}
	Status      string   {{.BackTick}}form:"status" json:"status"{{.BackTick}}
	Price       float64  {{.BackTick}}form:"price" json:"price"{{.BackTick}}
	CategoryID  uint     {{.BackTick}}form:"categoryId" json:"categoryId"{{.BackTick}}
	Tags        []string {{.BackTick}}form:"tags" json:"tags"{{.BackTick}}
}

// Rules defines validation rules for {{.ResourceNameLower}} creation
func (r *{{.RequestName}}) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		// TODO: Update validation rules according to your model
		"name":        fmt.Sprintf("%s|%s", contracts.Required, fmt.Sprintf(contracts.MaxLength, 255)),
		"description": fmt.Sprintf(contracts.MaxLength, 1000),
		"status":      fmt.Sprintf("in:%s", "ACTIVE,INACTIVE,PENDING"),
		"price":       fmt.Sprintf("%s|%s|%s", contracts.Required, contracts.Numeric, fmt.Sprintf(contracts.MinValue, 0)),
		"categoryId":  fmt.Sprintf("%s|%s", contracts.Required, fmt.Sprintf(contracts.Exists, "categories", "id")),
		"tags":        fmt.Sprintf("%s|%s", contracts.Array, fmt.Sprintf(contracts.ArrayMax, 10)),
		"tags.*":      fmt.Sprintf(contracts.MaxLength, 50),
	}
}

// Messages defines custom validation messages
func (r *{{.RequestName}}) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		// TODO: Update messages according to your fields
		"name.required":        "{{.ResourceName}} name is required",
		"name.max":             "{{.ResourceName}} name cannot exceed 255 characters",
		"description.max":      "Description cannot exceed 1000 characters",
		"status.in":            "Status must be one of: ACTIVE, INACTIVE, PENDING",
		"price.required":       "Price is required",
		"price.numeric":        "Price must be a valid number",
		"price.min":            "Price must be greater than or equal to 0",
		"categoryId.required":  "Category is required",
		"categoryId.exists":    "Selected category does not exist",
		"tags.array":           "Tags must be an array",
		"tags.max":             "Maximum 10 tags allowed",
		"tags.*.max":           "Each tag cannot exceed 50 characters",
	}
}

// Attributes defines custom attribute names
func (r *{{.RequestName}}) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// TODO: Update attribute mappings
		"categoryId": "category",
		"tags.*":     "tag",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *{{.RequestName}}) Authorize(ctx http.Context) error {
	// TODO: Update authorization logic
	// Example: Check if user can create {{.ResourceNameLower}}s
	// return facades.Gate().Allows("create.{{.ResourceNameLower}}s", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *{{.RequestName}}) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic
	// Example: Normalize name
	if r.Name != "" {
		r.Name = strings.TrimSpace(r.Name)
	}

	// Set default status if not provided
	if r.Status == "" {
		r.Status = "ACTIVE"
	}

	return nil
}

// PassedValidation is called after validation passes
func (r *{{.RequestName}}) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic
	return nil
}

// ToCreateData converts the request to create data map
func (r *{{.RequestName}}) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		// TODO: Update according to your fields
		"name":        r.Name,
		"description": r.Description,
		"status":      r.Status,
		"price":       r.Price,
		"categoryId":  r.CategoryID,
	}

	// Only include optional fields if they have values
	if len(r.Tags) > 0 {
		data["tags"] = r.Tags
	}

	return data
}
`

	t, err := template.New("createRequest").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = t.Execute(&result, map[string]string{
		"RequestName":        requestName,
		"ResourceName":       resourceName,
		"ResourceNameLower":  resourceNameLower,
		"BackTick":           "`",
	})

	return result.String(), err
}

func generateUpdateRequestContent(requestName string) (string, error) {
	resourceName := strings.TrimSuffix(strings.TrimSuffix(requestName, "Request"), "Update")
	resourceNameLower := strings.ToLower(resourceName)

	tmpl := `package requests

import (
	"fmt"
	"players/app/contracts"
	"strings"

	"github.com/goravel/framework/contracts/http"
)

// {{.RequestName}} handles {{.ResourceNameLower}} update validation
type {{.RequestName}} struct {
	// TODO: Add your fields here with pointers for optional updates
	Name        *string   {{.BackTick}}form:"name" json:"name"{{.BackTick}}
	Description *string   {{.BackTick}}form:"description" json:"description"{{.BackTick}}
	Status      *string   {{.BackTick}}form:"status" json:"status"{{.BackTick}}
	Price       *float64  {{.BackTick}}form:"price" json:"price"{{.BackTick}}
	CategoryID  *uint     {{.BackTick}}form:"categoryId" json:"categoryId"{{.BackTick}}
	Tags        *[]string {{.BackTick}}form:"tags" json:"tags"{{.BackTick}}
	ID          uint      {{.BackTick}}form:"-" json:"-"{{.BackTick}} // Set by controller
}

// Rules defines validation rules for {{.ResourceNameLower}} updates
func (r *{{.RequestName}}) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Only validate fields that are provided
	// TODO: Update validation rules according to your model
	if r.Name != nil {
		rules["name"] = fmt.Sprintf(contracts.MaxLength, 255)
	}
	if r.Description != nil {
		rules["description"] = fmt.Sprintf(contracts.MaxLength, 1000)
	}
	if r.Status != nil {
		rules["status"] = "in:ACTIVE,INACTIVE,PENDING"
	}
	if r.Price != nil {
		rules["price"] = fmt.Sprintf("%s|%s", contracts.Numeric, fmt.Sprintf(contracts.MinValue, 0))
	}
	if r.CategoryID != nil {
		rules["categoryId"] = fmt.Sprintf(contracts.Exists, "categories", "id")
	}
	if r.Tags != nil {
		rules["tags"] = fmt.Sprintf("%s|%s", contracts.Array, fmt.Sprintf(contracts.ArrayMax, 10))
		rules["tags.*"] = fmt.Sprintf(contracts.MaxLength, 50)
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *{{.RequestName}}) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		// TODO: Update messages according to your fields
		"name.max":             "{{.ResourceName}} name cannot exceed 255 characters",
		"description.max":      "Description cannot exceed 1000 characters",
		"status.in":            "Status must be one of: ACTIVE, INACTIVE, PENDING",
		"price.numeric":        "Price must be a valid number",
		"price.min":            "Price must be greater than or equal to 0",
		"categoryId.exists":    "Selected category does not exist",
		"tags.array":           "Tags must be an array",
		"tags.max":             "Maximum 10 tags allowed",
		"tags.*.max":           "Each tag cannot exceed 50 characters",
	}
}

// Attributes defines custom attribute names for updates
func (r *{{.RequestName}}) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// TODO: Update attribute mappings
		"categoryId": "category",
		"tags.*":     "tag",
	}
}

// Authorize determines if the user is authorized to update this {{.ResourceNameLower}}
func (r *{{.RequestName}}) Authorize(ctx http.Context) error {
	// TODO: Update authorization logic
	// Example: Check if user can update this specific {{.ResourceNameLower}}
	// return facades.Gate().Allows("update.{{.ResourceNameLower}}s", {{.ResourceNameLower}})
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *{{.RequestName}}) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic
	// Example: Normalize name if provided
	if r.Name != nil && *r.Name != "" {
		normalized := strings.TrimSpace(*r.Name)
		r.Name = &normalized
	}

	return nil
}

// PassedValidation is called after validation passes
func (r *{{.RequestName}}) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic
	return nil
}

// ToUpdateData converts the request to update data map
func (r *{{.RequestName}}) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only include fields that are provided (not nil)
	// TODO: Update according to your fields
	if r.Name != nil {
		data["name"] = *r.Name
	}
	if r.Description != nil {
		data["description"] = *r.Description
	}
	if r.Status != nil {
		data["status"] = *r.Status
	}
	if r.Price != nil {
		data["price"] = *r.Price
	}
	if r.CategoryID != nil {
		data["categoryId"] = *r.CategoryID
	}
	if r.Tags != nil {
		data["tags"] = *r.Tags
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *{{.RequestName}}) GetResourceID() interface{} {
	return r.ID
}
`

	t, err := template.New("updateRequest").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = t.Execute(&result, map[string]string{
		"RequestName":        requestName,
		"ResourceName":       resourceName,
		"ResourceNameLower":  resourceNameLower,
		"BackTick":           "`",
	})

	return result.String(), err
}