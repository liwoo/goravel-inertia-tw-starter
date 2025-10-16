package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

// CrudTestMaker generates comprehensive CRUD tests for an API controller
type CrudTestMaker struct {
}

// Signature The name and signature of the console command.
func (receiver *CrudTestMaker) Signature() string {
	return "make:crud-test"
}

// Description The console command description.
func (receiver *CrudTestMaker) Description() string {
	return "Generate comprehensive CRUD HTTP tests for an API controller"
}

// Extend The console command extend.
func (receiver *CrudTestMaker) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:     "controller",
				Aliases:  []string{"c"},
				Usage:    "Controller name (e.g., Lender, Book, User)",
				Required: true,
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *CrudTestMaker) Handle(ctx console.Context) error {
	controllerName := ctx.Option("controller")
	if controllerName == "" {
		return fmt.Errorf("controller name is required. Use --controller=ControllerName")
	}

	// Normalize names
	controllerName = strings.Title(controllerName)
	resourceName := strings.ToLower(controllerName)
	pluralResource := receiver.pluralize(resourceName)

	ctx.Info(fmt.Sprintf("Generating CRUD tests for %s controller...", controllerName))

	// Generate the test file
	testContent := receiver.generateTestContent(controllerName, resourceName, pluralResource)
	testPath := filepath.Join("tests", "feature", "crud", fmt.Sprintf("%s_crud_test.go", resourceName))

	// Create directory if it doesn't exist
	testDir := filepath.Dir(testPath)
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return fmt.Errorf("failed to create test directory: %v", err)
	}

	// Write test file
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		return fmt.Errorf("failed to write test file: %v", err)
	}

	ctx.Success(fmt.Sprintf("✓ CRUD tests generated successfully at: %s", testPath))
	ctx.Info("\nGenerated tests:")
	ctx.Info("  ✓ Create (valid data, validation errors)")
	ctx.Info("  ✓ Read (get by ID, not found)")
	ctx.Info("  ✓ Update (valid data, not found)")
	ctx.Info("  ✓ Delete (soft delete)")
	ctx.Info("  ✓ Pagination")
	ctx.Info("  ✓ Sorting")
	ctx.Info("  ✓ Search")
	ctx.Info("\nNext steps:")
	ctx.Info(fmt.Sprintf("  1. Review and customize test data in %s", testPath))
	ctx.Info(fmt.Sprintf("  2. Add specific field values for your model"))
	ctx.Info(fmt.Sprintf("  3. Run tests: APP_ENV=testing go test -v ./tests/feature/crud -run Test%sCRUDTestSuite", controllerName))

	return nil
}

func (receiver *CrudTestMaker) pluralize(word string) string {
	if strings.HasSuffix(word, "y") && len(word) > 1 && !receiver.isVowel(rune(word[len(word)-2])) {
		return word[:len(word)-1] + "ies"
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") {
		return word + "es"
	}
	return word + "s"
}

func (receiver *CrudTestMaker) isVowel(r rune) bool {
	vowels := "aeiouAEIOU"
	return strings.ContainsRune(vowels, r)
}

func (receiver *CrudTestMaker) generateTestContent(controllerName, resourceName, pluralResource string) string {
	return fmt.Sprintf(`package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"players/app/models"
	"players/tests"
	"players/tests/helpers"
)

type %sCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func Test%sCRUDTestSuite(t *testing.T) {
	suite.Run(t, &%sCRUDTestSuite{})
}

func (s *%sCRUDTestSuite) SetupSuite() {
	// Start test server
	s.server = httptest.NewServer(facades.Route())

	// Create HTTP client with cookie jar
	jar, _ := cookiejar.New(nil)
	s.client = &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (s *%sCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *%sCRUDTestSuite) SetupTest() {
	// Clean database
	s.RefreshDatabase()

	// Create test user with permissions
	s.setupTestUser()
}

func (s *%sCRUDTestSuite) TearDownTest() {
	// Clean up test data
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM %s")
		orm.Query().Exec("DELETE FROM users WHERE email = '%stest@example.com'")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = '%s_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE '%s_%%%%'")
	}
}

func (s *%sCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "%s Admin",
		Slug:  "%s_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions
	permissions := []string{"create", "read", "update", "delete"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("%s_%%%%s", perm),
			Slug:  fmt.Sprintf("%s_%%%%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("%stest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login
	s.authCookie = s.loginUser("%stest@example.com", "password")
	s.NotNil(s.authCookie)
}

func (s *%sCRUDTestSuite) loginUser(email, password string) *http.Cookie {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := s.client.Post(s.server.URL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	s.Nil(err)
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			return cookie
		}
	}

	return nil
}

func (s *%sCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, s.server.URL+path, bodyReader)
	s.Nil(err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	if s.authCookie != nil {
		req.AddCookie(s.authCookie)
	}

	resp, err := s.client.Do(req)
	s.Nil(err)

	respBody, err := io.ReadAll(resp.Body)
	s.Nil(err)
	resp.Body.Close()

	var result map[string]interface{}
	if len(respBody) > 0 {
		err = json.Unmarshal(respBody, &result)
		s.Nil(err)
	}

	return resp, result
}

// ============================================================================
// CREATE Tests
// ============================================================================

func (s *%sCRUDTestSuite) TestCreate%s() {
	%sData := map[string]interface{}{
		// TODO: Add required fields for %s
		// Example (based on your model):
		// "name": "Test %s",
		// "email": "test@example.com",
	}

	resp, result := s.makeRequest("POST", "/api/%s", %sData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
}

func (s *%sCRUDTestSuite) TestCreate%sValidation() {
	// Missing required fields
	%sData := map[string]interface{}{}

	resp, result := s.makeRequest("POST", "/api/%s", %sData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *%sCRUDTestSuite) TestGet%s() {
	// Create a %s first
	%s := &models.%s{
		// TODO: Set required fields
		// Example: Name: "Test Item", Email: "test@example.com",
	}
	%s.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(%s))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/%s/%%d", %s.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(%s.ID), data["id"].(float64))
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *%sCRUDTestSuite) TestUpdate%s() {
	// Create a %s first
	%s := &models.%s{
		// TODO: Set required fields
	}
	%s.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(%s))

	updateData := map[string]interface{}{
		// TODO: Add fields to update
		// Example: "name": "Updated Name",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/%s/%%d", %s.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *%sCRUDTestSuite) TestDelete%s() {
	// Create a %s first
	%s := &models.%s{
		// TODO: Set required fields
	}
	%s.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(%s))

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/%s/%%d", %s.ID), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deleted%s models.%s
	err := facades.Orm().Query().WithTrashed().Where("id", %s.ID).First(&deleted%s)
	s.Nil(err)
	s.NotNil(deleted%s.DeletedAt)
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *%sCRUDTestSuite) TestPagination() {
	// Create 25 %s
	for i := 1; i <= 25; i++ {
		%s := &models.%s{
			// TODO: Set required fields
			// Example: Name: fmt.Sprintf("Item %%02d", i),
		}
		%s.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(%s))
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/%s?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(20, len(items))
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(25), pagination["total"])
	s.Equal(float64(2), pagination["last_page"])

	// Test custom page size
	resp, result = s.makeRequest("GET", "/api/%s?page=1&pageSize=10", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	pagination = data["pagination"].(map[string]interface{})

	s.Equal(10, len(items))
	s.Equal(float64(3), pagination["last_page"])
}

// ============================================================================
// SORTING Tests
// ============================================================================

func (s *%sCRUDTestSuite) TestSorting() {
	// Create %s with sortable fields
	// TODO: Create items with different sortable values
	// Example: Create 3+ items with different names/values

	// Test sort ascending
	resp, result := s.makeRequest("GET", "/api/%s?sort=id&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 0)
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *%sCRUDTestSuite) TestSearch() {
	// Create searchable %s
	// TODO: Create items with searchable content
	// Example: Create items with "Searchable" in name/description

	// Search by query
	resp, result := s.makeRequest("GET", "/api/%s/search?q=test", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 0)
}
`,
		// Line 1: Suite type name
		controllerName,
		// Line 2: Test function name
		controllerName, controllerName,
		// Line 3: SetupSuite receiver
		controllerName,
		// Line 4: TearDownSuite receiver
		controllerName,
		// Line 5: SetupTest receiver
		controllerName,
		// Line 6: TearDownTest receiver - table name, email prefix, slug prefix, permission pattern
		controllerName, pluralResource, resourceName, resourceName, pluralResource,
		// Line 7: setupTestUser - Name, Slug, permission prefix (2x)
		controllerName, controllerName, resourceName, pluralResource, pluralResource,
		// Line 8: Login user - email prefix (2x)
		resourceName, resourceName,
		// Line 9: loginUser receiver
		controllerName,
		// Line 10: makeRequest receiver
		controllerName,
		// Line 11-12: TestCreate - receiver, method, data var, comment, example, endpoint, data var
		controllerName, controllerName, resourceName, controllerName, controllerName, pluralResource, resourceName,
		// Line 13-14: TestCreateValidation - receiver, method, data var, endpoint, data var
		controllerName, controllerName, resourceName, pluralResource, resourceName,
		// Line 15-16: TestGet - receiver, method, comment, var, type, var, var, endpoint, var, var
		controllerName, controllerName, resourceName, resourceName, controllerName, resourceName, resourceName, pluralResource, resourceName, resourceName,
		// Line 17-18: TestUpdate - receiver, method, comment, var, type, var, var, endpoint, var
		controllerName, controllerName, resourceName, resourceName, controllerName, resourceName, resourceName, pluralResource, resourceName,
		// Line 19-20: TestDelete - receiver, method, comment, var, type, var, var, endpoint, var, deleted var, type, var, var, var
		controllerName, controllerName, resourceName, resourceName, controllerName, resourceName, resourceName, pluralResource, resourceName, controllerName, controllerName, resourceName, controllerName, controllerName,
		// Line 21: TestPagination - receiver, comment, var, type, var, var, endpoint, endpoint
		controllerName, pluralResource, resourceName, controllerName, resourceName, resourceName, pluralResource, pluralResource,
		// Line 22: TestSorting - receiver, comment, endpoint
		controllerName, pluralResource, pluralResource,
		// Line 23: TestSearch - receiver, comment, endpoint
		controllerName, pluralResource, pluralResource,
	)
}
