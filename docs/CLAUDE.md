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

The codebase follows modern web development patterns with strong separation of concerns and defensive security practices.