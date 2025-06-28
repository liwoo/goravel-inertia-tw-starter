package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type MakeCrudE2E struct {
}

// Signature The name and signature of the console command.
func (receiver *MakeCrudE2E) Signature() string {
	return "make:crud-e2e {name} {--force}"
}

// Description The console command description.
func (receiver *MakeCrudE2E) Description() string {
	return "Generate complete CRUD system: model, migration, service, controllers, permissions, and UI components"
}

// Extend The console command extend.
func (receiver *MakeCrudE2E) Extend() command.Extend {
	return command.Extend{
		Category: "make",
	}
}

// Handle Execute the console command.
func (receiver *MakeCrudE2E) Handle(ctx console.Context) error {
	name := ctx.Argument(0)
	if name == "" {
		ctx.Error("Resource name is required")
		ctx.Info("Usage: go run . artisan make:crud-e2e Product")
		return errors.New("missing resource name")
	}

	forceOption := ctx.Option("force")
	force := forceOption != ""

	// Convert name to various formats
	resourceConfig := receiver.parseResourceName(name)
	
	ctx.Info(fmt.Sprintf("Generating complete CRUD system for: %s", resourceConfig.DisplayName))
	ctx.Info("=====================================")

	// Generate all components
	steps := []struct {
		name        string
		description string
		fn          func(console.Context, ResourceConfig, bool) error
	}{
		{"model", "Creating model", receiver.generateModel},
		{"migration", "Creating migration", receiver.generateMigration},
		{"service", "Creating service with contracts", receiver.generateService},
		{"requests", "Creating validation requests", receiver.generateRequests},
		{"controller", "Creating API controller", receiver.generateController},
		{"page-controller", "Creating page controller", receiver.generatePageController},
		{"routes", "Adding routes", receiver.generateRoutes},
		{"permissions", "Creating permissions", receiver.generatePermissions},
		{"ui-types", "Creating TypeScript types", receiver.generateUITypes},
		{"ui-components", "Creating React components", receiver.generateUIComponents},
		{"ui-pages", "Creating React pages", receiver.generateUIPages},
	}

	generatedFiles := []string{}
	
	for _, step := range steps {
		ctx.Info(fmt.Sprintf("🔨 %s...", step.description))
		
		if err := step.fn(ctx, resourceConfig, force); err != nil {
			ctx.Error(fmt.Sprintf("Failed to generate %s: %v", step.name, err))
			return err
		}
		
		ctx.Success(fmt.Sprintf("✓ %s generated successfully", step.description))
	}

	// Display summary
	ctx.Info("")
	ctx.Success("🎉 Complete CRUD system generated successfully!")
	ctx.Info("Generated files:")
	for _, file := range generatedFiles {
		ctx.Info(fmt.Sprintf("  • %s", file))
	}
	
	ctx.Info("")
	ctx.Info("Next steps:")
	ctx.Info("1. Run migration: go run . artisan migrate")
	ctx.Info("2. Seed permissions: go run . artisan seed --seeder=rbac")
	ctx.Info("3. Update your frontend routing")
	ctx.Info("4. Test the CRUD operations")

	return nil
}

// ResourceConfig holds all the naming variations for a resource
type ResourceConfig struct {
	// Input name variations
	Name            string // Product
	LowerName       string // product
	PluralName      string // Products
	LowerPluralName string // products
	SnakeName       string // product
	SnakePluralName string // products
	KebabName       string // product
	KebabPluralName string // products
	DisplayName     string // Product
	
	// Database
	TableName string // products
	
	// File paths
	ModelPath       string // app/models/product.go
	ServicePath     string // app/services/product_service.go
	ControllerPath  string // app/http/controllers/product_controller.go
	PageControllerPath string // app/http/controllers/product_page_controller.go
	RequestPath     string // app/http/requests/product_request.go
	MigrationPath   string // database/migrations/
	
	// Frontend paths
	UITypesPath     string // resources/js/types/product.ts
	UIComponentsPath string // resources/js/components/Products/
	UIPagesPath     string // resources/js/pages/Products/
}

// parseResourceName converts the input name to all required variations
func (receiver *MakeCrudE2E) parseResourceName(name string) ResourceConfig {
	name = strings.Title(strings.ToLower(name))
	lowerName := strings.ToLower(name)
	pluralName := receiver.pluralize(name)
	lowerPluralName := strings.ToLower(pluralName)
	
	return ResourceConfig{
		Name:            name,
		LowerName:       lowerName,
		PluralName:      pluralName,
		LowerPluralName: lowerPluralName,
		SnakeName:       receiver.toSnakeCase(name),
		SnakePluralName: receiver.toSnakeCase(pluralName),
		KebabName:       receiver.toKebabCase(name),
		KebabPluralName: receiver.toKebabCase(pluralName),
		DisplayName:     name,
		TableName:       receiver.toSnakeCase(pluralName),
		
		ModelPath:       fmt.Sprintf("app/models/%s.go", receiver.toSnakeCase(name)),
		ServicePath:     fmt.Sprintf("app/services/%s_service.go", receiver.toSnakeCase(name)),
		ControllerPath:  fmt.Sprintf("app/http/controllers/%s_controller.go", receiver.toSnakeCase(name)),
		PageControllerPath: fmt.Sprintf("app/http/controllers/%s_page_controller.go", receiver.toSnakeCase(name)),
		RequestPath:     fmt.Sprintf("app/http/requests/%s_request.go", receiver.toSnakeCase(name)),
		MigrationPath:   "database/migrations/",
		
		UITypesPath:     fmt.Sprintf("resources/js/types/%s.ts", lowerName),
		UIComponentsPath: fmt.Sprintf("resources/js/components/%s/", pluralName),
		UIPagesPath:     fmt.Sprintf("resources/js/pages/%s/", pluralName),
	}
}

// Helper functions for string transformations
func (receiver *MakeCrudE2E) pluralize(word string) string {
	if strings.HasSuffix(word, "y") {
		return word[:len(word)-1] + "ies"
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "sh") || strings.HasSuffix(word, "ch") {
		return word + "es"
	}
	return word + "s"
}

func (receiver *MakeCrudE2E) toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func (receiver *MakeCrudE2E) toKebabCase(str string) string {
	return strings.ReplaceAll(receiver.toSnakeCase(str), "_", "-")
}

// Generation functions
func (receiver *MakeCrudE2E) generateModel(ctx console.Context, config ResourceConfig, force bool) error {
	template := `package models

import (
	"github.com/goravel/framework/database/orm"
)

// {{.Name}} represents a {{.LowerName}} in the system
type {{.Name}} struct {
	orm.Model
	Name        string ` + "`" + `gorm:"not null" json:"name"` + "`" + `
	Description string ` + "`" + `gorm:"type:text" json:"description"` + "`" + `
	IsActive    bool   ` + "`" + `gorm:"default:true" json:"is_active"` + "`" + `
	
	// Add your custom fields here
	// Price       float64 ` + "`" + `gorm:"type:decimal(10,2)" json:"price"` + "`" + `
	// Category    string  ` + "`" + `gorm:"index" json:"category"` + "`" + `
	// CreatedByID *uint   ` + "`" + `gorm:"index" json:"created_by_id,omitempty"` + "`" + `
	// CreatedBy   *User   ` + "`" + `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"` + "`" + `
	
	orm.SoftDeletes
}

// TableName returns the table name for {{.Name}} model
func ({{.Name}}) TableName() string {
	return "{{.TableName}}"
}

// Validate performs model validation
func ({{.LowerName}} *{{.Name}}) Validate() error {
	if {{.LowerName}}.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
`

	return receiver.writeFileFromTemplate(config.ModelPath, template, config, force)
}

func (receiver *MakeCrudE2E) generateMigration(ctx console.Context, config ResourceConfig, force bool) error {
	timestamp := time.Now().Format("20060102150405")
	migrationFile := fmt.Sprintf("%s%s_create_%s_table.go", config.MigrationPath, timestamp, config.TableName)
	
	template := `package migrations

import (
	"github.com/goravel/framework/contracts/database/migration"
	"github.com/goravel/framework/contracts/database/schema"
)

func init() {
	Migrator.Register(&Create{{.PluralName}}Table{})
}

type Create{{.PluralName}}Table struct {
}

// Signature The name and signature of the migration.
func (m *Create{{.PluralName}}Table) Signature() string {
	return "{{.TableName}}_table"
}

// Up Run the migrations.
func (m *Create{{.PluralName}}Table) Up() error {
	return schema.Create("{{.TableName}}", func(table schema.Blueprint) error {
		table.ID()
		table.String("name").NotNull()
		table.Text("description").Nullable()
		table.Boolean("is_active").Default(true)
		
		// Add your custom columns here
		// table.Decimal("price", 10, 2).Nullable()
		// table.String("category").Index().Nullable()
		// table.UnsignedBigInteger("created_by_id").Index().Nullable()
		// table.Foreign("created_by_id").References("id").On("users").OnDelete("SET NULL")
		
		table.Timestamps()
		table.SoftDeletes()
		
		// Add indexes
		table.Index("name")
		table.Index("is_active")
		
		return nil
	})
}

// Down Reverse the migrations.
func (m *Create{{.PluralName}}Table) Down() error {
	return schema.DropIfExists("{{.TableName}}")
}
`

	return receiver.writeFileFromTemplate(migrationFile, template, config, force)
}

func (receiver *MakeCrudE2E) generateService(ctx console.Context, config ResourceConfig, force bool) error {
	template := `package services

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/facades"
	"players/app/contracts"
	"players/app/models"
)

// {{.Name}}Service handles {{.LowerName}} business logic with contract enforcement
type {{.Name}}Service struct {
	*contracts.BaseCrudService
}

// New{{.Name}}Service creates a new {{.LowerName}} service that implements all contracts
func New{{.Name}}Service() *{{.Name}}Service {
	service := &{{.Name}}Service{
		BaseCrudService: contracts.NewBaseCrudService("{{.LowerName}}", "id"),
	}

	// Register service with validation
	contracts.MustRegisterCrudService("{{.LowerPluralName}}", service)

	return service
}

// GetList with built-in pagination, sorting, filtering using GORM directly
// Implements CrudServiceContract interface
func (s *{{.Name}}Service) GetList(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Use base service validation
	if err := s.ValidateListRequest(&req); err != nil {
		return nil, err
	}
	s.SanitizeListRequest(&req)

	// Build query
	query := facades.Orm().Query().Model(&models.{{.Name}}{})

	// Apply search if provided using searchable fields
	if req.Search != "" {
		if err := s.ValidateSearchQuery(req.Search); err != nil {
			return nil, err
		}
		searchFields := s.GetSearchableFields()
		if len(searchFields) > 0 {
			conditions := make([]string, len(searchFields))
			values := make([]interface{}, len(searchFields))
			for i, field := range searchFields {
				conditions[i] = field + " LIKE ?"
				values[i] = "%" + req.Search + "%"
			}
			query = query.Where(strings.Join(conditions, " OR "), values...)
		}
	}

	// Apply sorting with field validation and mapping
	if req.Sort != "" && req.Direction != "" {
		if s.ValidateSortField(req.Sort) && s.ValidateSortDirection(req.Direction) {
			dbColumn, valid := s.MapSortField(req.Sort)
			if valid {
				orderClause := dbColumn + " " + strings.ToUpper(req.Direction)
				query = query.Order(orderClause)
			} else {
				// Use default sort
				defaultField, defaultDir := s.GetDefaultSort()
				query = query.Order(defaultField + " " + defaultDir)
			}
		} else {
			// Use default sort
			defaultField, defaultDir := s.GetDefaultSort()
			query = query.Order(defaultField + " " + defaultDir)
		}
	} else {
		// Default sorting
		defaultField, defaultDir := s.GetDefaultSort()
		query = query.Order(defaultField + " " + defaultDir)
	}

	// Get all {{.LowerPluralName}} with applied filters and sorting
	var all{{.PluralName}} []models.{{.Name}}
	if err := query.Find(&all{{.PluralName}}); err != nil {
		return nil, err
	}

	// Manual pagination
	total := int64(len(all{{.PluralName}}))
	offset := (req.Page - 1) * req.PageSize
	end := offset + req.PageSize

	if offset > len(all{{.PluralName}}) {
		offset = len(all{{.PluralName}})
	}
	if end > len(all{{.PluralName}}) {
		end = len(all{{.PluralName}})
	}

	var page{{.PluralName}} []models.{{.Name}}
	if offset < len(all{{.PluralName}}) {
		page{{.PluralName}} = all{{.PluralName}}[offset:end]
	}

	lastPage := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	// Convert to interface slice
	data := make([]interface{}, len(page{{.PluralName}}))
	for i, {{.LowerName}} := range page{{.PluralName}} {
		data[i] = {{.LowerName}}
	}

	return &contracts.PaginatedResult{
		Data:        data,
		Total:       total,
		PerPage:     req.PageSize,
		CurrentPage: req.Page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + len(page{{.PluralName}}),
		HasNext:     req.Page < lastPage,
		HasPrev:     req.Page > 1,
	}, nil
}

// GetListAdvanced with additional filters using GORM directly
// Implements CrudServiceContract interface
func (s *{{.Name}}Service) GetListAdvanced(req contracts.ListRequest, filters map[string]interface{}) (*contracts.PaginatedResult, error) {
	// Validate and sanitize request
	if err := s.ValidateListRequest(&req); err != nil {
		return nil, err
	}
	s.SanitizeListRequest(&req)

	// Validate filters
	validatedFilters, err := s.BuildFilterQuery(filters)
	if err != nil {
		return nil, err
	}

	// Create separate queries for count and data
	countQuery := facades.Orm().Query().Model(&models.{{.Name}}{})
	dataQuery := facades.Orm().Query().Model(&models.{{.Name}}{})

	// Apply search to both queries if provided
	if req.Search != "" {
		searchCondition := "name LIKE ?"
		searchValue := "%" + req.Search + "%"
		countQuery = countQuery.Where(searchCondition, searchValue)
		dataQuery = dataQuery.Where(searchCondition, searchValue)
	}

	// Apply validated filters to both queries
	for field, value := range validatedFilters {
		var condition string
		switch field {
		case "is_active":
			condition = "is_active = ?"
		case "name":
			condition = "name LIKE ?"
			value = "%" + fmt.Sprintf("%v", value) + "%"
		default:
			continue
		}
		countQuery = countQuery.Where(condition, value)
		dataQuery = dataQuery.Where(condition, value)
	}

	// Count total records
	var total int64
	if err := countQuery.Count(&total); err != nil {
		return nil, err
	}

	// Add sorting to data query only
	if req.Sort != "" {
		dataQuery = dataQuery.Order("id DESC") // Should parse req.Sort properly
	} else {
		dataQuery = dataQuery.Order("id DESC")
	}

	// Calculate pagination
	offset := (req.Page - 1) * req.PageSize
	lastPage := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	// Get paginated data
	var {{.LowerPluralName}} []models.{{.Name}}
	if err := dataQuery.Offset(offset).Limit(req.PageSize).Find(&{{.LowerPluralName}}); err != nil {
		return nil, err
	}

	// Convert to interface slice
	data := make([]interface{}, len({{.LowerPluralName}}))
	for i, {{.LowerName}} := range {{.LowerPluralName}} {
		data[i] = {{.LowerName}}
	}

	return &contracts.PaginatedResult{
		Data:        data,
		Total:       total,
		PerPage:     req.PageSize,
		CurrentPage: req.Page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + len({{.LowerPluralName}}),
		HasNext:     req.Page < lastPage,
		HasPrev:     req.Page > 1,
	}, nil
}

// GetByID - Implements CrudServiceContract interface
func (s *{{.Name}}Service) GetByID(id uint) (interface{}, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	return s.get{{.Name}}ByID(id)
}

// get{{.Name}}ByID is a helper method that returns the actual model type
func (s *{{.Name}}Service) get{{.Name}}ByID(id uint) (*models.{{.Name}}, error) {
	var {{.LowerName}} models.{{.Name}}
	if err := facades.Orm().Query().Model(&models.{{.Name}}{}).Where("id = ?", id).First(&{{.LowerName}}); err != nil {
		return nil, fmt.Errorf("{{.LowerName}} not found: %w", err)
	}

	return &{{.LowerName}}, nil
}

// Create - Implements CrudServiceContract interface
func (s *{{.Name}}Service) Create(data map[string]interface{}) (interface{}, error) {
	// Validate using validation rules
	if err := s.validateWithRules(data, false); err != nil {
		return nil, err
	}

	return s.create{{.Name}}(data)
}

// create{{.Name}} is a helper method that returns the actual model type
func (s *{{.Name}}Service) create{{.Name}}(data map[string]interface{}) (*models.{{.Name}}, error) {
	// Basic validation
	if err := s.validate{{.Name}}Data(data, false); err != nil {
		return nil, err
	}

	// Set default values if not provided
	if _, exists := data["is_active"]; !exists {
		data["is_active"] = true
	}

	// Create {{.LowerName}} struct from data
	{{.LowerName}} := models.{{.Name}}{
		Name:        data["name"].(string),
		IsActive:    data["is_active"].(bool),
	}

	if desc, ok := data["description"].(string); ok {
		{{.LowerName}}.Description = desc
	}

	// Create using GORM
	if err := facades.Orm().Query().Create(&{{.LowerName}}); err != nil {
		return nil, fmt.Errorf("failed to create {{.LowerName}}: %w", err)
	}

	return &{{.LowerName}}, nil
}

// Update - Implements CrudServiceContract interface
func (s *{{.Name}}Service) Update(id uint, data map[string]interface{}) (interface{}, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	// Validate using validation rules
	if err := s.validateWithRules(data, true); err != nil {
		return nil, err
	}

	return s.update{{.Name}}(id, data)
}

// update{{.Name}} is a helper method that returns the actual model type
func (s *{{.Name}}Service) update{{.Name}}(id uint, data map[string]interface{}) (*models.{{.Name}}, error) {
	// Check if {{.LowerName}} exists
	_, err := s.get{{.Name}}ByID(id)
	if err != nil {
		return nil, err
	}

	// Update using GORM
	var {{.LowerName}} models.{{.Name}}
	if _, err := facades.Orm().Query().Model(&{{.LowerName}}).Where("id = ?", id).Update(data); err != nil {
		return nil, fmt.Errorf("failed to update {{.LowerName}}: %w", err)
	}

	// Return updated {{.LowerName}}
	return s.get{{.Name}}ByID(id)
}

// Delete - Implements CrudServiceContract interface
func (s *{{.Name}}Service) Delete(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	// Check if {{.LowerName}} exists
	_, err := s.get{{.Name}}ByID(id)
	if err != nil {
		return err
	}

	// Delete using GORM (soft delete)
	if _, err := facades.Orm().Query().Model(&models.{{.Name}}{}).Where("id = ?", id).Delete(&models.{{.Name}}{}); err != nil {
		return fmt.Errorf("failed to delete {{.LowerName}}: %w", err)
	}

	return nil
}

// CONTRACT IMPLEMENTATIONS - Required by CompleteCrudService interface

// PaginationServiceContract implementation
func (s *{{.Name}}Service) GetPaginatedList(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	return s.GetList(req)
}

// SortableServiceContract implementation
func (s *{{.Name}}Service) GetSortableFields() []string {
	return []string{"id", "name", "description", "is_active", "createdAt", "updatedAt"}
}

func (s *{{.Name}}Service) ValidateSortField(field string) bool {
	sortableFields := s.GetSortableFields()
	for _, validField := range sortableFields {
		if field == validField {
			return true
		}
	}
	return false
}

func (s *{{.Name}}Service) MapSortField(frontendField string) (string, bool) {
	columnMapping := s.GetColumnMapping()
	if dbColumn, exists := columnMapping[frontendField]; exists {
		return dbColumn, true
	}
	return "", false
}

// FilterableServiceContract implementation
func (s *{{.Name}}Service) GetFilterableFields() []string {
	return []string{"name", "is_active"}
}

func (s *{{.Name}}Service) ValidateFilterField(field string) bool {
	filterableFields := s.GetFilterableFields()
	for _, validField := range filterableFields {
		if field == validField {
			return true
		}
	}
	return false
}

func (s *{{.Name}}Service) GetSearchableFields() []string {
	return []string{"name", "description"}
}

func (s *{{.Name}}Service) BuildFilterQuery(filters map[string]interface{}) (map[string]interface{}, error) {
	validatedFilters := make(map[string]interface{})

	for field, value := range filters {
		if !s.ValidateFilterField(field) {
			continue // Skip invalid fields
		}

		if !s.ValidateFilterValue(field, value) {
			continue // Skip invalid values
		}

		validatedFilters[field] = value
	}

	return validatedFilters, nil
}

// SearchableServiceContract implementation
func (s *{{.Name}}Service) Search(query string, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	if err := s.ValidateSearchQuery(query); err != nil {
		return nil, err
	}

	req.Search = query
	return s.GetList(req)
}

func (s *{{.Name}}Service) ValidateSearchQuery(query string) error {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return fmt.Errorf("search query must be at least 2 characters")
	}
	if len(query) > 100 {
		return fmt.Errorf("search query cannot exceed 100 characters")
	}
	return nil
}

// BulkOperationsContract implementation
func (s *{{.Name}}Service) BulkCreate(data []map[string]interface{}) ([]interface{}, error) {
	if err := s.ValidateBulkOperation([]uint{uint(len(data))}); err != nil {
		return nil, err
	}

	results := make([]interface{}, 0, len(data))
	for _, item := range data {
		result, err := s.Create(item)
		if err != nil {
			return nil, fmt.Errorf("bulk create failed: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *{{.Name}}Service) BulkUpdate(ids []uint, data map[string]interface{}) error {
	if err := s.ValidateBulkOperation(ids); err != nil {
		return err
	}

	for _, id := range ids {
		_, err := s.Update(id, data)
		if err != nil {
			return fmt.Errorf("bulk update failed for ID %d: %w", id, err)
		}
	}

	return nil
}

func (s *{{.Name}}Service) BulkDelete(ids []uint) error {
	if err := s.ValidateBulkOperation(ids); err != nil {
		return err
	}

	for _, id := range ids {
		if err := s.Delete(id); err != nil {
			return fmt.Errorf("bulk delete failed for ID %d: %w", id, err)
		}
	}

	return nil
}

// CrudServiceConfiguration implementation
func (s *{{.Name}}Service) GetModel() interface{} {
	return &models.{{.Name}}{}
}

func (s *{{.Name}}Service) GetValidationRules() map[string]interface{} {
	return map[string]interface{}{
		"name":        "required|string|max:255",
		"description": "string|max:1000",
		"is_active":   "boolean",
	}
}

func (s *{{.Name}}Service) GetColumnMapping() map[string]string {
	return map[string]string{
		"id":          "id",
		"name":        "name",
		"description": "description",
		"isActive":    "is_active",
		"createdAt":   "created_at",
		"updatedAt":   "updated_at",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
		"is_active":   "is_active",
	}
}

// HELPER METHODS

// validateWithRules uses the validation rules from the contract
func (s *{{.Name}}Service) validateWithRules(data map[string]interface{}, isUpdate bool) error {
	rules := s.GetValidationRules()

	// For updates, make required fields optional
	if isUpdate {
		for field, rule := range rules {
			if ruleStr, ok := rule.(string); ok {
				// Remove 'required|' from validation rules for updates
				if strings.HasPrefix(ruleStr, "required|") {
					rules[field] = strings.TrimPrefix(ruleStr, "required|")
				}
			}
		}
	}

	// Basic validation implementation (can be enhanced with proper validator)
	return s.validate{{.Name}}Data(data, isUpdate)
}

// validate{{.Name}}Data performs simple validation
func (s *{{.Name}}Service) validate{{.Name}}Data(data map[string]interface{}, isUpdate bool) error {
	// Required fields for creation
	if !isUpdate {
		requiredFields := []string{"name"}
		for _, field := range requiredFields {
			if value, exists := data[field]; !exists || value == "" {
				return fmt.Errorf("%s is required", field)
			}
		}
	}

	// Validate name if provided
	if name, ok := data["name"].(string); ok && name != "" {
		if len(name) < 2 || len(name) > 255 {
			return fmt.Errorf("name must be between 2 and 255 characters")
		}
	}

	// Validate description if provided
	if desc, exists := data["description"]; exists {
		if descStr, ok := desc.(string); ok {
			if len(descStr) > 1000 {
				return fmt.Errorf("description cannot exceed 1000 characters")
			}
		}
	}

	return nil
}
`

	return receiver.writeFileFromTemplate(config.ServicePath, template, config, force)
}

func (receiver *MakeCrudE2E) generateRequests(ctx console.Context, config ResourceConfig, force bool) error {
	template := `package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

// {{.Name}}CreateRequest handles validation for creating {{.LowerPluralName}}
type {{.Name}}CreateRequest struct {
	Name        string ` + "`" + `form:"name" json:"name"` + "`" + `
	Description string ` + "`" + `form:"description" json:"description"` + "`" + `
	IsActive    bool   ` + "`" + `form:"is_active" json:"is_active"` + "`" + `
}

// Authorize determines if the user can make this request
func (r *{{.Name}}CreateRequest) Authorize(ctx http.Context) error {
	// Add authorization logic here
	return nil
}

// Rules returns the validation rules for the request
func (r *{{.Name}}CreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "required|string|max:255|min:2",
		"description": "string|max:1000",
		"is_active":   "boolean",
	}
}

// Messages returns custom validation messages
func (r *{{.Name}}CreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required": "{{.Name}} name is required",
		"name.min":      "{{.Name}} name must be at least 2 characters",
		"name.max":      "{{.Name}} name cannot exceed 255 characters",
		"description.max": "Description cannot exceed 1000 characters",
	}
}

// Attributes returns custom attribute names for validation
func (r *{{.Name}}CreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "{{.Name}} Name",
		"description": "Description",
		"is_active":   "Active Status",
	}
}

// PrepareForValidation allows you to modify the data before validation
func (r *{{.Name}}CreateRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	// Set default values or modify data before validation
	if data.Get("is_active") == nil {
		data.Set("is_active", true)
	}
	return nil
}

// ToCreateData converts the request to data suitable for the service
func (r *{{.Name}}CreateRequest) ToCreateData() map[string]interface{} {
	return map[string]interface{}{
		"name":        r.Name,
		"description": r.Description,
		"is_active":   r.IsActive,
	}
}

// {{.Name}}UpdateRequest handles validation for updating {{.LowerPluralName}}
type {{.Name}}UpdateRequest struct {
	ID          uint   ` + "`" + `form:"id" json:"id"` + "`" + `
	Name        string ` + "`" + `form:"name" json:"name"` + "`" + `
	Description string ` + "`" + `form:"description" json:"description"` + "`" + `
	IsActive    bool   ` + "`" + `form:"is_active" json:"is_active"` + "`" + `
}

// Authorize determines if the user can make this request
func (r *{{.Name}}UpdateRequest) Authorize(ctx http.Context) error {
	// Add authorization logic here
	return nil
}

// Rules returns the validation rules for the request
func (r *{{.Name}}UpdateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "string|max:255|min:2",
		"description": "string|max:1000",
		"is_active":   "boolean",
	}
}

// Messages returns custom validation messages
func (r *{{.Name}}UpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.min":      "{{.Name}} name must be at least 2 characters",
		"name.max":      "{{.Name}} name cannot exceed 255 characters",
		"description.max": "Description cannot exceed 1000 characters",
	}
}

// Attributes returns custom attribute names for validation
func (r *{{.Name}}UpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "{{.Name}} Name",
		"description": "Description",
		"is_active":   "Active Status",
	}
}

// PrepareForValidation allows you to modify the data before validation
func (r *{{.Name}}UpdateRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	// Modify data before validation if needed
	return nil
}

// ToUpdateData converts the request to data suitable for the service
func (r *{{.Name}}UpdateRequest) ToUpdateData() map[string]interface{} {
	data := make(map[string]interface{})
	
	if r.Name != "" {
		data["name"] = r.Name
	}
	if r.Description != "" {
		data["description"] = r.Description
	}
	// Always include is_active for updates
	data["is_active"] = r.IsActive
	
	return data
}
`

	return receiver.writeFileFromTemplate(config.RequestPath, template, config, force)
}

func (receiver *MakeCrudE2E) generateController(ctx console.Context, config ResourceConfig, force bool) error {
	template := `package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/http/requests"
	"players/app/services"
)

// {{.Name}}Controller - API controller with contract enforcement
// Implements ResourceControllerContract interface
type {{.Name}}Controller struct {
	*contracts.BaseCrudController
	{{.LowerName}}Service *services.{{.Name}}Service
	authHelper  contracts.AuthHelper
}

// New{{.Name}}Controller creates a new {{.LowerName}} controller that implements all contracts
func New{{.Name}}Controller() *{{.Name}}Controller {
	controller := &{{.Name}}Controller{
		BaseCrudController: contracts.NewBaseCrudController("{{.LowerName}}"),
		{{.LowerName}}Service:        services.New{{.Name}}Service(),
		authHelper:         helpers.NewAuthHelper(),
	}

	// Register controller with validation
	contracts.MustRegisterCrudController("{{.LowerPluralName}}", controller)

	return controller
}

// Index GET /{{.LowerPluralName}} - Implements CrudControllerContract
func (c *{{.Name}}Controller) Index(ctx http.Context) http.Response {
	// Validate pagination request using contract
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Check permission
	if err := c.CheckPermission(ctx, "{{.LowerPluralName}}.viewAny", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Get {{.LowerPluralName}} using service
	result, err := c.{{.LowerName}}Service.GetList(*req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve {{.LowerPluralName}}: "+err.Error())
	}

	// Build standardized paginated response
	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, "{{.PluralName}} retrieved successfully")
}

// Show GET /{{.LowerPluralName}}/{id} - Implements CrudControllerContract (JSON for modals)
func (c *{{.Name}}Controller) Show(ctx http.Context) http.Response {
	// Validate ID parameter using contract
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid {{.LowerName}} ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Check authorization
	if err := c.CheckPermission(ctx, "{{.LowerPluralName}}.view", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Get the {{.LowerName}}
	{{.LowerName}}, err := c.{{.LowerName}}Service.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, "{{.LowerName}}", id)
	}

	return c.SuccessResponse(ctx, {{.LowerName}}, "{{.Name}} details retrieved successfully")
}

// Store POST /{{.LowerPluralName}} - Implements CrudControllerContract
func (c *{{.Name}}Controller) Store(ctx http.Context) http.Response {
	// Check authorization
	if err := c.CheckPermission(ctx, "{{.LowerPluralName}}.create", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Validate create request using contract
	data, err := c.ValidateCreateRequest(ctx)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Create the {{.LowerName}} using validated data
	{{.LowerName}}, err := c.{{.LowerName}}Service.Create(data)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to create {{.LowerName}}: "+err.Error())
	}

	return c.ResourceCreatedResponse(ctx, {{.LowerName}}, "{{.LowerName}}")
}

// Update PUT /{{.LowerPluralName}}/{id} - Implements CrudControllerContract
func (c *{{.Name}}Controller) Update(ctx http.Context) http.Response {
	// Validate ID parameter using contract
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid {{.LowerName}} ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Check if {{.LowerName}} exists
	_, err = c.{{.LowerName}}Service.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, "{{.LowerName}}", id)
	}

	// Check authorization
	if err := c.CheckPermission(ctx, "{{.LowerPluralName}}.update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Validate update request using contract
	data, err := c.ValidateUpdateRequest(ctx, id)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Update the {{.LowerName}} using validated data
	updated{{.Name}}, err := c.{{.LowerName}}Service.Update(id, data)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to update {{.LowerName}}: "+err.Error())
	}

	return c.ResourceUpdatedResponse(ctx, updated{{.Name}}, "{{.LowerName}}")
}

// Delete DELETE /{{.LowerPluralName}}/{id} - Implements CrudControllerContract
func (c *{{.Name}}Controller) Delete(ctx http.Context) http.Response {
	// Validate ID parameter using contract
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid {{.LowerName}} ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Check if {{.LowerName}} exists
	_, err = c.{{.LowerName}}Service.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, "{{.LowerName}}", id)
	}

	// Check authorization
	if err := c.CheckPermission(ctx, "{{.LowerPluralName}}.delete", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Delete the {{.LowerName}}
	err = c.{{.LowerName}}Service.Delete(id)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to delete {{.LowerName}}: "+err.Error())
	}

	return c.ResourceDeletedResponse(ctx, "{{.LowerName}}", id)
}

// CONTRACT IMPLEMENTATIONS - Required by ResourceControllerContract interface

// ValidationControllerContract implementation
func (c *{{.Name}}Controller) ValidateCreateRequest(ctx http.Context) (map[string]interface{}, error) {
	var createRequest requests.{{.Name}}CreateRequest
	errors, err := ctx.Request().ValidateRequest(&createRequest)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	if errors != nil {
		return nil, fmt.Errorf("validation errors: %v", errors.All())
	}

	return createRequest.ToCreateData(), nil
}

func (c *{{.Name}}Controller) ValidateUpdateRequest(ctx http.Context, id uint) (map[string]interface{}, error) {
	var updateRequest requests.{{.Name}}UpdateRequest
	updateRequest.ID = id // Set the ID for validation context

	errors, err := ctx.Request().ValidateRequest(&updateRequest)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	if errors != nil {
		return nil, fmt.Errorf("validation errors: %v", errors.All())
	}

	return updateRequest.ToUpdateData(), nil
}

func (c *{{.Name}}Controller) GetValidationRules() map[string]interface{} {
	return c.{{.LowerName}}Service.GetValidationRules()
}

// AuthorizationControllerContract implementation
func (c *{{.Name}}Controller) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequirePermission(ctx, permission)
	return err
}

func (c *{{.Name}}Controller) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	return permHelper.GetAuthenticatedUser(ctx)
}

func (c *{{.Name}}Controller) RequireAuthentication(ctx http.Context) error {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func (c *{{.Name}}Controller) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	permHelper := auth.GetPermissionHelper()
	return permHelper.BuildPermissionsMap(ctx, resourceType)
}
`

	return receiver.writeFileFromTemplate(config.ControllerPath, template, config, force)
}

func (receiver *MakeCrudE2E) generatePageController(ctx console.Context, config ResourceConfig, force bool) error {
	template := `package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/http/inertia"
	"players/app/services"
)

// {{.Name}}PageController handles the Inertia.js {{.Name}} management page
// Implements PageControllerContract interface
type {{.Name}}PageController struct {
	*contracts.BasePageController
	{{.LowerName}}Service *services.{{.Name}}Service
	authHelper  contracts.AuthHelper
}

// New{{.Name}}PageController creates a new {{.LowerName}} page controller that implements all contracts
func New{{.Name}}PageController() *{{.Name}}PageController {
	controller := &{{.Name}}PageController{
		BasePageController: contracts.NewBasePageController("{{.LowerName}}", "{{.PluralName}}/Index"),
		{{.LowerName}}Service:        services.New{{.Name}}Service(),
		authHelper:         helpers.NewAuthHelper(),
	}

	// Register page controller with validation
	contracts.MustRegisterPageController("{{.LowerPluralName}}_page", controller)

	return controller
}

// Index renders the {{.PluralName}} management page with data and permissions - Implements PageControllerContract
func (c *{{.Name}}PageController) Index(ctx http.Context) http.Response {
	// Validate page request using contract
	req, err := c.ValidatePageRequest(ctx)
	if err != nil {
		// Return error page or redirect with error
		req = &contracts.ListRequest{Page: 1, PageSize: 20}
		req.SetDefaults()
	}

	// Check permission
	if err := c.CheckPermission(ctx, "{{.LowerPluralName}}.viewAny", nil); err != nil {
		return inertia.Render(ctx, "Errors/403", map[string]interface{}{
			"message": "Access denied: You don't have permission to view {{.LowerPluralName}}",
		})
	}

	// Build permissions map using contract
	permissions := c.BuildPermissionsMap(ctx, "{{.LowerName}}")

	// Get {{.LowerPluralName}} data
	{{.LowerPluralName}}Result, err := c.{{.LowerName}}Service.GetList(*req)
	if err != nil {
		// Handle error gracefully, provide empty result
		{{.LowerPluralName}}Result = &contracts.PaginatedResult{
			Data:        []interface{}{},
			Total:       0,
			CurrentPage: 1,
			LastPage:    1,
			PerPage:     req.PageSize,
		}
	}

	// Get {{.LowerName}} statistics if user can view reports
	var stats map[string]interface{}
	if permissions["canViewReports"] {
		stats = c.get{{.Name}}Statistics()
	}

	// Build standardized page props using contract
	data := map[string]interface{}{
		"data":        {{.LowerPluralName}}Result.Data,
		"total":       {{.LowerPluralName}}Result.Total,
		"currentPage": {{.LowerPluralName}}Result.CurrentPage,
		"lastPage":    {{.LowerPluralName}}Result.LastPage,
		"perPage":     {{.LowerPluralName}}Result.PerPage,
		"from":        {{.LowerPluralName}}Result.From,
		"to":          {{.LowerPluralName}}Result.To,
		"hasNext":     {{.LowerPluralName}}Result.HasNext,
		"hasPrev":     {{.LowerPluralName}}Result.HasPrev,
	}

	filters := map[string]interface{}{
		"page":      req.Page,
		"pageSize":  req.PageSize,
		"search":    req.Search,
		"sort":      req.Sort,
		"direction": req.Direction,
		"filters":   req.Filters,
	}

	meta := map[string]interface{}{
		"stats": stats,
	}

	props := c.BuildPageProps(data, filters, permissions, meta)

	return inertia.Render(ctx, "{{.PluralName}}/Index", props)
}

// get{{.Name}}Statistics returns {{.LowerName}} statistics for the dashboard
func (c *{{.Name}}PageController) get{{.Name}}Statistics() map[string]interface{} {
	// Get status counts
	activeCount := c.get{{.Name}}CountByStatus(true)
	inactiveCount := c.get{{.Name}}CountByStatus(false)
	totalCount := activeCount + inactiveCount

	return map[string]interface{}{
		"total{{.PluralName}}":     totalCount,
		"active{{.PluralName}}":    activeCount,
		"inactive{{.PluralName}}":  inactiveCount,
	}
}

// get{{.Name}}CountByStatus is a helper function to get {{.LowerName}} count by status
func (c *{{.Name}}PageController) get{{.Name}}CountByStatus(isActive bool) int {
	// This is a simplified version. In a real implementation,
	// you'd want to add a method to the service to get counts efficiently
	req := contracts.ListRequest{
		PageSize: 1,
		Filters: map[string]interface{}{
			"is_active": isActive,
		},
	}

	result, err := c.{{.LowerName}}Service.GetListAdvanced(req, map[string]interface{}{
		"is_active": isActive,
	})
	if err != nil {
		return 0
	}

	return int(result.Total)
}

// CONTRACT IMPLEMENTATIONS - Required by PageControllerContract interface

// AuthorizationControllerContract implementation
func (c *{{.Name}}PageController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequirePermission(ctx, permission)
	return err
}

func (c *{{.Name}}PageController) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	return permHelper.GetAuthenticatedUser(ctx)
}

func (c *{{.Name}}PageController) RequireAuthentication(ctx http.Context) error {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func (c *{{.Name}}PageController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	permHelper := auth.GetPermissionHelper()
	return permHelper.BuildPermissionsMap(ctx, resourceType)
}
`

	return receiver.writeFileFromTemplate(config.PageControllerPath, template, config, force)
}

func (receiver *MakeCrudE2E) generateRoutes(ctx console.Context, config ResourceConfig, force bool) error {
	routeFile := fmt.Sprintf("routes/%s.go", config.LowerPluralName)
	
	template := `package routes

import (
	"github.com/goravel/framework/contracts/route"
	"players/app/http/controllers"
)

// {{.Name}}Routes registers all {{.LowerName}} related routes
func {{.Name}}Routes(router route.Route) {
	{{.LowerName}}Controller := controllers.New{{.Name}}Controller()
	{{.LowerName}}PageController := controllers.New{{.Name}}PageController()

	// API Routes
	apiGroup := router.Prefix("/api").Middleware("cors")
	{{.LowerName}}ApiGroup := apiGroup.Prefix("/{{.LowerPluralName}}")
	{
		{{.LowerName}}ApiGroup.Get("/", {{.LowerName}}Controller.Index)
		{{.LowerName}}ApiGroup.Get("/{id}", {{.LowerName}}Controller.Show)
		{{.LowerName}}ApiGroup.Post("/", {{.LowerName}}Controller.Store)
		{{.LowerName}}ApiGroup.Put("/{id}", {{.LowerName}}Controller.Update)
		{{.LowerName}}ApiGroup.Delete("/{id}", {{.LowerName}}Controller.Delete)
	}

	// Admin Web Routes (Inertia.js)
	adminGroup := router.Prefix("/admin").Middleware("web", "auth")
	{
		adminGroup.Get("/{{.LowerPluralName}}", {{.LowerName}}PageController.Index)
	}
}
`

	return receiver.writeFileFromTemplate(routeFile, template, config, force)
}

func (receiver *MakeCrudE2E) generatePermissions(ctx console.Context, config ResourceConfig, force bool) error {
	permissionFile := fmt.Sprintf("database/seeders/%s_permissions_seeder.go", config.LowerName)
	
	template := `package seeders

import (
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/models"
)

// {{.Name}}PermissionsSeeder seeds {{.LowerName}} permissions
type {{.Name}}PermissionsSeeder struct{}

// Signature implements the Seeder interface
func (s *{{.Name}}PermissionsSeeder) Signature() string {
	return "{{.LowerName}}_permissions"
}

// Run seeds {{.LowerName}} permissions
func (s *{{.Name}}PermissionsSeeder) Run() error {
	// Create {{.LowerName}} permissions
	permissions := []models.Permission{
		{Name: "View Any {{.PluralName}}", Slug: "{{.LowerPluralName}}.viewAny", Category: "{{.LowerPluralName}}", Action: "viewAny", Description: "View any {{.LowerPluralName}} in the system"},
		{Name: "View {{.PluralName}}", Slug: "{{.LowerPluralName}}.view", Category: "{{.LowerPluralName}}", Action: "view", Description: "View specific {{.LowerPluralName}}"},
		{Name: "Create {{.PluralName}}", Slug: "{{.LowerPluralName}}.create", Category: "{{.LowerPluralName}}", Action: "create", Description: "Create new {{.LowerPluralName}}"},
		{Name: "Update {{.PluralName}}", Slug: "{{.LowerPluralName}}.update", Category: "{{.LowerPluralName}}", Action: "update", Description: "Update existing {{.LowerPluralName}}"},
		{Name: "Delete {{.PluralName}}", Slug: "{{.LowerPluralName}}.delete", Category: "{{.LowerPluralName}}", Action: "delete", Description: "Delete {{.LowerPluralName}}"},
		{Name: "Manage {{.PluralName}}", Slug: "{{.LowerPluralName}}.manage", Category: "{{.LowerPluralName}}", Action: "manage", Description: "Full {{.LowerName}} management"},
		{Name: "Export {{.PluralName}}", Slug: "{{.LowerPluralName}}.export", Category: "{{.LowerPluralName}}", Action: "export", Description: "Export {{.LowerPluralName}} data"},
	}

	for _, permission := range permissions {
		var existing models.Permission
		err := facades.Orm().Query().Where("slug = ?", permission.Slug).First(&existing)
		if err != nil {
			// Permission doesn't exist, create it
			err = facades.Orm().Query().Create(&permission)
			if err != nil {
				return err
			}
		}
	}

	// Assign permissions to roles
	permissionService := auth.GetPermissionService()

	// Super Admin gets all permissions (already has *.* wildcard)
	
	// Admin permissions
	adminPerms := []string{
		"{{.LowerPluralName}}.viewAny", "{{.LowerPluralName}}.view", "{{.LowerPluralName}}.create", 
		"{{.LowerPluralName}}.update", "{{.LowerPluralName}}.delete", "{{.LowerPluralName}}.manage", "{{.LowerPluralName}}.export",
	}
	s.assignPermissionsToRole("admin", adminPerms, permissionService)

	// Librarian permissions
	librarianPerms := []string{
		"{{.LowerPluralName}}.viewAny", "{{.LowerPluralName}}.view", "{{.LowerPluralName}}.create", 
		"{{.LowerPluralName}}.update", "{{.LowerPluralName}}.export",
	}
	s.assignPermissionsToRole("librarian", librarianPerms, permissionService)

	// Moderator permissions
	moderatorPerms := []string{
		"{{.LowerPluralName}}.viewAny", "{{.LowerPluralName}}.view", "{{.LowerPluralName}}.create", "{{.LowerPluralName}}.update",
	}
	s.assignPermissionsToRole("moderator", moderatorPerms, permissionService)

	// Member permissions
	memberPerms := []string{
		"{{.LowerPluralName}}.viewAny", "{{.LowerPluralName}}.view",
	}
	s.assignPermissionsToRole("member", memberPerms, permissionService)

	// Guest permissions
	guestPerms := []string{
		"{{.LowerPluralName}}.viewAny", "{{.LowerPluralName}}.view",
	}
	s.assignPermissionsToRole("guest", guestPerms, permissionService)

	return nil
}

// assignPermissionsToRole assigns specific permissions to a role
func (s *{{.Name}}PermissionsSeeder) assignPermissionsToRole(roleSlug string, permissionSlugs []string, permissionService *auth.PermissionService) {
	for _, permSlug := range permissionSlugs {
		err := permissionService.GrantPermissionToRole(roleSlug, permSlug, nil)
		if err != nil {
			// Permission might already be granted, continue
			continue
		}
	}
}
`

	return receiver.writeFileFromTemplate(permissionFile, template, config, force)
}

func (receiver *MakeCrudE2E) generateUITypes(ctx console.Context, config ResourceConfig, force bool) error {
	template := `// TypeScript type definitions for {{.Name}}
export interface {{.Name}} {
  id: number;
  name: string;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface {{.Name}}ListResponse {
  data: {{.Name}}[];
  total: number;
  currentPage: number;
  lastPage: number;
  perPage: number;
  from: number;
  to: number;
  hasNext: boolean;
  hasPrev: boolean;
}

export interface {{.Name}}ListRequest {
  page?: number;
  pageSize?: number;
  search?: string;
  sort?: string;
  direction?: 'ASC' | 'DESC';
  filters?: Record<string, any>;
}

export interface {{.Name}}Stats {
  total{{.PluralName}}: number;
  active{{.PluralName}}: number;
  inactive{{.PluralName}}: number;
}

export interface {{.Name}}FormData {
  name: string;
  description: string;
  is_active: boolean;
}

export interface {{.Name}}BulkOperation {
  action: 'delete' | 'activate' | 'deactivate';
  ids: number[];
}

export interface {{.Name}}ExportOptions {
  format: 'csv' | 'json' | 'excel';
  fields?: string[];
  includeStats?: boolean;
  filters?: {{.Name}}ListRequest;
}

export interface {{.Name}}ImportData {
  file: File;
  format: 'csv' | 'json' | 'excel';
  skipErrors: boolean;
  updateExisting: boolean;
}

// Props interfaces for React components
export interface {{.Name}}IndexProps {
  data: {{.Name}}ListResponse;
  filters: {{.Name}}ListRequest;
  stats?: {{.Name}}Stats;
  permissions: {
    canCreate: boolean;
    canEdit: boolean;
    canDelete: boolean;
    canManage: boolean;
    canExport: boolean;
    canViewReports: boolean;
  };
}

export interface {{.Name}}FormProps {
  {{.LowerName}}?: {{.Name}};
  onSubmit: (data: {{.Name}}FormData) => void;
  onCancel: () => void;
  isLoading?: boolean;
}

export interface {{.Name}}DetailProps {
  {{.LowerName}}: {{.Name}};
  onEdit: () => void;
  onDelete: () => void;
  canEdit: boolean;
  canDelete: boolean;
}
`

	return receiver.writeFileFromTemplate(config.UITypesPath, template, config, force)
}

func (receiver *MakeCrudE2E) generateUIComponents(ctx console.Context, config ResourceConfig, force bool) error {
	// Create components directory
	os.MkdirAll(config.UIComponentsPath, 0755)

	// Generate column definitions
	columnsFile := filepath.Join(config.UIComponentsPath, fmt.Sprintf("%sColumns.tsx", config.Name))
	columnsTemplate := `import React from 'react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Edit, Eye, Trash2 } from 'lucide-react';
import { {{.Name}} } from '@/types/{{.LowerName}}';

// Column definitions for {{.Name}} table
export const {{.LowerName}}Columns = [
  {
    accessorKey: 'id',
    header: 'ID',
    cell: ({ row }: { row: { original: {{.Name}} } }) => (
      <span className="font-mono text-sm">#{row.original.id}</span>
    ),
  },
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }: { row: { original: {{.Name}} } }) => (
      <div className="font-medium">{row.original.name}</div>
    ),
  },
  {
    accessorKey: 'description',
    header: 'Description',
    cell: ({ row }: { row: { original: {{.Name}} } }) => (
      <div className="max-w-xs truncate text-muted-foreground">
        {row.original.description || 'No description'}
      </div>
    ),
  },
  {
    accessorKey: 'is_active',
    header: 'Status',
    cell: ({ row }: { row: { original: {{.Name}} } }) => (
      <Badge variant={row.original.is_active ? 'default' : 'secondary'}>
        {row.original.is_active ? 'Active' : 'Inactive'}
      </Badge>
    ),
  },
  {
    accessorKey: 'created_at',
    header: 'Created',
    cell: ({ row }: { row: { original: {{.Name}} } }) => (
      <span className="text-sm text-muted-foreground">
        {new Date(row.original.created_at).toLocaleDateString()}
      </span>
    ),
  },
  {
    id: 'actions',
    header: 'Actions',
    cell: ({ 
      row, 
      onView, 
      onEdit, 
      onDelete,
      canEdit,
      canDelete 
    }: { 
      row: { original: {{.Name}} };
      onView: ({{.LowerName}}: {{.Name}}) => void;
      onEdit: ({{.LowerName}}: {{.Name}}) => void;
      onDelete: ({{.LowerName}}: {{.Name}}) => void;
      canEdit: boolean;
      canDelete: boolean;
    }) => (
      <div className="flex items-center space-x-2">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onView(row.original)}
        >
          <Eye className="h-4 w-4" />
        </Button>
        {canEdit && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onEdit(row.original)}
          >
            <Edit className="h-4 w-4" />
          </Button>
        )}
        {canDelete && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onDelete(row.original)}
            className="text-destructive hover:text-destructive"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        )}
      </div>
    ),
  },
];

// Mobile-friendly columns
export const {{.LowerName}}ColumnsMobile = [
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }: { row: { original: {{.Name}} } }) => (
      <div>
        <div className="font-medium">{row.original.name}</div>
        <div className="text-sm text-muted-foreground">
          {row.original.description ? row.original.description.substring(0, 50) + '...' : 'No description'}
        </div>
        <div className="flex items-center space-x-2 mt-1">
          <Badge variant={row.original.is_active ? 'default' : 'secondary'} className="text-xs">
            {row.original.is_active ? 'Active' : 'Inactive'}
          </Badge>
          <span className="text-xs text-muted-foreground">
            #{row.original.id}
          </span>
        </div>
      </div>
    ),
  },
  {
    id: 'actions',
    header: 'Actions',
    cell: ({ 
      row, 
      onView, 
      onEdit, 
      onDelete,
      canEdit,
      canDelete 
    }: { 
      row: { original: {{.Name}} };
      onView: ({{.LowerName}}: {{.Name}}) => void;
      onEdit: ({{.LowerName}}: {{.Name}}) => void;
      onDelete: ({{.LowerName}}: {{.Name}}) => void;
      canEdit: boolean;
      canDelete: boolean;
    }) => (
      <div className="flex flex-col space-y-1">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onView(row.original)}
          className="h-8"
        >
          <Eye className="h-4 w-4" />
        </Button>
        {canEdit && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onEdit(row.original)}
            className="h-8"
          >
            <Edit className="h-4 w-4" />
          </Button>
        )}
        {canDelete && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onDelete(row.original)}
            className="h-8 text-destructive hover:text-destructive"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        )}
      </div>
    ),
  },
];

// Filter definitions
export const {{.LowerName}}Filters = [
  {
    key: 'is_active',
    label: 'Status',
    type: 'select',
    options: [
      { label: 'All', value: '' },
      { label: 'Active', value: 'true' },
      { label: 'Inactive', value: 'false' },
    ],
  },
];

// Quick filter buttons
export const {{.LowerName}}QuickFilters = [
  {
    key: 'all',
    label: 'All {{.PluralName}}',
    icon: React.createElement('span', { className: 'text-xs' }, '📋'),
    filters: {},
  },
  {
    key: 'active',
    label: 'Active',
    icon: React.createElement('span', { className: 'text-xs' }, '✅'),
    filters: { is_active: true },
  },
  {
    key: 'inactive',
    label: 'Inactive',
    icon: React.createElement('span', { className: 'text-xs' }, '❌'),
    filters: { is_active: false },
  },
];
`

	if err := receiver.writeFileFromTemplate(columnsFile, columnsTemplate, config, force); err != nil {
		return err
	}

	// Generate form components
	formsFile := filepath.Join(config.UIComponentsPath, fmt.Sprintf("%sForms.tsx", config.Name))
	formsTemplate := `import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { {{.Name}}, {{.Name}}FormData, {{.Name}}DetailProps } from '@/types/{{.LowerName}}';

// Create form component
export function {{.Name}}CreateForm({ 
  onSubmit, 
  onCancel, 
  isLoading = false 
}: {
  onSubmit: (data: {{.Name}}FormData) => void;
  onCancel: () => void;
  isLoading?: boolean;
}) {
  const [formData, setFormData] = useState<{{.Name}}FormData>({
    name: '',
    description: '',
    is_active: true,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Basic validation
    const newErrors: Record<string, string> = {};
    if (!formData.name.trim()) {
      newErrors.name = 'Name is required';
    }
    
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }
    
    setErrors({});
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Name *</Label>
        <Input
          id="name"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder="Enter {{.LowerName}} name"
          className={errors.name ? 'border-destructive' : ''}
        />
        {errors.name && (
          <p className="text-sm text-destructive">{errors.name}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <Textarea
          id="description"
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          placeholder="Enter {{.LowerName}} description"
          rows={3}
        />
      </div>

      <div className="flex items-center space-x-2">
        <Switch
          id="is_active"
          checked={formData.is_active}
          onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
        />
        <Label htmlFor="is_active">Active</Label>
      </div>

      <div className="flex justify-end space-x-2">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? 'Creating...' : 'Create {{.Name}}'}
        </Button>
      </div>
    </form>
  );
}

// Edit form component
export function {{.Name}}EditForm({ 
  {{.LowerName}}, 
  onSubmit, 
  onCancel, 
  isLoading = false 
}: {
  {{.LowerName}}: {{.Name}};
  onSubmit: (data: {{.Name}}FormData) => void;
  onCancel: () => void;
  isLoading?: boolean;
}) {
  const [formData, setFormData] = useState<{{.Name}}FormData>({
    name: {{.LowerName}}.name,
    description: {{.LowerName}}.description,
    is_active: {{.LowerName}}.is_active,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Basic validation
    const newErrors: Record<string, string> = {};
    if (!formData.name.trim()) {
      newErrors.name = 'Name is required';
    }
    
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }
    
    setErrors({});
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Name *</Label>
        <Input
          id="name"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder="Enter {{.LowerName}} name"
          className={errors.name ? 'border-destructive' : ''}
        />
        {errors.name && (
          <p className="text-sm text-destructive">{errors.name}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <Textarea
          id="description"
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          placeholder="Enter {{.LowerName}} description"
          rows={3}
        />
      </div>

      <div className="flex items-center space-x-2">
        <Switch
          id="is_active"
          checked={formData.is_active}
          onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
        />
        <Label htmlFor="is_active">Active</Label>
      </div>

      <div className="flex justify-end space-x-2">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? 'Updating...' : 'Update {{.Name}}'}
        </Button>
      </div>
    </form>
  );
}

// Detail view component
export function {{.Name}}DetailView({ 
  {{.LowerName}}, 
  onEdit, 
  onDelete, 
  canEdit, 
  canDelete 
}: {{.Name}}DetailProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{{{.LowerName}}.name}</CardTitle>
        <CardDescription>
          {{.Name}} ID: #{{{.LowerName}}.id}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <Label className="text-sm font-medium">Description</Label>
          <p className="text-sm text-muted-foreground mt-1">
            {{{.LowerName}}.description || 'No description provided'}
          </p>
        </div>

        <div>
          <Label className="text-sm font-medium">Status</Label>
          <p className="text-sm text-muted-foreground mt-1">
            {{{.LowerName}}.is_active ? 'Active' : 'Inactive'}
          </p>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <Label className="text-sm font-medium">Created</Label>
            <p className="text-sm text-muted-foreground mt-1">
              {new Date({{.LowerName}}.created_at).toLocaleDateString()}
            </p>
          </div>
          <div>
            <Label className="text-sm font-medium">Updated</Label>
            <p className="text-sm text-muted-foreground mt-1">
              {new Date({{.LowerName}}.updated_at).toLocaleDateString()}
            </p>
          </div>
        </div>

        <div className="flex justify-end space-x-2 pt-4">
          {canEdit && (
            <Button onClick={onEdit}>
              Edit {{.Name}}
            </Button>
          )}
          {canDelete && (
            <Button variant="destructive" onClick={onDelete}>
              Delete {{.Name}}
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
`

	return receiver.writeFileFromTemplate(formsFile, formsTemplate, config, force)
}

func (receiver *MakeCrudE2E) generateUIPages(ctx console.Context, config ResourceConfig, force bool) error {
	// Create pages directory
	os.MkdirAll(config.UIPagesPath, 0755)

	indexFile := filepath.Join(config.UIPagesPath, "Index.tsx")
	indexTemplate := `import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import { Download, Upload, Plus, RefreshCw } from 'lucide-react';
import { 
  {{.Name}}, 
  {{.Name}}IndexProps,
  {{.Name}}FormData 
} from '@/types/{{.LowerName}}';
import { CrudPage } from '@/components/Crud/CrudPage';
import { {{.LowerName}}Columns, {{.LowerName}}ColumnsMobile, {{.LowerName}}Filters, {{.LowerName}}QuickFilters } from '@/components/{{.PluralName}}/{{.Name}}Columns';
import { {{.Name}}CreateForm, {{.Name}}EditForm, {{.Name}}DetailView } from '@/components/{{.PluralName}}/{{.Name}}Forms';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';

export default function {{.PluralName}}Index({ 
  data, 
  filters, 
  stats,
  permissions 
}: {{.Name}}IndexProps) {
  const isMobile = useIsMobile();
  
  // Debug logging
  console.log('{{.PluralName}}Index - data:', data);
  console.log('{{.PluralName}}Index - filters:', filters);
  console.log('{{.PluralName}}Index - stats:', stats);
  console.log('{{.PluralName}}Index - permissions:', permissions);
  
  // Dialog states
  const [showImportDialog, setShowImportDialog] = useState(false);
  const [showExportDialog, setShowExportDialog] = useState(false);
  const [selected{{.PluralName}}, setSelected{{.PluralName}}] = useState<{{.Name}}[]>([]);

  // Handle bulk operations
  const handleBulkAction = async (action: string, selectedIds: number[]) => {
    if (selectedIds.length === 0) return;

    // Get selected {{.LowerName}} objects
    const selected = data.data.filter({{.LowerName}} => selectedIds.includes({{.LowerName}}.id));
    setSelected{{.PluralName}}(selected);

    const operations: Record<string, () => void> = {
      delete: () => handleBulkDelete(selectedIds),
      activate: () => handleBulkStatusUpdate(selectedIds, true),
      deactivate: () => handleBulkStatusUpdate(selectedIds, false),
      export: () => handleBulkExport(selectedIds),
    };

    const operation = operations[action];
    if (operation) {
      operation();
    }
  };

  const handleBulkDelete = ({{.LowerName}}Ids: number[]) => {
    const confirmMessage = ` + "`" + `Are you sure you want to delete ${{{.LowerName}}Ids.length} {{.LowerName}}(s)? This action cannot be undone.` + "`" + `;
    if (confirm(confirmMessage)) {
      router.delete('/api/{{.LowerPluralName}}/bulk', {
        data: { {{.LowerName}}Ids },
        onSuccess: () => {
          // Refresh will be handled by the parent
        },
      });
    }
  };

  const handleBulkStatusUpdate = ({{.LowerName}}Ids: number[], isActive: boolean) => {
    router.put('/api/{{.LowerPluralName}}/bulk/status', {
      {{.LowerName}}Ids,
      is_active: isActive,
    });
  };

  const handleBulkExport = ({{.LowerName}}Ids: number[]) => {
    const format = prompt('Export format (csv, json, excel):') || 'csv';
    const options = {
      format: format as any,
      filters: { ...filters, {{.LowerName}}Ids },
    };
    
    // Trigger download
    window.open(` + "`" + `/api/{{.LowerPluralName}}/export?${new URLSearchParams(options as any).toString()}` + "`" + `);
  };

  const handleRefresh = () => {
    router.reload({ only: ['data', 'stats'] });
  };

  return (
    <Admin title="{{.PluralName}}">
      <Head title="{{.PluralName}} - Management" />
      
      <div className="space-y-6">
        {/* Statistics Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total {{.PluralName}}</CardTitle>
                <span className="text-2xl">📋</span>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total{{.PluralName}}}</div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Active</CardTitle>
                <div className="h-4 w-4 bg-green-500 rounded-full" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">{stats.active{{.PluralName}}}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.total{{.PluralName}} > 0 && ` + "`" + `${Math.round((stats.active{{.PluralName}} / stats.total{{.PluralName}}) * 100)}% active` + "`" + `}
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Inactive</CardTitle>
                <div className="h-4 w-4 bg-gray-500 rounded-full" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-gray-600">{stats.inactive{{.PluralName}}}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.total{{.PluralName}} > 0 && ` + "`" + `${Math.round((stats.inactive{{.PluralName}} / stats.total{{.PluralName}}) * 100)}% inactive` + "`" + `}
                </p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Quick Filter Buttons */}
        <div className="flex flex-wrap gap-2">
          {{{.LowerName}}QuickFilters.map((filter) => (
            <Button
              key={filter.key}
              variant={JSON.stringify(filters) === JSON.stringify(filter.filters) ? 'default' : 'outline'}
              size="sm"
              onClick={() => {
                router.get('/admin/{{.LowerPluralName}}', filter.filters, {
                  preserveState: true,
                  preserveScroll: true,
                  only: ['data', 'filters', 'stats'],
                });
              }}
              className="flex items-center space-x-2"
            >
              {filter.icon}
              <span>{filter.label}</span>
            </Button>
          ))}
        </div>

        {/* Management Actions */}
        {permissions.canManage && (
          <div className="flex flex-wrap gap-2">
            <Button 
              variant="outline" 
              size="sm"
              onClick={() => setShowImportDialog(true)}
            >
              <Upload className="h-4 w-4 mr-2" />
              Import {{.PluralName}}
            </Button>
            <Button 
              variant="outline" 
              size="sm" 
              onClick={() => setShowExportDialog(true)}
            >
              <Download className="h-4 w-4 mr-2" />
              Export {{.PluralName}}
            </Button>
            <Button 
              variant="outline" 
              size="sm"
              onClick={handleRefresh}
            >
              <RefreshCw className="h-4 w-4 mr-2" />
              Refresh
            </Button>
          </div>
        )}

        {/* Main CRUD Component */}
        <CrudPage<{{.Name}}>
          data={data}
          filters={filters}
          title="{{.PluralName}}"
          resourceName="{{.LowerPluralName}}"
          columns={isMobile ? {{.LowerName}}ColumnsMobile : {{.LowerName}}Columns}
          customFilters={ {{.LowerName}}Filters}
          createForm={ {{.Name}}CreateForm}
          editForm={ {{.Name}}EditForm}
          detailView={ {{.Name}}DetailView}
          onBulkAction={handleBulkAction}
          onRefresh={handleRefresh}
          canCreate={permissions.canCreate}
          canEdit={permissions.canEdit}
          canDelete={permissions.canDelete}
          canView={true}
        />
      </div>
    </Admin>
  );
}
`

	return receiver.writeFileFromTemplate(indexFile, indexTemplate, config, force)
}

// Helper method to write file from template
func (receiver *MakeCrudE2E) writeFileFromTemplate(filePath, template string, config ResourceConfig, force bool) error {
	// Check if file exists and force is not set
	if !force {
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			return fmt.Errorf("file %s already exists (use --force to overwrite)", filePath)
		}
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Parse template
	parsedTemplate := receiver.parseTemplate(template, config)

	// Write file
	return os.WriteFile(filePath, []byte(parsedTemplate), 0644)
}

// Simple template parser (replace {{.Field}} with config values)
func (receiver *MakeCrudE2E) parseTemplate(template string, config ResourceConfig) string {
	replacements := map[string]string{
		"{{.Name}}":            config.Name,
		"{{.LowerName}}":       config.LowerName,
		"{{.PluralName}}":      config.PluralName,
		"{{.LowerPluralName}}": config.LowerPluralName,
		"{{.SnakeName}}":       config.SnakeName,
		"{{.SnakePluralName}}": config.SnakePluralName,
		"{{.KebabName}}":       config.KebabName,
		"{{.KebabPluralName}}": config.KebabPluralName,
		"{{.DisplayName}}":     config.DisplayName,
		"{{.TableName}}":       config.TableName,
	}

	result := template
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}