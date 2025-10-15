# Goravel Blog Application

[![CI](https://github.com/liwoo/goravel-inertia-tw-starter/actions/workflows/ci.yml/badge.svg)](https://github.com/liwoo/goravel-inertia-tw-starter/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-ready-blue.svg)](https://hub.docker.com/r/liwoo/goravel-inertia-tw-starter)
[![codecov](https://codecov.io/gh/liwoo/goravel-inertia-tw-starter/branch/main/graph/badge.svg)](https://codecov.io/gh/liwoo/goravel-inertia-tw-starter)

A modern web application built with Goravel (Go) and React, featuring JWT authentication, RBAC permissions, CRUD generators, dark mode, and a responsive UI.

## 📊 Build Status

| Branch | Status | Coverage |
|--------|--------|----------|
| main   | [![CI](https://github.com/liwoo/goravel-inertia-tw-starter/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/liwoo/goravel-inertia-tw-starter/actions) | [![codecov](https://codecov.io/gh/liwoo/goravel-inertia-tw-starter/branch/main/graph/badge.svg)](https://codecov.io/gh/liwoo/goravel-inertia-tw-starter) |
| develop | [![CI](https://github.com/liwoo/goravel-inertia-tw-starter/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/liwoo/goravel-inertia-tw-starter/actions) | [![codecov](https://codecov.io/gh/liwoo/goravel-inertia-tw-starter/branch/develop/graph/badge.svg)](https://codecov.io/gh/liwoo/goravel-inertia-tw-starter) |

## 🚀 Quick Start

```bash
# 1. Clone and install dependencies
git clone <repository-url>
cd blog
go mod download
npm install

# 2. Setup environment
cp .env.example .env
# Edit .env with your database credentials

# 3. Initialize database
go run . artisan migrate
go run . artisan key:generate
go run . artisan seed  # Seeds RBAC permissions and sample data

# 4. Create admin user
go run . artisan user:create-admin

# 5. Run the application (need 2 terminals)
# Terminal 1: Backend
air  # or: go run .

# Terminal 2: Frontend
npm run dev

# 6. Open http://localhost:3500
```

## 📋 Prerequisites

- **Go 1.18+** - [Download Go](https://golang.org/dl/)
- **Node.js 16+** - [Download Node.js](https://nodejs.org/)
- **Database** - One of:
  - SQLite (default, no setup needed)
  - MySQL 5.7+
  - PostgreSQL 12+
- **Air** (optional, for hot reload):
  ```bash
  go install github.com/cosmtrek/air@latest
  ```

## 🧪 Testing

### Quick Test Run

```bash
# Run all tests with proper environment setup
APP_ENV=testing go test ./tests/... -v

# Or use the convenience script
./run_tests.sh

# Run specific test suites
APP_ENV=testing go test ./tests/unit -v
APP_ENV=testing go test ./tests/integration -v
APP_ENV=testing go test ./tests/feature -v
```

## 🔍 Code Quality

### Linting

The project uses `golangci-lint` for code quality checks:

```bash
# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2

# Run linter
golangci-lint run

# Auto-fix issues
golangci-lint run --fix
```

### Pre-commit Hooks

Enable automatic code quality checks before each commit:

```bash
# Option 1: Use git hooks directory
./enable-hooks.sh

# Option 2: Use pre-commit framework
./setup-hooks.sh
```

The pre-commit hook will:
- Format Go code with `gofmt`
- Run `go mod tidy`
- Run `golangci-lint`
- Check for large files (>5MB)
- Prevent commits with linting errors

To skip hooks for a single commit:
```bash
git commit --no-verify
```

### Test Categories

| Category | Description | Path |
|----------|-------------|------|
| Unit | Business logic, isolated components | `tests/unit/` |
| Integration | Service layer, API endpoints | `tests/integration/` |
| Feature | Full user workflows, UI interactions | `tests/feature/` |

### Writing Tests

```go
// Example test structure
package tests

import (
    "testing"
    "github.com/stretchr/testify/suite"
    "players/tests"
)

type YourTestSuite struct {
    suite.Suite
    tests.TestCase  // Provides database helpers
}

func TestYourTestSuite(t *testing.T) {
    suite.Run(t, new(YourTestSuite))
}

func (s *YourTestSuite) SetupTest() {
    s.RefreshDatabase()  // Clean database for each test
}

func (s *YourTestSuite) TestExample() {
    // Your test logic
    s.Equal(expected, actual)
}
```

### Test Helpers

- **Authentication**: Use `helpers.SetupJWTUser()` for creating test users
- **Database**: Tests automatically use SQLite in-memory database
- **HTTP Testing**: Use `httptest.NewServer(facades.Route())` for API tests

### Common Test Commands

```bash
# Run with coverage
APP_ENV=testing go test ./tests/... -v -cover

# Run specific test
APP_ENV=testing go test -v ./tests/feature -run TestBookCRUD

# Run tests in watch mode (requires entr)
find . -name "*.go" | entr -c go test ./tests/... -v

# Clean test artifacts
rm -rf tests/*/database/
rm -rf tests/*/storage/
```

## 🔧 Detailed Setup

### First-Time Developer Setup

```bash
# 1. Install Go (if not installed)
# macOS: brew install go
# Ubuntu: sudo apt install golang-go
# Windows: Download from https://golang.org/dl/

# 2. Install Node.js (if not installed)
# macOS: brew install node
# Ubuntu: curl -fsSL https://deb.nodesource.com/setup_16.x | sudo -E bash - && sudo apt install nodejs
# Windows: Download from https://nodejs.org/

# 3. Install Air for hot reload (recommended)
go install github.com/cosmtrek/air@latest

# 4. Verify installations
go version      # Should show 1.18+
node --version  # Should show 16+
air -v          # Should show air version
```

### 1. Initial Setup

```bash
# Clone the repository
git clone <repository-url>
cd blog

# Install Go dependencies
go mod download
go mod tidy  # Clean up any issues

# Install frontend dependencies
npm install  # or: yarn install

# Setup environment
cp .env.example .env
```

### 2. Configure Database

For quickest setup, use SQLite (default):

```env
# .env file - SQLite (no setup needed)
DB_CONNECTION=sqlite
DB_DATABASE=database/database.sqlite
```

For MySQL/PostgreSQL:

```env
# MySQL
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=goravel_blog
DB_USERNAME=root
DB_PASSWORD=yourpassword

# PostgreSQL
DB_CONNECTION=postgresql
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=goravel_blog
DB_USERNAME=postgres
DB_PASSWORD=yourpassword
```

### 3. Initialize Application

```bash
# Generate application key
go run . artisan key:generate

# Run database migrations
go run . artisan migrate

# Seed RBAC permissions system
go run . artisan seed --seeder=rbac
```

### 4. Create Admin User

```bash
# Interactive admin user creation
go run . artisan user:create-admin

# Or create a regular user
go run . artisan user:create
```

### 5. Run Development Servers

You need two terminal windows:

**Terminal 1 - Backend Server:**
```bash
# With hot reload (recommended)
air

# Or standard Go run
go run . serve
```

**Terminal 2 - Frontend Assets:**
```bash
npm run dev
# or
yarn dev
```

The application will be available at `http://localhost:3500`

## 🏗️ Project Structure

```
├── app/                      # Application code
│   ├── auth/                 # Authentication & permissions
│   │   ├── permission_service.go
│   │   ├── permission_constants.go
│   │   └── permission_helper.go
│   ├── console/              # Artisan commands
│   │   └── commands/         # Custom commands
│   ├── contracts/            # Interfaces & contracts
│   ├── http/                 # HTTP layer
│   │   ├── controllers/      # API & page controllers
│   │   ├── middleware/       # HTTP middleware
│   │   └── requests/         # Request validation
│   ├── models/               # Database models
│   ├── providers/            # Service providers
│   └── services/             # Business logic
├── database/                 # Database files
│   ├── migrations/           # Schema migrations
│   └── seeders/              # Data seeders
├── resources/                # Frontend resources
│   ├── css/                  # Stylesheets
│   └── js/                   # React/TypeScript
│       ├── components/       # Reusable components
│       ├── contexts/         # React contexts
│       ├── pages/            # Page components
│       └── types/            # TypeScript types
├── routes/                   # Route definitions
│   ├── api.go                # API routes
│   ├── web.go                # Web routes
│   └── permissions.go        # Permission routes
└── docs/                     # Documentation
```

## 🛠️ Key Features

### Semi-Dynamic Permission System

The application features a powerful permission system:

- **Service-Action Format**: Permissions like `books_create`, `users_delete`
- **Auto-Detection**: Components automatically detect permissions
- **Server-Side Enforcement**: All controllers enforce permissions
- **Permission Matrix UI**: Visual role-permission management at `/admin/permissions`

### CRUD Resource Scaffolding

The application includes three powerful code generation commands that scaffold complete CRUD resources based on your models:

#### Quick Scaffolding Workflow

```bash
# 1. Create your model
go run . artisan make:model Lender

# 2. Generate service with model introspection
go run . artisan make:svc lender --model=Lender

# 3. Generate request validation classes
go run . artisan make:req lender --model=Lender

# 4. Generate controller with permission integration
go run . artisan make:ctrl lender --model=Lender

# 5. Complete setup
# - Add service to app/auth/permission_constants.go
# - Run: go run . artisan permissions:setup
# - Add routes to routes/api.go (instructions provided by make:ctrl)
```

#### Command Reference

##### `make:svc` - Service Generator

Generates a complete CRUD service based on your model with automatic field introspection.

```bash
# Basic usage
go run . artisan make:svc lender --model=Lender

# Flags
--model, -m    Model name to introspect (required)
```

**What it generates:**
- Service file in `app/services/` with builder pattern configuration
- Auto-detected search, sort, and filter fields from model
- Validation rules inferred from field types
- Scope filtering and soft delete support
- Commented GetFilterDefinitions() template for custom filters

**After generation:**
1. Add service constant to `app/auth/permission_constants.go`:
   ```go
   const (
       ServiceLenders ServiceRegistry = "lenders"  // Add this
   )
   ```
2. Run `go run . artisan permissions:setup` to sync permissions to database

##### `make:req` - Request Validation Generator

Generates Create and Update request validation classes based on your model.

```bash
# Basic usage
go run . artisan make:req lender --model=Lender

# Flags
--model, -m    Model name to introspect (required)
```

**What it generates:**
- Create request in `app/http/requests/` with non-pointer fields
- Update request with pointer fields for optional updates
- Validation rules based on field types
- Human-readable validation messages
- ToCreateData() and ToUpdateData() helper methods

**Field type handling:**
- Required fields (Create): Uses field type as-is from model
- Optional fields (Update): All fields become pointers for partial updates
- Already-pointer fields: Preserved as pointers in both requests

##### `make:ctrl` - Controller Generator

Generates a CRUD controller with permission-based authorization.

```bash
# Basic usage
go run . artisan make:ctrl lender --model=Lender

# Flags
--model, -m    Model name to introspect (required)
```

##### `make:audit` - Audit Fields Migration Generator

Generates a migration to add audit fields to an existing table.

```bash
# Basic usage (adds all audit fields including created_by)
go run . artisan make:audit lenders

# Exclude created_by field (if table already has it)
go run . artisan make:audit lenders --without-created-by

# Flags
--table, -t                Table name to add audit fields to
--without-created-by, -w   Exclude created_by field (default: false)
```

**What it generates:**
- Migration file in `database/migrations/`
- Automatically registers in `database/kernel.go`
- Adds all audit fields by default: `created_by`, `updated_by`, `deleted_by`, `ip_address`, `user_agent`
- Foreign key constraints to `users` table
- Indexes on audit fields for performance
- Complete `Down()` method for rollback

**Audit fields added (by default):**
- `created_by`: User who created the record
- `updated_by`: User who last updated the record
- `deleted_by`: User who soft-deleted the record
- `ip_address`: IP address from where the action was performed (45 chars for IPv6)
- `user_agent`: Browser/client user agent string

**After generation:**
```bash
# Run the migration
go run . artisan migrate

# To rollback
go run . artisan migrate:rollback
```

#### Controller Route Registration

After running `make:ctrl`, register routes in `routes/api.go` (detailed instructions provided):
```go
// 1. Import the controller
import "players/app/http/controllers/lenders"

// 2. Initialize controller
lenderController := lenders.NewLenderController()

// 3. Add optional auth routes (for public/scoped access - use for non-sensitive data)
router.Middleware(optionalAuth).Group(func(optionalAuthRouter route.Router) {
    optionalAuthRouter.Get("/lenders", lenderController.Index)
    optionalAuthRouter.Get("/lenders/search", lenderController.Search)
    optionalAuthRouter.Get("/lenders/filters", lenderController.FilterMetadata)
    optionalAuthRouter.Get("/lenders/{id}", lenderController.Show)
})

// 4. Add protected routes (always require authentication)
router.Middleware(jwtAuth).Group(func(protectedRouter route.Router) {
    protectedRouter.Post("/lenders", lenderController.Store)
    protectedRouter.Put("/lenders/{id}", lenderController.Update)
    protectedRouter.Delete("/lenders/{id}", lenderController.Delete)
})
```

**Middleware patterns:**
- **`optionalAuth`**: Allows public access but identifies authenticated users. Use for resources that support anonymous viewing with scoped results (e.g., books - anonymous users see all, authenticated users see filtered by permissions)
- **`jwtAuth`**: Requires authentication - blocks request if no valid token. Always use for write operations and sensitive data
- **Sensitive resources** (e.g., users): Use `jwtAuth` for ALL routes, including GET endpoints

#### Complete Example: Scaffolding a Lender Resource

```bash
# Step 1: Create the model
go run . artisan make:model Lender

# Edit app/models/lender.go
# package models
#
# import (
#     "github.com/goravel/framework/database/orm"
# )
#
# type Lender struct {
#     orm.Model
#     orm.SoftDeletes
#
#     FirstName      string  `json:"first_name" db:"first_name"`
#     LastName       string  `json:"last_name" db:"last_name"`
#     Phone          *string `json:"phone" db:"phone"`  // Optional
#     Email          string  `json:"email" db:"email"`
#     Address        *string `json:"address" db:"address"`  // Optional
#     IdentityNumber *string `json:"identity_number" db:"identity_number"`  // Optional
# }

# Step 2: Create migration
go run . artisan make:migration create_lenders_table

# Edit database/migrations/TIMESTAMP_create_lenders_table.go and run
go run . artisan migrate

# Step 3: Generate service with model introspection
go run . artisan make:svc lender --model=Lender
# ✓ Creates app/services/lender_service.go
# ✓ Auto-detects all fields from model
# ✓ Configures search, sort, filter fields
# ✓ Adds soft delete support

# Step 4: Register service permissions
# Edit app/auth/permission_constants.go
# Add: ServiceLenders ServiceRegistry = "lenders"

# Sync permissions to database
go run . artisan permissions:setup

# Step 5: Generate request validation classes
go run . artisan make:req lender --model=Lender
# ✓ Creates app/http/requests/lender_create_request.go
# ✓ Creates app/http/requests/lender_update_request.go
# ✓ Proper pointer handling for optional fields

# Step 6: Generate controller
go run . artisan make:ctrl lender --model=Lender
# ✓ Creates app/http/controllers/lenders/lender_controller.go
# ✓ Provides route registration instructions

# Step 7: Register routes
# Follow the instructions printed by make:ctrl command
# Edit routes/api.go and add the controller routes

# Step 8: Test your new CRUD resource
# Start servers and visit http://localhost:3500/api/lenders
```

#### Model Introspection Features

The scaffolding commands use Go AST parsing to intelligently analyze your models:

- **Field Detection**: Automatically extracts all struct fields
- **Type Inference**: Understands basic types (string, int, bool, etc.) and pointers
- **Tag Parsing**: Reads `json`, `db`, and `form` tags for field configuration
- **Soft Deletes**: Detects `orm.SoftDeletes` and configures service accordingly
- **Optional Fields**: Identifies pointer fields as optional (nullable in database)
- **Validation Rules**: Generates sensible validation rules based on field types

#### Legacy Commands

Individual component generators are still available for fine-grained control:

```bash
go run . artisan make:model ModelName
go run . artisan make:service ServiceName
go run . artisan make:repository RepositoryName
go run . artisan make:request RequestName
```

### Authentication & Authorization

- JWT-based authentication with HTTP-only cookies
- Role-Based Access Control (RBAC)
- Protected routes with middleware
- Global permission context in React

### Modern UI

- React with TypeScript
- Inertia.js for SPA-like experience
- Dark/Light theme toggle
- Responsive design with Tailwind CSS
- shadcn/ui component library

## 📝 Common Artisan Commands

### User Management
```bash
# Create admin user
go run . artisan user:create-admin

# Create regular user
go run . artisan user:create

# Assign role to user
go run . artisan role:assign user@example.com role-slug
```

### Database Operations
```bash
# Run migrations
go run . artisan migrate

# Rollback migrations
go run . artisan migrate:rollback

# Fresh migration (drop all tables and re-run)
go run . artisan migrate:fresh

# Run seeders
go run . artisan seed
go run . artisan seed --seeder=rbac
```

### CRUD Scaffolding
```bash
# Generate complete CRUD resource (recommended)
go run . artisan make:svc lender --model=Lender       # Service with model introspection
go run . artisan make:req lender --model=Lender       # Create/Update request validators
go run . artisan make:ctrl lender --model=Lender      # Controller with permissions

# Add audit fields to existing table
go run . artisan make:audit lenders                      # Adds all audit fields (including created_by)
go run . artisan make:audit lenders --without-created-by # Skip created_by if already exists

# Legacy: Generate individual components
go run . artisan make:model ModelName
go run . artisan make:service ServiceName
go run . artisan make:request RequestName
go run . artisan make:repository RepositoryName
```

### Permission Management
```bash
# Setup permissions (creates all service-action combinations)
go run . artisan permissions:setup

# Setup RBAC system
go run . artisan rbac:setup
```

## 🔒 Permission System Usage

### Backend - Page Controller
```go
func (c *BooksPageController) Index(ctx http.Context) http.Response {
    // Enforce permission check
    permHelper := auth.GetPermissionHelper()
    _, err := permHelper.RequireServicePermission(ctx, auth.ServiceBooks, auth.PermissionRead)
    if err != nil {
        return ctx.Response().Status(403).Json(map[string]interface{}{
            "error": "Forbidden",
        })
    }
    
    // Continue with rendering...
}
```

### Frontend - Auto Detection
```tsx
// CrudPage automatically detects permissions
<CrudPage
    resourceName="books"  // Auto-detects books_create, books_read, etc.
    title="Books Management"
    columns={bookColumns}
    data={data}
    filters={filters}
/>
```

### Frontend - Permission Hooks
```tsx
import { usePermissions } from '@/contexts/PermissionsContext';

function MyComponent() {
    const { canPerformAction, isSuperAdmin } = usePermissions();
    
    if (canPerformAction('books', 'create')) {
        // Show create button
    }
}
```

## 🧪 Development Workflow

### 1. Creating a New CRUD Resource

```bash
# Step 1: Create and define your model
go run . artisan make:model Feature
# Edit app/models/feature.go and add your fields

# Step 2: Create and run migration
go run . artisan make:migration create_features_table
# Edit the migration file
go run . artisan migrate

# Step 3: Scaffold complete CRUD resource
go run . artisan make:svc feature --model=Feature      # Service
go run . artisan make:req feature --model=Feature      # Request validators
go run . artisan make:ctrl feature --model=Feature     # Controller

# Step 4: Register permissions
# Edit app/auth/permission_constants.go and add:
# ServiceFeatures ServiceRegistry = "features"
go run . artisan permissions:setup

# Step 5: Register routes
# Follow instructions printed by make:ctrl command
# Edit routes/api.go

# Step 6: Restart servers
# Ctrl+C to stop both terminals, then restart:
# Terminal 1: air (or go run .)
# Terminal 2: npm run dev
```

### 2. Managing Permissions

1. Visit `/admin/permissions` as super admin
2. Create/edit roles
3. Assign permissions using the matrix grid
4. Test with different user accounts

### 3. Testing

```bash
# Create test users
go run . artisan user:create test@example.com testpass
go run . artisan role:assign test@example.com member

# Test API endpoints
curl -X GET "http://localhost:3500/api/books"

# Check server logs for permission debugging
# Look for: DEBUG HasPermission: user 1 has permissions: [books_create, books_read]
```

## 🐛 Troubleshooting

### Test Issues
```bash
# Tests failing with "panic: test timed out"
# - Increase timeout: go test -timeout 30s
# - Check for infinite loops or deadlocks

# Database/migration errors in tests
# - Tests use SQLite by default, no setup needed
# - Ensure APP_ENV=testing is set
# - Clean test databases: rm -rf tests/*/database/

# "sql: Scan error" with dates
# - Check migration uses DateTime() not String() for date fields
# - Ensure model uses *time.Time for nullable dates

# Mock errors in tests
# - Use real facades for integration tests
# - For unit tests, see helpers in tests/helpers/

# Permission errors in tests
# - Use helpers.SetupJWTUser() for test users with roles
# - Check CLAUDE.md for test examples
```

### Database Issues
```bash
# Connection errors
# - Check .env database credentials
# - Ensure database server is running
# - For SQLite: touch database/database.sqlite

# Migration errors
go run . artisan migrate:rollback
go run . artisan migrate:fresh

# Foreign key constraint errors
# - Check migration order in database/kernel.go
# - Ensure related tables exist first
```

### Permission Issues
```bash
# Permissions not working
# 1. Enable debug logging: LOG_LEVEL=debug
# 2. Check logs for "CheckScopedPermission" entries
# 3. Verify permission format: service_action (e.g., books_create)
# 4. Re-seed permissions:
go run . artisan seed --seeder=rbac

# User can't access features
# - Check user role: go run . artisan user:show email@example.com
# - Verify permissions: visit /admin/permissions as admin
# - Check scoped permissions (by_all, by_me, by_my_role)
```

### Development Server Issues
```bash
# Backend not starting
# - Check port 3500 is free: lsof -i :3500
# - Verify Go modules: go mod tidy
# - Check .env file exists and is valid

# Frontend not building
# - Clear node_modules: rm -rf node_modules && npm install
# - Check Node version: node --version (needs 16+)
# - Clear Vite cache: rm -rf node_modules/.vite

# Hot reload not working
# - Install Air: go install github.com/cosmtrek/air@latest
# - Check .air.toml configuration
# - Use 'air' instead of 'go run .'
```

### Common API Errors
```go
// "unexpected end of JSON input" in validation
// Replace ValidateRequest with manual binding:
var request requests.YourRequest
if err := ctx.Request().Bind(&request); err != nil {
    return ctx.Response().Json(http.StatusBadRequest, map[string]string{
        "error": "Invalid request format",
    })
}

// JWT token issues
// - Token stored in HTTP-only cookie named "token"
// - Check cookie domain matches your URL
// - For tests, use helpers.SetupJWTUser()
```

## 📚 Documentation

Detailed documentation available in the `docs/` directory:

- **[CRUD Resource Scaffolding](#crud-resource-scaffolding)** - Quick-start guide for `make:svc`, `make:req`, `make:ctrl` commands
- [Permission System Guide](docs/PERMISSION_SYSTEM_GUIDE.md) - Complete permission system documentation
- [CRUD E2E Guide](docs/CRUD_E2E_GUIDE.md) - Step-by-step CRUD implementation
- [Artisan Commands](docs/ARTISAN_COMMANDS.md) - All available commands
- [RBAC Implementation](docs/RBAC_IMPLEMENTATION.md) - Role-based access control details

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Follow the existing code style and conventions
4. Write tests for new features
5. Commit your changes: `git commit -m 'Add feature: description'`
6. Push to the branch: `git push origin feature-name`
7. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Built with [Goravel](https://www.goravel.dev/) - The Laravel-inspired Go framework
- UI components from [shadcn/ui](https://ui.shadcn.com/)
- Icons from [Lucide](https://lucide.dev/)