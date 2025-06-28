-- Clear existing data
DELETE FROM role_permissions;
DELETE FROM user_roles;
DELETE FROM permissions;
DELETE FROM roles;

-- Insert roles
INSERT INTO roles (name, slug, description, level, is_active, created_at, updated_at) VALUES 
('Super Administrator', 'super-admin', 'Full system access with all permissions', 100, 1, datetime('now'), datetime('now')),
('Administrator', 'admin', 'Administrative access to most features', 80, 1, datetime('now'), datetime('now')),
('Librarian', 'librarian', 'Full book management access', 60, 1, datetime('now'), datetime('now')),
('Moderator', 'moderator', 'Limited administrative access', 40, 1, datetime('now'), datetime('now')),
('Member', 'member', 'Regular user with borrowing privileges', 20, 1, datetime('now'), datetime('now')),
('Guest', 'guest', 'Basic read-only access', 10, 1, datetime('now'), datetime('now'));

-- Insert permissions
INSERT INTO permissions (name, slug, description, category, resource, action, is_active, requires_ownership, can_delegate, created_at, updated_at) VALUES 
-- Books permissions
('Create Books', 'books_create', 'Create new books', 'books', 'books', 'create', 1, 0, 0, datetime('now'), datetime('now')),
('Read Books', 'books_read', 'View books', 'books', 'books', 'read', 1, 0, 0, datetime('now'), datetime('now')),
('Update Books', 'books_update', 'Update existing books', 'books', 'books', 'update', 1, 0, 0, datetime('now'), datetime('now')),
('Delete Books', 'books_delete', 'Delete books', 'books', 'books', 'delete', 1, 0, 0, datetime('now'), datetime('now')),
('Export Books', 'books_export', 'Export books data', 'books', 'books', 'export', 1, 0, 0, datetime('now'), datetime('now')),
('Bulk Update Books', 'books_bulk_update', 'Bulk update books', 'books', 'books', 'bulk_update', 1, 0, 0, datetime('now'), datetime('now')),
('Bulk Delete Books', 'books_bulk_delete', 'Bulk delete books', 'books', 'books', 'bulk_delete', 1, 0, 0, datetime('now'), datetime('now')),

-- Users permissions
('Create Users', 'users_create', 'Create new users', 'users', 'users', 'create', 1, 0, 0, datetime('now'), datetime('now')),
('Read Users', 'users_read', 'View users', 'users', 'users', 'read', 1, 0, 0, datetime('now'), datetime('now')),
('Update Users', 'users_update', 'Update existing users', 'users', 'users', 'update', 1, 0, 0, datetime('now'), datetime('now')),
('Delete Users', 'users_delete', 'Delete users', 'users', 'users', 'delete', 1, 0, 0, datetime('now'), datetime('now')),
('Export Users', 'users_export', 'Export users data', 'users', 'users', 'export', 1, 0, 0, datetime('now'), datetime('now')),
('Bulk Update Users', 'users_bulk_update', 'Bulk update users', 'users', 'users', 'bulk_update', 1, 0, 0, datetime('now'), datetime('now')),
('Bulk Delete Users', 'users_bulk_delete', 'Bulk delete users', 'users', 'users', 'bulk_delete', 1, 0, 0, datetime('now'), datetime('now')),

-- Roles permissions
('Create Roles', 'roles_create', 'Create new roles', 'roles', 'roles', 'create', 1, 0, 0, datetime('now'), datetime('now')),
('Read Roles', 'roles_read', 'View roles', 'roles', 'roles', 'read', 1, 0, 0, datetime('now'), datetime('now')),
('Update Roles', 'roles_update', 'Update existing roles', 'roles', 'roles', 'update', 1, 0, 0, datetime('now'), datetime('now')),
('Delete Roles', 'roles_delete', 'Delete roles', 'roles', 'roles', 'delete', 1, 0, 0, datetime('now'), datetime('now')),
('Export Roles', 'roles_export', 'Export roles data', 'roles', 'roles', 'export', 1, 0, 0, datetime('now'), datetime('now')),
('Bulk Update Roles', 'roles_bulk_update', 'Bulk update roles', 'roles', 'roles', 'bulk_update', 1, 0, 0, datetime('now'), datetime('now')),
('Bulk Delete Roles', 'roles_bulk_delete', 'Bulk delete roles', 'roles', 'roles', 'bulk_delete', 1, 0, 0, datetime('now'), datetime('now')),

-- Permissions permissions
('Create Permissions', 'permissions_create', 'Create new permissions', 'permissions', 'permissions', 'create', 1, 0, 0, datetime('now'), datetime('now')),
('Read Permissions', 'permissions_read', 'View permissions', 'permissions', 'permissions', 'read', 1, 0, 0, datetime('now'), datetime('now')),
('Update Permissions', 'permissions_update', 'Update existing permissions', 'permissions', 'permissions', 'update', 1, 0, 0, datetime('now'), datetime('now')),
('Delete Permissions', 'permissions_delete', 'Delete permissions', 'permissions', 'permissions', 'delete', 1, 0, 0, datetime('now'), datetime('now')),
('Export Permissions', 'permissions_export', 'Export permissions data', 'permissions', 'permissions', 'export', 1, 0, 0, datetime('now'), datetime('now')),
('Bulk Update Permissions', 'permissions_bulk_update', 'Bulk update permissions', 'permissions', 'permissions', 'bulk_update', 1, 0, 0, datetime('now'), datetime('now')),
('Bulk Delete Permissions', 'permissions_bulk_delete', 'Bulk delete permissions', 'permissions', 'permissions', 'bulk_delete', 1, 0, 0, datetime('now'), datetime('now')),

-- System permissions
('System Manage', 'system_manage', 'Full system management', 'system', 'system', 'manage', 1, 0, 0, datetime('now'), datetime('now')),

-- Reports permissions
('Read Reports', 'reports_read', 'View reports and analytics', 'reports', 'reports', 'read', 1, 0, 0, datetime('now'), datetime('now')),
('Create Reports', 'reports_create', 'Create custom reports', 'reports', 'reports', 'create', 1, 0, 0, datetime('now'), datetime('now')),
('Export Reports', 'reports_export', 'Export reports', 'reports', 'reports', 'export', 1, 0, 0, datetime('now'), datetime('now'));

-- Assign all permissions to super-admin role
INSERT INTO role_permissions (role_id, permission_id, is_active, created_at, updated_at)
SELECT r.id, p.id, 1, datetime('now'), datetime('now')
FROM roles r, permissions p
WHERE r.slug = 'super-admin';

-- Assign specific permissions to other roles
-- Admin gets most permissions except system management
INSERT INTO role_permissions (role_id, permission_id, is_active, created_at, updated_at)
SELECT r.id, p.id, 1, datetime('now'), datetime('now')
FROM roles r, permissions p
WHERE r.slug = 'admin' 
AND p.category IN ('books', 'users', 'roles', 'reports')
AND p.action NOT IN ('delete', 'bulk_delete');

-- Librarian gets book management permissions
INSERT INTO role_permissions (role_id, permission_id, is_active, created_at, updated_at)
SELECT r.id, p.id, 1, datetime('now'), datetime('now')
FROM roles r, permissions p
WHERE r.slug = 'librarian' 
AND p.category = 'books';

-- Member gets read and borrow permissions for books
INSERT INTO role_permissions (role_id, permission_id, is_active, created_at, updated_at)
SELECT r.id, p.id, 1, datetime('now'), datetime('now')
FROM roles r, permissions p
WHERE r.slug = 'member' 
AND p.category = 'books'
AND p.action = 'read';

-- Guest gets only read permissions for books
INSERT INTO role_permissions (role_id, permission_id, is_active, created_at, updated_at)
SELECT r.id, p.id, 1, datetime('now'), datetime('now')
FROM roles r, permissions p
WHERE r.slug = 'guest' 
AND p.category = 'books'
AND p.action = 'read';

-- Assign current admin user (ID=1) to super-admin role
INSERT INTO user_roles (user_id, role_id, assigned_at, is_active, notes, created_at, updated_at)
SELECT 1, r.id, datetime('now'), 1, 'Initial super admin assignment', datetime('now'), datetime('now')
FROM roles r
WHERE r.slug = 'super-admin';