package feature

import (
	"testing"
	"time"

	"books-database/app/auth"
	"books-database/app/models"
	"books-database/tests"
	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"
)

type ScopedPermissionsSimpleTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestScopedPermissionsSimpleTestSuite(t *testing.T) {
	suite.Run(t, new(ScopedPermissionsSimpleTestSuite))
}

func (s *ScopedPermissionsSimpleTestSuite) SetupTest() {
	s.RefreshDatabase()
}

// Test that the permission helper correctly recognizes scoped permissions
func (s *ScopedPermissionsSimpleTestSuite) TestPermissionHelperRecognizesScopes() {
	// Create a test user
	password, _ := facades.Hash().Make("password")
	user := &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: password,
		IsActive: true,
	}
	err := facades.Orm().Query().Create(user)
	s.NoError(err)

	// Create a role
	role := &models.Role{
		Name:     "Test Role",
		Slug:     "test-role",
		IsActive: true,
	}
	err = facades.Orm().Query().Create(role)
	s.NoError(err)

	// Create permissions
	permissions := []models.Permission{
		{Slug: "books_read", Name: "Read Books", IsActive: true},
		{Slug: "books_create", Name: "Create Books", IsActive: true},
		{Slug: "books_update", Name: "Update Books", IsActive: true},
		{Slug: "books_delete", Name: "Delete Books", IsActive: true},
	}

	for i := range permissions {
		err = facades.Orm().Query().Create(&permissions[i])
		s.NoError(err)
	}

	// Assign user to role
	userRole := &models.UserRole{
		UserID:     user.ID,
		RoleID:     role.ID,
		AssignedAt: time.Now(),
		IsActive:   true,
	}
	err = facades.Orm().Query().Create(userRole)
	s.NoError(err)

	// Assign scoped permissions to role
	rolePermissions := []models.RolePermission{
		{RoleID: role.ID, PermissionID: permissions[0].ID, Scope: "by_all"},     // read by_all
		{RoleID: role.ID, PermissionID: permissions[1].ID, Scope: "by_all"},     // create by_all
		{RoleID: role.ID, PermissionID: permissions[2].ID, Scope: "by_me"},      // update by_me
		{RoleID: role.ID, PermissionID: permissions[3].ID, Scope: "by_my_role"}, // delete by_my_role
	}

	for i := range rolePermissions {
		err = facades.Orm().Query().Create(&rolePermissions[i])
		s.NoError(err)
	}

	// Test permission service
	permService := auth.GetPermissionService()

	// Reload user with roles
	err = facades.Orm().Query().Where("id", user.ID).With("Roles").First(user)
	s.NoError(err)

	// Test that HasPermission does NOT automatically grant base permissions for scoped versions
	// This ensures proper scope validation happens in CheckScopedPermission
	s.True(permService.HasPermission(user, "books_read"), "User should have books_read permission (by_all grants base)")
	s.True(permService.HasPermission(user, "books_create"), "User should have books_create permission (by_all grants base)")
	s.False(permService.HasPermission(user, "books_update"), "User should NOT have books_update permission (only has by_me)")
	s.False(permService.HasPermission(user, "books_delete"), "User should NOT have books_delete permission (only has by_my_role)")

	// Test that GetUserPermissions returns scoped versions
	userPerms := permService.GetUserPermissions(user)
	s.Contains(userPerms, "books_read_by_all", "Should have scoped read permission")
	s.Contains(userPerms, "books_create_by_all", "Should have scoped create permission")
	s.Contains(userPerms, "books_update_by_me", "Should have scoped update permission")
	s.Contains(userPerms, "books_delete_by_my_role", "Should have scoped delete permission")
}

// Test BuildPermissionsMap with scoped permissions
func (s *ScopedPermissionsSimpleTestSuite) TestBuildPermissionsMapWithScopes() {
	// Create test data similar to above
	password, _ := facades.Hash().Make("password")
	user := &models.User{
		Name:     "Test User 2",
		Email:    "test2@example.com",
		Password: password,
		IsActive: true,
	}
	err := facades.Orm().Query().Create(user)
	s.NoError(err)

	role := &models.Role{
		Name:     "Editor Role",
		Slug:     "editor",
		IsActive: true,
	}
	err = facades.Orm().Query().Create(role)
	s.NoError(err)

	// Create permissions
	perms := []models.Permission{
		{Slug: "books_read", Name: "Read Books", IsActive: true},
		{Slug: "books_create", Name: "Create Books", IsActive: true},
		{Slug: "books_update", Name: "Update Books", IsActive: true},
	}

	for i := range perms {
		err = facades.Orm().Query().Create(&perms[i])
		s.NoError(err)
	}

	// Assign role to user
	userRole := &models.UserRole{
		UserID:     user.ID,
		RoleID:     role.ID,
		AssignedAt: time.Now(),
		IsActive:   true,
	}
	err = facades.Orm().Query().Create(userRole)
	s.NoError(err)

	// Assign scoped permissions
	rolePerms := []models.RolePermission{
		{RoleID: role.ID, PermissionID: perms[0].ID, Scope: "by_all"},
		{RoleID: role.ID, PermissionID: perms[1].ID, Scope: "by_me"},
		{RoleID: role.ID, PermissionID: perms[2].ID, Scope: "by_me"},
	}

	for i := range rolePerms {
		err = facades.Orm().Query().Create(&rolePerms[i])
		s.NoError(err)
	}

	// Test permission helper
	// Note: In real usage, this would use HTTP context
	// For testing, we verify the permission service directly
	permService := auth.GetPermissionService()

	// Reload user
	err = facades.Orm().Query().Where("id", user.ID).With("Roles").First(user)
	s.NoError(err)

	// Test that HasPermission correctly handles scoped permissions
	// by_all scope grants base permission, by_me and by_my_role do not
	s.True(permService.HasPermission(user, "books_read"), "Should have read permission (by_all grants base)")
	s.False(permService.HasPermission(user, "books_create"), "Should NOT have create permission (only has by_me)")
	s.False(permService.HasPermission(user, "books_update"), "Should NOT have update permission (only has by_me)")
	s.False(permService.HasPermission(user, "books_delete"), "Should not have delete permission (not assigned)")
}

// Test different scope combinations
func (s *ScopedPermissionsSimpleTestSuite) TestMultipleScopeScenarios() {
	// Create users with different roles and scopes
	password, _ := facades.Hash().Make("password")

	// Create roles
	adminRole := &models.Role{Name: "Admin", Slug: "admin", IsActive: true}
	editorRole := &models.Role{Name: "Editor", Slug: "editor", IsActive: true}
	memberRole := &models.Role{Name: "Member", Slug: "member", IsActive: true}

	err := facades.Orm().Query().Create(adminRole)
	s.NoError(err)
	err = facades.Orm().Query().Create(editorRole)
	s.NoError(err)
	err = facades.Orm().Query().Create(memberRole)
	s.NoError(err)

	// Create users
	admin := &models.User{Name: "Admin", Email: "admin@test.com", Password: password, IsActive: true}
	editor := &models.User{Name: "Editor", Email: "editor@test.com", Password: password, IsActive: true}
	member := &models.User{Name: "Member", Email: "member@test.com", Password: password, IsActive: true}

	err = facades.Orm().Query().Create(admin)
	s.NoError(err)
	err = facades.Orm().Query().Create(editor)
	s.NoError(err)
	err = facades.Orm().Query().Create(member)
	s.NoError(err)

	// Assign roles
	facades.Orm().Query().Create(&models.UserRole{UserID: admin.ID, RoleID: adminRole.ID, AssignedAt: time.Now(), IsActive: true})
	facades.Orm().Query().Create(&models.UserRole{UserID: editor.ID, RoleID: editorRole.ID, AssignedAt: time.Now(), IsActive: true})
	facades.Orm().Query().Create(&models.UserRole{UserID: member.ID, RoleID: memberRole.ID, AssignedAt: time.Now(), IsActive: true})

	// Create permissions
	readPerm := &models.Permission{Slug: "books_read", Name: "Read Books", IsActive: true}
	createPerm := &models.Permission{Slug: "books_create", Name: "Create Books", IsActive: true}
	updatePerm := &models.Permission{Slug: "books_update", Name: "Update Books", IsActive: true}
	deletePerm := &models.Permission{Slug: "books_delete", Name: "Delete Books", IsActive: true}

	facades.Orm().Query().Create(readPerm)
	facades.Orm().Query().Create(createPerm)
	facades.Orm().Query().Create(updatePerm)
	facades.Orm().Query().Create(deletePerm)

	// Admin: all permissions with by_all scope
	facades.Orm().Query().Create(&models.RolePermission{RoleID: adminRole.ID, PermissionID: readPerm.ID, Scope: "by_all"})
	facades.Orm().Query().Create(&models.RolePermission{RoleID: adminRole.ID, PermissionID: createPerm.ID, Scope: "by_all"})
	facades.Orm().Query().Create(&models.RolePermission{RoleID: adminRole.ID, PermissionID: updatePerm.ID, Scope: "by_all"})
	facades.Orm().Query().Create(&models.RolePermission{RoleID: adminRole.ID, PermissionID: deletePerm.ID, Scope: "by_all"})

	// Editor: mixed scopes
	facades.Orm().Query().Create(&models.RolePermission{RoleID: editorRole.ID, PermissionID: readPerm.ID, Scope: "by_all"})
	facades.Orm().Query().Create(&models.RolePermission{RoleID: editorRole.ID, PermissionID: createPerm.ID, Scope: "by_all"})
	facades.Orm().Query().Create(&models.RolePermission{RoleID: editorRole.ID, PermissionID: updatePerm.ID, Scope: "by_me"})
	facades.Orm().Query().Create(&models.RolePermission{RoleID: editorRole.ID, PermissionID: deletePerm.ID, Scope: "by_my_role"})

	// Member: restrictive scopes
	facades.Orm().Query().Create(&models.RolePermission{RoleID: memberRole.ID, PermissionID: readPerm.ID, Scope: "by_all"})
	facades.Orm().Query().Create(&models.RolePermission{RoleID: memberRole.ID, PermissionID: createPerm.ID, Scope: "by_me"})
	facades.Orm().Query().Create(&models.RolePermission{RoleID: memberRole.ID, PermissionID: updatePerm.ID, Scope: "by_me"})
	// No delete permission for member

	// Test permissions
	permService := auth.GetPermissionService()

	// Reload users with roles
	facades.Orm().Query().Where("id", admin.ID).With("Roles").First(admin)
	facades.Orm().Query().Where("id", editor.ID).With("Roles").First(editor)
	facades.Orm().Query().Where("id", member.ID).With("Roles").First(member)

	// Test admin permissions
	s.True(permService.HasPermission(admin, "books_read"), "Admin should have read")
	s.True(permService.HasPermission(admin, "books_create"), "Admin should have create")
	s.True(permService.HasPermission(admin, "books_update"), "Admin should have update")
	s.True(permService.HasPermission(admin, "books_delete"), "Admin should have delete")

	// Test editor permissions
	s.True(permService.HasPermission(editor, "books_read"), "Editor should have read (by_all)")
	s.True(permService.HasPermission(editor, "books_create"), "Editor should have create (by_all)")
	s.False(permService.HasPermission(editor, "books_update"), "Editor should NOT have update (only by_me)")
	s.False(permService.HasPermission(editor, "books_delete"), "Editor should NOT have delete (only by_my_role)")

	// Test member permissions
	s.True(permService.HasPermission(member, "books_read"), "Member should have read (by_all)")
	s.False(permService.HasPermission(member, "books_create"), "Member should NOT have create (only by_me)")
	s.False(permService.HasPermission(member, "books_update"), "Member should NOT have update (only by_me)")
	s.False(permService.HasPermission(member, "books_delete"), "Member should NOT have delete")

	// Verify scoped permissions are loaded correctly
	adminPerms := permService.GetUserPermissions(admin)
	s.Contains(adminPerms, "books_read_by_all")
	s.Contains(adminPerms, "books_update_by_all")

	editorPerms := permService.GetUserPermissions(editor)
	s.Contains(editorPerms, "books_read_by_all")
	s.Contains(editorPerms, "books_update_by_me")
	s.Contains(editorPerms, "books_delete_by_my_role")

	memberPerms := permService.GetUserPermissions(member)
	s.Contains(memberPerms, "books_read_by_all")
	s.Contains(memberPerms, "books_create_by_me")
	s.Contains(memberPerms, "books_update_by_me")
}

func (s *ScopedPermissionsSimpleTestSuite) TearDownTest() {
	// Cleanup handled by RefreshDatabase
}
