package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type MakeApp struct {
}

// Signature The name and signature of the console command.
func (r *MakeApp) Signature() string {
	return "make:app"
}

// Description The console command description.
func (r *MakeApp) Description() string {
	return "Customize the application name and package name for your fork"
}

// Extend The console command extend.
func (r *MakeApp) Extend() command.Extend {
	return command.Extend{
		Category: "make",
	}
}

// Handle Execute the console command.
func (r *MakeApp) Handle(ctx console.Context) error {
	ctx.Info("┌─────────────────────────────────────────────────┐")
	ctx.Info("│  🚀 Application Customization Tool             │")
	ctx.Info("└─────────────────────────────────────────────────┘")
	ctx.NewLine()

	// Get current package name from go.mod
	currentPackage, err := r.getCurrentPackageName()
	if err != nil {
		return fmt.Errorf("failed to read current package name: %w", err)
	}

	ctx.Info(fmt.Sprintf("Current package name: %s", currentPackage))
	ctx.NewLine()

	// Get new package name
	ctx.Info("Enter the new package name (e.g., 'myapp', 'github.com/username/myapp'):")
	ctx.Info("This will be used as the Go module name in go.mod")
	newPackage := r.promptForInput(ctx, "Package name")
	if newPackage == "" {
		ctx.Error("Package name cannot be empty")
		return fmt.Errorf("package name is required")
	}

	// Validate package name
	if !r.isValidPackageName(newPackage) {
		ctx.Error("Invalid package name. Use lowercase letters, numbers, hyphens, and slashes only.")
		return fmt.Errorf("invalid package name")
	}

	ctx.NewLine()

	// Get new app name
	ctx.Info("Enter the human-readable application name (e.g., 'My Awesome App'):")
	ctx.Info("This will be used in the .env file as APP_NAME")
	newAppName := r.promptForInput(ctx, "App name")
	if newAppName == "" {
		ctx.Error("App name cannot be empty")
		return fmt.Errorf("app name is required")
	}

	ctx.NewLine()

	// Confirm changes
	ctx.Warning("⚠️  This will update the following:")
	ctx.Line(fmt.Sprintf("  • go.mod module name: %s → %s", currentPackage, newPackage))
	ctx.Line(fmt.Sprintf("  • APP_NAME in .env files: → %s", newAppName))
	ctx.Line(fmt.Sprintf("  • All import statements referencing '%s'", currentPackage))
	ctx.NewLine()

	ctx.Info("Do you want to continue? (yes/no):")
	confirmation := r.promptForInput(ctx, "Confirm")
	if strings.ToLower(confirmation) != "yes" && strings.ToLower(confirmation) != "y" {
		ctx.Warning("Operation cancelled")
		return nil
	}

	ctx.NewLine()
	ctx.Info("🔄 Starting customization...")
	ctx.NewLine()

	// Perform updates
	if err := r.updateGoMod(ctx, currentPackage, newPackage); err != nil {
		return err
	}

	if err := r.updateImportStatements(ctx, currentPackage, newPackage); err != nil {
		return err
	}

	if err := r.updateEnvFiles(ctx, newAppName); err != nil {
		return err
	}

	// Run go mod tidy
	ctx.Info("Running 'go mod tidy'...")
	if err := r.runGoModTidy(ctx); err != nil {
		ctx.Warning("Warning: go mod tidy failed, but customization was successful")
		ctx.Warning(fmt.Sprintf("Error: %v", err))
	} else {
		ctx.Success("✓ go mod tidy completed")
	}

	ctx.NewLine()
	ctx.Success("✨ Application customization completed successfully!")
	ctx.NewLine()
	ctx.Info("Next steps:")
	ctx.Line("  1. Review the changes with 'git diff'")
	ctx.Line("  2. Update your .env file if needed")
	ctx.Line("  3. Test your application: 'go run .'")
	ctx.Line("  4. Commit the changes: 'git add . && git commit -m \"Customize app name and package\"'")

	return nil
}

// getCurrentPackageName reads the current package name from go.mod
func (r *MakeApp) getCurrentPackageName() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	// Parse module name from go.mod
	re := regexp.MustCompile(`(?m)^module\s+(.+)$`)
	matches := re.FindSubmatch(data)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find module name in go.mod")
	}

	return strings.TrimSpace(string(matches[1])), nil
}

// promptForInput prompts the user for input
func (r *MakeApp) promptForInput(ctx console.Context, prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// isValidPackageName validates the package name format
func (r *MakeApp) isValidPackageName(name string) bool {
	// Allow lowercase letters, numbers, hyphens, underscores, dots, and slashes
	re := regexp.MustCompile(`^[a-z0-9._/-]+$`)
	return re.MatchString(name)
}

// updateGoMod updates the module name in go.mod
func (r *MakeApp) updateGoMod(ctx console.Context, oldPackage, newPackage string) error {
	ctx.Info("Updating go.mod...")

	data, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	// Replace module name
	re := regexp.MustCompile(`(?m)^module\s+` + regexp.QuoteMeta(oldPackage) + `$`)
	newData := re.ReplaceAll(data, []byte("module "+newPackage))

	if err := os.WriteFile("go.mod", newData, 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	ctx.Success("✓ Updated go.mod")
	return nil
}

// updateImportStatements updates all import statements in Go files
func (r *MakeApp) updateImportStatements(ctx console.Context, oldPackage, newPackage string) error {
	ctx.Info("Updating import statements in Go files...")

	// Find all .go files
	var goFiles []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor and hidden directories
		if info.IsDir() {
			name := info.Name()
			if name == "vendor" || name == ".git" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .go files
		if filepath.Ext(path) == ".go" {
			goFiles = append(goFiles, path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	updatedCount := 0
	for _, file := range goFiles {
		updated, err := r.updateImportsInFile(file, oldPackage, newPackage)
		if err != nil {
			ctx.Warning(fmt.Sprintf("Warning: failed to update %s: %v", file, err))
			continue
		}
		if updated {
			updatedCount++
		}
	}

	ctx.Success(fmt.Sprintf("✓ Updated imports in %d files", updatedCount))
	return nil
}

// updateImportsInFile updates imports in a single file
func (r *MakeApp) updateImportsInFile(filename, oldPackage, newPackage string) (bool, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}

	// Replace import statements
	// Match both single and double quoted imports
	oldPackageEscaped := regexp.QuoteMeta(oldPackage)

	// Replace in import statements: "oldPackage/..." or `oldPackage/...`
	re1 := regexp.MustCompile(`"` + oldPackageEscaped + `(/[^"]*)?"|` + "`" + oldPackageEscaped + `(/[^` + "`" + `]*)?` + "`")

	newData := re1.ReplaceAllFunc(data, func(match []byte) []byte {
		matchStr := string(match)
		quote := matchStr[0:1]

		// Extract the path after the old package
		path := strings.TrimPrefix(matchStr, quote+oldPackage)
		path = strings.TrimSuffix(path, quote)

		return []byte(quote + newPackage + path + quote)
	})

	// Check if anything changed
	if string(data) == string(newData) {
		return false, nil
	}

	if err := os.WriteFile(filename, newData, 0644); err != nil {
		return false, err
	}

	return true, nil
}

// updateEnvFiles updates APP_NAME in .env files
func (r *MakeApp) updateEnvFiles(ctx console.Context, newAppName string) error {
	ctx.Info("Updating .env files...")

	envFiles := []string{".env", ".env.example", ".env.testing"}
	updatedCount := 0

	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); os.IsNotExist(err) {
			continue
		}

		if err := r.updateEnvFile(envFile, newAppName); err != nil {
			ctx.Warning(fmt.Sprintf("Warning: failed to update %s: %v", envFile, err))
			continue
		}
		updatedCount++
	}

	if updatedCount > 0 {
		ctx.Success(fmt.Sprintf("✓ Updated %d .env file(s)", updatedCount))
	} else {
		ctx.Warning("⚠ No .env files found to update")
	}

	return nil
}

// updateEnvFile updates APP_NAME in a single .env file
func (r *MakeApp) updateEnvFile(filename, newAppName string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Replace APP_NAME value
	re := regexp.MustCompile(`(?m)^APP_NAME=.*$`)
	newData := re.ReplaceAll(data, []byte("APP_NAME="+newAppName))

	if err := os.WriteFile(filename, newData, 0644); err != nil {
		return err
	}

	return nil
}

// runGoModTidy runs 'go mod tidy' command
func (r *MakeApp) runGoModTidy(ctx console.Context) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
