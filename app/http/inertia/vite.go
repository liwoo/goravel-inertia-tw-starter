package inertia

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/goravel/framework/facades"
)

// ManifestEntry represents a single entry in Vite's manifest.json
type ManifestEntry struct {
	File    string   `json:"file"`
	Src     string   `json:"src"`     // Original source file path, e.g., "resources/js/app.js"
	IsEntry bool     `json:"isEntry"` // True if this is an entry point defined in vite.config.js input
	CSS     []string `json:"css"`     // CSS files that this entry depends on
	Assets  []string `json:"assets"`  // Other static assets (images, fonts) imported by this entry
}

// ViteManifest represents the structure of Vite's manifest.json
// It's a map where keys are usually the source file paths (like ManifestEntry.Src)
type ViteManifest map[string]ManifestEntry

// ViteHelper handles Vite asset generation
type ViteHelper struct {
	isDev         bool
	manifest      ViteManifest
	publicPath    string // Base path for built assets, e.g., "/build"
	devServerURL  string // URL for the Vite development server, e.g., "http://localhost:5173"
}

// NewViteHelper creates a new ViteHelper instance.
// Configuration is loaded from facades.Config() with defaults.
func NewViteHelper() *ViteHelper {
	isDev := facades.Config().GetString("app.env", "production") != "production"
	manifestPath := facades.Config().GetString("vite.manifest_path", "public/build/manifest.json")
	publicPath := facades.Config().GetString("vite.public_path", "/build")
	devServerURL := facades.Config().GetString("vite.dev_server_url", "http://localhost:5173")

	vh := &ViteHelper{
		isDev:         isDev,
		publicPath:    publicPath,
		devServerURL:  devServerURL,
	}

	// Load manifest in production mode
	if !isDev {
		vh.loadManifest(manifestPath)
	} else {
		log.Println("VITE: Running in development mode. Assets will be served by Vite dev server.")
	}

	return vh
}

// loadManifest loads the Vite manifest file into the ViteHelper instance.
func (vh *ViteHelper) loadManifest(manifestPath string) {
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		log.Printf("VITE: Warning: Vite manifest not found at %s. Ensure 'npm run build' has been executed.", manifestPath)
		vh.manifest = make(ViteManifest) // Initialize to empty map to prevent nil panics
		return
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Printf("VITE: Error reading manifest file '%s': %v", manifestPath, err)
		vh.manifest = make(ViteManifest) // Initialize to empty map
		return
	}

	err = json.Unmarshal(data, &vh.manifest)
	if err != nil {
		log.Printf("VITE: Error parsing manifest file '%s': %v", manifestPath, err)
		vh.manifest = make(ViteManifest) // Initialize to empty map
		return
	}
	log.Println("VITE: Manifest loaded successfully for production.")
}

// Vite generates the appropriate HTML <script> or <link> tags for the given asset.
// assetPath should be the original source path of the asset (e.g., "resources/js/app.tsx").
func (vh *ViteHelper) Vite(assetPath string) template.HTML {
	if vh.isDev {
		return vh.devAsset(assetPath)
	}
	return vh.prodAsset(assetPath)
}

// devAsset generates HTML tags for assets in development mode (served by Vite dev server).
func (vh *ViteHelper) devAsset(assetPath string) template.HTML {
	url := fmt.Sprintf("%s/%s", vh.devServerURL, strings.TrimPrefix(assetPath, "/"))
	if strings.HasSuffix(assetPath, ".css") {
		return template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s">`, url))
	}
	return template.HTML(fmt.Sprintf(`<script type="module" src="%s"></script>`, url))
}

// prodAsset generates HTML tags for assets in production mode (using the manifest).
func (vh *ViteHelper) prodAsset(assetPath string) template.HTML {
	if vh.manifest == nil {
		log.Printf("VITE: Manifest is nil for asset '%s'. This shouldn't happen if initialized correctly.", assetPath)
		return template.HTML(fmt.Sprintf("<!-- Vite manifest not loaded for asset %s -->", assetPath))
	}

	manifestEntry, exists := vh.manifest[assetPath]
	if !exists {
		log.Printf("VITE: Asset '%s' not found in manifest. Ensure it's an entry point in vite.config.js and 'npm run build' has completed.", assetPath)
		return template.HTML(fmt.Sprintf("<!-- Asset %s not found in manifest -->", assetPath))
	}

	var html strings.Builder

	// Add CSS files listed as dependencies for this entry
	for _, cssFile := range manifestEntry.CSS {
		cssURL := fmt.Sprintf("%s/%s", strings.TrimSuffix(vh.publicPath, "/"), strings.TrimPrefix(cssFile, "/"))
		html.WriteString(fmt.Sprintf(`<link rel="stylesheet" href="%s">`, cssURL))
	}

	// Add the main asset file (JS or CSS)
	assetURL := fmt.Sprintf("%s/%s", strings.TrimSuffix(vh.publicPath, "/"), strings.TrimPrefix(manifestEntry.File, "/"))
	if strings.HasSuffix(manifestEntry.File, ".css") {
		html.WriteString(fmt.Sprintf(`<link rel="stylesheet" href="%s">`, assetURL))
	} else {
		html.WriteString(fmt.Sprintf(`<script type="module" src="%s"></script>`, assetURL))
	}

	return template.HTML(html.String())
}

// CreateTemplateFuncs creates a template.FuncMap containing the "vite" function.
func (vh *ViteHelper) CreateTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"vite": vh.Vite,
	}
}
