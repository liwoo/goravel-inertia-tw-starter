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

type MakeRepositoryCommand struct{}

func (receiver *MakeRepositoryCommand) Signature() string {
	return "make:repository {name : The name of the repository}"
}

func (receiver *MakeRepositoryCommand) Description() string {
	return "Create a new repository class with searchable capabilities"
}

func (receiver *MakeRepositoryCommand) Extend() command.Extend {
	return command.Extend{}
}

func (receiver *MakeRepositoryCommand) Handle(ctx console.Context) error {
	name := ctx.Argument(0)
	if name == "" {
		return fmt.Errorf("repository name is required")
	}

	// Ensure name ends with "Repository"
	if !strings.HasSuffix(name, "Repository") {
		name += "Repository"
	}

	repoPath := filepath.Join("app", "repositories")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return err
	}

	filename := filepath.Join(repoPath, strings.ToLower(strings.TrimSuffix(name, "Repository"))+"_repository.go")

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("repository %s already exists", name)
	}

	// Generate repository content
	content, err := generateRepositoryContent(name)
	if err != nil {
		return err
	}

	// Write file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return err
	}

	ctx.Info(fmt.Sprintf("Repository [%s] created successfully", filename))
	return nil
}

func generateRepositoryContent(repositoryName string) (string, error) {
	resourceName := strings.TrimSuffix(repositoryName, "Repository")
	modelName := strings.Title(resourceName)
	tableName := strings.ToLower(resourceName) + "s"

	tmpl := `package repositories

import (
	"players/app/contracts"
	"players/app/models"
)

// {{.RepositoryName}} handles {{.ResourceNameLower}}-specific data access
type {{.RepositoryName}} struct {
	*SearchableBaseRepository
}

// New{{.RepositoryName}} creates a new {{.ResourceNameLower}} repository
func New{{.RepositoryName}}() contracts.SearchableRepository {
	return &{{.RepositoryName}}{
		SearchableBaseRepository: NewSearchableRepository(&models.{{.ModelName}}{}, "{{.TableName}}"),
	}
}

// TODO: Add domain-specific repository methods below

// GetByName retrieves a {{.ResourceNameLower}} by name
func (r *{{.RepositoryName}}) GetByName(name string) (*models.{{.ModelName}}, error) {
	var {{.ResourceNameLower}} models.{{.ModelName}}
	err := r.builder.Where("name", "=", name).First(&{{.ResourceNameLower}})
	if err != nil {
		return nil, err
	}
	return &{{.ResourceNameLower}}, nil
}

// GetByStatus retrieves {{.ResourceNameLower}}s by status
func (r *{{.RepositoryName}}) GetByStatus(status string) ([]models.{{.ModelName}}, error) {
	var {{.ResourceNameLower}}s []models.{{.ModelName}}
	err := r.builder.Where("status", "=", status).Find(&{{.ResourceNameLower}}s)
	if err != nil {
		return nil, err
	}
	return {{.ResourceNameLower}}s, nil
}

// GetActive retrieves all active {{.ResourceNameLower}}s
func (r *{{.RepositoryName}}) GetActive() ([]models.{{.ModelName}}, error) {
	return r.GetByStatus("ACTIVE")
}

// UpdateStatus updates the status of a {{.ResourceNameLower}}
func (r *{{.RepositoryName}}) UpdateStatus(id interface{}, status string) error {
	data := map[string]interface{}{
		"status": status,
	}
	_, err := r.Update(id, data)
	return err
}

// CountByStatus counts {{.ResourceNameLower}}s by status
func (r *{{.RepositoryName}}) CountByStatus(status string) (int64, error) {
	conditions := map[string]interface{}{
		"status": status,
	}
	return r.Count(conditions)
}

// TODO: Add more domain-specific methods as needed
// Examples:
// - GetByCategoryID(categoryID uint) ([]models.{{.ModelName}}, error)
// - GetPopular(limit int) ([]models.{{.ModelName}}, error)
// - GetRecentlyCreated(days int) ([]models.{{.ModelName}}, error)
// - BulkUpdateStatus(ids []interface{}, status string) error
`

	t, err := template.New("repository").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = t.Execute(&result, map[string]string{
		"RepositoryName":     repositoryName,
		"ModelName":          modelName,
		"ResourceName":       resourceName,
		"ResourceNameLower":  strings.ToLower(resourceName),
		"TableName":          tableName,
	})

	return result.String(), err
}