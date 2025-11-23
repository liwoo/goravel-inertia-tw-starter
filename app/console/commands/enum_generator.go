package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

// EnumGenerator generates TypeScript enum types from Go enums
type EnumGenerator struct {
}

// Signature The name and signature of the console command.
func (receiver *EnumGenerator) Signature() string {
	return "make:ts-enums"
}

// Description The console command description.
func (receiver *EnumGenerator) Description() string {
	return "Generate TypeScript enum types from Go enum definitions"
}

// Extend The console command extend.
func (receiver *EnumGenerator) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "source",
				Aliases: []string{"s"},
				Usage:   "Source directory to scan for Go enums (default: app/http/requests)",
				Value:   "app/http/requests",
			},
			&command.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output directory for TypeScript files (default: resources/js/types)",
				Value:   "resources/js/types",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *EnumGenerator) Handle(ctx console.Context) error {
	sourceDir := ctx.Option("source")
	outputDir := ctx.Option("output")

	if sourceDir == "" {
		sourceDir = "app/http/requests"
	}
	if outputDir == "" {
		outputDir = "resources/js/types"
	}

	ctx.Info(fmt.Sprintf("Scanning %s for Go enum types...", sourceDir))

	// Use the enhanced type generator
	gen := &EnhancedTypeGenerator{}
	enums, err := gen.ScanDirectoryForEnums(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to scan for enums: %w", err)
	}

	if len(enums) == 0 {
		ctx.Warning("No enum types found in the source directory")
		return nil
	}

	ctx.Info(fmt.Sprintf("Found %d enum type(s)", len(enums)))

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate a TypeScript file for each enum
	for _, enum := range enums {
		fileName := strings.ToLower(gen.toSnakeCase(enum.TypeName)) + ".ts"
		filePath := filepath.Join(outputDir, fileName)

		tsContent := gen.GenerateEnumTypeScript(enum)

		if err := os.WriteFile(filePath, []byte(tsContent), 0644); err != nil {
			ctx.Error(fmt.Sprintf("Failed to write %s: %v", fileName, err))
			continue
		}

		ctx.Success(fmt.Sprintf("✓ Generated %s (%d values)", fileName, len(enum.Values)))
	}

	ctx.Info("\n" + strings.Repeat("=", 60))
	ctx.Info("✅ TypeScript enum files generated successfully!")
	ctx.Info(strings.Repeat("=", 60))
	ctx.Info("\nYou can now import these enums in your TypeScript code:")
	ctx.Info("  import { GenderType, GENDER_TYPE_OPTIONS } from '@/types/gender_type';")
	ctx.Info("\nThe generated files include:")
	ctx.Info("  - Union type definitions (e.g., type GenderType = 'MALE' | 'FEMALE')")
	ctx.Info("  - Option arrays for forms (e.g., GENDER_TYPE_OPTIONS)")
	ctx.Info(strings.Repeat("=", 60))

	return nil
}
