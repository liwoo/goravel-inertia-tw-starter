package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type SwaggerGenerator struct {
}

// Signature The name and signature of the console command.
func (r *SwaggerGenerator) Signature() string {
	return "swagger:generate"
}

// Description The console command description.
func (r *SwaggerGenerator) Description() string {
	return "Generate Swagger documentation from code"
}

// Extend The console command extend.
func (r *SwaggerGenerator) Extend() command.Extend {
	return command.Extend{
		Category: "swagger",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "open",
				Aliases: []string{"o"},
				Usage:   "Open Swagger UI in browser after generation",
			},
		},
	}
}

// Handle Execute the console command.
func (r *SwaggerGenerator) Handle(ctx console.Context) error {
	ctx.Info("Generating Swagger documentation...")
	ctx.NewLine()

	// Check if swag is installed
	if !r.isSwagInstalled() {
		ctx.Warning("Installing swag CLI tool...")
		if err := r.installSwag(ctx); err != nil {
			return fmt.Errorf("failed to install swag: %v", err)
		}
		ctx.Success("✓ Swag CLI installed")
		ctx.NewLine()
	}

	// Create docs directory if it doesn't exist
	docsDir := filepath.Join("docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return fmt.Errorf("failed to create docs directory: %v", err)
	}

	// Run swag init
	ctx.Info("Running swag init...")
	cmd := exec.Command("swag", "init",
		"--generalInfo", "main.go",
		"--dir", ".",
		"--parseDependency",
		"--parseInternal",
		"--output", "docs")

	output, err := cmd.CombinedOutput()
	if err != nil {
		ctx.Error(string(output))
		return fmt.Errorf("failed to generate swagger docs: %v", err)
	}

	ctx.Success("✓ Swagger documentation generated successfully!")
	ctx.NewLine()

	// Print generated files
	ctx.Info("Generated files:")
	ctx.Line("  - docs/docs.go")
	ctx.Line("  - docs/swagger.json")
	ctx.Line("  - docs/swagger.yaml")
	ctx.NewLine()

	// Print access URL
	ctx.Info("Access Swagger UI at:")
	ctx.Line("  http://localhost:3000/api/swagger/index.html")
	ctx.NewLine()

	// Print helpful tips
	ctx.Comment("Tips:")
	ctx.Line("  - Swagger will auto-discover your models via json tags")
	ctx.Line("  - Update main.go to change API title, version, or description")
	ctx.Line("  - Re-run this command after making changes to your API")
	ctx.NewLine()

	return nil
}

// isSwagInstalled checks if swag CLI is installed
func (r *SwaggerGenerator) isSwagInstalled() bool {
	cmd := exec.Command("swag", "--version")
	return cmd.Run() == nil
}

// installSwag installs the swag CLI tool
func (r *SwaggerGenerator) installSwag(ctx console.Context) error {
	cmd := exec.Command("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
