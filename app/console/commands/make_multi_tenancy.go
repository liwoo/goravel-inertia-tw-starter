package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type MakeMultiTenancy struct{}

func (r *MakeMultiTenancy) Signature() string {
	return "make:multi-tenancy"
}

func (r *MakeMultiTenancy) Description() string {
	return "Scaffold multi-tenancy support into the application (non-reversible)"
}

func (r *MakeMultiTenancy) Extend() command.Extend {
	return command.Extend{
		Category: "make",
	}
}

func (r *MakeMultiTenancy) Handle(ctx console.Context) error {
	ctx.Info("┌─────────────────────────────────────────────────┐")
	ctx.Info("│  Multi-Tenancy Scaffolding                      │")
	ctx.Info("└─────────────────────────────────────────────────┘")
	ctx.NewLine()
	ctx.Warning("This command will permanently modify your application to support multi-tenancy.")
	ctx.Warning("This is a non-reversible operation. Make sure you have committed your current changes.")
	ctx.NewLine()

	// Check if already scaffolded
	if _, err := os.Stat("app/tenant/context.go"); err == nil {
		ctx.Error("Multi-tenancy appears to already be scaffolded (app/tenant/context.go exists).")
		return fmt.Errorf("multi-tenancy already scaffolded")
	}

	// Detect module name from go.mod
	moduleName, err := r.detectModuleName()
	if err != nil {
		return fmt.Errorf("failed to detect module name from go.mod: %w", err)
	}
	ctx.Info(fmt.Sprintf("Detected module: %s", moduleName))
	ctx.NewLine()

	timestamp := time.Now().Format("20060102150405")

	var created []string
	var modified []string

	// Step 1: Config
	ctx.Info("[1/15] Creating tenancy config...")
	if err := r.createFile("config/tenancy.go", mtTplConfig(moduleName)); err != nil {
		return r.fail(ctx, "config/tenancy.go", err)
	}
	created = append(created, "config/tenancy.go")

	// Step 2: Tenant context package
	ctx.Info("[2/15] Creating tenant context package...")
	if err := os.MkdirAll("app/tenant", 0755); err != nil {
		return err
	}
	if err := r.createFile("app/tenant/context.go", mtTplTenantContext(moduleName)); err != nil {
		return r.fail(ctx, "app/tenant/context.go", err)
	}
	created = append(created, "app/tenant/context.go")

	// Step 3: Models
	ctx.Info("[3/15] Creating Tenant and UserTenant models...")
	if err := r.createFile("app/models/tenant.go", mtTplTenantModel()); err != nil {
		return r.fail(ctx, "app/models/tenant.go", err)
	}
	created = append(created, "app/models/tenant.go")
	if err := r.createFile("app/models/user_tenant.go", mtTplUserTenantModel()); err != nil {
		return r.fail(ctx, "app/models/user_tenant.go", err)
	}
	created = append(created, "app/models/user_tenant.go")

	// Step 4: Modify BaseAuditableModel
	ctx.Info("[4/15] Adding TenantID to BaseAuditableModel...")
	if err := r.modifyBaseAuditableModel(); err != nil {
		return r.fail(ctx, "app/models/base_auditable.go", err)
	}
	modified = append(modified, "app/models/base_auditable.go")

	// Step 5: Migrations
	ctx.Info("[5/15] Creating migrations...")
	migFiles := []struct {
		suffix  string
		content string
	}{
		{"_create_tenants_table", mtTplMigrationTenants(timestamp + "01")},
		{"_create_user_tenants_table", mtTplMigrationUserTenants(timestamp + "02")},
		{"_add_tenant_id_to_all_tables", mtTplMigrationAddTenantID(timestamp + "03")},
	}
	for i, mig := range migFiles {
		ts := fmt.Sprintf("%s0%d", timestamp, i+1)
		filename := fmt.Sprintf("database/migrations/%s%s.go", ts, mig.suffix)
		if err := r.createFile(filename, mig.content); err != nil {
			return r.fail(ctx, filename, err)
		}
		created = append(created, filename)
	}

	// Register migrations in database/kernel.go
	if err := r.modifyDatabaseKernel(timestamp); err != nil {
		return r.fail(ctx, "database/kernel.go", err)
	}
	modified = append(modified, "database/kernel.go")

	// Step 6: Seeder
	ctx.Info("[6/15] Creating tenant seeder...")
	if err := r.createFile("database/seeders/tenant_seeder.go", mtTplTenantSeeder(moduleName)); err != nil {
		return r.fail(ctx, "database/seeders/tenant_seeder.go", err)
	}
	created = append(created, "database/seeders/tenant_seeder.go")
	if err := r.modifyDatabaseSeeder(); err != nil {
		return r.fail(ctx, "database/seeders/database_seeder.go", err)
	}
	modified = append(modified, "database/seeders/database_seeder.go")

	// Step 7: Tenant middleware
	ctx.Info("[7/15] Creating tenant middleware...")
	if err := r.createFile("app/http/middleware/tenant.go", mtTplTenantMiddleware(moduleName)); err != nil {
		return r.fail(ctx, "app/http/middleware/tenant.go", err)
	}
	created = append(created, "app/http/middleware/tenant.go")

	// Step 8: Modify JwtAuth middleware
	ctx.Info("[8/15] Modifying JWT auth middleware for tenant verification...")
	if err := r.modifyJwtAuthMiddleware(moduleName); err != nil {
		ctx.Warning(fmt.Sprintf("Could not auto-modify jwt_auth.go: %v", err))
		ctx.Warning("You will need to manually add tenant membership verification to JwtAuth middleware.")
	} else {
		modified = append(modified, "app/http/middleware/jwt_auth.go")
	}

	// Step 9: Refactor routes
	ctx.Info("[9/15] Refactoring routes for tenant support...")
	if err := r.modifyWebRoutes(moduleName); err != nil {
		return r.fail(ctx, "routes/web.go", err)
	}
	modified = append(modified, "routes/web.go")
	if err := r.modifyRouteServiceProvider(moduleName); err != nil {
		return r.fail(ctx, "app/providers/route_service_provider.go", err)
	}
	modified = append(modified, "app/providers/route_service_provider.go")

	// Step 10: Modify GenericCrudService
	ctx.Info("[10/15] Adding tenant filtering to GenericCrudService...")
	if err := r.modifyGenericCrudService(moduleName); err != nil {
		return r.fail(ctx, "app/contracts/generic_crud_service.go", err)
	}
	modified = append(modified, "app/contracts/generic_crud_service.go")

	// Step 10b: Add WithTenantAwareness to ServiceBuilder
	if err := r.modifyServiceBuilder(); err != nil {
		return r.fail(ctx, "app/contracts/service_builder.go", err)
	}
	modified = append(modified, "app/contracts/service_builder.go")

	// Step 11: Modify CrudController
	ctx.Info("[11/15] Adding tenant injection to CrudController...")
	if err := r.modifyCrudController(moduleName); err != nil {
		return r.fail(ctx, "app/contracts/crud_controller.go", err)
	}
	modified = append(modified, "app/contracts/crud_controller.go")

	// Step 12: Permission constants
	ctx.Info("[12/15] Adding tenant permission constants...")
	if err := r.modifyPermissionConstants(); err != nil {
		return r.fail(ctx, "app/auth/permission_constants.go", err)
	}
	modified = append(modified, "app/auth/permission_constants.go")

	// Step 13: Tenant service & controller
	ctx.Info("[13/15] Creating tenant service, controller, and requests...")
	if err := r.createFile("app/services/tenant_service.go", mtTplTenantService(moduleName)); err != nil {
		return r.fail(ctx, "app/services/tenant_service.go", err)
	}
	created = append(created, "app/services/tenant_service.go")
	if err := os.MkdirAll("app/http/controllers/tenants", 0755); err != nil {
		return err
	}
	if err := r.createFile("app/http/controllers/tenants/tenant_controller.go", mtTplTenantController(moduleName)); err != nil {
		return r.fail(ctx, "app/http/controllers/tenants/tenant_controller.go", err)
	}
	created = append(created, "app/http/controllers/tenants/tenant_controller.go")
	if err := r.createFile("app/http/controllers/tenants/tenant_page_controller.go", mtTplTenantPageController(moduleName)); err != nil {
		return r.fail(ctx, "app/http/controllers/tenants/tenant_page_controller.go", err)
	}
	created = append(created, "app/http/controllers/tenants/tenant_page_controller.go")
	if err := r.createFile("app/http/requests/tenant_create_request.go", mtTplTenantCreateRequest()); err != nil {
		return r.fail(ctx, "app/http/requests/tenant_create_request.go", err)
	}
	created = append(created, "app/http/requests/tenant_create_request.go")
	if err := r.createFile("app/http/requests/tenant_update_request.go", mtTplTenantUpdateRequest()); err != nil {
		return r.fail(ctx, "app/http/requests/tenant_update_request.go", err)
	}
	created = append(created, "app/http/requests/tenant_update_request.go")

	// Step 14: Register tenant routes
	if err := r.modifyApiRoutes(moduleName); err != nil {
		return r.fail(ctx, "routes/api.go", err)
	}
	modified = append(modified, "routes/api.go")

	// Step 15: Enable tenant awareness on existing services
	ctx.Info("[14/15] Enabling tenant awareness on existing services...")
	svcFiles := []string{
		"app/services/book_service.go",
		"app/services/lender_service.go",
		"app/services/application_service.go",
		"app/services/config_service.go",
	}
	for _, svcFile := range svcFiles {
		if err := r.enableTenantAwareness(svcFile); err != nil {
			ctx.Warning(fmt.Sprintf("Could not modify %s: %v", svcFile, err))
		} else {
			modified = append(modified, svcFile)
		}
	}

	// Step 16: .env update
	ctx.Info("[15/15] Updating .env files...")
	r.appendToEnv(".env", "\n# Multi-tenancy\nDEFAULT_TENANT_SLUG=main\n")
	r.appendToEnv(".env.example", "\n# Multi-tenancy\nDEFAULT_TENANT_SLUG=main\n")

	// Print summary
	ctx.NewLine()
	ctx.Info("┌─────────────────────────────────────────────────┐")
	ctx.Info("│  Multi-tenancy scaffolded successfully!         │")
	ctx.Info("└─────────────────────────────────────────────────┘")
	ctx.NewLine()

	ctx.Info("Created files:")
	for _, f := range created {
		ctx.Success(fmt.Sprintf("  + %s", f))
	}
	ctx.NewLine()

	ctx.Info("Modified files:")
	for _, f := range modified {
		ctx.Warning(fmt.Sprintf("  ~ %s", f))
	}
	ctx.NewLine()

	ctx.Info("Next steps:")
	ctx.Line("  1. Run migrations:     go run . artisan migrate")
	ctx.Line("  2. Run seeders:        go run . artisan db:seed")
	ctx.Line("  3. Build frontend:     pnpm install && pnpm build")
	ctx.Line("  4. Add frontend tenant context (see resources/js/contexts/)")
	ctx.Line("  5. Review generated code and customize as needed")
	ctx.NewLine()

	return nil
}

// Helper methods

func (r *MakeMultiTenancy) detectModuleName() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("module name not found in go.mod")
}

func (r *MakeMultiTenancy) createFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func (r *MakeMultiTenancy) fail(ctx console.Context, file string, err error) error {
	ctx.Error(fmt.Sprintf("Failed to process %s: %v", file, err))
	return err
}

func (r *MakeMultiTenancy) readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *MakeMultiTenancy) writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func (r *MakeMultiTenancy) appendToEnv(path, content string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(content)
}

// insertAfter inserts newContent after the first occurrence of marker in the file
func (r *MakeMultiTenancy) insertAfter(filePath, marker, newContent string) error {
	content, err := r.readFile(filePath)
	if err != nil {
		return err
	}
	if !strings.Contains(content, marker) {
		return fmt.Errorf("marker not found in %s: %s", filePath, marker)
	}
	content = strings.Replace(content, marker, marker+newContent, 1)
	return r.writeFile(filePath, content)
}

// insertBefore inserts newContent before the first occurrence of marker in the file
func (r *MakeMultiTenancy) insertBefore(filePath, marker, newContent string) error {
	content, err := r.readFile(filePath)
	if err != nil {
		return err
	}
	if !strings.Contains(content, marker) {
		return fmt.Errorf("marker not found in %s: %s", filePath, marker)
	}
	content = strings.Replace(content, marker, newContent+marker, 1)
	return r.writeFile(filePath, content)
}

// replaceInFile replaces old with new in the file
func (r *MakeMultiTenancy) replaceInFile(filePath, old, new string) error {
	content, err := r.readFile(filePath)
	if err != nil {
		return err
	}
	if !strings.Contains(content, old) {
		return fmt.Errorf("text not found in %s", filePath)
	}
	content = strings.Replace(content, old, new, 1)
	return r.writeFile(filePath, content)
}
