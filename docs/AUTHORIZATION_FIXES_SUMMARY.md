# Authorization and Single Role Fixes Summary

## Issues Fixed

### 1. 403 Error on Books Page ✅
**Problem**: User with `books_read_by_all` permission was getting 403 error when accessing /admin/books
**Root Cause**: Page controller was checking for `books_view` permission, but the system uses `books_read`
**Solution**: 
- Updated `RequireServicePermission` method to check for any scoped variant of the permission
- Changed page controller to use `PermissionRead` instead of `PermissionView`

### 2. Improved Error Messages ✅
**Problem**: 403 errors were showing generic "Forbidden" message without details
**Solution**: Updated `renderForbidden` method to:
- Include the actual error message from permission check
- Provide context about which resource is being accessed
- Properly handle both Inertia and JSON responses

### 3. Single Active Role Constraint ✅
**Problem**: Users could have multiple active roles at the same time
**Solution**: 
- Created migration to add unique constraint on active roles
- Updated `AssignRole` method to automatically deactivate existing roles
- Added logic to reactivate previously assigned roles if reassigning

## Technical Changes

### Backend Updates

1. **Permission Helper** (`app/auth/permission_helper.go`)
   ```go
   // Now checks for any scoped variant
   scopedPermissions := []string{
       basePermission,
       basePermission + "_by_all",
       basePermission + "_by_my_role", 
       basePermission + "_by_me",
   }
   ```

2. **Permission Service** (`app/auth/permission_service.go`)
   ```go
   // Deactivates all existing active roles before assigning new one
   _, err = facades.Orm().Query().Model(&models.UserRole{}).
       Where("user_id = ? AND is_active = ?", user.ID, true).
       Update("is_active", false)
   ```

3. **Generic Page Controller** (`app/contracts/generic_page_controller.go`)
   - Changed from `PermissionView` to `PermissionRead`
   - Enhanced error rendering with detailed messages

4. **Database Migration** (`20250723100000_ensure_single_active_role_per_user.go`)
   - Adds unique constraint ensuring one active role per user
   - Handles SQLite compatibility issues

## Usage Notes

### Permission Checking
The system now properly recognizes scoped permissions:
- `books_read_by_all` grants access to book pages
- `books_update_by_me` allows editing only own books
- `books_delete_by_my_role` allows deleting books by users with same role

### Role Assignment
When assigning a new role:
- Previous active roles are automatically deactivated
- If reassigning a previously held role, it's reactivated
- Database constraint prevents multiple active roles

### Error Handling
Users now see clear error messages:
- "insufficient permissions: books_read required" instead of generic "Forbidden"
- Proper Inertia error pages with context
- Resource type included in error responses

## Testing Recommendations

1. **Test Permission Access**
   - Login with user having `books_read_by_all`
   - Verify access to /admin/books page
   - Check data filtering based on scope

2. **Test Role Assignment**
   - Assign a role to a user
   - Assign a different role
   - Verify only one role is active

3. **Test Error Messages**
   - Access restricted pages
   - Verify clear error messages appear