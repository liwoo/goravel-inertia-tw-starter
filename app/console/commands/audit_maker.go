package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type AuditMaker struct {
}

// Signature The name and signature of the console command.
func (r *AuditMaker) Signature() string {
	return "make:audit"
}

// Description The console command description.
func (r *AuditMaker) Description() string {
	return "Generate a migration to add audit fields to a table"
}

// Extend The console command extend.
func (r *AuditMaker) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "table",
				Aliases: []string{"t"},
				Usage:   "The table name to add audit fields to",
			},
			&command.BoolFlag{
				Name:    "without-created-by",
				Aliases: []string{"w"},
				Usage:   "Exclude created_by field (if table already has it)",
			},
		},
	}
}

// Handle Execute the console command.
func (r *AuditMaker) Handle(ctx console.Context) error {
	// Get the table name from argument or flag
	tableName := ctx.Argument(0)
	if tableName == "" {
		tableName = ctx.Option("table")
	}
	if tableName == "" {
		return fmt.Errorf("table name is required")
	}

	// Include created_by by default unless --without-created-by is specified
	withCreatedBy := !ctx.OptionBool("without-created-by")

	// Generate timestamp for migration filename
	timestamp := time.Now().Format("20060102150405")

	// Generate migration name
	migrationName := fmt.Sprintf("add_audit_fields_to_%s_table", tableName)
	filename := fmt.Sprintf("%s_%s.go", timestamp, migrationName)

	// Create migrations directory if it doesn't exist
	migrationsDir := filepath.Join("database", "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return err
	}

	// Generate full path
	fullPath := filepath.Join(migrationsDir, filename)

	// Generate migration content
	content := r.generateMigrationContent(timestamp, migrationName, tableName, withCreatedBy)

	// Write file
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return err
	}

	ctx.Success(fmt.Sprintf("Migration created: %s", fullPath))

	// Register migration in database/kernel.go
	structName := r.migrationNameToStructName(timestamp, migrationName)
	if err := r.registerMigration(structName); err != nil {
		ctx.Error(fmt.Sprintf("Failed to register migration: %v", err))
		ctx.NewLine()
		ctx.Comment("Please manually add the following line to database/kernel.go:")
		ctx.Line(fmt.Sprintf("  &migrations.%s{},", structName))
	} else {
		ctx.Success("Migration registered successfully")
	}
	ctx.NewLine()

	ctx.Info("Audit fields that will be added:")
	if withCreatedBy {
		ctx.Line("  - created_by (UnsignedBigInteger, nullable)")
	}
	ctx.Line("  - updated_by (UnsignedBigInteger, nullable)")
	ctx.Line("  - deleted_by (UnsignedBigInteger, nullable)")
	ctx.Line("  - deleted_at (Timestamp, nullable)")
	ctx.Line("  - ip_address (String, 45 chars, nullable)")
	ctx.Line("  - user_agent (Text, nullable)")
	ctx.NewLine()

	ctx.Comment("Run the migration:")
	ctx.Line("  go run . artisan migrate")

	return nil
}

// generateMigrationContent generates the migration file content
func (r *AuditMaker) generateMigrationContent(timestamp, migrationName, tableName string, withCreatedBy bool) string {
	// Generate struct name from migration name
	structName := r.migrationNameToStructName(timestamp, migrationName)

	// Generate singular form for comments
	singularName := strings.TrimSuffix(tableName, "s")
	if singularName == tableName {
		singularName = tableName
	}

	// Build created_by field section
	createdBySection := ""
	createdByForeignSection := ""
	createdByIndexSection := ""
	createdByDropForeignSection := ""
	createdByDropIndexSection := ""
	createdByDropColumnSection := ""

	if withCreatedBy {
		createdBySection = fmt.Sprintf("\t\ttable.UnsignedBigInteger(\"created_by\").Nullable().Comment(\"User who created this %s\")\n", singularName)
		createdByForeignSection = "\t\ttable.Foreign(\"created_by\").References(\"id\").On(\"users\")\n"
		createdByIndexSection = "\t\ttable.Index(\"created_by\")\n"
		createdByDropForeignSection = "\t\ttable.DropForeign(\"created_by\")\n"
		createdByDropIndexSection = "\t\ttable.DropIndex(\"created_by\")\n"
		createdByDropColumnSection = "\t\ttable.DropColumn(\"created_by\")\n"
	}

	return fmt.Sprintf(`package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type %s struct {
}

// Signature The unique signature for the migration.
func (r *%s) Signature() string {
	return "%s_%s"
}

// Up Run the migrations.
func (r *%s) Up() error {
	return facades.Schema().Table("%s", func(table schema.Blueprint) {
%s		table.UnsignedBigInteger("updated_by").Nullable().Comment("User who last updated this %s")
		table.UnsignedBigInteger("deleted_by").Nullable().Comment("User who deleted this %s")
		table.Timestamp("deleted_at").Nullable().Comment("Timestamp when this %s was deleted")
		table.String("ip_address", 45).Nullable().Comment("IP address from where the action was performed")
		table.Text("user_agent").Nullable().Comment("User agent from where the action was performed")

%s		table.Foreign("updated_by").References("id").On("users")
		table.Foreign("deleted_by").References("id").On("users")

%s		table.Index("updated_by")
		table.Index("deleted_by")
		table.Index("deleted_at")
	})
}

// Down Reverse the migrations.
func (r *%s) Down() error {
	return facades.Schema().Table("%s", func(table schema.Blueprint) {
%s		table.DropForeign("updated_by")
		table.DropForeign("deleted_by")

%s		table.DropIndex("updated_by")
		table.DropIndex("deleted_by")
		table.DropIndex("deleted_at")

%s		table.DropColumn("updated_by")
		table.DropColumn("deleted_by")
		table.DropColumn("deleted_at")
		table.DropColumn("ip_address")
		table.DropColumn("user_agent")
	})
}
`,
		structName,
		structName, timestamp, migrationName,
		structName, tableName,
		createdBySection, singularName, singularName, singularName,
		createdByForeignSection,
		createdByIndexSection,
		structName, tableName,
		createdByDropForeignSection,
		createdByDropIndexSection,
		createdByDropColumnSection,
	)
}

// migrationNameToStructName converts migration name to struct name
// e.g., "20060102150405_add_audit_fields_to_books_table" -> "M20060102150405AddAuditFieldsToBooksTable"
func (r *AuditMaker) migrationNameToStructName(timestamp, migrationName string) string {
	// Start with M + timestamp
	structName := "M" + timestamp

	// Split by underscores and capitalize each word
	parts := strings.Split(migrationName, "_")
	for _, part := range parts {
		if part != "" {
			structName += strings.Title(part)
		}
	}

	return structName
}

// registerMigration adds the migration to database/kernel.go
func (r *AuditMaker) registerMigration(structName string) error {
	kernelPath := filepath.Join("database", "kernel.go")

	// Read the kernel file
	content, err := os.ReadFile(kernelPath)
	if err != nil {
		return fmt.Errorf("failed to read kernel.go: %w", err)
	}

	kernelContent := string(content)

	// Check if already registered
	migrationLine := fmt.Sprintf("&migrations.%s{}", structName)
	if strings.Contains(kernelContent, migrationLine) {
		return fmt.Errorf("migration already registered")
	}

	// Find the Migrations() function and its closing brace
	// Look for "func (kernel Kernel) Migrations()" then find the closing of that function
	migrationsStart := strings.Index(kernelContent, "func (kernel Kernel) Migrations()")
	if migrationsStart == -1 {
		return fmt.Errorf("could not find Migrations() function in kernel.go")
	}

	// Find the end of the Migrations slice (before the closing }\n})
	// Look for the pattern after migrationsStart
	searchArea := kernelContent[migrationsStart:]
	insertionMarker := "\t}\n}"
	relativeInsertionIndex := strings.Index(searchArea, insertionMarker)
	if relativeInsertionIndex == -1 {
		return fmt.Errorf("could not find insertion point in Migrations() function")
	}

	// Calculate the absolute insertion index
	insertionIndex := migrationsStart + relativeInsertionIndex

	// Insert the new migration before the closing brace
	newMigration := fmt.Sprintf("\t\t&migrations.%s{},\n", structName)
	updatedContent := kernelContent[:insertionIndex] +
		newMigration +
		kernelContent[insertionIndex:]

	// Write the updated content back
	if err := os.WriteFile(kernelPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write kernel.go: %w", err)
	}

	return nil
}
