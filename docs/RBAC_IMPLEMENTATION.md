# Role-Based Access Control (RBAC) Implementation

This document describes the comprehensive RBAC system implemented for robust permission management.

## Overview

The RBAC system provides:
- **Hierarchical Roles**: Roles can inherit from parent roles
- **Granular Permissions**: Fine-grained control over resources and actions
- **Caching**: High-performance permission checking with intelligent caching
- **Flexible**: Easy to extend with new roles and permissions
- **Database-driven**: No code changes needed for new permissions
- **Contract Integration**: Seamlessly works with controller contracts

## System Architecture

### Database Schema

```
Users (1) ←→ (M) UserRoles (M) ←→ (1) Roles
Roles (1) ←→ (M) RolePermissions (M) ←→ (1) Permissions
Roles (1) ←→ (M) Roles (parent-child hierarchy)
```

### Models

#### User Model
```go
type User struct {
    orm.Model
    Name     string `gorm:"not null"`
    Email    string `gorm:"uniqueIndex;not null"`
    Password string `gorm:"not null"`
    
    // Legacy role field (backward compatibility)
    Role string `gorm:"default:'USER'"`
    
    // User status
    IsActive      bool `gorm:"default:true"`
    EmailVerified bool `gorm:"default:false"`
    
    // Many-to-many relationships
    Roles []Role `gorm:"many2many:user_roles"`
    
    orm.SoftDeletes
}
```

#### Role Model
```go
type Role struct {
    orm.Model
    Name        string `gorm:"uniqueIndex;not null"`
    Slug        string `gorm:"uniqueIndex;not null"`
    Description string
    Level       int    `gorm:"default:0"` // Higher = more permissions
    IsActive    bool   `gorm:"default:true"`
    
    // Hierarchical structure
    ParentID *uint
    Parent   *Role
    Children []Role
    
    // Relationships
    Users       []User       `gorm:"many2many:user_roles"`
    Permissions []Permission `gorm:"many2many:role_permissions"`
    
    orm.SoftDeletes
}
```

#### Permission Model
```go
type Permission struct {
    orm.Model
    Name        string `gorm:"uniqueIndex;not null"`
    Slug        string `gorm:"uniqueIndex;not null"`
    Description string
    Category    string `gorm:"index;not null"` // e.g., "books", "users"
    Resource    string `gorm:"index"`          // Specific resource type
    Action      string `gorm:"index;not null"` // e.g., "create", "read"
    IsActive    bool   `gorm:"default:true"`
    
    // Permission metadata
    RequiresOwnership bool `gorm:"default:false"`
    CanDelegate       bool `gorm:"default:false"`
    
    // Relationships
    Roles []Role `gorm:"many2many:role_permissions"`
    
    orm.SoftDeletes
}
```

## Default Role Hierarchy

```
Super Admin (Level 100)
└── Admin (Level 80)
    └── Librarian (Level 60)
        └── Moderator (Level 40)
            └── Member (Level 20)
                └── Guest (Level 10)
```

### Role Permissions

| Role | Permissions |
|------|-------------|
| **Super Admin** | All permissions (system-wide access) |
| **Admin** | Books: full management, Users: view/create/update/manage, Roles: view/assign, Reports: view/export/create |
| **Librarian** | Books: full management + export, Users: view only, Reports: view/export |
| **Moderator** | Books: view/create/update/borrow/return, Users: view, Reports: view |
| **Member** | Books: view/borrow/return |
| **Guest** | Books: view only |

## Permission Categories

### Books Management
- `books.viewAny` - View any books in the system
- `books.view` - View specific books
- `books.create` - Create new books
- `books.update` - Update existing books
- `books.delete` - Delete books
- `books.manage` - Full book management
- `books.borrow` - Borrow books
- `books.return` - Return books
- `books.export` - Export books data

### User Management
- `users.viewAny` - View any users
- `users.view` - View specific users
- `users.create` - Create new users
- `users.update` - Update existing users
- `users.delete` - Delete users
- `users.manage` - Full user management
- `users.impersonate` - Impersonate other users

### Role Management
- `roles.viewAny` - View any roles
- `roles.view` - View specific roles
- `roles.create` - Create new roles
- `roles.update` - Update existing roles
- `roles.delete` - Delete roles
- `roles.manage` - Full role management
- `roles.assign` - Assign roles to users

### System Administration
- `system.backup` - Create system backups
- `system.restore` - Restore from backup
- `system.configure` - Configure system settings
- `system.monitor` - Monitor system performance
- `system.logs` - View system logs
- `system.manage` - Full system management

### Reports & Analytics
- `reports.view` - View reports and analytics
- `reports.export` - Export reports
- `reports.create` - Create custom reports
- `reports.manage` - Full report management

## Usage Examples

### Basic Permission Checking

```go
// In controllers
func (c *BookController) Store(ctx http.Context) http.Response {
    // Check permission using contract
    if err := c.CheckPermission(ctx, "books.create", nil); err != nil {
        return c.ForbiddenResponse(ctx, "Access denied")
    }
    
    // Continue with creation...
}
```

### Using Permission Helper

```go
// Direct permission checking
permHelper := helpers.GetPermissionHelper()

// Require specific permission
user, err := permHelper.RequirePermission(ctx, "books.manage")
if err != nil {
    // Handle unauthorized access
}

// Check permission (returns bool)
canDelete := permHelper.CheckPermission(ctx, "books.delete")

// Check resource access
canAccess := permHelper.CheckResourceAccess(ctx, "update", "books", bookID)
```

### Using Permission Service

```go
permService := services.GetPermissionService()

// Check user permission
hasPermission := permService.HasPermission(user, "books.create")

// Check role
hasRole := permService.HasRole(user, "admin")

// Check resource access with ownership
canAccess := permService.CanAccessResource(user, "update", "books", bookID)

// Assign role to user
err := permService.AssignRole(user, "librarian", assignedBy)

// Create new role
role, err := permService.CreateRole("Editor", "editor", "Content editor", 30, "member")

// Grant permission to role
err := permService.GrantPermissionToRole("editor", "books.create", grantedBy)
```

### Frontend Permission Maps

Controllers automatically generate permission maps for frontend:

```go
// In any controller implementing contracts
permissions := c.BuildPermissionsMap(ctx, "books")

// Returns:
{
    "canView":   true,
    "canCreate": false,
    "canEdit":   true,
    "canDelete": false,
    "canManage": false,
    "canBorrow": true,
    "canReturn": true,
    "isAdmin":   false,
    "isSuperAdmin": false
}
```

### Wildcard Permissions

The system supports wildcard permissions:

```go
// Grant all book permissions
permService.GrantPermissionToRole("librarian", "books.*", grantedBy)

// Grant specific action across all resources
permService.GrantPermissionToRole("viewer", "*.view", grantedBy)

// Grant everything (super admin)
permService.GrantPermissionToRole("super-admin", "*.*", grantedBy)
```

### Role Hierarchy & Inheritance

```go
// Roles automatically inherit permissions from parent roles
// If "member" has "books.view" and "admin" inherits from "member",
// then "admin" automatically gets "books.view" + their own permissions

// Check role hierarchy
adminRole := user.GetHighestRole()
canManageUser := user.CanManageUser(otherUser) // Based on role levels
```

## Advanced Features

### Ownership-Based Permissions

```go
// Some permissions can require ownership
permission := models.Permission{
    Slug: "books.update.own",
    RequiresOwnership: true, // User must own the resource
}

// The system will check ownership before granting access
canUpdate := permService.CanAccessResource(user, "update", "books", bookID)
// This will check both permission AND ownership
```

### Temporary Role Assignments

```go
// Assign role with expiration
userRole := models.UserRole{
    UserID:    user.ID,
    RoleID:    role.ID,
    ExpiresAt: &futureTime, // Role expires
    IsActive:  true,
}
```

### Permission Caching

The system automatically caches permissions for performance:

```go
// Permissions are cached for 15 minutes by default
// Cache is automatically invalidated when:
// - User roles change
// - Role permissions change
// - User is deleted/deactivated

// Manual cache refresh
permService.RefreshCache()
```

## Database Setup

### Run RBAC Seeder

```bash
# Seed the database with default roles and permissions
go run . artisan seed --seeder=rbac
```

### Manual Role/Permission Creation

```go
// Create a new role
role, err := permService.CreateRole(
    "Content Manager", 
    "content-manager", 
    "Manages content and books", 
    50, 
    "member" // parent role
)

// Create a new permission
permission, err := permService.CreatePermission(
    "Manage Categories",
    "categories.manage",
    "categories",
    "manage",
    "",
    "Full category management access"
)

// Grant permission to role
err = permService.GrantPermissionToRole("content-manager", "categories.manage", grantedBy)
```

## Integration with Controller Contracts

The RBAC system seamlessly integrates with your existing controller contracts:

```go
// All controllers automatically get permission checking
type BookController struct {
    *contracts.BaseCrudController
    // ... other fields
}

// Contract methods automatically use RBAC
func (c *BookController) Index(ctx http.Context) http.Response {
    // This automatically checks permissions based on the resource type
    req, err := c.ValidatePaginationRequest(ctx)
    // ... rest of implementation
}
```

## Security Best Practices

1. **Principle of Least Privilege**: Users start with minimal permissions
2. **Role Hierarchy**: Higher roles can manage lower roles
3. **Permission Inheritance**: Reduces duplication and ensures consistency
4. **Audit Trail**: All role assignments and permission grants are logged
5. **Caching**: Improves performance while maintaining security
6. **Wildcard Control**: Carefully managed wildcard permissions
7. **Ownership Validation**: Resource-specific ownership checks

## Benefits

1. **Scalable**: Easy to add new roles and permissions without code changes
2. **Flexible**: Supports complex permission hierarchies and inheritance
3. **Performant**: Intelligent caching reduces database queries
4. **Secure**: Built-in security best practices and audit trails
5. **Contract Integrated**: Works seamlessly with existing controller contracts
6. **User-Friendly**: Clear permission maps for frontend development
7. **Database-Driven**: Configure permissions without code deployment

The RBAC system provides enterprise-grade permission management while maintaining simplicity and performance!