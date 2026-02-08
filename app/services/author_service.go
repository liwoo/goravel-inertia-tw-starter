package services

import (
	"fmt"

	"books-database/app/contracts"
	"books-database/app/models"

	"github.com/goravel/framework/facades"
)

// AuthorService implements author-specific business logic using the builder pattern
type AuthorService struct {
	contracts.CrudServiceContract // Embedded - automatically exposes all methods!
}

// NewAuthorService creates a new author service using the builder pattern
func NewAuthorService() *AuthorService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Author]("authors", "id").
		WithSearchFields("first_name", "last_name", "email", "nationality").
		WithSortFields("id", "first_name", "last_name", "email", "nationality", "status", "created_at", "updated_at").
		WithFilterFields("status", "nationality").
		WithValidationRules(map[string]interface{}{
			"firstName":   "required|string|max:100",
			"lastName":    "required|string|max:100",
			"email":       "string|max:255",
			"website":     "string|max:255",
			"bio":         "string",
			"birthDate":   "string|max:20",
			"nationality": "string|max:100",
			"photoUrl":    "string|max:500",
			"status":      "string|in:ACTIVE,INACTIVE",
		}).
		WithRelations("Creator", "Updater", "Books").
		WithDefaultSort("created_at", "DESC").
		WithSoftDeletes().
		WithScopeFiltering("authors", "created_by").
		WithBeforeCreate(func(data map[string]interface{}) error {
			// The created_by field should already be set by the controller
			// Just ensure it's properly formatted if present
			if createdBy, exists := data["created_by"]; exists && createdBy != nil {
				switch v := createdBy.(type) {
				case float64:
					data["created_by"] = uint(v)
				case int:
					data["created_by"] = uint(v)
				case uint:
					// Already correct type
				default:
					if fmt.Sprintf("%v", v) != "" {
						// Keep the value as is, let GORM handle the conversion
					}
				}
			}

			// Set default status if not provided
			if _, exists := data["status"]; !exists {
				data["status"] = "ACTIVE"
			}

			return nil
		}).
		WithTenantAwareness().
		Build()

	authorServiceInstance := &AuthorService{
		CrudServiceContract: service,
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, authorServiceInstance, "AuthorService")

	return authorServiceInstance
}

// Override GetColumnMapping to include author-specific mappings
// NOTE: Only map fields used in sort/filter query params.
// Do NOT map fields that come from ToCreateData/ToUpdateData, as applyFieldMapping
// would break setFieldsRecursively (which matches against model json tags).
// MapSortField auto-converts camelCase→snake_case, so explicit mapping is only
// needed for non-standard conversions.
func (s *AuthorService) GetColumnMapping() map[string]string {
	mapping := s.CrudServiceContract.GetColumnMapping()
	mapping["createdAt"] = "created_at"
	mapping["updatedAt"] = "updated_at"
	return mapping
}

// GetFilterDefinitions returns filter definitions for the authors resource
func (s *AuthorService) GetFilterDefinitions() []contracts.FilterDefinition {
	return []contracts.FilterDefinition{
		contracts.NewFilterDefinition(
			"status",
			"Author Status",
			contracts.FilterTypeEnum,
			&[]string{
				"ACTIVE",
				"INACTIVE",
			},
		),
		contracts.NewFilterDefinition(
			"nationality",
			"Nationality",
			contracts.FilterTypeString,
			nil,
		),
	}
}

// GetAuthorStatistics returns statistics about authors
func (s *AuthorService) GetAuthorStatistics() (map[string]interface{}, error) {
	var stats struct {
		TotalAuthors    int64
		ActiveAuthors   int64
		InactiveAuthors int64
	}

	// Get total authors (GORM automatically handles deleted_at IS NULL for soft delete models)
	stats.TotalAuthors, _ = facades.Orm().Query().Model(&models.Author{}).Count()

	// Get active authors
	stats.ActiveAuthors, _ = facades.Orm().Query().Model(&models.Author{}).Where("status = ?", "ACTIVE").Count()

	// Get inactive authors
	stats.InactiveAuthors, _ = facades.Orm().Query().Model(&models.Author{}).Where("status = ?", "INACTIVE").Count()

	return map[string]interface{}{
		"totalAuthors":    stats.TotalAuthors,
		"activeAuthors":   stats.ActiveAuthors,
		"inactiveAuthors": stats.InactiveAuthors,
	}, nil
}
