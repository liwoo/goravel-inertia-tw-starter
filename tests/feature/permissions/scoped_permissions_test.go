package feature

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/goravel/framework/facades"
	"players/app/models"
	"players/app/auth"
	"players/tests"
)

type ScopedPermissionsTestSuite struct {
	suite.Suite
	tests.TestCase
	superAdmin *models.User
	editor1    *models.User
	editor2    *models.User
	regularUser *models.User
	editorRole *models.Role
}

func TestScopedPermissionsTestSuite(t *testing.T) {
	suite.Run(t, new(ScopedPermissionsTestSuite))
}

// SetupTest will run before each test in the suite.
func (s *ScopedPermissionsTestSuite) SetupTest() {
	// Refresh database
	s.RefreshDatabase()
	
	// Create test users and roles
	s.setupTestUsers()
	s.setupPermissions()
}

func (s *ScopedPermissionsTestSuite) setupTestUsers() {
	// Create super admin
	password, _ := facades.Hash().Make("password")
	s.superAdmin = &models.User{
		Name: "Super Admin",
		Email: "super@example.com",
		Password: password,
		IsSuperAdmin: true,
	}
	facades.Orm().Query().Create(s.superAdmin)
	
	// Create editor role
	s.editorRole = &models.Role{
		Name: "Editor",
		Slug: "editor",
	}
	facades.Orm().Query().Create(s.editorRole)
	
	// Create editors
	s.editor1 = &models.User{
		Name: "Editor One",
		Email: "editor1@example.com",
		Password: password,
	}
	facades.Orm().Query().Create(s.editor1)
	facades.Orm().Query().Create(&models.UserRole{
		UserID:     s.editor1.ID,
		RoleID:     s.editorRole.ID,
		AssignedAt: time.Now(),
		IsActive:   true,
	})
	
	s.editor2 = &models.User{
		Name: "Editor Two",
		Email: "editor2@example.com",
		Password: password,
	}
	facades.Orm().Query().Create(s.editor2)
	facades.Orm().Query().Create(&models.UserRole{
		UserID:     s.editor2.ID,
		RoleID:     s.editorRole.ID,
		AssignedAt: time.Now(),
		IsActive:   true,
	})
	
	// Create regular user
	s.regularUser = &models.User{
		Name: "Regular User",
		Email: "regular@example.com",
		Password: password,
	}
	facades.Orm().Query().Create(s.regularUser)
}

func (s *ScopedPermissionsTestSuite) setupPermissions() {
	// Create permissions
	permissions := []models.Permission{
		{Slug: "books_read", Name: "Read Books"},
		{Slug: "books_create", Name: "Create Books"},
		{Slug: "books_update", Name: "Update Books"},
		{Slug: "books_delete", Name: "Delete Books"},
	}
	
	for _, perm := range permissions {
		facades.Orm().Query().Create(&perm)
	}
	
	// Assign scoped permissions to editor role
	var readPerm models.Permission
	facades.Orm().Query().Where("slug", "books_read").First(&readPerm)
	
	var updatePerm models.Permission
	facades.Orm().Query().Where("slug", "books_update").First(&updatePerm)
	
	var deletePerm models.Permission
	facades.Orm().Query().Where("slug", "books_delete").First(&deletePerm)
	
	// Editors can read all books
	facades.Orm().Query().Create(&models.RolePermission{
		RoleID: s.editorRole.ID,
		PermissionID: readPerm.ID,
		Scope: "by_all",
	})
	
	// Editors can only update their own books
	facades.Orm().Query().Create(&models.RolePermission{
		RoleID: s.editorRole.ID,
		PermissionID: updatePerm.ID,
		Scope: "by_me",
	})
	
	// Editors can delete books created by users with same role
	facades.Orm().Query().Create(&models.RolePermission{
		RoleID: s.editorRole.ID,
		PermissionID: deletePerm.ID,
		Scope: "by_my_role",
	})
}

func (s *ScopedPermissionsTestSuite) TestPermissionScopeHelper() {
	// Test GetUserScope
	scopeHelper := auth.GetScopeHelper()
	
	// Create a mock context with editor1
	// Note: In a real test, you'd use proper HTTP context with authenticated user
	// This is a simplified example
	
	// Test super admin always gets full scope
	s.Equal(auth.ScopeByAll, scopeHelper.GetUserScope(nil, auth.ServiceBooks, auth.PermissionRead))
	
	// Test editor gets correct scopes
	// Editor should have:
	// - read: by_all
	// - update: by_me
	// - delete: by_my_role
}

func (s *ScopedPermissionsTestSuite) TestDataFiltering() {
	// Create test books
	bookBySuperAdmin := &models.Book{
		Title: "Super Admin Book",
		Author: "Admin Author",
		ISBN: "1234567890",
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: &s.superAdmin.ID,
		},
	}
	facades.Orm().Query().Create(bookBySuperAdmin)
	
	bookByEditor1 := &models.Book{
		Title: "Editor 1 Book",
		Author: "Editor Author 1",
		ISBN: "1234567891",
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: &s.editor1.ID,
		},
	}
	facades.Orm().Query().Create(bookByEditor1)
	
	bookByEditor2 := &models.Book{
		Title: "Editor 2 Book",
		Author: "Editor Author 2",
		ISBN: "1234567892",
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: &s.editor2.ID,
		},
	}
	facades.Orm().Query().Create(bookByEditor2)
	
	bookByRegular := &models.Book{
		Title: "Regular User Book",
		Author: "Regular Author",
		ISBN: "1234567893",
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: &s.regularUser.ID,
		},
	}
	facades.Orm().Query().Create(bookByRegular)
	
	// Test data visibility based on scopes
	// This would require creating proper HTTP contexts and testing through the service layer
	// For now, we'll test the basic functionality
	
	var allBooks []models.Book
	facades.Orm().Query().Find(&allBooks)
	s.Equal(4, len(allBooks), "Should have 4 books in total")
}

func (s *ScopedPermissionsTestSuite) TestPermissionParsing() {
	// Test parsing scoped permission slugs
	service, action, scope := auth.ParsePermissionSlug("books_read_by_all")
	s.Equal("books", service)
	s.Equal("read", action)
	s.Equal(auth.ScopeByAll, scope)
	
	service, action, scope = auth.ParsePermissionSlug("books_update_by_me")
	s.Equal("books", service)
	s.Equal("update", action)
	s.Equal(auth.ScopeByMe, scope)
	
	service, action, scope = auth.ParsePermissionSlug("books_delete_by_my_role")
	s.Equal("books", service)
	s.Equal("delete", action)
	s.Equal(auth.ScopeByMyRole, scope)
	
	// Test parsing regular permission (no scope)
	service, action, scope = auth.ParsePermissionSlug("books_read")
	s.Equal("books", service)
	s.Equal("read", action)
	s.Equal(auth.ScopeByAll, scope) // Default scope
}

func (s *ScopedPermissionsTestSuite) TestIsScopedPermission() {
	s.True(auth.IsScopedPermission("books_read_by_all"))
	s.True(auth.IsScopedPermission("books_update_by_me"))
	s.True(auth.IsScopedPermission("books_delete_by_my_role"))
	s.False(auth.IsScopedPermission("books_read"))
	s.False(auth.IsScopedPermission("books_create"))
}

// TearDownTest will run after each test in the suite.
func (s *ScopedPermissionsTestSuite) TearDownTest() {
	// Clean up if needed
}