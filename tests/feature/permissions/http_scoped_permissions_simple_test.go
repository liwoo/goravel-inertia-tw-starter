package feature

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
)

func TestSimpleHTTPRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(facades.Route())
	defer server.Close()
	
	// Create HTTP client
	client := &http.Client{}
	
	// Make a simple request to root
	resp, err := client.Get(server.URL + "/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}