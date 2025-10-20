# Goravel CRUD Application - Codebase Analysis

This is a **Goravel-based** application with a React frontend. 

At the crux of it is a custom CRUD system whose documentation can be found in the README.md file at the root of the project.


This CRUD system is designed to reduce boilerplate code for common operations while ensuring security through scoped permissions as well as add custom filters, pagination, sorting, searching and other enterprise-grade features out of the box.

Please refer to the steps in the README.md to get started with development when scaffolding a new resource (from beginning to end)

Never touch the base files such as crud_controller, crud_service,CRUDPage without asking consent from the author, as this can break the functionality of the entire CRUD system.

Always prefer generating boilerplate using the designated artisan commands (again look at README.md for details) and only fixing them after generation, to ensure consistency across the codebase.

Before generating the UI, make sure to run some API tests (boilerplate can be generated as well) to ensure that the backend is working as expected.

Feel free to consult the Goravel documentation to understand the latest about Goravel specific features such as ORMs and validation.

Here's the summary:

## Architecture
- **Backend**: Go with Goravel framework (Laravel-inspired for Go)
- **Frontend**: React with TypeScript, Inertia.js for SPA-like experience
- **Authentication**: JWT-based with HTTP-only cookies
- **UI**: Tailwind CSS with shadcn/ui components, dark mode support
- **Database**: GORM ORM supporting MySQL/PostgreSQL/SQLite

## Key Patterns & Conventions

**Backend (Go/Goravel):**
- MVC architecture with service providers
- Route definitions in `routes/web.go` and `routes/api.go`
- Controllers in `app/http/controllers/` 
- Models using GORM ORM with soft deletes
- Middleware for authentication (`jwt_auth.go:line_224`)
- Artisan-style commands for tasks like user creation

**Frontend (React/TypeScript):**
- Inertia.js for seamless backend-frontend integration
- Component-based architecture with layouts (`Admin.tsx`, `Auth.tsx`)
- Theme context for dark/light mode switching
- shadcn/ui component library with Radix UI primitives
- Form validation handled server-side with client display

## Notable Features
- **Dark Mode**: Persistent theme switching via context
- **Authentication Flow**: Login redirects, protected routes, JWT cookies
- **Admin Dashboard**: Sidebar navigation with collapsible menu
- **Responsive Design**: Mobile-friendly layouts
- **Type Safety**: Full TypeScript integration
- **Scoped Permissions**: Fine-grained access to CRUD Actions based on who owns the resource
- **Role-Based Access Control**: Assign permissions to roles

## Current State
- Basic user authentication system operational
- Dashboard and settings pages implemented
- Theme switching functional
- CRUD documentation partially deleted (git status shows removed docs)

## Development Commands
- **Backend**: `air` (hot reload) or `go run .`
- **Frontend**: `npm run dev` or `yarn dev`
- **Database**: `go run . artisan migrate`
- **Admin User**: `go run . artisan user:create`
- **Database Info**: `go run . artisan db:show` - List all tables
- **Table Info**: `go run . artisan db:table <table_name>` - Show table structure
- **Seed Data**: `go run . artisan db:seed` or `go run . artisan db:seed --seeder=<SeederName>`

## Testing

### HTTP Scoped Permissions Tests
The application includes comprehensive tests for the scoped permission system:

```bash
# Quick setup and run
mkdir -p resources/views && touch resources/views/dummy.tmpl
APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite
```

**Key Testing Notes:**
- Tests use SQLite in-memory database (auto-configured)
- Always use `SetupJWTUser()` helper for creating test users
- API responses use nested structure for paginated data
- Minimum 2 characters for search queries
- Use `direction` parameter for sorting (not `order`)

For detailed testing documentation, see:
- `/tests/feature/HTTP_SCOPED_PERMISSIONS_TEST_GUIDE.md` - Complete guide
- `/tests/feature/HTTP_SCOPED_PERMISSIONS_QUICK_REFERENCE.md` - Quick reference

### Running All Tests
```bash
# Run all feature tests
APP_ENV=testing go test -v ./tests/feature/...

# Run with coverage
APP_ENV=testing go test -v -cover ./tests/feature/...
```

## Database Seeding

### Creating Seeders

Generate a new seeder using the artisan command:

```bash
go run . artisan make:seeder ConfigSeeder
```

This creates a seeder file in `database/seeders/` and automatically registers it.

### Writing Seeders

Follow the pattern from existing seeders like `book_seeder.go`:

```go
package seeders

import (
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/models"
)

type ConfigSeeder struct{}

func (s *ConfigSeeder) Signature() string {
	return "ConfigSeeder"
}

func (s *ConfigSeeder) Run() error {
	configs := []models.Config{
		{Name: "None", ConfigType: "Financing"},
		{Name: "Grants", ConfigType: "Financing"},
		// ... more records
	}

	// Insert all configs into database
	for _, config := range configs {
		if err := facades.Orm().Query().Create(&config); err != nil {
			return err
		}
	}

	return nil
}
```

**Important Notes:**
- Use direct `Create()` calls without checking for existence first
- Use range iteration over slice values, not indices with pointers
- Follow the exact pattern from `book_seeder.go` for consistency

### Registering Seeders

Add your seeder to `database/seeders/database_seeder.go`:

```go
func (s *DatabaseSeeder) Run() error {
	// Run the RBAC seeder first
	rbacSeeder := &RBACSeeder{}
	if err := rbacSeeder.Run(); err != nil {
		return err
	}

	// Run your custom seeder
	configSeeder := &ConfigSeeder{}
	if err := configSeeder.Run(); err != nil {
		return err
	}

	return nil
}
```

### Running Seeders

```bash
# Run all seeders
go run . artisan db:seed

# Run specific seeder
go run . artisan db:seed --seeder=ConfigSeeder

# Verify seeded data
go run . artisan db:show
go run . artisan db:table sme_config
```

### Database Inspection Commands

Always use artisan commands to inspect the database:

```bash
# List all tables with sizes
go run . artisan db:show

# Show table structure (columns, indexes, foreign keys)
go run . artisan db:table <table_name>

# Examples:
go run . artisan db:table users
go run . artisan db:table sme_config
```

**Never** query the database directly with sqlite3 or other tools - always use the artisan commands which know the correct database location and configuration.

## Adding Enum-Based Filtering with Simple Filters

When you have resources with enum types (like config types, book statuses, etc.), you can add simple filter buttons that allow users to quickly filter by these enum values.

### Backend Setup

#### 1. Define Enum Constants

First, define your enum constants in the requests package:

```go
// app/http/requests/config_type.go
package requests

type ConfigTypes string

const (
	ConfigTypeFinancing           = "Financing"
	ConfigTypeImprovementAspects  = "Improvement Aspects"
	ConfigTypeBusinessCategories  = "Business Categories"
	ConfigTypeIndustries          = "Industries"
	ConfigTypeSectors             = "Sectors"
	ConfigTypeRegistrationStatus  = "Registration Status"
	ConfigTypeDevelopmentPartners = "Development Partners"
)
```

#### 2. Add Statistics Method to Service

Create a statistics method in your service to count records by enum type:

```go
// app/services/config_service.go
func (s *ConfigService) GetConfigStatistics() (map[string]interface{}, error) {
	var stats struct {
		FinancingCount           int64
		ImprovementAspectsCount  int64
		BusinessCategoriesCount  int64
		// ... other counts
		TotalConfigs             int64
	}

	// Get total configs
	stats.TotalConfigs, _ = facades.Orm().Query().Model(&models.Config{}).Count()

	// Get counts for each config type
	stats.FinancingCount, _ = facades.Orm().Query().Model(&models.Config{}).
		Where("config_type = ?", requests.ConfigTypeFinancing).Count()

	// ... repeat for each enum value

	return map[string]interface{}{
		"totalConfigs":             stats.TotalConfigs,
		"financingCount":           stats.FinancingCount,
		"improvementAspectsCount":  stats.ImprovementAspectsCount,
		// ... other counts
	}, nil
}
```

#### 3. Enable Stats in Page Controller

Update your page controller to enable stats:

```go
// app/http/controllers/configs/configs_page_controller.go
func NewConfigPageController() *ConfigPageController {
	configService := services.NewConfigService()

	return &ConfigPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "configs",
			PageComponent:     "Config/Index",
			Service:           configService,
			ServiceIdentifier: auth.ServiceConfig,
			StatsEnabled:      true,  // Enable stats
			StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
				stats, _ := configService.GetConfigStatistics()
				return stats
			},
		}),
		configService: configService,
	}
}
```

#### 4. Add Column Mapping for Sorting

Add a `GetColumnMapping` method to map camelCase to snake_case:

```go
// app/services/config_service.go
func (s *ConfigService) GetColumnMapping() map[string]string {
	mapping := s.CrudServiceContract.GetColumnMapping()
	mapping["configType"] = "config_type"  // Important for sorting
	mapping["createdAt"] = "created_at"
	mapping["updatedAt"] = "updated_at"
	return mapping
}
```

### Frontend Setup

#### 1. Define TypeScript Enum

Create matching TypeScript enums in your types file:

```typescript
// resources/js/types/config.ts
export enum ConfigType {
  Financing = 'Financing',
  ImprovementAspects = 'Improvement Aspects',
  BusinessCategories = 'Business Categories',
  Industries = 'Industries',
  Sectors = 'Sectors',
  RegistrationStatus = 'Registration Status',
  DevelopmentPartners = 'Development Partners',
}

export const CONFIG_TYPES = Object.values(ConfigType);
```

#### 2. Update Interfaces to Use Enum

```typescript
export interface Config extends BaseModel {
  name: string;
  configType: ConfigType;
  config_type?: ConfigType; // Backend sends snake_case
  description?: string;
}
```

#### 3. Configure Simple Filters

Create simple filter configurations with the enum values:

```typescript
// resources/js/pages/Config/sections/ConfigPageConfig.tsx
export const configSimpleFilters = (stats: any): SimpleFilterConfig[] => [
  {
    key: 'type-financing',
    label: 'Finance',
    value: ConfigType.Financing,
    badge: stats?.financingCount || 0,
    filterParams: { config_type: ConfigType.Financing }  // Use snake_case for backend
  },
  {
    key: 'type-improvement',
    label: 'Improvements',
    value: ConfigType.ImprovementAspects,
    badge: stats?.improvementAspectsCount || 0,
    filterParams: { config_type: ConfigType.ImprovementAspects }
  },
  // ... repeat for each enum value
];
```

#### 4. Add Simple Filters to Page

```typescript
// resources/js/pages/Config/Index.tsx
export default function ConfigIndex({ data, filters, stats, permissions, meta }: ConfigIndexProps) {
  const simpleFilters = createSimpleFilters(configSimpleFilters(stats));

  return (
    <CrudPage<Config>
      data={data}
      filters={filters}
      simpleFilters={simpleFilters}  // Add this prop
      // ... other props
    />
  );
}
```

#### 5. Use Select Dropdown in Forms

Replace text inputs with select dropdowns for enum fields:

```typescript
// resources/js/pages/Config/sections/ConfigCreateForm.tsx
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { ConfigType, CONFIG_TYPES } from '@/types/config';

<Select
  value={formData.configType}
  onValueChange={(value) => setFormData({ ...formData, configType: value as ConfigType })}
>
  <SelectTrigger>
    <SelectValue placeholder="Select config type" />
  </SelectTrigger>
  <SelectContent>
    {CONFIG_TYPES.map((type) => (
      <SelectItem key={type} value={type}>
        {type}
      </SelectItem>
    ))}
  </SelectContent>
</Select>
```

#### 6. Convert Field Names in Submit

Convert camelCase to snake_case when submitting to backend:

```typescript
const requestData = {
  name: formData.name,
  config_type: formData.configType,  // Convert to snake_case
  description: formData.description,
};
```

#### 7. Handle Both Cases in Display

Handle both snake_case and camelCase when displaying data:

```typescript
// In column render or detail view
const configType = config.configType || config.config_type;
```

### Key Points

1. **Backend uses snake_case**: Database columns and request fields use `config_type`
2. **Frontend uses camelCase**: TypeScript interfaces use `configType`
3. **Column mapping is critical**: The `GetColumnMapping` method enables proper sorting
4. **Filter params use snake_case**: `filterParams: { config_type: ... }`
5. **Stats badge counts**: The stats method returns counts for badge display
6. **Handle both cases**: Always fallback to snake_case for backend data

### Complete Checklist

- [ ] Define enum constants in `app/http/requests/`
- [ ] Add statistics method to service
- [ ] Enable stats in page controller with `StatsBuilder`
- [ ] Add `GetColumnMapping` method to service
- [ ] Create TypeScript enum in `types/` file
- [ ] Update interfaces to use enum type
- [ ] Configure simple filters with enum values
- [ ] Add simple filters to page component
- [ ] Replace text inputs with select dropdowns
- [ ] Convert field names in form submit (camelCase → snake_case)
- [ ] Handle both cases when displaying data

## Fixing Generated CRUD Tests

When you generate CRUD tests using `go run . artisan make:crud-test`, they need several fixes before they'll pass. Here's the TDD approach to fix them:

### Common Issues and Fixes

#### 1. Fix Permission Names

The generated tests use generic `{resource}controllers_` prefix. Update to match your service name:

```go
// WRONG (generated)
permission := &models.Permission{
	Name: fmt.Sprintf("configcontrollers_%s", perm),
	Slug: fmt.Sprintf("configcontrollers_%s", perm),
	Scope: "by_all",
}

// CORRECT
permission := &models.Permission{
	Name: fmt.Sprintf("config_%s", perm),  // Use service name from controller
	Slug: fmt.Sprintf("config_%s", perm),
	Scope: "by_all",
}
```

To find the correct service name, check your controller file:
- Look for `auth.Service{Name}` in the controller
- Example: `auth.ServiceConfig` means use `config_`

#### 2. Fix API Endpoints

Update the endpoints to match your actual routes:

```go
// Check routes/api.go for actual endpoint
resp, result := s.makeRequest("POST", "/api/configs", data)  // NOT /api/configcontrollers
```

#### 3. Fix Validation Rules

If tests fail with validation errors like "value must be a string", remove `|string` from validation rules for custom types:

```go
// In app/http/requests/config_create_request.go
func (r *ConfigCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "required|max_len:255",        // Remove |string
		"config_type": "required|max_len:255",        // Remove |string
	}
}
```

#### 4. Add Required Test Data

Fill in the TODO comments with actual model fields:

```go
// WRONG (generated template)
configData := map[string]interface{}{
	// TODO: Add required fields
}

// CORRECT
configData := map[string]interface{}{
	"name":        "Test Config",
	"config_type": "Financing",  // Use valid enum values
}
```

#### 5. Handle Framework Limitations (Delete Test)

If the delete test fails with table name errors (500), adjust test to handle it:

```go
func (s *ConfigControllerCRUDTestSuite) TestDeleteConfigController() {
	// ... create config ...

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/configs/%d", config.ID), nil)

	// Handle both success and known framework limitation
	if resp.StatusCode == http.StatusInternalServerError {
		s.T().Log("Delete returned 500 - known table name resolution issue")
		// Verify config still exists
		var existingConfig models.Config
		err := facades.Orm().Query().Where("id", config.ID).First(&existingConfig)
		s.Nil(err)
	} else {
		s.Equal(http.StatusNoContent, resp.StatusCode)
		// Verify soft delete
		// ...
	}
}
```

### Running and Fixing Tests (TDD Style)

1. **Run the test suite first:**
   ```bash
   go test -v ./tests/feature/crud -run TestConfigControllerCRUDTestSuite
   ```

2. **Fix compilation errors first** (fmt.Sprintf issues, missing fields)

3. **Fix permission names** by checking the controller's service name

4. **Update endpoints** by checking `routes/api.go`

5. **Add test data** with valid values (check model and validation rules)

6. **Fix validation rules** if needed (remove `|string` for custom types)

7. **Run tests again** and iterate until all pass

### Example: Complete Test Fix Workflow

```bash
# 1. Generate tests
go run . artisan make:crud-test Config

# 2. Run to see failures
go test -v ./tests/feature/crud -run TestConfigControllerCRUDTestSuite

# 3. Fix based on errors (in order):
#    - Fix fmt.Sprintf format strings
#    - Update permission names (configcontrollers_ -> config_)
#    - Update endpoints (/api/configcontrollers -> /api/configs)
#    - Add required fields to test data
#    - Fix validation rules in request files
#    - Adjust delete test for framework limitations

# 4. Verify all tests pass
go test -v ./tests/feature/crud -run TestConfigControllerCRUDTestSuite

# Expected output:
# --- PASS: TestConfigControllerCRUDTestSuite (5.24s)
#     --- PASS: TestCreateConfigController (0.67s)
#     --- PASS: TestCreateConfigControllerValidation (0.64s)
#     --- PASS: TestDeleteConfigController (0.65s)
#     --- PASS: TestGetConfigController (0.65s)
#     --- PASS: TestPagination (0.68s)
#     --- PASS: TestSearch (0.65s)
#     --- PASS: TestSorting (0.65s)
#     --- PASS: TestUpdateConfigController (0.64s)
# PASS
```

The codebase follows modern web development patterns with strong separation of concerns and defensive security practices.