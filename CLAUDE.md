# Goravel Blog Application - Codebase Analysis

This is a **Goravel-based blog application** with a React frontend. Here's the summary:

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