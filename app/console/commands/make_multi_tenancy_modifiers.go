package commands

import (
	"fmt"
	"strings"
)

// ──────────────────────────────────────────────
// Step 4: Modify BaseAuditableModel
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) modifyBaseAuditableModel() error {
	content, err := r.readFile("app/models/base_auditable.go")
	if err != nil {
		return err
	}

	// Check if already modified
	if strings.Contains(content, "TenantID") {
		return nil // Already has TenantID
	}

	// Add TenantID field to BaseAuditableModel struct after the audit fields comment
	old := `	// Audit fields
	CreatedBy *uint  ` + "`" + `gorm:"index" json:"created_by,omitempty"` + "`"
	new := `	// Tenant isolation
	TenantID *uint   ` + "`" + `gorm:"index" json:"tenant_id,omitempty"` + "`" + `
	Tenant   *Tenant ` + "`" + `gorm:"foreignKey:TenantID" json:"tenant,omitempty" swaggerignore:"true"` + "`" + `

	// Audit fields
	CreatedBy *uint  ` + "`" + `gorm:"index" json:"created_by,omitempty"` + "`"

	content = strings.Replace(content, old, new, 1)

	// Add GetTenantID and SetTenantID methods before the closing of the file
	methods := `
// GetTenantID returns the tenant ID
func (b *BaseAuditableModel) GetTenantID() *uint {
	return b.TenantID
}

// SetTenantID sets the tenant ID
func (b *BaseAuditableModel) SetTenantID(tenantID *uint) {
	b.TenantID = tenantID
}
`
	content = content + methods

	return r.writeFile("app/models/base_auditable.go", content)
}

// ──────────────────────────────────────────────
// Step 5: Modify database/kernel.go
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) modifyDatabaseKernel(timestamp string) error {
	content, err := r.readFile("database/kernel.go")
	if err != nil {
		return err
	}

	// Find the last migration entry and add new ones after it
	// Look for the closing of Migrations() return slice
	lastEntry := "&migrations.M20260207000003AddAuthorIdToBooksTable{},"
	if !strings.Contains(content, lastEntry) {
		// Try to find any migration entry pattern
		lines := strings.Split(content, "\n")
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.Contains(lines[i], "&migrations.M") && strings.Contains(lines[i], "{},") {
				lastEntry = strings.TrimSpace(lines[i])
				break
			}
		}
	}

	t1 := timestamp + "01"
	t2 := timestamp + "02"
	t3 := timestamp + "03"

	newMigrations := fmt.Sprintf(`
		&migrations.M%sCreateTenantsTable{},
		&migrations.M%sCreateUserTenantsTable{},
		&migrations.M%sAddTenantIdToAllTables{},`, t1, t2, t3)

	content = strings.Replace(content, lastEntry, lastEntry+newMigrations, 1)

	return r.writeFile("database/kernel.go", content)
}

// ──────────────────────────────────────────────
// Step 6: Modify database seeder
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) modifyDatabaseSeeder() error {
	content, err := r.readFile("database/seeders/database_seeder.go")
	if err != nil {
		return err
	}

	// Add tenant seeder after book seeder
	old := `	// Run the book seeder
	bookSeeder := &BookSeeder{}
	if err := bookSeeder.Run(); err != nil {
		return err
	}

	return nil`

	new := `	// Run the book seeder
	bookSeeder := &BookSeeder{}
	if err := bookSeeder.Run(); err != nil {
		return err
	}

	// Run the tenant seeder
	tenantSeeder := &TenantSeeder{}
	if err := tenantSeeder.Run(); err != nil {
		return err
	}

	return nil`

	content = strings.Replace(content, old, new, 1)
	return r.writeFile("database/seeders/database_seeder.go", content)
}

// ──────────────────────────────────────────────
// Step 8: Modify JwtAuth middleware
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) modifyJwtAuthMiddleware(moduleName string) error {
	content, err := r.readFile("app/http/middleware/jwt_auth.go")
	if err != nil {
		return err
	}

	// Add import for tenant package
	if !strings.Contains(content, "tenant") {
		old := `"github.com/goravel/framework/facades"`
		new := fmt.Sprintf(`"%s/app/models"
	"%s/app/tenant"

	"github.com/goravel/framework/facades"`, moduleName, moduleName)
		content = strings.Replace(content, old, new, 1)
	}

	// Add tenant verification after successful JWT parse
	// Look for ctx.Request().Next() and add tenant check before it
	marker := `ctx.Request().Next()`
	if strings.Contains(content, marker) {
		tenantCheck := `// Verify user belongs to current tenant
		currentTenant := tenant.GetFromContext(ctx)
		if currentTenant != nil {
			var user models.User
			if err := facades.Auth(ctx).User(&user); err == nil && user.ID > 0 {
				if !user.IsSuperAdminUser() {
					var ut models.UserTenant
					facades.Orm().Query().
						Where("user_id = ? AND tenant_id = ? AND is_active = ?", user.ID, currentTenant.ID, true).
						First(&ut)
					if ut.ID == 0 {
						ctx.Request().AbortWithStatusJson(403, map[string]string{
							"error": "You do not have access to this tenant",
						})
						return
					}
				}
			}
		}

		`

		// Only insert before the first ctx.Request().Next()
		content = strings.Replace(content, marker, tenantCheck+marker, 1)
	}

	return r.writeFile("app/http/middleware/jwt_auth.go", content)
}

// ──────────────────────────────────────────────
// Step 9: Modify routes
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) modifyWebRoutes(moduleName string) error {
	content, err := r.readFile("routes/web.go")
	if err != nil {
		return err
	}

	// Change Web() to Web(router route.Router)
	content = strings.Replace(content, "func Web() {", "func Web(router route.Router) {", 1)

	// Replace facades.Route().X(...) with router.X(...)
	// But NOT the GlobalMiddleware or Static calls - those move to RouteServiceProvider
	// Remove static and global middleware lines
	linesToRemove := []string{
		`	facades.Route().Static("/images", "./public/images")`,
		`	facades.Route().Static("/css", "./public/css")`,
		`	facades.Route().Static("/js", "./public/js")`,
		`	// Register the Inertia middleware globally`,
		`	facades.Route().GlobalMiddleware(inertiaMiddleware)`,
	}
	for _, line := range linesToRemove {
		content = strings.Replace(content, line+"\n", "", 1)
	}

	// Replace facades.Route(). calls with router.
	content = strings.Replace(content, "facades.Route().Get(", "router.Get(", -1)
	content = strings.Replace(content, "facades.Route().Post(", "router.Post(", -1)
	content = strings.Replace(content, "facades.Route().Put(", "router.Put(", -1)
	content = strings.Replace(content, "facades.Route().Delete(", "router.Delete(", -1)
	content = strings.Replace(content, "facades.Route().Middleware(", "router.Middleware(", -1)

	// Export the inertia middleware function so RouteServiceProvider can reference it
	content = strings.Replace(content, "func inertiaMiddleware(", "func InertiaMiddleware(", 1)

	// Remove unused facades import if no other facades calls remain
	if !strings.Contains(content, "facades.") {
		content = strings.Replace(content, `	"github.com/goravel/framework/facades"
`, "", 1)
	}

	// Add tenant page controller import and route
	// Add import
	if !strings.Contains(content, "tenants") {
		old := fmt.Sprintf(`"%s/app/http/controllers/lenders"`, moduleName)
		new := fmt.Sprintf(`"%s/app/http/controllers/lenders"
	tenants_ctrl "%s/app/http/controllers/tenants"`, moduleName, moduleName)
		content = strings.Replace(content, old, new, 1)
	}

	// Add page controller initialization
	if !strings.Contains(content, "tenantPageController") {
		old := `applicationsPageController := applications.NewApplicationPageController()`
		new := `applicationsPageController := applications.NewApplicationPageController()
	tenantPageController := tenants_ctrl.NewTenantPageController()`
		content = strings.Replace(content, old, new, 1)
	}

	// Add tenant admin route
	if !strings.Contains(content, `"/admin/tenants"`) {
		old := `		// User management pages (super admin only)
		router.Get("/admin/users", userPageController.Index)`
		new := `		// Tenant management (super admin only)
		router.Get("/admin/tenants", tenantPageController.Index)

		// User management pages (super admin only)
		router.Get("/admin/users", userPageController.Index)`
		content = strings.Replace(content, old, new, 1)
	}

	return r.writeFile("routes/web.go", content)
}

func (r *MakeMultiTenancy) modifyRouteServiceProvider(moduleName string) error {
	content, err := r.readFile("app/providers/route_service_provider.go")
	if err != nil {
		return err
	}

	// Add middleware import
	if !strings.Contains(content, fmt.Sprintf(`"%s/app/http/middleware"`, moduleName)) {
		old := fmt.Sprintf(`	"%s/routes"`, moduleName)
		new := fmt.Sprintf(`	"%s/app/http/middleware"
	"%s/routes"`, moduleName, moduleName)
		content = strings.Replace(content, old, new, 1)
	}

	// Replace the Boot method
	old := `	// Add routes
	routes.Web()

	// API routes will be prefixed with /api
	facades.Route().Prefix("api").Group(func(apiRouter route.Router) {
		routes.Api(apiRouter)
	})`

	new := `	// Static files and inertia middleware (global, run once)
	facades.Route().Static("/images", "./public/images")
	facades.Route().Static("/css", "./public/css")
	facades.Route().Static("/js", "./public/js")
	facades.Route().GlobalMiddleware(routes.InertiaMiddleware)

	// Main routes (no tenant prefix - backward compatible)
	routes.Web(facades.Route())

	// API routes prefixed with /api
	facades.Route().Prefix("api").Group(func(apiRouter route.Router) {
		routes.Api(apiRouter)
	})

	// Tenant-prefixed routes: /t/{tenant}/...
	tenantMw := middleware.TenantFromSlug()
	facades.Route().Prefix("t/{tenant}").Middleware(tenantMw).Group(func(tenantRouter route.Router) {
		routes.Web(tenantRouter)
	})
	facades.Route().Prefix("t/{tenant}/api").Middleware(tenantMw).Group(func(tenantApiRouter route.Router) {
		routes.Api(tenantApiRouter)
	})`

	content = strings.Replace(content, old, new, 1)

	return r.writeFile("app/providers/route_service_provider.go", content)
}

// ──────────────────────────────────────────────
// Step 10: Modify GenericCrudService
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) modifyGenericCrudService(moduleName string) error {
	content, err := r.readFile("app/contracts/generic_crud_service.go")
	if err != nil {
		return err
	}

	// Add import for tenant package
	if !strings.Contains(content, "tenant") {
		old := `"smedi-sme-db/app/auth"`
		if !strings.Contains(content, old) {
			old = fmt.Sprintf(`"%s/app/auth"`, moduleName)
		}
		new := fmt.Sprintf(`"%s/app/auth"
	"%s/app/tenant"`, moduleName, moduleName)
		content = strings.Replace(content, old, new, 1)
	}

	// Add tenantAware field to struct
	old := `	// Custom query builders
	customQuery`
	new := `	// Tenant awareness
	tenantAware bool

	// Custom query builders
	customQuery`
	content = strings.Replace(content, old, new, 1)

	// Add tenant scope application in GetList after scope filtering
	// Find the block after scope filtering in GetList
	listMarker := `	// Apply search
	if req.Search != "" {`
	tenantFilter := `	// Apply tenant scope filtering
	if s.tenantAware && req.Context != nil {
		if tenantID := tenant.GetIDFromContext(req.Context); tenantID != nil {
			query = query.Where(s.tableName+".tenant_id = ?", *tenantID)
		}
	}

`
	content = strings.Replace(content, listMarker, tenantFilter+listMarker, 1)

	// Add tenant scope to count query (find the second occurrence of scope filtering in GetList)
	countMarker := `	// Apply permission-based scope filtering if enabled`
	// Count how many occurrences there are
	parts := strings.SplitN(content, countMarker, 3)
	if len(parts) >= 3 {
		// Insert tenant filter before the second occurrence's search block
		// The count query section follows the second scope filtering
		content = parts[0] + countMarker + parts[1] + countMarker + parts[2]
	}

	// Also add to count query - find the count query's search application
	countSearchMarker := `	// Apply search for count`
	if strings.Contains(content, countSearchMarker) {
		tenantCountFilter := `	// Apply tenant scope filtering for count
	if s.tenantAware && req.Context != nil {
		if tenantID := tenant.GetIDFromContext(req.Context); tenantID != nil {
			countQuery = countQuery.Where(s.tableName+".tenant_id = ?", *tenantID)
		}
	}

`
		content = strings.Replace(content, countSearchMarker, tenantCountFilter+countSearchMarker, 1)
	}

	// Add tenant scope to GetByIDWithContext
	getByIDMarker := `	// Apply scope filtering if enabled
	if s.enableScopeFiltering && ctx != nil {`
	if strings.Contains(content, getByIDMarker) {
		tenantGetByID := `	// Apply tenant scope filtering
	if s.tenantAware && ctx != nil {
		if tenantID := tenant.GetIDFromContext(ctx); tenantID != nil {
			query = query.Where(s.tableName+".tenant_id = ?", *tenantID)
		}
	}

	`
		content = strings.Replace(content, getByIDMarker, tenantGetByID+getByIDMarker, 1)
	}

	// Add EnableTenantAwareness method near the other Enable methods
	enableMethod := `
// EnableTenantAwareness enables automatic tenant_id filtering on queries
func (s *GenericCrudService[T]) EnableTenantAwareness() *GenericCrudService[T] {
	s.tenantAware = true
	return s
}
`
	// Add after EnableScopeFiltering
	scopeMethod := `func (s *GenericCrudService[T]) EnableScopeFiltering`
	idx := strings.Index(content, scopeMethod)
	if idx >= 0 {
		// Find the end of the EnableScopeFiltering method (next blank line after the closing brace)
		endIdx := strings.Index(content[idx:], "\n\n")
		if endIdx >= 0 {
			insertPos := idx + endIdx + 2
			content = content[:insertPos] + enableMethod + content[insertPos:]
		}
	}

	return r.writeFile("app/contracts/generic_crud_service.go", content)
}

// ──────────────────────────────────────────────
// Step 10b: Modify ServiceBuilder
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) modifyServiceBuilder() error {
	content, err := r.readFile("app/contracts/service_builder.go")
	if err != nil {
		return err
	}

	if strings.Contains(content, "WithTenantAwareness") {
		return nil // Already modified
	}

	// Add WithTenantAwareness method after WithScopeFiltering method
	// Find the closing brace of WithScopeFiltering by looking for the next func declaration after it
	marker := `func (b *ServiceBuilderComplete[T]) WithBeforeCreate`
	newMethod := `func (b *ServiceBuilderComplete[T]) WithTenantAwareness() *ServiceBuilderComplete[T] {
	b.builder.service.tenantAware = true
	return b
}

`
	content = strings.Replace(content, marker, newMethod+marker, 1)

	return r.writeFile("app/contracts/service_builder.go", content)
}

// ──────────────────────────────────────────────
// Step 11: Modify CrudController
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) modifyCrudController(moduleName string) error {
	content, err := r.readFile("app/contracts/crud_controller.go")
	if err != nil {
		return err
	}

	// Add import for tenant package
	if !strings.Contains(content, "/app/tenant") {
		// Try to add after an existing module import
		old := fmt.Sprintf(`"%s/app/auth"`, moduleName)
		if strings.Contains(content, old) {
			new := fmt.Sprintf(`"%s/app/auth"
	"%s/app/tenant"`, moduleName, moduleName)
			content = strings.Replace(content, old, new, 1)
		} else {
			// Add as a new import group before the framework imports
			old = `	"github.com/goravel/framework/contracts/http"`
			new := fmt.Sprintf(`	"%s/app/tenant"

	"github.com/goravel/framework/contracts/http"`, moduleName)
			content = strings.Replace(content, old, new, 1)
		}
	}

	// Add tenant_id injection in Store method before the beforeStore hook
	storeMarker := `	// Run before hook if set
	if c.beforeStore != nil {`
	if strings.Contains(content, storeMarker) {
		tenantInject := `	// Auto-inject tenant_id from request context
	if tenantID := tenant.GetIDFromContext(ctx); tenantID != nil {
		data["tenant_id"] = *tenantID
	}

`
		content = strings.Replace(content, storeMarker, tenantInject+storeMarker, 1)
	}

	// Strip tenant_id from update data in Update method
	updateMarker := `	// Run before hook if set (with mapped data)
	if c.beforeUpdate != nil {`
	if !strings.Contains(content, updateMarker) {
		updateMarker = `	// Run before hook if set
	if c.beforeUpdate != nil {`
	}
	if strings.Contains(content, updateMarker) {
		tenantStrip := `	// Prevent changing tenant_id via update
	delete(data, "tenant_id")

`
		content = strings.Replace(content, updateMarker, tenantStrip+updateMarker, 1)
	}

	return r.writeFile("app/contracts/crud_controller.go", content)
}

// ──────────────────────────────────────────────
// Step 12: Modify permission constants
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) modifyPermissionConstants() error {
	content, err := r.readFile("app/auth/permission_constants.go")
	if err != nil {
		return err
	}

	// Add ServiceTenants constant
	if strings.Contains(content, "ServiceTenants") {
		return nil // Already added
	}

	// Add after the last ServiceXxx constant
	old := `	ServiceAuthors      ServiceRegistry = "authors"`
	new := `	ServiceAuthors      ServiceRegistry = "authors"
	ServiceTenants      ServiceRegistry = "tenants"`
	content = strings.Replace(content, old, new, 1)

	// Add to GetAllServiceRegistries
	old = `		ServiceAuthors,
	}`
	new = `		ServiceAuthors,
		ServiceTenants,
	}`
	content = strings.Replace(content, old, new, 1)

	// Add display name
	old = `	case ServiceAuthors:
		return "Authors Management"`
	new = `	case ServiceAuthors:
		return "Authors Management"
	case ServiceTenants:
		return "Tenant Management"`
	content = strings.Replace(content, old, new, 1)

	// Add service actions
	old = `	default:
		return GetAllCorePermissionActions()
	}`
	new = `	case ServiceTenants:
		return []CorePermissionAction{
			PermissionCreate,
			PermissionRead,
			PermissionUpdate,
			PermissionDelete,
			PermissionView,
			PermissionManage,
		}
	default:
		return GetAllCorePermissionActions()
	}`
	content = strings.Replace(content, old, new, 1)

	return r.writeFile("app/auth/permission_constants.go", content)
}

// ──────────────────────────────────────────────
// Step 14: Modify API routes
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) modifyApiRoutes(moduleName string) error {
	content, err := r.readFile("routes/api.go")
	if err != nil {
		return err
	}

	// Add tenant controller import
	if !strings.Contains(content, "tenants") {
		old := fmt.Sprintf(`"%s/app/http/controllers/lenders"`, moduleName)
		if !strings.Contains(content, old) {
			// Try alternate import pattern
			old = `"books-database/app/http/controllers/lenders"`
		}
		new := old + fmt.Sprintf(`
	tenants_api "%s/app/http/controllers/tenants"`, moduleName)
		content = strings.Replace(content, old, new, 1)
	}

	// Add tenant controller initialization
	if !strings.Contains(content, "tenantController") {
		// Find a controller initialization to add after
		marker := `lenderController := lenders.NewLenderController()`
		if strings.Contains(content, marker) {
			content = strings.Replace(content, marker,
				marker+"\n\ttenantController := tenants_api.NewTenantController()", 1)
		}
	}

	// Add tenant routes in the protected section
	if !strings.Contains(content, `"tenants"`) || !strings.Contains(content, "tenantController.Index") {
		// Find the end of the protected routes group to add tenant routes
		// Look for a good insertion point inside the protected group
		marker := `		// Users management (super admin only)`
		if strings.Contains(content, marker) {
			tenantRoutes := `		// Tenant management (super admin only)
		protectedRouter.Prefix("tenants").Group(func(tenantRouter route.Router) {
			tenantRouter.Get("/", tenantController.Index)
			tenantRouter.Get("/search", tenantController.Search)
			tenantRouter.Get("/filters", tenantController.FilterMetadata)
			tenantRouter.Get("/{id}", tenantController.Show)
			tenantRouter.Post("/", tenantController.Store)
			tenantRouter.Put("/{id}", tenantController.Update)
			tenantRouter.Delete("/{id}", tenantController.Delete)
			tenantRouter.Post("/{id}/users", tenantController.AssignUser)
			tenantRouter.Delete("/{id}/users/{userId}", tenantController.RemoveUser)
			tenantRouter.Get("/{id}/users", tenantController.GetUsers)
		})

`
			content = strings.Replace(content, marker, tenantRoutes+marker, 1)
		}
	}

	return r.writeFile("routes/api.go", content)
}

// ──────────────────────────────────────────────
// Step 15: Enable tenant awareness on services
// ──────────────────────────────────────────────

func (r *MakeMultiTenancy) enableTenantAwareness(filePath string) error {
	content, err := r.readFile(filePath)
	if err != nil {
		return err
	}

	// Check if already has tenant awareness
	if strings.Contains(content, "WithTenantAwareness") {
		return nil
	}

	// Find Build() — it may appear as ".Build()" on the same line or
	// "\t\tBuild()" on its own line in chained method calls
	if strings.Contains(content, "\t\tBuild()") {
		content = strings.Replace(content, "\t\tBuild()", "\t\tWithTenantAwareness().\n\t\tBuild()", 1)
	} else if strings.Contains(content, ".Build()") {
		content = strings.Replace(content, ".Build()", ".WithTenantAwareness().\n\t\tBuild()", 1)
	} else {
		return fmt.Errorf("Build() not found in %s", filePath)
	}

	return r.writeFile(filePath, content)
}
