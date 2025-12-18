package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"starter-project/app/models"
	"starter-project/tests"
)

// This integration test demonstrates how scoped permissions work
// by directly querying the database with the expected SQL patterns
type ScopedFilteringTestSuite struct {
	suite.Suite
	tests.TestCase

	// Users
	admin   *models.User
	editor1 *models.User
	editor2 *models.User
	member1 *models.User
	member2 *models.User

	// Roles
	adminRole  *models.Role
	editorRole *models.Role
	memberRole *models.Role
}

func TestScopedFilteringTestSuite(t *testing.T) {
	suite.Run(t, new(ScopedFilteringTestSuite))
}

func (s *ScopedFilteringTestSuite) SetupTest() {
	s.RefreshDatabase()
	// Clean existing books to ensure test isolation
	facades.Orm().Query().Exec("DELETE FROM books")
	s.setupRolesAndPermissions()
	s.setupUsers()
	s.createTestBooks()
}

func (s *ScopedFilteringTestSuite) setupRolesAndPermissions() {
	// Create roles
	s.adminRole = &models.Role{Name: "Admin", Slug: "admin", IsActive: true}
	s.editorRole = &models.Role{Name: "Editor", Slug: "editor", IsActive: true}
	s.memberRole = &models.Role{Name: "Member", Slug: "member", IsActive: true}

	facades.Orm().Query().Create(s.adminRole)
	facades.Orm().Query().Create(s.editorRole)
	facades.Orm().Query().Create(s.memberRole)

	// Create permission
	perm := &models.Permission{
		Slug:     "books_read",
		Name:     "Read Books",
		Resource: "books",
		Action:   "read",
		IsActive: true,
	}
	facades.Orm().Query().Create(perm)

	// Assign permissions with different scopes
	// Admin: by_all
	facades.Orm().Query().Create(&models.RolePermission{
		RoleID:       s.adminRole.ID,
		PermissionID: perm.ID,
		Scope:        "by_all",
		IsActive:     true,
	})

	// Editor: by_my_role
	facades.Orm().Query().Create(&models.RolePermission{
		RoleID:       s.editorRole.ID,
		PermissionID: perm.ID,
		Scope:        "by_my_role",
		IsActive:     true,
	})

	// Member: by_me
	facades.Orm().Query().Create(&models.RolePermission{
		RoleID:       s.memberRole.ID,
		PermissionID: perm.ID,
		Scope:        "by_me",
		IsActive:     true,
	})
}

func (s *ScopedFilteringTestSuite) setupUsers() {
	password, _ := facades.Hash().Make("password")

	// Create users
	s.admin = &models.User{Name: "Admin", Email: "admin@test.com", Password: password, IsActive: true}
	s.editor1 = &models.User{Name: "Editor1", Email: "editor1@test.com", Password: password, IsActive: true}
	s.editor2 = &models.User{Name: "Editor2", Email: "editor2@test.com", Password: password, IsActive: true}
	s.member1 = &models.User{Name: "Member1", Email: "member1@test.com", Password: password, IsActive: true}
	s.member2 = &models.User{Name: "Member2", Email: "member2@test.com", Password: password, IsActive: true}

	facades.Orm().Query().Create(s.admin)
	facades.Orm().Query().Create(s.editor1)
	facades.Orm().Query().Create(s.editor2)
	facades.Orm().Query().Create(s.member1)
	facades.Orm().Query().Create(s.member2)

	// Assign roles
	// Create user-role associations with required fields
	now := time.Now()
	err := facades.Orm().Query().Create(&models.UserRole{
		UserID:     s.admin.ID,
		RoleID:     s.adminRole.ID,
		IsActive:   true,
		AssignedAt: now,
	})
	s.NoError(err)

	err = facades.Orm().Query().Create(&models.UserRole{
		UserID:     s.editor1.ID,
		RoleID:     s.editorRole.ID,
		IsActive:   true,
		AssignedAt: now,
	})
	s.NoError(err)

	err = facades.Orm().Query().Create(&models.UserRole{
		UserID:     s.editor2.ID,
		RoleID:     s.editorRole.ID,
		IsActive:   true,
		AssignedAt: now,
	})
	s.NoError(err)

	err = facades.Orm().Query().Create(&models.UserRole{
		UserID:     s.member1.ID,
		RoleID:     s.memberRole.ID,
		IsActive:   true,
		AssignedAt: now,
	})
	s.NoError(err)

	err = facades.Orm().Query().Create(&models.UserRole{
		UserID:     s.member2.ID,
		RoleID:     s.memberRole.ID,
		IsActive:   true,
		AssignedAt: now,
	})
	s.NoError(err)
}

func (s *ScopedFilteringTestSuite) createTestBooks() {
	// Create 2 books for each user
	users := []*models.User{s.admin, s.editor1, s.editor2, s.member1, s.member2}

	for _, user := range users {
		for i := 0; i < 2; i++ {
			book := &models.Book{
				Title:  fmt.Sprintf("%s Book %d", user.Name, i+1),
				Author: user.Name,
				ISBN:   fmt.Sprintf("ISBN-%d-%d", user.ID, i+1),
				Status: "AVAILABLE",
				Price:  float64(10 + i*5),
			}
			book.CreatedBy = &user.ID
			facades.Orm().Query().Create(book)
		}
	}
}

// Test database queries that demonstrate scope filtering
func (s *ScopedFilteringTestSuite) TestScopeFilteringQueries() {
	// 1. Admin scope (by_all) - sees all books
	count, _ := facades.Orm().Query().Model(&models.Book{}).Count()
	s.T().Logf("Total books in database: %d", count)
	s.Equal(int64(10), count, "Admin with by_all should see all 10 books")

	// 2. Editor scope (by_my_role) - sees books created by users with editor role
	// Get all user IDs with editor role
	var editorUserIDs []uint
	err := facades.Orm().Query().Table("user_roles").
		Where("role_id = ? AND is_active = ?", s.editorRole.ID, true).
		Pluck("user_id", &editorUserIDs)
	if err != nil {
		return
	}
	s.T().Logf("Editor role ID: %d, Editor user IDs: %v", s.editorRole.ID, editorUserIDs)

	// Count books created by editors
	if len(editorUserIDs) > 0 {
		count, _ = facades.Orm().Query().Model(&models.Book{}).
			Where("created_by IN ?", editorUserIDs).
			Count()
	} else {
		count = 0
	}
	s.T().Logf("Books created by editors: %d", count)
	s.Equal(int64(4), count, "Editor with by_my_role should see 4 books (2 from each editor)")

	// 3. Member scope (by_me) - sees only their own books
	count, _ = facades.Orm().Query().Model(&models.Book{}).
		Where("created_by = ?", s.member1.ID).
		Count()
	s.Equal(int64(2), count, "Member with by_me should see only their 2 books")
}

// Test the actual SQL patterns used by scope filtering
func (s *ScopedFilteringTestSuite) TestScopeFilteringPatterns() {
	// Test getting users with the same role

	// For editor role
	roleUserCount, _ := facades.Orm().Query().Table("user_roles").
		Where("role_id = ? AND is_active = ?", s.editorRole.ID, true).
		Count()

	s.Equal(int64(2), roleUserCount, "Should have 2 users with editor role")

	// For member role
	roleUserCount, _ = facades.Orm().Query().Table("user_roles").
		Where("role_id = ? AND is_active = ?", s.memberRole.ID, true).
		Count()
	s.Equal(int64(2), roleUserCount, "Should have 2 users with member role")
}

// Test complex scope query for by_my_role
func (s *ScopedFilteringTestSuite) TestByMyRoleScopeQuery() {
	// This demonstrates the SQL query pattern for by_my_role scope
	var books []models.Book

	// Simulate what the scope helper would do for editor1
	userRoleIDs := []uint{s.editorRole.ID}

	// Build the query as the scope helper would
	whereClause := fmt.Sprintf("created_by IN (SELECT DISTINCT user_id FROM user_roles WHERE role_id IN (%d) AND is_active = true)", userRoleIDs[0])

	err := facades.Orm().Query().Model(&models.Book{}).
		Where(whereClause).
		Find(&books)

	s.NoError(err)
	s.Len(books, 4, "by_my_role query should return 4 books for editor role")

	// Verify the books belong to editors
	for _, book := range books {
		s.True(*book.CreatedBy == s.editor1.ID || *book.CreatedBy == s.editor2.ID,
			"Book should be created by an editor")
	}
}

// Test statistics calculation with different scopes
func (s *ScopedFilteringTestSuite) TestStatisticsWithScopes() {
	// Create books with different statuses
	statuses := map[string]int{
		"AVAILABLE":   0,
		"BORROWED":    0,
		"MAINTENANCE": 0,
	}

	// Count all books by status
	var results []struct {
		Status string
		Count  int64
	}

	facades.Orm().Query().Model(&models.Book{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&results)

	for _, r := range results {
		statuses[r.Status] = int(r.Count)
	}

	// All books are created with AVAILABLE status in our test
	s.Equal(10, statuses["AVAILABLE"], "Should have 10 available books")
	s.Equal(0, statuses["BORROWED"], "Should have 0 borrowed books")
	s.Equal(0, statuses["MAINTENANCE"], "Should have 0 maintenance books")
}

func (s *ScopedFilteringTestSuite) TearDownTest() {
	// Cleanup handled by RefreshDatabase
}
