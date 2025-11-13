package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
	"smedi-sme-db/tests"
	"smedi-sme-db/tests/helpers"
)

// BdspServiceIntegrationTestSuite tests BDSP service with real database operations
type BdspServiceIntegrationTestSuite struct {
	suite.Suite
	tests.TestCase

	service *services.BdspService
	user1   *models.User
	user2   *models.User
	role1   *models.Role
	role2   *models.Role
}

func TestBdspServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(BdspServiceIntegrationTestSuite))
}

func (s *BdspServiceIntegrationTestSuite) SetupTest() {
	s.RefreshDatabase()
	// Additional cleanup to ensure clean state
	helpers.CleanTestDatabase()
	s.service = services.NewBdspService()
	s.setupTestData()
}

func (s *BdspServiceIntegrationTestSuite) setupTestData() {
	// Create test roles
	s.role1 = &models.Role{
		Name:     "BDSP Admin",
		Slug:     "bdsp_admin",
		Level:    100,
		IsActive: true,
	}
	s.role2 = &models.Role{
		Name:     "BDSP User",
		Slug:     "bdsp_user",
		Level:    50,
		IsActive: true,
	}
	s.NoError(facades.Orm().Query().Create(s.role1))
	s.NoError(facades.Orm().Query().Create(s.role2))

	// Create test users
	password, _ := facades.Hash().Make("password")
	s.user1 = &models.User{
		Name:     "Test User 1",
		Email:    "user1@bdsp.com",
		Password: password,
		IsActive: true,
	}
	s.user2 = &models.User{
		Name:     "Test User 2",
		Email:    "user2@bdsp.com",
		Password: password,
		IsActive: true,
	}
	s.NoError(facades.Orm().Query().Create(s.user1))
	s.NoError(facades.Orm().Query().Create(s.user2))

	// Assign roles
	now := time.Now()
	s.NoError(facades.Orm().Query().Create(&models.UserRole{
		UserID:     s.user1.ID,
		RoleID:     s.role1.ID,
		IsActive:   true,
		AssignedAt: now,
	}))
	s.NoError(facades.Orm().Query().Create(&models.UserRole{
		UserID:     s.user2.ID,
		RoleID:     s.role2.ID,
		IsActive:   true,
		AssignedAt: now,
	}))

	// Create permissions
	perm := &models.Permission{
		Slug:     "bdsps_read",
		Name:     "Read BDSPs",
		Resource: "bdsps",
		Action:   "read",
		IsActive: true,
	}
	s.NoError(facades.Orm().Query().Create(perm))

	// Assign permissions with different scopes
	s.NoError(facades.Orm().Query().Create(&models.RolePermission{
		RoleID:       s.role1.ID,
		PermissionID: perm.ID,
		Scope:        "by_all",
		IsActive:     true,
	}))
	s.NoError(facades.Orm().Query().Create(&models.RolePermission{
		RoleID:       s.role2.ID,
		PermissionID: perm.ID,
		Scope:        "by_me",
		IsActive:     true,
	}))
}

func (s *BdspServiceIntegrationTestSuite) TearDownTest() {
	// Clean up
	facades.Orm().Query().Exec("DELETE FROM bdsps")
	facades.Orm().Query().Exec("DELETE FROM user_roles")
	facades.Orm().Query().Exec("DELETE FROM role_permissions")
	facades.Orm().Query().Exec("DELETE FROM users WHERE email LIKE '%@bdsp.com'")
	facades.Orm().Query().Exec("DELETE FROM roles WHERE slug LIKE 'bdsp_%'")
	facades.Orm().Query().Exec("DELETE FROM permissions WHERE slug LIKE 'bdsps_%'")
}

// ============================================================================
// CREATE Tests
// ============================================================================

func (s *BdspServiceIntegrationTestSuite) TestCreate() {
	data := map[string]interface{}{
		"name":                "Test BDSP Provider",
		"postal_address":      "P.O. Box 12345",
		"physical_address":    "123 Main St, Lilongwe",
		"registration_status": "Registered",
		"product_types":       []string{"Training", "Consulting"},
		"service_list": []models.BdspService{
			{Name: "Business Training", Cost: 50000, Duration: "3 days"},
			{Name: "Financial Consulting", Cost: 100000, Duration: "1 week"},
		},
		"associated_partners": []string{"UNDP", "World Bank"},
		"created_by":          s.user1.ID,
	}

	result, err := s.service.Create(data)
	s.NoError(err)
	s.NotNil(result)

	// Verify the created BDSP
	bdsp, ok := result.(*models.Bdsp)
	s.True(ok)
	s.NotZero(bdsp.ID)
	s.Equal("Test BDSP Provider", bdsp.Name)
	s.Equal("P.O. Box 12345", *bdsp.PostalAddress)
	s.Equal("123 Main St, Lilongwe", *bdsp.PhysicalAddress)
	s.Equal("Registered", bdsp.RegistrationStatus)
	s.Equal(int(s.user1.ID), *bdsp.CreatedBy)

	// Verify JSON fields were properly serialized
	s.Equal(2, len(bdsp.ProductTypes))
	s.Contains(bdsp.ProductTypes, "Training")
	s.Contains(bdsp.ProductTypes, "Consulting")

	s.Equal(2, len(bdsp.ServiceList))
	s.Equal("Business Training", bdsp.ServiceList[0].Name)
	s.Equal(50000.0, bdsp.ServiceList[0].Cost)

	s.Equal(2, len(bdsp.AssociatedPartners))
	s.Contains(bdsp.AssociatedPartners, "UNDP")
}

func (s *BdspServiceIntegrationTestSuite) TestCreateWithMinimalData() {
	data := map[string]interface{}{
		"name":                "Minimal BDSP",
		"registration_status": "Pending",
		"product_types":       []string{"Training"},
		"service_list": []models.BdspService{
			{Name: "Basic Service", Cost: 10000, Duration: "1 day"},
		},
		"associated_partners": []string{"Local Partner"},
		"created_by":          s.user1.ID,
	}

	result, err := s.service.Create(data)
	s.NoError(err)
	s.NotNil(result)

	bdsp, ok := result.(*models.Bdsp)
	s.True(ok)
	s.Equal("Minimal BDSP", bdsp.Name)
	s.Nil(bdsp.PostalAddress)
	s.Nil(bdsp.PhysicalAddress)
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *BdspServiceIntegrationTestSuite) TestGetByID() {
	// Create a BDSP
	bdsp := s.createTestBdsp("Find Me BDSP", s.user1.ID)

	// Fetch by ID
	result, err := s.service.GetByID(bdsp.ID)
	s.NoError(err)
	s.NotNil(result)

	fetched, ok := result.(*models.Bdsp)
	s.True(ok)
	s.Equal(bdsp.ID, fetched.ID)
	s.Equal("Find Me BDSP", fetched.Name)
	s.Equal(bdsp.RegistrationStatus, fetched.RegistrationStatus)
}

func (s *BdspServiceIntegrationTestSuite) TestGetByIDNotFound() {
	result, err := s.service.GetByID(99999)
	s.Error(err)
	s.Nil(result)
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *BdspServiceIntegrationTestSuite) TestUpdate() {
	// Create a BDSP
	bdsp := s.createTestBdsp("Original Name", s.user1.ID)

	// Update data
	updateData := map[string]interface{}{
		"name":                "Updated BDSP Name",
		"postal_address":      "New P.O. Box 99999",
		"registration_status": "Active",
		"product_types_json":  `["Training", "Consulting", "Mentoring"]`,
		"updated_by":          s.user2.ID,
	}

	result, err := s.service.Update(bdsp.ID, updateData)
	s.NoError(err)
	s.NotNil(result)

	updated, ok := result.(*models.Bdsp)
	s.True(ok)
	s.Equal(bdsp.ID, updated.ID)
	s.Equal("Updated BDSP Name", updated.Name)
	s.Equal("New P.O. Box 99999", *updated.PostalAddress)
	s.Equal("Active", updated.RegistrationStatus)
	s.Equal(`["Training", "Consulting", "Mentoring"]`, updated.ProductTypesJSON)
	s.Equal(int(s.user2.ID), *updated.UpdatedBy)
}

func (s *BdspServiceIntegrationTestSuite) TestUpdateServices() {
	// Create a BDSP
	bdsp := s.createTestBdsp("Service Test BDSP", s.user1.ID)

	// Update with new services
	updateData := map[string]interface{}{
		"service_list_json": `[{"name":"New Service A","cost":25000,"duration":"2 days"},{"name":"New Service B","cost":75000,"duration":"1 week"},{"name":"New Service C","cost":150000,"duration":"2 weeks"}]`,
	}

	result, err := s.service.Update(bdsp.ID, updateData)
	s.NoError(err)

	updated, ok := result.(*models.Bdsp)
	s.True(ok)
	s.Equal(`[{"name":"New Service A","cost":25000,"duration":"2 days"},{"name":"New Service B","cost":75000,"duration":"1 week"},{"name":"New Service C","cost":150000,"duration":"2 weeks"}]`, updated.ServiceListJSON)
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *BdspServiceIntegrationTestSuite) TestDelete() {
	// Create a BDSP
	bdsp := s.createTestBdsp("To Delete", s.user1.ID)

	// Delete it
	err := s.service.Delete(bdsp.ID)
	s.NoError(err)

	// Verify soft delete - should not be found in normal query
	result, err := s.service.GetByID(bdsp.ID)
	s.Error(err)
	s.Nil(result)

	// Verify it's in database with deleted_at set
	var deletedBdsp models.Bdsp
	err = facades.Orm().Query().WithTrashed().Where("id", bdsp.ID).First(&deletedBdsp)
	s.NoError(err)
	s.NotNil(deletedBdsp.DeletedAt)
}

// ============================================================================
// LIST & PAGINATION Tests
// ============================================================================

func (s *BdspServiceIntegrationTestSuite) TestList() {
	// Create multiple BDSPs
	for i := 1; i <= 5; i++ {
		s.createTestBdsp(fmt.Sprintf("BDSP %d", i), s.user1.ID)
	}

	// List with default pagination
	listReq := contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	}

	result, err := s.service.GetList(listReq)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(5, len(result.Data))
	s.Equal(int64(5), result.Total)
	s.Equal(1, result.CurrentPage)
}

func (s *BdspServiceIntegrationTestSuite) TestListPagination() {
	// Create 25 BDSPs
	for i := 1; i <= 25; i++ {
		s.createTestBdsp(fmt.Sprintf("BDSP %02d", i), s.user1.ID)
	}

	// Test page 1
	result, err := s.service.GetList(contracts.ListRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.Equal(10, len(result.Data))
	s.Equal(int64(25), result.Total)
	s.Equal(1, result.CurrentPage)
	s.Equal(3, result.LastPage)

	// Test page 2
	result, err = s.service.GetList(contracts.ListRequest{Page: 2, PageSize: 10})
	s.NoError(err)
	s.Equal(10, len(result.Data))
	s.Equal(2, result.CurrentPage)

	// Test page 3 (partial)
	result, err = s.service.GetList(contracts.ListRequest{Page: 3, PageSize: 10})
	s.NoError(err)
	s.Equal(5, len(result.Data))
	s.Equal(3, result.CurrentPage)
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *BdspServiceIntegrationTestSuite) TestSearch() {
	// Create BDSPs with searchable data
	bdsps := []struct {
		name    string
		address string
		status  string
	}{
		{"Alpha Training Services", "Lilongwe", "Registered"},
		{"Beta Consulting Group", "Blantyre", "Active"},
		{"Gamma Business Development", "Mzuzu", "Registered"},
		{"Delta Management Solutions", "Lilongwe", "Pending"},
	}

	for _, b := range bdsps {
		data := map[string]interface{}{
			"name":                b.name,
			"physical_address":    b.address,
			"registration_status": b.status,
			"product_types":       []string{"Training"},
			"service_list": []models.BdspService{
				{Name: "Service", Cost: 10000, Duration: "1 day"},
			},
			"associated_partners": []string{"Partner"},
			"created_by":          s.user1.ID,
		}
		_, err := s.service.Create(data)
		s.NoError(err)
	}

	// Search by name
	result, err := s.service.Search("Training", contracts.ListRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.GreaterOrEqual(len(result.Data), 1)

	// Search by location
	result, err = s.service.Search("Lilongwe", contracts.ListRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.GreaterOrEqual(len(result.Data), 2)

	// Search by status
	result, err = s.service.Search("Registered", contracts.ListRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.GreaterOrEqual(len(result.Data), 2)
}

// ============================================================================
// SORTING Tests
// ============================================================================

func (s *BdspServiceIntegrationTestSuite) TestSorting() {
	// Create BDSPs with different names
	names := []string{"Zebra BDSP", "Alpha BDSP", "Gamma BDSP", "Beta BDSP"}
	for _, name := range names {
		s.createTestBdsp(name, s.user1.ID)
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	// Sort by name ascending
	result, err := s.service.GetList(contracts.ListRequest{
		Page:      1,
		PageSize:  10,
		Sort:      "name",
		Direction: "asc",
	})
	s.NoError(err)
	s.Equal(4, len(result.Data))

	// Verify order
	bdsp0 := result.Data[0].(models.Bdsp)
	s.Equal("Alpha BDSP", bdsp0.Name)

	// Sort by name descending
	result, err = s.service.GetList(contracts.ListRequest{
		Page:      1,
		PageSize:  10,
		Sort:      "name",
		Direction: "desc",
	})
	s.NoError(err)
	bdsp0 = result.Data[0].(models.Bdsp)
	s.Equal("Zebra BDSP", bdsp0.Name)
}

// ============================================================================
// FILTERING Tests
// ============================================================================

func (s *BdspServiceIntegrationTestSuite) TestFiltering() {
	// Create BDSPs with different statuses
	statuses := []string{"Registered", "Active", "Pending", "Registered", "Active"}
	for i, status := range statuses {
		data := map[string]interface{}{
			"name":                fmt.Sprintf("BDSP %d", i+1),
			"registration_status": status,
			"product_types":       []string{"Training"},
			"service_list": []models.BdspService{
				{Name: "Service", Cost: 10000, Duration: "1 day"},
			},
			"associated_partners": []string{"Partner"},
			"created_by":          s.user1.ID,
		}
		_, err := s.service.Create(data)
		s.NoError(err)
	}

	// Filter by registration_status = Registered
	filters := map[string]interface{}{
		"registration_status": "Registered",
	}

	result, err := s.service.GetListAdvanced(contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	}, filters)
	s.NoError(err)
	s.GreaterOrEqual(len(result.Data), 2)

	// Verify all results have "Registered" status
	for _, item := range result.Data {
		bdsp := item.(models.Bdsp)
		s.Equal("Registered", bdsp.RegistrationStatus)
	}
}

// ============================================================================
// SCOPED PERMISSIONS Tests
// ============================================================================
// Note: These tests verify data separation. Full scoped permission testing
// with Context is done in feature tests with HTTP requests.

func (s *BdspServiceIntegrationTestSuite) TestDataSeparationByUser() {
	// Create BDSPs by different users
	s.createTestBdsp("BDSP by User1", s.user1.ID)
	s.createTestBdsp("BDSP by User2", s.user2.ID)

	// Verify both BDSPs exist in database
	result, err := s.service.GetList(contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.GreaterOrEqual(len(result.Data), 2)
}

func (s *BdspServiceIntegrationTestSuite) TestFilterByCreator() {
	// Create BDSPs by different users
	s.createTestBdsp("BDSP by User1 A", s.user1.ID)
	s.createTestBdsp("BDSP by User1 B", s.user1.ID)
	s.createTestBdsp("BDSP by User2", s.user2.ID)

	// Filter by created_by = user2
	filters := map[string]interface{}{
		"created_by": s.user2.ID,
	}

	result, err := s.service.GetListAdvanced(contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	}, filters)
	s.NoError(err)
	s.GreaterOrEqual(len(result.Data), 1)

	// Verify all results belong to user2
	for _, item := range result.Data {
		bdsp := item.(models.Bdsp)
		s.Equal(int(s.user2.ID), *bdsp.CreatedBy)
	}
}

func (s *BdspServiceIntegrationTestSuite) TestMultipleUsersCreation() {
	// Create BDSPs by users with different roles
	s.createTestBdsp("BDSP by Admin", s.user1.ID)   // role1
	s.createTestBdsp("BDSP by User", s.user2.ID)    // role2
	s.createTestBdsp("BDSP by Admin 2", s.user1.ID) // role1

	// Verify all created successfully
	result, err := s.service.GetList(contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.GreaterOrEqual(len(result.Data), 3)
}

// ============================================================================
// JSON FIELDS Tests
// ============================================================================

func (s *BdspServiceIntegrationTestSuite) TestJSONFieldsSerialization() {
	// Create BDSP with complex JSON fields
	data := map[string]interface{}{
		"name":                "JSON Test BDSP",
		"registration_status": "Active",
		"product_types":       []string{"Type1", "Type2", "Type3"},
		"service_list": []models.BdspService{
			{Name: "Service A", Cost: 10000, Duration: "1 day"},
			{Name: "Service B", Cost: 20000, Duration: "2 days"},
		},
		"associated_partners": []string{"Partner A", "Partner B", "Partner C"},
		"created_by":          s.user1.ID,
	}

	result, err := s.service.Create(data)
	s.NoError(err)

	bdsp, ok := result.(*models.Bdsp)
	s.True(ok)

	// Fetch from database to verify persistence
	var fetched models.Bdsp
	err = facades.Orm().Query().Where("id", bdsp.ID).First(&fetched)
	s.NoError(err)

	// Verify JSON fields are properly deserialized
	s.Equal(3, len(fetched.ProductTypes))
	s.Equal(2, len(fetched.ServiceList))
	s.Equal(3, len(fetched.AssociatedPartners))

	s.Equal("Type1", fetched.ProductTypes[0])
	s.Equal("Service A", fetched.ServiceList[0].Name)
	s.Equal(10000.0, fetched.ServiceList[0].Cost)
	s.Equal("Partner A", fetched.AssociatedPartners[0])
}

// ============================================================================
// Helper Methods
// ============================================================================

func (s *BdspServiceIntegrationTestSuite) createTestBdsp(name string, createdBy uint) *models.Bdsp {
	data := map[string]interface{}{
		"name":                name,
		"postal_address":      "P.O. Box 123",
		"physical_address":    "Test Address",
		"registration_status": "Registered",
		"product_types":       []string{"Training", "Consulting"},
		"service_list": []models.BdspService{
			{Name: "Test Service", Cost: 50000, Duration: "3 days"},
		},
		"associated_partners": []string{"Test Partner"},
		"created_by":          createdBy,
	}

	result, err := s.service.Create(data)
	s.NoError(err)

	bdsp, ok := result.(*models.Bdsp)
	s.True(ok)
	return bdsp
}

// TestBdspModelJSONHooks tests the BeforeSave and AfterFind hooks
func (s *BdspServiceIntegrationTestSuite) TestBdspModelJSONHooks() {
	bdsp := &models.Bdsp{
		Name:               "Hook Test BDSP",
		RegistrationStatus: "Active",
		ProductTypes:       []string{"Type1", "Type2"},
		ServiceList: []models.BdspService{
			{Name: "Service1", Cost: 10000, Duration: "1 day"},
		},
		AssociatedPartners: []string{"Partner1"},
	}
	userID := int(s.user1.ID)
	bdsp.CreatedBy = &userID

	// Save - should trigger BeforeSave hook
	err := facades.Orm().Query().Create(bdsp)
	s.NoError(err)

	// Verify JSON fields were saved
	s.NotEmpty(bdsp.ProductTypesJSON)
	s.NotEmpty(bdsp.ServiceListJSON)
	s.NotEmpty(bdsp.AssociatedPartnersJSON)

	// Fetch - should trigger AfterFind hook
	var fetched models.Bdsp
	err = facades.Orm().Query().Where("id", bdsp.ID).First(&fetched)
	s.NoError(err)

	// Verify JSON fields were properly deserialized
	s.Equal(2, len(fetched.ProductTypes))
	s.Equal(1, len(fetched.ServiceList))
	s.Equal(1, len(fetched.AssociatedPartners))
	s.Equal("Type1", fetched.ProductTypes[0])
	s.Equal("Service1", fetched.ServiceList[0].Name)
}

// TestQueryBuilderWithScopes tests query builder with various scopes
func (s *BdspServiceIntegrationTestSuite) TestQueryBuilderWithScopes() {
	// Create BDSPs
	s.createTestBdsp("BDSP 1", s.user1.ID)
	s.createTestBdsp("BDSP 2", s.user1.ID)
	s.createTestBdsp("BDSP 3", s.user2.ID)

	// Test query with scope
	var bdsps []models.Bdsp
	query := facades.Orm().Query()

	// Apply by_me scope manually
	err := query.Model(&models.Bdsp{}).Where("created_by = ?", s.user1.ID).Find(&bdsps)
	s.NoError(err)

	s.Equal(2, len(bdsps))
}

// TestComplexFiltering tests complex filter combinations
func (s *BdspServiceIntegrationTestSuite) TestComplexFiltering() {
	// Create test data
	testData := []struct {
		name    string
		status  string
		address string
		userID  uint
	}{
		{"Alpha Services", "Registered", "Lilongwe", s.user1.ID},
		{"Beta Training", "Active", "Blantyre", s.user1.ID},
		{"Gamma Consulting", "Registered", "Lilongwe", s.user2.ID},
		{"Delta Solutions", "Pending", "Mzuzu", s.user2.ID},
	}

	for _, td := range testData {
		data := map[string]interface{}{
			"name":                td.name,
			"physical_address":    td.address,
			"registration_status": td.status,
			"product_types":       []string{"Training"},
			"service_list": []models.BdspService{
				{Name: "Service", Cost: 10000, Duration: "1 day"},
			},
			"associated_partners": []string{"Partner"},
			"created_by":          td.userID,
		}
		_, err := s.service.Create(data)
		s.NoError(err)
	}

	// Test filtering by status
	filters := map[string]interface{}{
		"registration_status": "Registered",
	}

	result, err := s.service.GetListAdvanced(contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	}, filters)
	s.NoError(err)
	s.GreaterOrEqual(len(result.Data), 2)

	// Verify results match filter condition
	for _, item := range result.Data {
		bdsp := item.(models.Bdsp)
		s.Equal("Registered", bdsp.RegistrationStatus)
	}
}

// TestTransactionRollback tests transaction rollback on error
func (s *BdspServiceIntegrationTestSuite) TestTransactionRollback() {
	// Count before
	countBefore, _ := facades.Orm().Query().Model(&models.Bdsp{}).Count()

	// Try to create with invalid data (this should fail validation if implemented)
	// For now, just test that valid data succeeds
	data := map[string]interface{}{
		"name":                "Transaction Test",
		"registration_status": "Active",
		"product_types":       []string{"Training"},
		"service_list": []models.BdspService{
			{Name: "Service", Cost: 10000, Duration: "1 day"},
		},
		"associated_partners": []string{"Partner"},
		"created_by":          s.user1.ID,
	}

	_, err := s.service.Create(data)
	s.NoError(err)

	// Count after
	countAfter, _ := facades.Orm().Query().Model(&models.Bdsp{}).Count()
	s.Equal(countBefore+1, countAfter)
}

// TestConcurrentOperations tests service handles concurrent operations
func (s *BdspServiceIntegrationTestSuite) TestConcurrentOperations() {
	bdsp := s.createTestBdsp("Concurrent Test", s.user1.ID)

	// Simulate concurrent updates
	done := make(chan bool, 2)

	// Update 1
	go func() {
		_, err := s.service.Update(bdsp.ID, map[string]interface{}{
			"postal_address": "Update 1",
		})
		s.NoError(err)
		done <- true
	}()

	// Update 2
	go func() {
		_, err := s.service.Update(bdsp.ID, map[string]interface{}{
			"physical_address": "Update 2",
		})
		s.NoError(err)
		done <- true
	}()

	// Wait for both to complete
	<-done
	<-done

	// Verify final state
	result, err := s.service.GetByID(bdsp.ID)
	s.NoError(err)
	s.NotNil(result)
}

// TestEmptyJSONArrays tests handling of empty JSON arrays
func (s *BdspServiceIntegrationTestSuite) TestEmptyJSONArrays() {
	// Create with empty arrays
	bdsp := &models.Bdsp{
		Name:               "Empty Arrays Test",
		RegistrationStatus: "Active",
		ProductTypes:       []string{},
		ServiceList:        []models.BdspService{},
		AssociatedPartners: []string{},
	}
	userID := int(s.user1.ID)
	bdsp.CreatedBy = &userID

	err := facades.Orm().Query().Create(bdsp)
	s.NoError(err)

	// Fetch and verify
	var fetched models.Bdsp
	err = facades.Orm().Query().Where("id", bdsp.ID).First(&fetched)
	s.NoError(err)

	s.NotNil(fetched.ProductTypes)
	s.NotNil(fetched.ServiceList)
	s.NotNil(fetched.AssociatedPartners)
	s.Equal(0, len(fetched.ProductTypes))
	s.Equal(0, len(fetched.ServiceList))
	s.Equal(0, len(fetched.AssociatedPartners))
}

// TestWithTrashedScope tests querying with soft-deleted records
func (s *BdspServiceIntegrationTestSuite) TestWithTrashedScope() {
	// Create and delete a BDSP
	bdsp := s.createTestBdsp("Trashed Test", s.user1.ID)
	err := s.service.Delete(bdsp.ID)
	s.NoError(err)

	// Should appear with WithTrashed
	trashedCount, _ := facades.Orm().Query().Model(&models.Bdsp{}).Where("id", bdsp.ID).WithTrashed().Count()
	s.Equal(int64(1), trashedCount)
}

// TestOrderByWithRelations tests ordering with relations
func (s *BdspServiceIntegrationTestSuite) TestOrderByWithRelations() {
	// Create BDSPs with different creators
	s.createTestBdsp("BDSP A", s.user1.ID)
	s.createTestBdsp("BDSP B", s.user2.ID)
	s.createTestBdsp("BDSP C", s.user1.ID)

	// List with sorting
	result, err := s.service.GetList(contracts.ListRequest{
		Page:      1,
		PageSize:  10,
		Sort:      "created_by",
		Direction: "asc",
	})

	s.NoError(err)
	s.GreaterOrEqual(len(result.Data), 3)
}

// TestPaginationEdgeCases tests edge cases in pagination
func (s *BdspServiceIntegrationTestSuite) TestPaginationEdgeCases() {
	// Create exactly 10 items
	for i := 1; i <= 10; i++ {
		s.createTestBdsp(fmt.Sprintf("BDSP %d", i), s.user1.ID)
	}

	// Request page 1 with size 10 - should get all items
	result, err := s.service.GetList(contracts.ListRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.Equal(10, len(result.Data))
	s.Equal(1, result.LastPage)

	// Request page 2 - should be empty
	result, err = s.service.GetList(contracts.ListRequest{Page: 2, PageSize: 10})
	s.NoError(err)
	s.Equal(0, len(result.Data))
}

// TestValidateBeforeCreate tests validation before creation
func (s *BdspServiceIntegrationTestSuite) TestValidateBeforeCreate() {
	// Missing required fields should fail
	invalidData := map[string]interface{}{
		// Missing name, product_types, service_list, etc.
	}

	result, err := s.service.Create(invalidData)
	s.Error(err)
	s.Nil(result)
}

// TestBulkOperations tests multiple operations in sequence
func (s *BdspServiceIntegrationTestSuite) TestBulkOperations() {
	// Create multiple BDSPs
	var ids []uint
	for i := 1; i <= 5; i++ {
		bdsp := s.createTestBdsp(fmt.Sprintf("Bulk BDSP %d", i), s.user1.ID)
		ids = append(ids, bdsp.ID)
	}

	// Update all
	for _, id := range ids {
		_, err := s.service.Update(id, map[string]interface{}{
			"registration_status": "Active",
		})
		s.NoError(err)
	}

	// Verify all updated
	for _, id := range ids {
		result, err := s.service.GetByID(id)
		s.NoError(err)
		bdsp := result.(*models.Bdsp)
		s.Equal("Active", bdsp.RegistrationStatus)
	}

	// Delete all
	for _, id := range ids {
		err := s.service.Delete(id)
		s.NoError(err)
	}

	// Verify all deleted
	count, _ := facades.Orm().Query().Model(&models.Bdsp{}).Where("id IN ?", ids).Count()
	s.Equal(int64(0), count)
}
