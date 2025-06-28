# Semi-Dynamic Permission System Guide

This document explains the semi-dynamic permission system with automatic permission detection, global context integration, and server-side enforcement.

## Overview

The permission system uses a **service_action** format where:
- **Services** are the main entities (books, users, roles, etc.)
- **Actions** are the operations (create, read, update, delete, etc.)
- **Permissions** are combinations like `books_create`, `users_delete`

## Key Features

- ✅ **Auto-detection**: Components automatically detect permissions based on resource names
- ✅ **Global Context**: Permissions available throughout the React app via context
- ✅ **Server-side Enforcement**: All page controllers enforce permissions before rendering
- ✅ **Real-time Updates**: Permission changes reflect immediately (no caching issues)
- ✅ **Type Safety**: Full TypeScript support for permission checks

## Quick Start

### 1. Backend - Page Controller with Permission Check

```go
// app/http/controllers/books_page_controller.go
func (c *BooksPageController) Index(ctx http.Context) http.Response {
    // Server-side permission check - returns 403 if unauthorized
    permHelper := auth.GetPermissionHelper()
    _, err := permHelper.RequireServicePermission(ctx, auth.ServiceBooks, auth.PermissionRead)
    if err != nil {
        return ctx.Response().Status(403).Json(map[string]interface{}{
            "error": "Forbidden",
            "message": "You don't have permission to access this page",
        })
    }
    
    // Continue with normal page rendering...
}
```

### 2. Frontend - Auto Permission Detection

```tsx
// No need to pass permission props manually!
// CrudPage auto-detects permissions based on resourceName
<CrudPage
    resourceName="books"  // Automatically checks books_create, books_read, etc.
    title="Books Management"
    columns={bookColumns}
    data={data}
    filters={filters}
    // No canCreate, canEdit, canDelete props needed!
/>
```

### 3. Navigation - Automatic Permission Filtering

```tsx
// app-sidebar.tsx configuration
const navigationConfig = {
    navMain: [
        {
            title: "Books",
            url: "/admin/books",
            icon: BookIcon,
            requiredService: "books",
            requiredAction: "read" as const,  // Menu item only shows if user has books_read
        },
    ]
}
```

## Permission Context API

### Using Permission Hooks

```tsx
import { usePermissions } from '@/contexts/PermissionsContext';

function MyComponent() {
    const { canPerformAction, isSuperAdmin } = usePermissions();
    
    // Check specific permission
    if (canPerformAction('books', 'create')) {
        // User can create books
    }
    
    // Check super admin status
    if (isSuperAdmin()) {
        // User is super admin
    }
}
```

### Available Context Methods

```typescript
interface PermissionsContextType {
    user: User | null;
    permissions: Record<string, ServicePermissions>;
    hasPermission: (permission: string) => boolean;
    hasServicePermission: (service: string, action: string) => boolean;
    canPerformAction: (service: string, action: CoreAction) => boolean;
    isSuperAdmin: () => boolean;
    isAdmin: () => boolean;
}
```

## Core Concepts

### Services (Entities)
- `books` - Book management
- `users` - User management  
- `roles` - Role management
- `permissions` - Permission management
- `reports` - Reports and analytics
- `system` - System administration

### Actions (Operations)
- `create` - Create new records
- `read` - Read/list records
- `update` - Update existing records
- `delete` - Delete records
- `export` - Export data
- `bulk_update` - Bulk update operations
- `bulk_delete` - Bulk delete operations
- `view` - View/list (UI display for read permission)
- `manage` - Full management (all operations)

### Permission Format
Permissions are stored as: `{service}_{action}`
- `books_create` - Can create books
- `users_update` - Can update users
- `reports_view` - Can view reports

## Implementation Guide

### 1. Page Controllers (Backend)

Every page controller MUST enforce permissions:

```go
type BooksPageController struct {
    *contracts.BasePageController
    bookService *services.BookService
    authHelper  contracts.AuthHelper
}

// Required: Identify the service this controller manages
func (c *BooksPageController) GetServiceIdentifier() auth.ServiceRegistry {
    return auth.ServiceBooks
}

// Required: Check permissions before rendering
func (c *BooksPageController) Index(ctx http.Context) http.Response {
    // Enforce server-side permission check
    permHelper := auth.GetPermissionHelper()
    _, err := permHelper.RequireServicePermission(ctx, auth.ServiceBooks, auth.PermissionRead)
    if err != nil {
        return ctx.Response().Status(403).Json(map[string]interface{}{
            "error": "Forbidden",
            "message": "You don't have permission to access this page",
        })
    }
    
    // Permissions are automatically included in global props
    // No need to manually pass them
    return inertia.Render(ctx, "Books/Index", props)
}
```

### 2. Global Permission Loading (Backend)

Permissions are automatically loaded in `app/http/inertia/inertia.go`:

```go
// Automatically builds permissions for all services
allPermissions := make(map[string]map[string]bool)
allServices := auth.GetAllServiceRegistries()

for _, service := range allServices {
    servicePerms := permHelper.BuildPermissionsMap(ctx, string(service))
    allPermissions[string(service)] = servicePerms
}

// Included in every page render
sharedProps["auth"] = map[string]interface{}{
    "user": userWithPermissions,
    "permissions": allPermissions,
}
```

### 3. Frontend Components

Components auto-detect permissions - no manual props needed:

```tsx
// CrudPage.tsx - Automatic permission detection
export function CrudPage<T>({ resourceName, ...props }) {
    const { canPerformAction } = usePermissions();
    
    // Auto-detect all permissions based on resourceName
    const canCreate = canPerformAction(resourceName, 'create');
    const canEdit = canPerformAction(resourceName, 'update');
    const canDelete = canPerformAction(resourceName, 'delete');
    const canView = canPerformAction(resourceName, 'read');
    
    // UI automatically adjusts based on permissions
    return (
        <>
            {canCreate && <Button>Add New</Button>}
            {/* Rest of UI respects permissions */}
        </>
    );
}
```

### 4. Navigation Auto-Filtering

Menu items automatically hide based on permissions:

```tsx
// app-sidebar.tsx
export function AppSidebar() {
    const { canPerformAction } = usePermissions();
    
    // Filter navigation based on permissions
    const filteredNavMain = navigationConfig.navMain.filter(item => {
        if (item.requiredService && item.requiredAction) {
            return canPerformAction(item.requiredService, item.requiredAction);
        }
        return true;
    });
    
    // Render only accessible menu items
}
```

## Database Schema

### Permissions Table
```sql
permissions (
    id, name, slug, description, category, resource, action, is_active
)
-- Example: slug = 'books_create'
```

### Role-Permission Pivot
```sql
role_permissions (
    role_id, permission_id, granted_at, granted_by_id, is_active
)
```

### User-Role Pivot
```sql
user_roles (
    user_id, role_id, assigned_at, assigned_by_id, is_active
)
```

## Debugging Permissions

### Enable Debug Logging

The permission system includes comprehensive debug logging:

```go
// In app/auth/permission_service.go
DEBUG HasPermission: user 1 has permissions: [books_create, books_read]
DEBUG HasPermission: checking permission: books_update
DEBUG loadUserPermissions: user has 2 roles
DEBUG loadUserPermissions: role master has 4 permissions
```

### Common Debug Points

1. **Check User Authentication**:
```go
fmt.Printf("DEBUG: User authentication status: %+v\n", c.GetCurrentUser(ctx) != nil)
```

2. **Check Loaded Permissions**:
```go
fmt.Printf("DEBUG: Permissions for 'books' resource: %+v\n", permissions)
```

3. **Check Role Loading**:
```go
fmt.Printf("DEBUG: User %d loaded with %d roles\n", userWithRoles.ID, len(userWithRoles.Roles))
```

## Setup & Commands

### Setup Permissions
```bash
# Creates all standard service-action permission combinations
go run . artisan permissions:setup
```

### Assign Role to User
```bash
go run . artisan role:assign <user-email> <role-slug>
```

### Create Admin User
```bash
go run . artisan user:create-admin
```

### Seed RBAC System
```bash
go run . artisan seed --seeder=rbac
```

## Security Best Practices

### 1. Always Enforce Server-Side

```go
// ✅ ALWAYS check permissions server-side
_, err := permHelper.RequireServicePermission(ctx, auth.ServiceBooks, auth.PermissionRead)
if err != nil {
    return ctx.Response().Status(403).Json(errorResponse)
}

// ❌ NEVER trust frontend-only checks
```

### 2. Use Type-Safe Constants

```go
// ✅ Use constants from auth package
permHelper.RequireServicePermission(ctx, auth.ServiceBooks, auth.PermissionCreate)

// ❌ Avoid string literals
permHelper.RequirePermission(ctx, "books.create")
```

### 3. Consistent Permission Format

```go
// ✅ Correct format: service_action
BuildPermissionSlug(ServiceBooks, PermissionCreate) // Returns: "books_create"

// ❌ Wrong formats: service.action, service:action
```

## Common Issues & Solutions

### Issue: Permissions Not Loading
```go
// Solution: Ensure roles are preloaded
err = facades.Orm().Query().
    Where("id = ?", user.ID).
    With("Roles.Permissions").  // Critical: Load both roles AND permissions
    First(&userWithRoles)
```

### Issue: Menu Items Not Hiding
```tsx
// Solution: Check navigation config
{
    title: "Books",
    requiredService: "books",     // Must match service name
    requiredAction: "read",       // Must match action name
}
```

### Issue: Permission Changes Not Reflecting
```go
// Solution: Permission cache is disabled for real-time updates
// In permission_service.go, cache is bypassed:
permissions := s.loadUserPermissions(user)  // Always loads fresh
```

### Issue: 403 on Page Access
```go
// Solution: Ensure user has required permission
// Check debug logs for which permission is being checked
DEBUG HasPermission: checking permission: books_read
DEBUG HasPermission: user 1 has permissions: [books_create]  // Missing books_read!
```

## Migration from Manual Permissions

If upgrading from manual permission props:

### Before (Manual):
```tsx
<CrudPage
    canCreate={permissions.canCreate}
    canEdit={permissions.canEdit}
    canDelete={permissions.canDelete}
    // ... other props
/>
```

### After (Automatic):
```tsx
<CrudPage
    resourceName="books"  // That's it! Permissions auto-detected
    // ... other props
/>
```

## Testing Permissions

### 1. Create Test Users
```bash
# Create test users with different roles
go run . artisan user:create test@example.com password123
go run . artisan role:assign test@example.com member
```

### 2. Test Permission Matrix
1. Visit `/admin/permissions` as super admin
2. Create/edit roles and assign permissions
3. Test access with different user accounts

### 3. Verify Server-Side Enforcement
```bash
# Try accessing protected endpoints without permission
curl -X GET "http://localhost:3500/admin/books" -H "Cookie: your-session-cookie"
# Should return 403 if no books_read permission
```

## Summary

The permission system provides:
- ✅ Automatic permission detection in components
- ✅ Server-side enforcement on all pages
- ✅ Global permission context for React
- ✅ Real-time permission updates
- ✅ Type-safe permission checks
- ✅ Automatic navigation filtering
- ✅ Comprehensive debug logging
- ✅ Role-based access control (RBAC)
- ✅ Permission matrix UI management

Remember: **Always enforce permissions server-side** - frontend checks are for UX only!