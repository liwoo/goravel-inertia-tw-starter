package providers

import (
	// "html/template" // No longer needed here

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
	"github.com/petaki/inertia-go"
	inertia_http "players/app/http/inertia"
)

// InertiaServiceProvider is responsible for setting up the Inertia.js integration
type InertiaServiceProvider struct {
	// viteFunc any // No longer storing here
}

// Register registers the Inertia.js service
func (provider *InertiaServiceProvider) Register(app foundation.Application) {
	facades.Log().Info("InertiaServiceProvider: Register method STARTED")
	// Get application URL from config
	appURL := facades.Config().GetString("app.url", "http://localhost:8000")
	
	// Root template path
	rootTemplate := "resources/views/app.tmpl"
	
	// Asset version (can be updated based on your assets)
	version := "1.0.0"

	facades.Log().Debug("InertiaServiceProvider: Initializing Inertia manager")
	// Create Inertia manager
	inertiaManager := inertia.New(appURL, rootTemplate, version)

	facades.Log().Debug("InertiaServiceProvider: Creating ViteHelper instance for Inertia manager")
	// Create Vite helper instance
	viteHelper := inertia_http.NewViteHelper()
	// Get template funcs from ViteHelper
	templateFuncs := viteHelper.CreateTemplateFuncs()

	facades.Log().Debug("InertiaServiceProvider: Attempting to share 'vite' function with Inertia manager")
	// Set custom template functions for Inertia manager
	if viteFuncFromHelper, ok := templateFuncs["vite"]; ok {
		inertiaManager.ShareFunc("vite", viteFuncFromHelper)
		facades.Log().Info("InertiaServiceProvider: 'vite' function shared with Inertia manager")
	} else {
		facades.Log().Error("InertiaServiceProvider: 'vite' function NOT FOUND in ViteHelper's FuncMap for Inertia manager")
	}

	// Add a simple test function
	inertiaManager.ShareFunc("hello", func() string { return "Hello from Inertia FuncMap!" })
	facades.Log().Debug("InertiaServiceProvider: 'hello' function shared")
	
	// Share global view data
	inertiaManager.ShareViewData("appName", facades.Config().GetString("app.name", "Goravel"))
	facades.Log().Debug("InertiaServiceProvider: Global view data 'appName' shared")
	
	// Register the Inertia manager as a singleton
	facades.App().Singleton("inertia", func(app foundation.Application) (any, error) {
		return inertiaManager, nil
	})
	facades.Log().Info("InertiaServiceProvider: Inertia manager registered as singleton")
	facades.Log().Info("InertiaServiceProvider: Register method FINISHED")
}

// Boot performs any bootstrapping needed for Inertia.js
func (provider *InertiaServiceProvider) Boot(app foundation.Application) {
	facades.Log().Info("InertiaServiceProvider: Boot method STARTED - No view functions to register here anymore.")
	// Global view functions are now registered in config/http.go
	facades.Log().Info("InertiaServiceProvider: Boot method FINISHED")
}
