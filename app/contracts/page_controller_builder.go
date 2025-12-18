package contracts

import (
	"github.com/goravel/framework/contracts/http"
	"starter-project/app/auth"
)

// PageControllerBuilder uses a step-by-step builder pattern to ensure all required methods are implemented
type PageControllerBuilder struct {
	controller       *GenericPageController
	permissionAction auth.CorePermissionAction
	propsFunc        func(result *PaginatedResult, req *ListRequest, permissions map[string]bool, stats map[string]interface{}) PageProps
}

// Step 1: Required - Set service
type PageControllerBuilderWithService struct {
	builder *PageControllerBuilder
}

// Step 2: Required - Set page permission
type PageControllerBuilderWithPermission struct {
	builder *PageControllerBuilder
}

// Step 3: Required - Set props configuration
type PageControllerBuilderWithProps struct {
	builder *PageControllerBuilder
}

// Step 4: Optional configurations and build
type PageControllerBuilderComplete struct {
	builder *PageControllerBuilder
}

// NewPageControllerBuilder creates a new page controller builder that enforces required setup
func NewPageControllerBuilder(resourceName string, pageComponent string) *PageControllerBuilder {
	config := GenericPageConfig{
		ResourceType:  resourceName,
		PageComponent: pageComponent,
	}
	controller := NewGenericPageController(config)

	return &PageControllerBuilder{
		controller: controller,
	}
}

// Step 1: Set service (REQUIRED)
func (b *PageControllerBuilder) WithService(service CrudServiceContract) *PageControllerBuilderWithService {
	if service == nil {
		panic("Service cannot be nil")
	}
	b.controller.service = service
	return &PageControllerBuilderWithService{builder: b}
}

// Step 2: Set page permission (REQUIRED)
func (b *PageControllerBuilderWithService) WithPagePermission(
	permissionAction auth.CorePermissionAction,
) *PageControllerBuilderWithPermission {
	// Store the permission action for use in the custom permission check
	b.builder.permissionAction = permissionAction
	return &PageControllerBuilderWithPermission{builder: b.builder}
}

// Step 3: Set props configuration (REQUIRED)
func (b *PageControllerBuilderWithPermission) WithPropsConfiguration(
	getProps func(result *PaginatedResult, req *ListRequest, permissions map[string]bool, stats map[string]interface{}) PageProps,
) *PageControllerBuilderComplete {
	if getProps == nil {
		panic("GetProps function cannot be nil")
	}
	b.builder.propsFunc = getProps
	return &PageControllerBuilderComplete{builder: b.builder}
}

// Optional configurations
func (b *PageControllerBuilderComplete) WithCustomPermissionCheck(
	check func(ctx http.Context) error,
) *PageControllerBuilderComplete {
	b.builder.controller.customPermissionCheck = check
	return b
}

func (b *PageControllerBuilderComplete) WithStatisticsEnabled(
	enabled bool,
	statsBuilder func(controller *GenericPageController) map[string]interface{},
) *PageControllerBuilderComplete {
	b.builder.controller.statsEnabled = enabled
	b.builder.controller.statsBuilder = statsBuilder
	return b
}

func (b *PageControllerBuilderComplete) WithServiceIdentifier(identifier auth.ServiceRegistry) *PageControllerBuilderComplete {
	b.builder.controller.serviceIdentifier = identifier
	return b
}

func (b *PageControllerBuilderComplete) WithRequireSuperAdmin(require bool) *PageControllerBuilderComplete {
	b.builder.controller.requireSuperAdmin = require
	return b
}

// Build ensures the controller implements the required interface
func (b *PageControllerBuilderComplete) Build() PageControllerContract {
	// Create a wrapper that ensures all interface methods are implemented
	return &pageControllerContractWrapper{
		GenericPageController: b.builder.controller,
		permissionAction:      b.builder.permissionAction,
		propsFunc:             b.builder.propsFunc,
	}
}

// pageControllerContractWrapper ensures the controller implements PageControllerContract
type pageControllerContractWrapper struct {
	*GenericPageController
	permissionAction auth.CorePermissionAction
	propsFunc        func(result *PaginatedResult, req *ListRequest, permissions map[string]bool, stats map[string]interface{}) PageProps
}

// Implement PageControllerContract methods that need customization
// The rest are inherited from GenericPageController

func (w *pageControllerContractWrapper) GetProps(result *PaginatedResult, req *ListRequest, permissions map[string]bool, stats map[string]interface{}) PageProps {
	if w.propsFunc != nil {
		return w.propsFunc(result, req, permissions, stats)
	}
	// Default implementation
	return PageProps{
		Data: PaginatedData{
			Data:        result.Data,
			CurrentPage: result.CurrentPage,
			LastPage:    result.LastPage,
			PerPage:     result.PerPage,
			Total:       result.Total,
			From:        (result.CurrentPage-1)*result.PerPage + 1,
			To:          min(result.CurrentPage*result.PerPage, int(result.Total)),
			HasNext:     result.CurrentPage < result.LastPage,
			HasPrev:     result.CurrentPage > 1,
		},
		Permissions: convertToPermissionsMap(permissions),
		Stats:       stats,
	}
}

// Helper function to convert map[string]bool to PermissionsMap
func convertToPermissionsMap(permissions map[string]bool) PermissionsMap {
	pm := PermissionsMap{
		Custom: make(map[string]bool),
	}

	// Map common permissions
	pm.CanView = permissions["canView"]
	pm.CanCreate = permissions["canCreate"]
	pm.CanEdit = permissions["canEdit"]
	pm.CanDelete = permissions["canDelete"]
	pm.CanManage = permissions["canManage"]
	pm.CanExport = permissions["canExport"]
	pm.CanBulkUpdate = permissions["canBulkUpdate"]
	pm.CanBulkDelete = permissions["canBulkDelete"]
	pm.IsAdmin = permissions["isAdmin"]
	pm.IsSuperAdmin = permissions["isSuperAdmin"]

	// Add any other permissions to custom
	for k, v := range permissions {
		switch k {
		case "canView", "canCreate", "canEdit", "canDelete", "canManage",
			"canExport", "canBulkUpdate", "canBulkDelete", "isAdmin", "isSuperAdmin":
			// Already mapped
		default:
			pm.Custom[k] = v
		}
	}

	return pm
}

// Min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Example usage that enforces compile-time checks:
//
// controller := NewPageControllerBuilder("books", "Books/BookList").
//     WithService(bookService).                           // REQUIRED
//     WithPagePermission(PermissionRead).                 // REQUIRED
//     WithPropsConfiguration(customPropsFunc).            // REQUIRED
//     WithStatisticsSettings(true, statsFunc).            // Optional
//     WithBeforeRender(auditFunc).                        // Optional
//     Build()
//
// If you skip any required step, you get a compile error!
