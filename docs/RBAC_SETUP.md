# RBAC Setup Guide

## 🚀 Quick Start with Your Existing Admin User

You now have enhanced commands that integrate RBAC with your existing admin user creation workflow.

## Setup Options

### Option 1: Complete RBAC Setup (Recommended)
This sets up the entire RBAC system and upgrades any existing users:

```bash
# Set up RBAC system with roles, permissions, and upgrade existing users
go run . artisan rbac:setup
```

**What this does:**
- ✅ Creates all roles (super-admin, admin, librarian, moderator, member, guest)
- ✅ Creates all permissions (30+ granular permissions)
- ✅ Assigns permissions to roles based on hierarchy
- ✅ Upgrades existing users with RBAC roles based on their legacy roles
- ✅ Shows a complete summary

### Option 2: Manual Setup (Step by Step)

#### Step 1: Create RBAC System
```bash
# Create roles and permissions only
go run . artisan seed --seeder=rbac
```

#### Step 2: Create New Admin User (Enhanced)
```bash
# Create admin user with automatic super-admin role assignment
go run . artisan user:create-admin
```

#### Step 3: Assign Roles to Existing Users
```bash
# Assign role to specific user
go run . artisan rbac:assign user@example.com super-admin
go run . artisan rbac:assign another@example.com librarian
```

## Available Roles

| Role | Level | Permissions | Description |
|------|-------|-------------|-------------|
| `super-admin` | 100 | All permissions | Full system access |
| `admin` | 80 | Administrative access | Books + Users + Reports management |
| `librarian` | 60 | Book management | Full book operations + reports |
| `moderator` | 40 | Limited admin | Books create/update + basic user view |
| `member` | 20 | User privileges | Book browsing + borrow/return |
| `guest` | 10 | Read-only | View books only |

## Available Commands

### Enhanced Admin Creation
```bash
# Creates user with both legacy ADMIN role and super-admin RBAC role
go run . artisan user:create-admin
```

### RBAC Management
```bash
# Complete RBAC setup
go run . artisan rbac:setup

# Assign role to user
go run . artisan rbac:assign email@domain.com role-name

# Available roles: super-admin, admin, librarian, moderator, member, guest
```

### Database Seeding
```bash
# Seed everything (includes RBAC)
go run . artisan migrate:fresh --seed

# Seed only RBAC
go run . artisan seed --seeder=rbac
```

## What Happens When You Create an Admin User

Your enhanced `user:create-admin` command now:

1. **Creates the user** with legacy `ADMIN` role (backward compatibility)
2. **Assigns super-admin RBAC role** automatically
3. **Sets user as active and verified**
4. **Provides clear feedback** about role assignments

### Example Output:
```bash
$ go run . artisan user:create-admin

Enter admin name: John Administrator
Enter admin email: admin@example.com
Enter admin password: [hidden]

✓ Admin user 'John Administrator' (admin@example.com) created successfully!
✓ Super-admin role assigned to 'John Administrator' successfully!
ℹ User has both legacy ADMIN role and super-admin RBAC role for maximum compatibility.
```

## Permission System Integration

Your controllers automatically use the RBAC system:

```go
// This automatically checks RBAC permissions
func (c *BookController) Store(ctx http.Context) http.Response {
    // Checks if user has "books.create" permission
    if err := c.CheckPermission(ctx, "books.create", nil); err != nil {
        return c.ForbiddenResponse(ctx, "Access denied")
    }
    // ... rest of implementation
}
```

## Troubleshooting

### If Role Assignment Fails
```bash
# Make sure RBAC system is set up first
go run . artisan seed --seeder=rbac

# Then create admin user
go run . artisan user:create-admin

# Or assign role manually
go run . artisan rbac:assign admin@example.com super-admin
```

### Check User Permissions
The system provides automatic permission maps for your frontend:
```go
// In your controllers
permissions := c.BuildPermissionsMap(ctx, "books")
// Returns: {"canCreate": true, "canEdit": true, "canDelete": true, ...}
```

## Migration from Legacy System

Your existing users with `Role = "ADMIN"` or `Role = "USER"` will be automatically upgraded:

- `ADMIN` users → `super-admin` RBAC role
- `USER` users → `member` RBAC role
- Legacy roles are preserved for backward compatibility

## Next Steps

1. **Run the setup**: `go run . artisan rbac:setup`
2. **Create admin users**: `go run . artisan user:create-admin`
3. **Your controllers already work** - permissions are automatically enforced!
4. **Frontend gets permission maps** - use them to show/hide UI elements

The RBAC system is now fully integrated with your existing admin creation workflow! 🎉