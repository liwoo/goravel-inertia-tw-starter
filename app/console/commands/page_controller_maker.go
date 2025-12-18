package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

// PageControllerMaker generates page controllers for Inertia.js views
type PageControllerMaker struct {
}

// Signature The name and signature of the console command.
func (receiver *PageControllerMaker) Signature() string {
	return "make:page-ctrl"
}

// Description The console command description.
func (receiver *PageControllerMaker) Description() string {
	return "Generate a page controller for Inertia.js views"
}

// Extend The console command extend.
func (receiver *PageControllerMaker) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:     "controller",
				Aliases:  []string{"c"},
				Usage:    "Controller name (e.g., Lender, Book, User)",
				Required: true,
			},
			&command.BoolFlag{
				Name:    "stats",
				Aliases: []string{"s"},
				Usage:   "Include statistics support",
				Value:   false,
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *PageControllerMaker) Handle(ctx console.Context) error {
	controllerName := ctx.Option("controller")
	if controllerName == "" {
		return fmt.Errorf("controller name is required. Use --controller=ControllerName")
	}

	enableStats := ctx.OptionBool("stats")

	// Normalize names
	controllerName = strings.Title(controllerName)
	resourceName := strings.ToLower(controllerName)
	pluralResource := receiver.pluralize(resourceName)

	ctx.Info(fmt.Sprintf("Generating page controller for %s...", controllerName))

	// Generate the controller file
	controllerContent := receiver.generateControllerContent(
		controllerName,
		resourceName,
		pluralResource,
		enableStats,
	)

	// Create controller directory
	controllerDir := filepath.Join("app", "http", "controllers", pluralResource)
	if err := os.MkdirAll(controllerDir, 0755); err != nil {
		return fmt.Errorf("failed to create controller directory: %v", err)
	}

	// Write controller file
	controllerPath := filepath.Join(controllerDir, fmt.Sprintf("%s_page_controller.go", pluralResource))
	if err := os.WriteFile(controllerPath, []byte(controllerContent), 0644); err != nil {
		return fmt.Errorf("failed to write controller file: %v", err)
	}

	ctx.Success(fmt.Sprintf("✓ Page controller generated successfully at: %s", controllerPath))
	ctx.Info("\nNext steps:")
	ctx.Info(fmt.Sprintf("  1. Create the Inertia page component at: resources/js/Pages/%s/Index.tsx", controllerName))
	ctx.Info(fmt.Sprintf("  2. Add route in routes/web.go:"))
	ctx.Info(fmt.Sprintf("     %sPageController := %s.New%sPageController()", pluralResource, pluralResource, controllerName))
	ctx.Info(fmt.Sprintf("     router.Get(\"/%s\", %sPageController.Index)", pluralResource, pluralResource))
	if enableStats {
		ctx.Info(fmt.Sprintf("  3. Implement Get%sStatistics() method in your service", controllerName))
	}

	return nil
}

func (receiver *PageControllerMaker) pluralize(word string) string {
	if strings.HasSuffix(word, "y") && len(word) > 1 && !receiver.isVowel(rune(word[len(word)-2])) {
		return word[:len(word)-1] + "ies"
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") {
		return word + "es"
	}
	return word + "s"
}

func (receiver *PageControllerMaker) isVowel(r rune) bool {
	vowels := "aeiouAEIOU"
	return strings.ContainsRune(vowels, r)
}

func (receiver *PageControllerMaker) generateControllerContent(
	controllerName, resourceName, pluralResource string,
	enableStats bool,
) string {
	// Generate auth service identifier
	authServiceName := fmt.Sprintf("Service%s", controllerName+"s")

	var statsConfig string
	if enableStats {
		statsConfig = fmt.Sprintf(`
			StatsEnabled:      true,
			StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
				stats, _ := %sService.Get%sStatistics()
				return stats
			},`, resourceName, controllerName)
	} else {
		statsConfig = `
			StatsEnabled:      false,`
	}

	return fmt.Sprintf(`package %s

import (
	"starter-project/app/auth"
	"starter-project/app/contracts"
	"starter-project/app/services"
)

// %sPageController handles the %s page
type %sPageController struct {
	*contracts.GenericPageController
	%sService *services.%sService
}

// New%sPageController creates a new %s page controller
func New%sPageController() *%sPageController {
	%sService := services.New%sService()

	return &%sPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "%s",
			PageComponent:     "%s/Index",
			Service:           %sService,
			ServiceIdentifier: auth.%s,%s
		}),
		%sService: %sService,
	}
}
`,
		pluralResource,
		controllerName, pluralResource,
		controllerName,
		resourceName, controllerName,
		controllerName, pluralResource,
		controllerName, controllerName,
		resourceName, controllerName,
		controllerName,
		pluralResource,
		controllerName,
		resourceName,
		authServiceName,
		statsConfig,
		resourceName, resourceName,
	)
}
