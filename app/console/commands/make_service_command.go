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

type MakeServiceCommand struct{}

func (receiver *MakeServiceCommand) Signature() string {
	return "make:service {name : The name of the service}"
}

func (receiver *MakeServiceCommand) Description() string {
	return "Create a new CRUD service class with repository pattern"
}

func (receiver *MakeServiceCommand) Extend() command.Extend {
	return command.Extend{}
}

func (receiver *MakeServiceCommand) Handle(ctx console.Context) error {
	name := ctx.Argument(0)
	if name == "" {
		return fmt.Errorf("service name is required")
	}

	// Ensure name ends with "Service"
	if !strings.HasSuffix(name, "Service") {
		name += "Service"
	}

	servicePath := filepath.Join("app", "services")
	if err := os.MkdirAll(servicePath, 0755); err != nil {
		return err
	}

	filename := filepath.Join(servicePath, strings.ToLower(strings.TrimSuffix(name, "Service"))+"_service.go")

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("service %s already exists", name)
	}

	// Generate service content
	content, err := generateServiceContent(name)
	if err != nil {
		return err
	}

	// Write file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return err
	}

	ctx.Info(fmt.Sprintf("Service [%s] created successfully", filename))
	ctx.Info("Don't forget to create the corresponding repository:")
	ctx.Info(fmt.Sprintf("  go run . artisan make:repository %sRepository", strings.TrimSuffix(name, "Service")))
	return nil
}

func generateServiceContent(serviceName string) (string, error) {
	resourceName := strings.TrimSuffix(serviceName, "Service")
	modelName := strings.Title(resourceName)
	tableName := strings.ToLower(resourceName) + "s"
	repositoryName := resourceName + "Repository"

	tmpl := `package services

import (
	"fmt"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/models"
	"players/app/repositories"
	"strconv"
)

// {{.ServiceName}} - enhanced with repository pattern
type {{.ServiceName}} struct {
	repository contracts.SearchableRepository
	authHelper contracts.AuthHelper
}

// New{{.ServiceName}} creates a new {{.ResourceName}} service with repository
func New{{.ServiceName}}() *{{.ServiceName}} {
	return &{{.ServiceName}}{
		repository: repositories.New{{.RepositoryName}}(),
		authHelper: helpers.NewAuthHelper(),
	}
}

// GetList with built-in pagination, sorting, filtering using repository
func (s *{{.ServiceName}}) GetList(req helpers.ListRequest) (*contracts.PaginatedResult, error) {
	req.SetDefaults()

	// Start with a fresh query
	repo := s.repository.Reset()

	// Apply search if provided
	if req.Search != "" {
		{{.ResourceNameLower}} := models.{{.ModelName}}{}
		repo = repo.Search(req.Search, {{.ResourceNameLower}}.SearchFields())
	}

	// Apply sorting
	if req.Sort != "" {
		// TODO: Parse sort string properly (e.g., "title ASC")
		repo = repo.OrderBy("id", "DESC") // Default, should parse req.Sort
	}

	// Execute paginated query
	return repo.Paginate(req.Page, req.PageSize)
}

// GetListAdvanced with additional filters using repository
func (s *{{.ServiceName}}) GetListAdvanced(req helpers.ListRequest, filters map[string]interface{}) (*contracts.PaginatedResult, error) {
	req.SetDefaults()

	// Start with a fresh query
	repo := s.repository.Reset()

	// Apply search if provided
	if req.Search != "" {
		{{.ResourceNameLower}} := models.{{.ModelName}}{}
		repo = repo.SearchWithFilters(req.Search, {{.ResourceNameLower}}.SearchFields(), filters)
	} else {
		// Apply filters without search
		for field, value := range filters {
			switch field {
			// TODO: Add specific filters for {{.ModelName}}
			case "status", "name":
				repo = repo.Where(field, "=", value)
			case "min_price":
				repo = repo.Where("price", ">=", value)
			case "max_price":
				repo = repo.Where("price", "<=", value)
			}
		}
	}

	// Apply sorting
	if req.Sort != "" {
		repo = repo.OrderBy("id", "DESC") // Should parse req.Sort properly
	}

	// Execute paginated query
	return repo.Paginate(req.Page, req.PageSize)
}

// GetByID - using repository
func (s *{{.ServiceName}}) GetByID(id uint) (*models.{{.ModelName}}, error) {
	result, err := s.repository.Find(id)
	if err != nil {
		return nil, fmt.Errorf("{{.ResourceNameLower}} not found: %w", err)
	}

	{{.ResourceNameLower}}, ok := result.(*models.{{.ModelName}})
	if !ok {
		return nil, fmt.Errorf("invalid {{.ResourceNameLower}} data")
	}

	return {{.ResourceNameLower}}, nil
}

// Create - using repository with validation
func (s *{{.ServiceName}}) Create(data map[string]interface{}) (*models.{{.ModelName}}, error) {
	// Basic validation
	if err := s.validate{{.ModelName}}Data(data, false); err != nil {
		return nil, err
	}

	// TODO: Set default values if needed
	// if _, exists := data["status"]; !exists {
	//     data["status"] = "ACTIVE"
	// }

	// Create using repository
	result, err := s.repository.Create(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create {{.ResourceNameLower}}: %w", err)
	}

	{{.ResourceNameLower}}, ok := result.(*models.{{.ModelName}})
	if !ok {
		return nil, fmt.Errorf("invalid {{.ResourceNameLower}} data")
	}

	return {{.ResourceNameLower}}, nil
}

// Update - using repository with validation
func (s *{{.ServiceName}}) Update(id uint, data map[string]interface{}) (*models.{{.ModelName}}, error) {
	// Check if {{.ResourceNameLower}} exists
	_, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate update data
	if err := s.validate{{.ModelName}}Data(data, true); err != nil {
		return nil, err
	}

	// Update using repository
	result, err := s.repository.Update(id, data)
	if err != nil {
		return nil, fmt.Errorf("failed to update {{.ResourceNameLower}}: %w", err)
	}

	{{.ResourceNameLower}}, ok := result.(*models.{{.ModelName}})
	if !ok {
		return nil, fmt.Errorf("invalid {{.ResourceNameLower}} data")
	}

	return {{.ResourceNameLower}}, nil
}

// Delete - using repository
func (s *{{.ServiceName}}) Delete(id uint) error {
	// Check if {{.ResourceNameLower}} exists
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// Delete using repository
	err = s.repository.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete {{.ResourceNameLower}}: %w", err)
	}

	return nil
}

// GetRepository returns the repository (for testing or advanced usage)
func (s *{{.ServiceName}}) GetRepository() contracts.Repository {
	return s.repository
}

// SetRepository sets the repository (for dependency injection)
func (s *{{.ServiceName}}) SetRepository(repo contracts.Repository) {
	if searchableRepo, ok := repo.(contracts.SearchableRepository); ok {
		s.repository = searchableRepo
	}
}

// validate{{.ModelName}}Data performs validation
func (s *{{.ServiceName}}) validate{{.ModelName}}Data(data map[string]interface{}, isUpdate bool) error {
	// TODO: Add validation logic specific to {{.ModelName}}
	
	// Required fields for creation
	if !isUpdate {
		requiredFields := []string{"name"} // TODO: Update with actual required fields
		for _, field := range requiredFields {
			if value, exists := data[field]; !exists || value == "" {
				return fmt.Errorf("%s is required", field)
			}
		}
	}

	// TODO: Add specific validation rules
	// Example: Validate name length
	if name, ok := data["name"].(string); ok && name != "" {
		if len(name) > 255 {
			return fmt.Errorf("name cannot exceed 255 characters")
		}
	}

	// Example: Validate numeric fields
	if price, exists := data["price"]; exists {
		switch v := price.(type) {
		case float64:
			if v < 0 {
				return fmt.Errorf("price cannot be negative")
			}
		case string:
			if p, err := strconv.ParseFloat(v, 64); err != nil || p < 0 {
				return fmt.Errorf("invalid price format")
			}
		}
	}

	return nil
}
`

	t, err := template.New("service").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = t.Execute(&result, map[string]string{
		"ServiceName":        serviceName,
		"ModelName":          modelName,
		"ResourceName":       resourceName,
		"ResourceNameLower":  strings.ToLower(resourceName),
		"TableName":          tableName,
		"RepositoryName":     repositoryName,
	})

	return result.String(), err
}