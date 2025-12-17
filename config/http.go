package config

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin/render"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	ginfacades "github.com/goravel/gin/facades"
	inertia_http "smedi-sme-db/app/http/inertia"
)

func init() {
	config := facades.Config()
	config.Add("http", map[string]any{
		// HTTP Driver
		"default": "gin",
		// HTTP Drivers
		"drivers": map[string]any{
			"gin": map[string]any{
				// Optional, default is 4096 KB
				"body_limit":   4096,
				"header_limit": 4096,
				"route": func() (route.Route, error) {
					return ginfacades.Route("gin"), nil
				},
				// Custom template with marshal function for Inertia.js and Vite function
				"template": func() (render.HTMLRender, error) {
					// Create a new template. It's good practice to give it a name.
					tmpl := template.New("app")

					// Initialize the FuncMap
					funcMap := template.FuncMap{
						"marshal": func(v interface{}) template.JS {
							a, _ := json.Marshal(v)
							return template.JS(a)
						},
					}

					// Create Vite helper instance and add its functions
					facades.Log().Debug("[config/http.go] Creating ViteHelper for global template functions.")
					viteHelper := inertia_http.NewViteHelper()
					viteTemplateFuncs := viteHelper.CreateTemplateFuncs()
					if viteFunc, ok := viteTemplateFuncs["vite"]; ok {
						funcMap["vite"] = viteFunc
						facades.Log().Info("[config/http.go] 'vite' function successfully added to global FuncMap.")
					} else {
						facades.Log().Error("[config/http.go] 'vite' function NOT FOUND in ViteHelper's FuncMap. It will not be available globally.")
					}

					// Apply the FuncMap to the template
					tmpl = tmpl.Funcs(funcMap)

					// Parse templates from the views directory
					pattern := filepath.Join("resources", "views", "*.tmpl")

					// Parse the templates (this will now use the tmpl instance with funcMap applied)
					parsedTmpl, err := tmpl.ParseGlob(pattern)
					if err != nil {
						// In test environment, we can skip template parsing errors
						// Check both config and environment variable directly
						isTestMode := facades.Config().GetString("app.env") == "testing" ||
							os.Getenv("APP_ENV") == "testing"
						if isTestMode {
							log.Printf("[config/http.go] Skipping template parsing error in test mode: %v", err)
							// Return an empty template for tests
							return &render.HTMLProduction{Template: tmpl}, nil
						}
						log.Printf("[config/http.go] Error parsing templates: %v", err) // Use log for setup errors
						return nil, err
					}

					// Return a new HTML renderer with the parsed templates and options
					// Pass the already parsed template (parsedTmpl) to gin.NewTemplate
					return &render.HTMLProduction{Template: parsedTmpl}, nil
				},
			},
		},
		// HTTP URL
		"url": config.Env("APP_URL", "http://localhost"),
		// HTTP Host
		"host": config.Env("APP_HOST", "127.0.0.1"),
		// HTTP Port
		"port": config.Env("APP_PORT", "3500"),
		// HTTP Timeout, default is 3 seconds (use higher value for testing with bcrypt)
		"request_timeout": config.Env("HTTP_REQUEST_TIMEOUT", 3),

		// HTTP Client Configuration
		"client": map[string]any{
			"base_url":                config.GetString("HTTP_CLIENT_BASE_URL"),
			"timeout":                 config.GetDuration("HTTP_CLIENT_TIMEOUT"),
			"max_idle_conns":          config.GetInt("HTTP_CLIENT_MAX_IDLE_CONNS"),
			"max_idle_conns_per_host": config.GetInt("HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST"),
			"max_conns_per_host":      config.GetInt("HTTP_CLIENT_MAX_CONN_PER_HOST"),
			"idle_conn_timeout":       config.GetDuration("HTTP_CLIENT_IDLE_CONN_TIMEOUT"),
		},

		// HTTPS Configuration
		"tls": map[string]any{
			// HTTPS Host
			"host": config.Env("APP_HOST", "127.0.0.1"),
			// HTTPS Port
			"port": config.Env("APP_PORT", "3500"),
			// SSL Certificate, you can put the certificate in /public folder
			"ssl": map[string]any{
				// ca.pem
				"cert": "",
				// ca.key
				"key": "",
			},
		},
	})
}
