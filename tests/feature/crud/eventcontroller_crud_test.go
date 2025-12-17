package crud

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
	"github.com/goravel/framework/support/carbon"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/models"
	"smedi-sme-db/tests"
	"smedi-sme-db/tests/helpers"
)

type EventControllerCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestEventControllerCRUDTestSuite(t *testing.T) {
	suite.Run(t, &EventControllerCRUDTestSuite{})
}

func (s *EventControllerCRUDTestSuite) SetupSuite() {
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

func (s *EventControllerCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *EventControllerCRUDTestSuite) SetupTest() {
	// Clean database
	s.RefreshDatabase()

	// Create test user with permissions
	s.setupTestUser()
}

func (s *EventControllerCRUDTestSuite) TearDownTest() {
	// Clean up test data
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM events")
		orm.Query().Exec("DELETE FROM users WHERE email = 'eventtest@example.com'")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = 'event_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'events_%%'")
	}
}

func (s *EventControllerCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "Event Admin",
		Slug:  "event_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions
	permissions := []string{"create", "read", "update", "delete"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("events_%s", perm),
			Slug:  fmt.Sprintf("events_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("eventtest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login
	s.authCookie = s.loginUser("eventtest@example.com", "password")
	s.NotNil(s.authCookie)
}

func (s *EventControllerCRUDTestSuite) loginUser(email, password string) *http.Cookie {
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

func (s *EventControllerCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
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

func (s *EventControllerCRUDTestSuite) TestCreateEventController() {
	eventData := map[string]interface{}{
		"title":           "SME Networking Event",
		"description":     "A networking event for small and medium enterprises",
		"date":            "2025-12-15",
		"venue":           "Kigali Convention Centre",
		"partners":        []string{"MINICOM", "RDB", "PSF"},
		"district":        "Gasabo",
		"attending_smes":  []int{},
		"notes":           "Event includes lunch and refreshments",
	}

	resp, result := s.makeRequest("POST", "/api/events", eventData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
	s.Equal("SME Networking Event", data["title"])
	s.Equal("Kigali Convention Centre", data["venue"])
}

func (s *EventControllerCRUDTestSuite) TestCreateEventControllerValidation() {
	// Missing required fields
	eventData := map[string]interface{}{}

	resp, result := s.makeRequest("POST", "/api/events", eventData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *EventControllerCRUDTestSuite) TestGetEventController() {
	// Create an event first
	createdByInt := int(s.testUser.ID)
	notes := "Limited to 50 participants"
	event := &models.Event{
		Title:         "Business Training Workshop",
		Description:   "Workshop on business development",
		Date:          *carbon.NewDateTime(carbon.Parse("2025-11-20")),
		Venue:         "Serena Hotel",
		Partners:      []string{"BDF", "IFC"},
		District:      "Nyarugenge",
		AttendingSmes: []int{},
		Notes:         &notes,
		CreatedBy:     &createdByInt,
	}
	s.Nil(facades.Orm().Query().Create(event))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/events/%d", event.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(event.ID), data["id"].(float64))
	s.Equal("Business Training Workshop", data["title"])
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *EventControllerCRUDTestSuite) TestUpdateEventController() {
	// Create an event first
	createdByInt := int(s.testUser.ID)
	event := &models.Event{
		Title:         "Innovation Summit",
		Description:   "Annual innovation summit",
		Date:          *carbon.NewDateTime(carbon.Parse("2025-10-10")),
		Venue:         "Kigali Arena",
		Partners:      []string{"MINICT"},
		District:      "Gasabo",
		AttendingSmes: []int{},
		CreatedBy:     &createdByInt,
	}
	s.Nil(facades.Orm().Query().Create(event))

	updateData := map[string]interface{}{
		"title": "Innovation Summit 2025",
		"venue": "Kigali Convention Centre",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/events/%d", event.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("Innovation Summit 2025", data["title"])
	s.Equal("Kigali Convention Centre", data["venue"])
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *EventControllerCRUDTestSuite) TestDeleteEventController() {
	// Create an event first
	createdByInt := int(s.testUser.ID)
	event := &models.Event{
		Title:         "Expo 2025",
		Description:   "Business expo",
		Date:          *carbon.NewDateTime(carbon.Parse("2025-09-01")),
		Venue:         "IPRC Kigali",
		Partners:      []string{"PSF"},
		District:      "Kicukiro",
		AttendingSmes: []int{},
		CreatedBy:     &createdByInt,
	}
	s.Nil(facades.Orm().Query().Create(event))

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/events/%d", event.ID), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deletedEvent models.Event
	err := facades.Orm().Query().WithTrashed().Where("id", event.ID).First(&deletedEvent)
	s.Nil(err)
	s.NotNil(deletedEvent.DeletedAt)
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *EventControllerCRUDTestSuite) TestPagination() {
	// Create 25 events
	createdByInt := int(s.testUser.ID)
	for i := 1; i <= 25; i++ {
		event := &models.Event{
			Title:         fmt.Sprintf("Event %02d", i),
			Description:   fmt.Sprintf("Description for event %02d", i),
			Date:          *carbon.NewDateTime(carbon.Parse("2025-12-01")),
			Venue:         fmt.Sprintf("Venue %02d", i),
			Partners:      []string{"Partner1"},
			District:      "Gasabo",
			AttendingSmes: []int{},
			CreatedBy:     &createdByInt,
		}
		s.Nil(facades.Orm().Query().Create(event))
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/events?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(20, len(items))
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(25), pagination["total"])
	s.Equal(float64(2), pagination["last_page"])

	// Test custom page size
	resp, result = s.makeRequest("GET", "/api/events?page=1&pageSize=10", nil)
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

func (s *EventControllerCRUDTestSuite) TestSorting() {
	// Create events with different titles for sorting
	createdByInt := int(s.testUser.ID)
	events := []string{"Zebra Event", "Alpha Event", "Beta Event"}
	for _, title := range events {
		event := &models.Event{
			Title:         title,
			Description:   "Test event",
			Date:          *carbon.NewDateTime(carbon.Parse("2025-12-01")),
			Venue:         "Test Venue",
			Partners:      []string{"Test"},
			District:      "Gasabo",
			AttendingSmes: []int{},
			CreatedBy:     &createdByInt,
		}
		s.Nil(facades.Orm().Query().Create(event))
	}

	// Test sort ascending by title
	resp, result := s.makeRequest("GET", "/api/events?sort=title&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	// Verify first item is "Alpha Event"
	firstItem := items[0].(map[string]interface{})
	s.Equal("Alpha Event", firstItem["title"])
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *EventControllerCRUDTestSuite) TestSearch() {
	// Create searchable events
	createdByInt := int(s.testUser.ID)
	searchableEvent := &models.Event{
		Title:         "Searchable Tech Conference",
		Description:   "A searchable tech conference for startups",
		Date:          *carbon.NewDateTime(carbon.Parse("2025-11-15")),
		Venue:         "Impact Hub",
		Partners:      []string{"Google", "Microsoft"},
		District:      "Gasabo",
		AttendingSmes: []int{},
		CreatedBy:     &createdByInt,
	}
	s.Nil(facades.Orm().Query().Create(searchableEvent))

	otherEvent := &models.Event{
		Title:         "Business Meetup",
		Description:   "Regular business meetup",
		Date:          *carbon.NewDateTime(carbon.Parse("2025-11-20")),
		Venue:         "Norrsken House",
		Partners:      []string{"PSF"},
		District:      "Kicukiro",
		AttendingSmes: []int{},
		CreatedBy:     &createdByInt,
	}
	s.Nil(facades.Orm().Query().Create(otherEvent))

	// Search by query
	resp, result := s.makeRequest("GET", "/api/events/search?q=searchable", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 1)

	// Verify search results contain the searchable event
	found := false
	for _, item := range items {
		eventData := item.(map[string]interface{})
		if eventData["title"] == "Searchable Tech Conference" {
			found = true
			break
		}
	}
	s.True(found, "Should find the searchable event")
}
