package commands

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type RequestMaker struct {
}

// Signature The name and signature of the console command.
func (r *RequestMaker) Signature() string {
	return "make:req"
}

// Description The console command description.
func (r *RequestMaker) Description() string {
	return "Generate Create and Update request classes based on a model"
}

// Extend The console command extend.
func (r *RequestMaker) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "The model to base the requests on",
			},
			&command.BoolFlag{
				Name:    "create-only",
				Aliases: []string{"c"},
				Usage:   "Generate only the Create request",
			},
			&command.BoolFlag{
				Name:    "update-only",
				Aliases: []string{"u"},
				Usage:   "Generate only the Update request",
			},
		},
	}
}

// Handle Execute the console command.
func (r *RequestMaker) Handle(ctx console.Context) error {
	// Get the request name from argument
	requestName := ctx.Argument(0)
	if requestName == "" {
		return fmt.Errorf("request name is required")
	}

	// Get the model name from flag or derive from request name
	modelName := ctx.Option("model")
	if modelName == "" {
		// Derive model name from request name
		// e.g., "LenderRequest" -> "Lender", "lender" -> "Lender"
		modelName = strings.TrimSuffix(requestName, "Request")
		modelName = strings.Title(strings.ToLower(modelName))
	}

	// Ensure proper capitalization
	requestBaseName := strings.TrimSuffix(requestName, "Request")
	requestBaseName = strings.Title(strings.ToLower(requestBaseName))

	// Determine which requests to generate
	createOnly := ctx.OptionBool("create-only")
	updateOnly := ctx.OptionBool("update-only")

	// Create requests directory if it doesn't exist
	requestsDir := filepath.Join("app", "http", "requests")
	if err := os.MkdirAll(requestsDir, 0755); err != nil {
		return err
	}

	// Find and introspect the model
	modelPath := filepath.Join("app", "models", strings.ToLower(modelName)+".go")
	modelInfo, err := r.introspectModel(modelPath, modelName)
	if err != nil {
		ctx.Warning(fmt.Sprintf("Could not introspect model %s: %v", modelName, err))
		ctx.Warning("Generating requests with default configuration...")
		modelInfo = &RequestModelInfo{
			Name:   modelName,
			Fields: []RequestFieldInfo{},
		}
	}

	filesCreated := []string{}

	// Generate Create request unless update-only is specified
	if !updateOnly {
		createFilename := filepath.Join(requestsDir, strings.ToLower(requestBaseName)+"_create_request.go")

		// Check if file already exists
		if _, err := os.Stat(createFilename); err == nil {
			ctx.Warning(fmt.Sprintf("Create request %s already exists, skipping...", createFilename))
		} else {
			content := r.generateCreateRequest(requestBaseName, modelInfo)
			if err := os.WriteFile(createFilename, []byte(content), 0644); err != nil {
				return err
			}
			filesCreated = append(filesCreated, createFilename)
		}
	}

	// Generate Update request unless create-only is specified
	if !createOnly {
		updateFilename := filepath.Join(requestsDir, strings.ToLower(requestBaseName)+"_update_request.go")

		// Check if file already exists
		if _, err := os.Stat(updateFilename); err == nil {
			ctx.Warning(fmt.Sprintf("Update request %s already exists, skipping...", updateFilename))
		} else {
			content := r.generateUpdateRequest(requestBaseName, modelInfo)
			if err := os.WriteFile(updateFilename, []byte(content), 0644); err != nil {
				return err
			}
			filesCreated = append(filesCreated, updateFilename)
		}
	}

	// Show success messages
	if len(filesCreated) == 0 {
		ctx.Warning("No files were created (all files already exist)")
		return nil
	}

	for _, file := range filesCreated {
		ctx.Success(fmt.Sprintf("Request created: %s", file))
	}

	ctx.NewLine()
	ctx.Info("Next steps:")
	ctx.Line("1. Customize the validation rules in the Rules() method")
	ctx.Line("2. Add custom validation messages in Messages() method")
	ctx.Line("3. Implement authorization logic in Authorize() method if needed")
	ctx.Line("4. Use the requests in your controller")

	return nil
}

type RequestModelInfo struct {
	Name   string
	Fields []RequestFieldInfo
}

type RequestFieldInfo struct {
	Name     string
	Type     string
	JSONTag  string
	FormTag  string
	Optional bool // true if pointer type
}

// introspectModel parses the model file and extracts field information
func (r *RequestMaker) introspectModel(modelPath, modelName string) (*RequestModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, modelPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	modelInfo := &RequestModelInfo{
		Name:   modelName,
		Fields: []RequestFieldInfo{},
	}

	// Find the struct definition
	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok || typeSpec.Name.Name != modelName {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Extract fields
		for _, field := range structType.Fields.List {
			// Skip embedded fields
			if field.Names == nil || len(field.Names) == 0 {
				continue
			}

			for _, name := range field.Names {
				// Skip internal fields
				if name.Name == "ID" || strings.HasPrefix(name.Name, "Created") ||
					strings.HasPrefix(name.Name, "Updated") || strings.HasPrefix(name.Name, "Deleted") {
					continue
				}

				fieldInfo := RequestFieldInfo{
					Name: name.Name,
				}

				// Determine if field is optional (pointer)
				if starExpr, ok := field.Type.(*ast.StarExpr); ok {
					fieldInfo.Optional = true
					fieldInfo.Type = "*" + r.getTypeName(starExpr.X)
				} else {
					fieldInfo.Type = r.getTypeName(field.Type)
				}

				// Extract tags
				if field.Tag != nil {
					tag := strings.Trim(field.Tag.Value, "`")
					fieldInfo.JSONTag = r.extractTag(tag, "json")
					fieldInfo.FormTag = r.extractTag(tag, "form")

					// Use json tag if no form tag
					if fieldInfo.FormTag == "" && fieldInfo.JSONTag != "" {
						fieldInfo.FormTag = fieldInfo.JSONTag
					}
					// Use field name in snake_case if no tags
					if fieldInfo.JSONTag == "" {
						fieldInfo.JSONTag = r.toSnakeCase(name.Name)
					}
					if fieldInfo.FormTag == "" {
						fieldInfo.FormTag = r.toSnakeCase(name.Name)
					}
				} else {
					// No tags, use field name in snake_case
					fieldInfo.JSONTag = r.toSnakeCase(name.Name)
					fieldInfo.FormTag = r.toSnakeCase(name.Name)
				}

				// Skip fields with "-" tag
				if fieldInfo.JSONTag == "-" {
					continue
				}

				modelInfo.Fields = append(modelInfo.Fields, fieldInfo)
			}
		}

		return false
	})

	return modelInfo, nil
}

// getTypeName extracts the type name from an AST expression
func (r *RequestMaker) getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", r.getTypeName(t.X), t.Sel.Name)
	case *ast.ArrayType:
		return "[]" + r.getTypeName(t.Elt)
	default:
		return "interface{}"
	}
}

// extractTag extracts a specific tag value from a struct tag string
func (r *RequestMaker) extractTag(tagString, tagName string) string {
	parts := strings.Split(tagString, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, tagName+":") {
			value := strings.TrimPrefix(part, tagName+":")
			return strings.Trim(value, `"`)
		}
	}
	return ""
}

// toSnakeCase converts CamelCase to snake_case
func (r *RequestMaker) toSnakeCase(s string) string {
	var result strings.Builder
	for i, char := range s {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}

// generateCreateRequest generates the Create request content
func (r *RequestMaker) generateCreateRequest(baseName string, model *RequestModelInfo) string {
	resourceName := strings.ToLower(baseName)

	// Generate struct fields
	structFields := r.generateStructFields(model.Fields, false)

	// Generate validation rules
	validationRules := r.generateValidationRules(model.Fields, false)

	// Generate validation messages
	validationMessages := r.generateValidationMessages(model.Fields, resourceName)

	// Generate ToCreateData fields
	createDataFields := r.generateCreateDataFields(model.Fields)

	return fmt.Sprintf(`package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// %sCreateRequest handles %s creation validation
type %sCreateRequest struct {
%s
}

// Rules defines validation rules for %s creation
func (r *%sCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
%s
	}
}

// Messages defines custom validation messages
func (r *%sCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
%s
	}
}

// Attributes defines custom attribute names
func (r *%sCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *%sCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create %s
	// return facades.Gate().Allows("create.%ss", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *%sCreateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic
	// Example: Normalize, trim, or set default values
	return nil
}

// PassedValidation is called after validation passes
func (r *%sCreateRequest) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic if needed
	return nil
}

// ToCreateData converts the request to create data map
func (r *%sCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
%s
	}

	return data
}
`,
		baseName, resourceName, baseName,
		structFields,
		resourceName, baseName,
		validationRules,
		baseName,
		validationMessages,
		baseName,
		baseName, resourceName, resourceName,
		baseName,
		baseName,
		baseName,
		createDataFields,
	)
}

// generateUpdateRequest generates the Update request content
func (r *RequestMaker) generateUpdateRequest(baseName string, model *RequestModelInfo) string {
	resourceName := strings.ToLower(baseName)

	// Generate struct fields (with pointers for optional updates)
	structFields := r.generateStructFields(model.Fields, true)

	// Generate validation rules (conditional based on presence)
	validationRules := r.generateUpdateValidationRules(model.Fields)

	// Generate validation messages
	validationMessages := r.generateValidationMessages(model.Fields, resourceName)

	// Generate ToUpdateData fields
	updateDataFields := r.generateUpdateDataFields(model.Fields)

	return fmt.Sprintf(`package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// %sUpdateRequest handles %s update validation
type %sUpdateRequest struct {
%s
	ID uint `+"`"+`form:"-" json:"-"`+"`"+` // Set by controller
}

// Rules defines validation rules for %s updates
func (r *%sUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

%s

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *%sUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
%s
	}
}

// Attributes defines custom attribute names for updates
func (r *%sUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to update this %s
func (r *%sUpdateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user can update this specific %s
	// return facades.Gate().Allows("update.%ss", %s)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *%sUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic for updates
	// Example: Normalize data if provided
	return nil
}

// PassedValidation is called after validation passes
func (r *%sUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *%sUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

%s

	return data
}

// GetResourceID returns the resource ID for update
func (r *%sUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
`,
		baseName, resourceName, baseName,
		structFields,
		resourceName, baseName,
		validationRules,
		baseName,
		validationMessages,
		baseName,
		resourceName, baseName,
		resourceName, resourceName, resourceName,
		baseName,
		baseName,
		baseName,
		updateDataFields,
		baseName,
	)
}

// generateStructFields generates struct field definitions
func (r *RequestMaker) generateStructFields(fields []RequestFieldInfo, isUpdate bool) string {
	if len(fields) == 0 {
		return "\t// TODO: Add your fields here"
	}

	lines := []string{}
	for _, field := range fields {
		goType := field.Type
		if isUpdate {
			// For updates, ensure all fields are pointers (optional)
			// If the field is already optional (pointer), keep it as is
			// Otherwise, make it a pointer
			if !field.Optional {
				goType = "*" + goType
			}
		}

		line := fmt.Sprintf("\t%s %s `form:\"%s\" json:\"%s\"`",
			field.Name, goType, field.FormTag, field.JSONTag)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// generateValidationRules generates validation rules for create requests
func (r *RequestMaker) generateValidationRules(fields []RequestFieldInfo, isUpdate bool) string {
	if len(fields) == 0 {
		return "\t\t// TODO: Add your validation rules here"
	}

	lines := []string{}
	for _, field := range fields {
		rule := r.generateFieldValidationRule(field)
		if rule != "" {
			lines = append(lines, fmt.Sprintf("\t\t\"%s\": \"%s\",", field.JSONTag, rule))
		}
	}

	if len(lines) == 0 {
		return "\t\t// TODO: Add your validation rules here"
	}

	return strings.Join(lines, "\n")
}

// generateUpdateValidationRules generates conditional validation rules for update requests
func (r *RequestMaker) generateUpdateValidationRules(fields []RequestFieldInfo) string {
	if len(fields) == 0 {
		return "\t// TODO: Add your conditional validation rules here"
	}

	lines := []string{}
	for _, field := range fields {
		rule := r.generateFieldValidationRule(field)
		if rule != "" {
			lines = append(lines, fmt.Sprintf("\t// Only validate %s if provided", field.Name))
			lines = append(lines, fmt.Sprintf("\tif r.%s != nil {", field.Name))
			lines = append(lines, fmt.Sprintf("\t\trules[\"%s\"] = \"%s\"", field.JSONTag, rule))
			lines = append(lines, "\t}")
		}
	}

	if len(lines) == 0 {
		return "\t// TODO: Add your conditional validation rules here"
	}

	return strings.Join(lines, "\n")
}

// generateFieldValidationRule generates a validation rule for a single field
func (r *RequestMaker) generateFieldValidationRule(field RequestFieldInfo) string {
	rules := []string{}

	// Add required rule for non-optional fields
	if !field.Optional {
		rules = append(rules, "required")
	}

	// Type-specific rules
	switch field.Type {
	case "string":
		rules = append(rules, "string", "max:255")
	case "int", "uint", "int64", "uint64":
		rules = append(rules, "numeric")
	case "float64", "float32":
		rules = append(rules, "numeric")
	case "bool":
		rules = append(rules, "boolean")
	case "[]string":
		rules = append(rules, "array")
	}

	return strings.Join(rules, "|")
}

// generateValidationMessages generates validation messages
func (r *RequestMaker) generateValidationMessages(fields []RequestFieldInfo, resourceName string) string {
	if len(fields) == 0 {
		return "\t\t// TODO: Add your custom validation messages here"
	}

	lines := []string{}
	for _, field := range fields {
		fieldDisplay := strings.ReplaceAll(field.JSONTag, "_", " ")
		fieldDisplay = strings.Title(fieldDisplay)

		if !field.Optional {
			lines = append(lines, fmt.Sprintf("\t\t\"%s.required\": \"%s is required\",", field.JSONTag, fieldDisplay))
		}

		switch field.Type {
		case "string":
			lines = append(lines, fmt.Sprintf("\t\t\"%s.max\": \"%s cannot exceed 255 characters\",", field.JSONTag, fieldDisplay))
		case "int", "uint", "int64", "uint64", "float64", "float32":
			lines = append(lines, fmt.Sprintf("\t\t\"%s.numeric\": \"%s must be a valid number\",", field.JSONTag, fieldDisplay))
		}
	}

	if len(lines) == 0 {
		return "\t\t// TODO: Add your custom validation messages here"
	}

	return strings.Join(lines, "\n")
}

// generateCreateDataFields generates ToCreateData field mappings
func (r *RequestMaker) generateCreateDataFields(fields []RequestFieldInfo) string {
	if len(fields) == 0 {
		return "\t\t// TODO: Add your fields here"
	}

	lines := []string{}
	for _, field := range fields {
		lines = append(lines, fmt.Sprintf("\t\t\"%s\": r.%s,", field.JSONTag, field.Name))
	}

	return strings.Join(lines, "\n")
}

// generateUpdateDataFields generates ToUpdateData field mappings with nil checks
func (r *RequestMaker) generateUpdateDataFields(fields []RequestFieldInfo) string {
	if len(fields) == 0 {
		return "\t// TODO: Add your fields here with nil checks"
	}

	lines := []string{}
	for _, field := range fields {
		lines = append(lines, fmt.Sprintf("\t// Only include %s if provided", field.Name))
		lines = append(lines, fmt.Sprintf("\tif r.%s != nil {", field.Name))
		lines = append(lines, fmt.Sprintf("\t\tdata[\"%s\"] = *r.%s", field.JSONTag, field.Name))
		lines = append(lines, "\t}")
	}

	return strings.Join(lines, "\n")
}
