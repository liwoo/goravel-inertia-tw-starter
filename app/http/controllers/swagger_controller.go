package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type SwaggerController struct {
}

func NewSwaggerController() *SwaggerController {
	return &SwaggerController{}
}

// ServeSwagger serves Swagger UI and documentation files
func (c *SwaggerController) ServeSwagger(ctx http.Context) http.Response {
	// Get the path after /swagger/
	path := ctx.Request().Route("any")

	// Default to index.html if path is empty
	if path == "" || path == "/" {
		path = "index.html"
	}

	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")

	// Serve swagger files from docs directory
	if strings.HasPrefix(path, "swagger.") {
		return c.serveSwaggerFile(ctx, path)
	}

	// Serve Swagger UI HTML
	if path == "index.html" || path == "" {
		return c.serveSwaggerUI(ctx)
	}

	return ctx.Response().String(404, "Not found")
}

// serveSwaggerFile serves swagger.json or swagger.yaml
func (c *SwaggerController) serveSwaggerFile(ctx http.Context, filename string) http.Response {
	filePath := filepath.Join("docs", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ctx.Response().String(404, "Swagger documentation not found. Run: go run . artisan swagger:generate")
	}

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ctx.Response().String(500, "Error reading swagger file")
	}

	// Set appropriate content type
	contentType := "application/json"
	if strings.HasSuffix(filename, ".yaml") {
		contentType = "text/yaml"
	}

	ctx.Response().Header("Content-Type", contentType)
	return ctx.Response().String(200, string(content))
}

// serveSwaggerUI serves the Swagger UI HTML page
func (c *SwaggerController) serveSwaggerUI(ctx http.Context) http.Response {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Library Management API - Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "/api/swagger/swagger.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

	ctx.Response().Header("Content-Type", "text/html; charset=utf-8")
	return ctx.Response().String(200, html)
}

// ServeSwaggerJSON is a convenience method to directly serve the JSON file
func (c *SwaggerController) ServeSwaggerJSON(ctx http.Context) http.Response {
	return c.serveSwaggerFile(ctx, "swagger.json")
}

// ServeSwaggerYAML is a convenience method to directly serve the YAML file
func (c *SwaggerController) ServeSwaggerYAML(ctx http.Context) http.Response {
	return c.serveSwaggerFile(ctx, "swagger.yaml")
}

// ServeSwaggerDocs serves the docs.go file for debugging
func (c *SwaggerController) ServeSwaggerDocs(ctx http.Context) http.Response {
	filePath := filepath.Join("docs", "docs.go")

	file, err := os.Open(filePath)
	if err != nil {
		return ctx.Response().String(404, "docs.go not found")
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return ctx.Response().String(500, "Error reading docs.go")
	}

	ctx.Response().Header("Content-Type", "text/plain; charset=utf-8")
	return ctx.Response().String(200, string(content))
}
