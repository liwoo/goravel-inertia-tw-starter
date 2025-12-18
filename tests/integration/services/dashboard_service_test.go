package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"starter-project/app/models"
	"starter-project/app/services"
	"starter-project/tests"
)

type DashboardServiceTestSuite struct {
	suite.Suite
	tests.TestCase

	dashboardService *services.DashboardService
	testUser         *models.User
}

func TestDashboardServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardServiceTestSuite))
}

func (s *DashboardServiceTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.dashboardService = services.NewDashboardService()

	// Clean existing test data
	facades.Orm().Query().Exec("DELETE FROM books")
	facades.Orm().Query().Exec("DELETE FROM users WHERE email LIKE '%@test.com'")

	// Create test user
	password, _ := facades.Hash().Make("password")
	s.testUser = &models.User{
		Name:     "Test User",
		Email:    fmt.Sprintf("dashboard-test-%d@test.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.testUser)
}

func (s *DashboardServiceTestSuite) TearDownTest() {
	// Clean up test data
	if s.testUser != nil && s.testUser.ID > 0 {
		facades.Orm().Query().Exec("DELETE FROM books WHERE created_by = ?", s.testUser.ID)
		facades.Orm().Query().Exec("DELETE FROM users WHERE id = ?", s.testUser.ID)
	}
}
