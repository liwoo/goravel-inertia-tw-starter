# Goravel Blog Application

A modern web application built with Goravel (Go) and React, featuring JWT authentication, RBAC permissions, CRUD generators, dark mode, and a responsive UI.

## 🚀 Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd blog
go mod download
npm install

# Configure environment
cp .env.example .env
# Edit .env with your database credentials

# Setup database and permissions
go run . artisan migrate
go run . artisan key:generate
go run . artisan seed #Seeds RBAC permissions as well as Books
go run . artisan user:create-admin

# Run the application
# Terminal 1: Backend
air  # or: go run . serve

# Terminal 2: Frontend
npm run dev
```

Visit `http://localhost:3500` and login with your admin credentials.

## 📋 Prerequisites

- Go 1.18 or higher
- Node.js 16 or higher
- NPM or Yarn
- MySQL/PostgreSQL/SQLite database
- Air (optional, for hot reload): `go install github.com/cosmtrek/air@latest`

## 🔧 Detailed Setup

### 1. Initial Setup

```bash
# Clone the repository
git clone <repository-url>
cd blog

# Install dependencies
go mod download
npm install  # or: yarn install

# Setup environment
cp .env.example .env
```

### 2. Configure Database

Edit `.env` file with your database credentials:

```env
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=goravel_blog
DB_USERNAME=root
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

### CRUD Generator

Generate complete CRUD systems with one command:

```bash
# Generate full CRUD for a Product resource
go run . artisan make:crud-e2e Product

# This creates:
# - Model with soft deletes
# - Migration with indexes
# - Service with contracts
# - Controllers (API + Page)
# - Request validation
# - TypeScript types
# - React components
# - Permissions
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

### CRUD Generation
```bash
# Generate complete CRUD system
go run . artisan make:crud-e2e ResourceName

# Generate individual components
go run . artisan make:model ModelName
go run . artisan make:service ServiceName
go run . artisan make:controller ControllerName
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

### 1. Creating a New Feature

```bash
# Generate CRUD for your feature
go run . artisan make:crud-e2e Feature

# Run migrations
go run . artisan migrate

# Seed permissions
go run . artisan seed --seeder=rbac

# Restart servers
# Ctrl+C to stop, then restart both backend and frontend
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

### Database Issues
```bash
# Connection errors
# - Check .env database credentials
# - Ensure database server is running
# - For SQLite, ensure database file exists

# Migration errors
go run . artisan migrate:rollback
go run . artisan migrate:fresh
```

### Permission Issues
```bash
# Permissions not working
# 1. Check debug logs in console
# 2. Verify permission format: service_action (e.g., books_create)
# 3. Re-seed permissions:
go run . artisan seed --seeder=rbac

# User can't access features
# - Check user has correct role
# - Verify role has required permissions in /admin/permissions
```

### Development Server Issues
```bash
# Backend not reloading
# - Restart with Ctrl+C then run again
# - Use 'air' for hot reload

# Frontend not updating
# - Check npm run dev is running
# - Clear browser cache
# - Check browser console for errors
```

### Common Validation Errors
```go
// If you get "unexpected end of JSON input" in validation
// Replace ValidateRequest with manual binding:
if err := ctx.Request().Bind(&request); err != nil {
    return nil, err
}
// Then add manual validation
```

## 📚 Documentation

Detailed documentation available in the `docs/` directory:

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