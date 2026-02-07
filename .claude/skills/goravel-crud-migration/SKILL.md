---
name: goravel-crud-migration
description: Create database migration and audit fields for a new Goravel entity. Use when setting up a new database table, adding fields, or starting CRUD scaffolding.
argument-hint: "[table_name]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel CRUD Migration

Create the database migration and audit fields for entity `$ARGUMENTS`.

## Step 1: Create the Migration

```bash
go run . artisan make:migration create_$ARGUMENTS_table
```

## Step 2: Edit the Migration

Open the generated file in `database/migrations/` and define the table schema.

### Field Type Reference

| Go/DB Type | Migration Method | Example |
|---|---|---|
| String | `table.String("field", 255)` | `table.String("title", 255)` |
| Text | `table.Text("field")` | `table.Text("description")` |
| Integer | `table.Integer("field")` | `table.Integer("quantity")` |
| Big Integer | `table.BigInteger("field")` | `table.BigInteger("amount")` |
| Float/Decimal | `table.Decimal("field", 10, 2)` | `table.Decimal("price", 10, 2)` |
| Boolean | `table.Boolean("field")` | `table.Boolean("is_active")` |
| Date | `table.Date("field")` | `table.Date("event_date")` |
| Timestamp | `table.TimestampTz("field")` | `table.TimestampTz("published_at")` |
| JSON | `table.Json("field")` | `table.Json("tags")` |
| Nullable | `.Nullable()` suffix | `table.Text("notes").Nullable()` |

### Required in Every Migration

- Add `table.SoftDeletesTz()` for soft deletes
- Foreign keys: `table.BigInteger("parent_id").Unsigned()` + foreign key constraint

## Step 3: Run the Migration

```bash
go run . artisan migrate
```

**IMPORTANT**: This runs migrations on the **dev database**. If you've only been running tests (which use testcontainers), the dev DB may not have the new tables. Always run `go run . artisan migrate` against the dev DB before attempting to use the UI.

## Step 4: Add Audit Fields

```bash
go run . artisan make:audit --table=$ARGUMENTS
```

### Post-Audit Fix

Review the generated audit migration and **remove** `table.Timestamp("deleted_at")` if present (SoftDeletesTz already creates it).

## Step 5: Run Audit Migration

```bash
go run . artisan migrate
```

## Step 6: Register in database/kernel.go

Add the new migration(s) to the `Migrations()` function in `database/kernel.go`:

```go
func (kernel Kernel) Migrations() []schema.Migration {
    return []schema.Migration{
        // ... existing migrations
        &migrations.CreateYourEntityTable{},
        &migrations.AddAuditFieldsToYourEntityTable{},
    }
}
```

## Verify

After registering in `database/kernel.go`, confirm the project compiles:

```bash
go build ./...
```

If there are import errors or type mismatches in the migration, fix them before proceeding.

## Post-Migration Checklist

- [ ] Migrations registered in `database/kernel.go`
- [ ] `go build ./...` passes
- [ ] `go run . artisan migrate` run against **dev database** (not just test containers)
- [ ] If adding FK to existing table, verify existing records handle nullable FK gracefully

**Common mistake**: Only running migrations in test (via testcontainers) but forgetting to run on the dev DB. This causes 500 errors when accessing the UI because the table doesn't exist in dev.

## Next Step

Run `/goravel-crud-model` to generate the model from the table.
