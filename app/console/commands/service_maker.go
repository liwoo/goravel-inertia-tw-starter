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

type ServiceMaker struct {
}

// Signature The name and signature of the console command.
func (r *ServiceMaker) Signature() string {
	return "make:svc"
}

// Description The console command description.
func (r *ServiceMaker) Description() string {
	return "Create a new service based on a model"
}

// Extend The console command extend.
func (r *ServiceMaker) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "The model to base the service on",
			},
		},
	}
}

// Handle Execute the console command.
func (r *ServiceMaker) Handle(ctx console.Context) error {
	// Get the service name from argument
	serviceName := ctx.Argument(0)
	if serviceName == "" {
		return fmt.Errorf("service name is required")
	}

	// Get the model name from flag or derive from service name
	modelName := ctx.Option("model")
	if modelName == "" {
		// Derive model name from service name
		// e.g., "LenderService" -> "Lender", "lender" -> "Lender"
		modelName = strings.TrimSuffix(serviceName, "Service")
		modelName = strings.Title(strings.ToLower(modelName))
	}

	// Ensure service name is properly capitalized and ends with "Service"
	serviceName = strings.TrimSuffix(serviceName, "Service")
	serviceName = strings.Title(strings.ToLower(serviceName)) + "Service"

	// Handle folder structure (e.g., "acme/lender" -> "app/services/acme")
	var serviceDir string
	var servicePackage string
	if strings.Contains(serviceName, "/") {
		parts := strings.Split(serviceName, "/")
		folderPath := strings.Join(parts[:len(parts)-1], "/")
		serviceName = parts[len(parts)-1]
		serviceDir = filepath.Join("app", "services", folderPath)
		servicePackage = parts[len(parts)-2] // Use the immediate parent folder as package
	} else {
		serviceDir = filepath.Join("app", "services")
		servicePackage = "services"
	}

	// Create service directory if it doesn't exist
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return err
	}

	// Generate filename
	filename := filepath.Join(serviceDir, strings.ToLower(strings.TrimSuffix(serviceName, "Service"))+"_service.go")

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("service %s already exists", filename)
	}

	// Find and introspect the model
	modelPath := filepath.Join("app", "models", strings.ToLower(modelName)+".go")
	modelInfo, err := r.introspectModel(modelPath, modelName)
	if err != nil {
		ctx.Warning(fmt.Sprintf("Could not introspect model %s: %v", modelName, err))
		ctx.Warning("Generating service with default configuration...")
		modelInfo = &ModelInfo{
			Name:          modelName,
			TableName:     strings.ToLower(modelName) + "s",
			Fields:        []FieldInfo{},
			HasSoftDelete: false,
		}
	}

	// Generate service content
	content := r.generateServiceContent(serviceName, servicePackage, modelInfo)

	// Write file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return err
	}

	ctx.Success(fmt.Sprintf("Service [%s] created successfully", filename))
	ctx.NewLine()

	// Generate registration instructions
	r.printRegistrationInstructions(ctx, modelInfo.TableName, serviceName)

	return nil
}

type ModelInfo struct {
	Name          string
	TableName     string
	Fields        []FieldInfo
	HasSoftDelete bool
}

type FieldInfo struct {
	Name     string
	Type     string
	JSONTag  string
	DBTag    string
	Optional bool // true if pointer type
}

// introspectModel parses the model file and extracts field information
func (r *ServiceMaker) introspectModel(modelPath, modelName string) (*ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, modelPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	modelInfo := &ModelInfo{
		Name:   modelName,
		Fields: []FieldInfo{},
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

		// Check for embedded types
		for _, field := range structType.Fields.List {
			if field.Names == nil || len(field.Names) == 0 {
				// Embedded field
				if ident, ok := field.Type.(*ast.SelectorExpr); ok {
					if ident.Sel.Name == "SoftDeletes" {
						modelInfo.HasSoftDelete = true
					}
				}
				continue
			}

			for _, name := range field.Names {
				fieldInfo := FieldInfo{
					Name: name.Name,
				}

				// Determine if field is optional (pointer)
				if starExpr, ok := field.Type.(*ast.StarExpr); ok {
					fieldInfo.Optional = true
					fieldInfo.Type = r.getTypeName(starExpr.X)
				} else {
					fieldInfo.Type = r.getTypeName(field.Type)
				}

				// Extract tags
				if field.Tag != nil {
					tag := strings.Trim(field.Tag.Value, "`")
					fieldInfo.JSONTag = r.extractTag(tag, "json")
					fieldInfo.DBTag = r.extractTag(tag, "db")
				}

				modelInfo.Fields = append(modelInfo.Fields, fieldInfo)
			}
		}

		return false
	})

	// Find TableName method if exists
	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != "TableName" {
			return true
		}

		// Extract table name from return statement
		if funcDecl.Body != nil {
			for _, stmt := range funcDecl.Body.List {
				if returnStmt, ok := stmt.(*ast.ReturnStmt); ok {
					if len(returnStmt.Results) > 0 {
						if basicLit, ok := returnStmt.Results[0].(*ast.BasicLit); ok {
							modelInfo.TableName = strings.Trim(basicLit.Value, `"`)
						}
					}
				}
			}
		}

		return false
	})

	// Default table name if not found
	if modelInfo.TableName == "" {
		modelInfo.TableName = strings.ToLower(modelName) + "s"
	}

	return modelInfo, nil
}

// getTypeName extracts the type name from an AST expression
func (r *ServiceMaker) getTypeName(expr ast.Expr) string {
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
func (r *ServiceMaker) extractTag(tagString, tagName string) string {
	parts := strings.Split(tagString, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, tagName+":") {
			value := strings.TrimPrefix(part, tagName+":")
			return strings.Trim(value, `"`)
		}
	}
	return ""
}

// generateServiceContent generates the service file content
func (r *ServiceMaker) generateServiceContent(serviceName, packageName string, model *ModelInfo) string {
	// Generate search fields from string fields
	searchFields := []string{}
	for _, field := range model.Fields {
		dbTag := field.DBTag
		if dbTag == "" {
			dbTag = strings.ToLower(field.Name)
		}
		if field.Type == "string" && dbTag != "" {
			searchFields = append(searchFields, dbTag)
		}
	}

	// Generate sort fields from all non-slice fields
	sortFields := []string{"id", "created_at", "updated_at"}
	for _, field := range model.Fields {
		dbTag := field.DBTag
		if dbTag == "" {
			dbTag = strings.ToLower(field.Name)
		}
		if dbTag != "" && !strings.HasPrefix(field.Type, "[]") {
			sortFields = append(sortFields, dbTag)
		}
	}

	// Generate filter fields (enum-like string fields and numeric fields)
	filterFields := []string{}
	for _, field := range model.Fields {
		dbTag := field.DBTag
		if dbTag == "" {
			dbTag = strings.ToLower(field.Name)
		}
		if (field.Type == "string" || field.Type == "int" || field.Type == "uint") && dbTag != "" {
			filterFields = append(filterFields, dbTag)
		}
	}

	// Generate validation rules
	validationRules := map[string]string{}
	for _, field := range model.Fields {
		jsonTag := field.JSONTag
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		rule := r.generateValidationRule(field)
		if rule != "" {
			validationRules[jsonTag] = rule
		}
	}

	// Build the service content
	serviceVarName := strings.ToLower(string(serviceName[0])) + serviceName[1:]

	content := fmt.Sprintf(`package %s

import (
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
)

// %s implements business logic for %s using the builder pattern
type %s struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// New%s creates a new %s service using the builder pattern
func New%s() *%s {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.%s]("%s", "id").
		WithSearchFields(%s). // Fields that will be searchable via the search query parameter
		WithSortFields(%s). // Fields that can be used for sorting results
		WithFilterFields(%s). // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{ // Validation rules for create/update operations
%s
		}).
		WithDefaultSort("created_at", "DESC"). // Default sorting when none specified
		WithScopeFiltering("%s", "created_by"). // Enable permission-based filtering
%s
		Build() // Returns a fully configured CrudServiceContract

	%sInstance := &%s{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, %sInstance, "%s")

	return %sInstance
}

// GetFilterDefinitions returns filter definitions for the %s resource
// This method is OPTIONAL - only implement if you need custom filter UI components
// Uncomment and customize the implementation below if needed:
//
// func (s *%s) GetFilterDefinitions() []contracts.FilterDefinition {
// 	return []contracts.FilterDefinition{
// 		// String field example
// 		contracts.NewFilterDefinition(
// 			"field_name",           // Field name in database
// 			"Display Name",         // Human-readable label
// 			contracts.FilterTypeString,
// 			nil,                    // nil = use all string operators (equals, contains, starts_with, etc.)
// 		),
//
// 		// Enum field example (dropdown)
// 		contracts.NewFilterDefinition(
// 			"status",
// 			"Status",
// 			contracts.FilterTypeEnum,
// 			&[]string{"ACTIVE", "INACTIVE", "PENDING"}, // Available options
// 		),
//
// 		// Number field example
// 		contracts.NewFilterDefinition(
// 			"price",
// 			"Price",
// 			contracts.FilterTypeNumber,
// 			nil,                    // nil = use all number operators (equals, greater_than, less_than, etc.)
// 		),
//
// 		// Date field example
// 		contracts.NewFilterDefinition(
// 			"created_at",
// 			"Created Date",
// 			contracts.FilterTypeDate,
// 			nil,                    // nil = use all date operators (equals, before, after, between, etc.)
// 		),
//
// 		// Boolean field example
// 		contracts.NewFilterDefinition(
// 			"is_active",
// 			"Active Status",
// 			contracts.FilterTypeBoolean,
// 			nil,
// 		),
// 	}
// }

// Add domain-specific methods below this line
// Examples:
// - GetByStatus(status string) ([]*models.%s, error)
// - GetActive() ([]*models.%s, error)
// - Custom business logic methods specific to %s
`,
		packageName,
		serviceName, model.TableName, serviceName,
		serviceName, model.Name, serviceName, serviceName,
		model.Name, model.TableName,
		r.formatStringSlice(searchFields),
		r.formatStringSlice(sortFields),
		r.formatStringSlice(filterFields),
		r.formatValidationRules(validationRules),
		model.TableName,
		r.generateSoftDeleteLine(model.HasSoftDelete),
		serviceVarName, serviceName,
		serviceVarName, serviceName,
		serviceVarName,
		model.TableName, serviceName,
		model.Name, model.Name, model.TableName,
	)

	return content
}

// generateValidationRule generates a validation rule string for a field
func (r *ServiceMaker) generateValidationRule(field FieldInfo) string {
	rules := []string{}

	// Required rule for non-optional fields
	if !field.Optional {
		rules = append(rules, "required")
	}

	// Type-specific rules
	switch field.Type {
	case "string":
		rules = append(rules, "string|max:255")
	case "int", "uint", "int64", "uint64":
		rules = append(rules, "numeric")
	case "float64", "float32":
		rules = append(rules, "numeric")
	case "bool":
		rules = append(rules, "boolean")
	case "[]string":
		rules = append(rules, "array")
	}

	if len(rules) == 0 {
		return ""
	}

	return strings.Join(rules, "|")
}

// formatStringSlice formats a string slice as Go code
func (r *ServiceMaker) formatStringSlice(items []string) string {
	if len(items) == 0 {
		return ""
	}

	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf(`"%s"`, item)
	}

	return strings.Join(quoted, ", ")
}

// formatValidationRules formats validation rules as Go code
func (r *ServiceMaker) formatValidationRules(rules map[string]string) string {
	if len(rules) == 0 {
		return "\t\t\t// Add your validation rules here"
	}

	lines := []string{}
	for field, rule := range rules {
		lines = append(lines, fmt.Sprintf("\t\t\t\"%s\": \"%s\",", field, rule))
	}

	return strings.Join(lines, "\n")
}

// generateSoftDeleteLine generates the WithSoftDeletes line if needed
func (r *ServiceMaker) generateSoftDeleteLine(hasSoftDelete bool) string {
	if hasSoftDelete {
		return "\t\tWithSoftDeletes(). // Enable soft delete support"
	}
	return ""
}

// printRegistrationInstructions outputs the registration instructions
func (r *ServiceMaker) printRegistrationInstructions(ctx console.Context, tableName, serviceName string) {
	ctx.Info("To complete the service registration, follow these steps:")
	ctx.NewLine()

	// Step 1: Permission constants
	ctx.Line("1. Add the service to app/auth/permission_constants.go:")
	ctx.NewLine()
	ctx.Comment("   In the ServiceRegistry constants section, add:")
	ctx.Line(fmt.Sprintf("   Service%s ServiceRegistry = \"%s\"", strings.Title(tableName), tableName))
	ctx.NewLine()
	ctx.Comment("   In GetAllServiceRegistries() function, add:")
	ctx.Line(fmt.Sprintf("   Service%s,", strings.Title(tableName)))
	ctx.NewLine()
	ctx.Comment("   In GetServiceDisplayName() switch statement, add:")
	ctx.Line(fmt.Sprintf("   case Service%s:", strings.Title(tableName)))
	ctx.Line(fmt.Sprintf("       return \"%s Management\"", strings.Title(tableName)))
	ctx.NewLine()
	ctx.Comment("   In GetServiceActions() switch statement, add:")
	ctx.Line(fmt.Sprintf("   case Service%s:", strings.Title(tableName)))
	ctx.Line("       return []CorePermissionAction{")
	ctx.Line("           PermissionCreate,")
	ctx.Line("           PermissionRead,")
	ctx.Line("           PermissionUpdate,")
	ctx.Line("           PermissionDelete,")
	ctx.Line("           PermissionView,")
	ctx.Line("       }")
	ctx.NewLine()

	// Step 2: Run permissions setup
	ctx.Line("2. Run the permissions setup command to create the permissions:")
	ctx.NewLine()
	ctx.Line("   go run . artisan permissions:setup")
	ctx.NewLine()
}
