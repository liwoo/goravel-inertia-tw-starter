---
name: fake-data
description: Create a Go database seeder for a Goravel entity with 25+ realistic records. Covers all status/enum values, diverse data for sorting/pagination testing, and proper nullable field handling.
argument-hint: "[EntityName]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Fake Data Seeder

**Agent**: Chikondi Banda — Senior Backend Engineer

Create a database seeder for `$ARGUMENTS`.

## Prerequisites

- Model exists in `app/models/`
- Migration has been run on dev DB (`go run . artisan migrate`)

## Step 1: Read the Model

Read `app/models/<entity_name>.go` to understand:
- All fields and their types
- Which fields are required vs nullable (`*string`, `*time.Time`, etc.)
- Enum/status values (check request validators in `app/http/requests/`)
- Relationships (FK fields)

## Step 2: Create the Seeder File

Create `database/seeders/<entity_name>_seeder.go`:

```go
package seeders

import (
	"github.com/goravel/framework/facades"
	"books-database/app/models"
	"time"
)

type EntitySeeder struct {
}

// Signature The name and signature of the seeder.
func (s *EntitySeeder) Signature() string {
	return "EntitySeeder"
}

// Run executes the seeder logic.
func (s *EntitySeeder) Run() error {
	// Check if table already has data (idempotent)
	var count int64
	count, err := facades.Orm().Query().Model(&models.Entity{}).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		facades.Log().Infof("EntitySeeder: Skipping - table already has %d records", count)
		return nil
	}

	entities := []models.Entity{
		// ... 25-50 records ...
	}

	for _, entity := range entities {
		if err := facades.Orm().Query().Create(&entity); err != nil {
			return err
		}
	}

	return nil
}
```

## Data Generation Rules

### Minimum Records: 25

Generate **at least 25 records** to ensure meaningful testing of:
- Pagination (default page size is typically 10-20)
- Sorting (needs diverse values to verify order)
- Search (needs varied text to test matching)
- Filters (needs records in each status/category)

### Status/Enum Distribution

Distribute records across **all** enum values:

```
// For a 2-value enum (ACTIVE/INACTIVE) with 30 records:
// ~20 ACTIVE, ~10 INACTIVE

// For a 4-value enum (AVAILABLE/BORROWED/MAINTENANCE/RESERVED) with 50 records:
// ~20 AVAILABLE, ~12 BORROWED, ~10 MAINTENANCE, ~8 RESERVED
```

**Every enum value must have at least 3 records** so filter tabs show meaningful results.

### Data Diversity for Sorting

Ensure fields have diverse values that produce a clear sort order:

```go
// GOOD — names that sort differently
{FirstName: "Alice", ...},
{FirstName: "Zara", ...},
{FirstName: "Michael", ...},

// BAD — all similar names
{FirstName: "John", ...},
{FirstName: "John", ...},
{FirstName: "Jane", ...},
```

For **date fields**, spread across multiple years:
```go
parseDate("2020-03-15"),
parseDate("2022-07-22"),
parseDate("2024-01-10"),
```

For **numeric fields** (price, amount), use varied values:
```go
Price: 9.99,
Price: 24.50,
Price: 149.99,
```

### Nullable Fields

Mix `nil` and populated values for optional fields:

```go
// ~70% populated, ~30% nil for optional fields
{Bio: strPtr("A detailed biography..."), Email: strPtr("alice@example.com")},
{Bio: nil, Email: strPtr("bob@example.com")},
{Bio: strPtr("Another bio"), Email: nil},
{Bio: nil, Email: nil},
```

### Helper Functions

```go
// String pointer helper
func strPtr(s string) *string { return &s }

// Date pointer helper (reuse if already defined in another seeder)
func parseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	return &t
}
```

**Important**: Check if `parseDate` or `strPtr` already exist in another seeder file before redefining them (Go won't allow duplicate function names in the same package).

### FK Relationships

If the entity references another entity via FK:

```go
// Option A: Set FK to known seeded IDs (if parent seeder runs first)
{AuthorID: uintPtr(1), ...},
{AuthorID: uintPtr(2), ...},
{AuthorID: nil, ...},  // Some records without FK (tests nullable FK)

// Option B: Look up parent IDs dynamically
var authors []models.Author
facades.Orm().Query().Find(&authors)
// Use authors[0].ID, authors[1].ID, etc.
```

### Realistic Data

Use **contextually appropriate** data, not "Test Entity 1":

| Entity Type | Good Data | Bad Data |
|-------------|-----------|----------|
| Author | "Chinua Achebe", "Chimamanda Adichie" | "Author 1", "Test Author" |
| Book | "Things Fall Apart", "Half of a Yellow Sun" | "Book 1", "Test Book" |
| Product | "Wireless Mouse", "Standing Desk" | "Product A", "Item 1" |

Group records by category/theme with comments for readability:

```go
entities := []models.Entity{
    // Category A
    {Name: "...", ...},
    {Name: "...", ...},

    // Category B
    {Name: "...", ...},
    {Name: "...", ...},
}
```

## Step 3: Register in DatabaseSeeder

Edit `database/seeders/database_seeder.go` to include the new seeder:

```go
func (s *DatabaseSeeder) Run() error {
    // ... existing seeders ...

    // Run the entity seeder
    entitySeeder := &EntitySeeder{}
    if err := entitySeeder.Run(); err != nil {
        return err
    }

    return nil
}
```

## Step 4: Register in kernel.go

Edit `database/kernel.go` to add the seeder:

```go
func (kernel Kernel) Seeders() []seeder.Seeder {
    return []seeder.Seeder{
        &seeders.DatabaseSeeder{},
        &seeders.BookSeeder{},
        &seeders.RBACSeeder{},
        &seeders.ConfigSeeder{},
        &seeders.EntitySeeder{},  // Add here
    }
}
```

## Step 5: Run the Seeder

```bash
# Run all seeders
go run . artisan db:seed

# Or run specific seeder
go run . artisan db:seed --seeder=EntitySeeder
```

## Verify

```bash
# 1. Project compiles
go build ./...

# 2. Seeder runs without errors
go run . artisan db:seed --seeder=EntitySeeder

# 3. Verify records exist
# Check via API or database client
```

## Idempotency (CRITICAL)

The seeder MUST be **idempotent** — running it multiple times should not create duplicate records.

Two patterns:

### Pattern A: Count Check (simple, used by ConfigSeeder)
```go
var count int64
count, _ = facades.Orm().Query().Model(&models.Entity{}).Count()
if count > 0 {
    return nil  // Skip if any records exist
}
```

### Pattern B: Upsert by Unique Field (for re-seeding)
```go
for _, entity := range entities {
    var existing models.Entity
    err := facades.Orm().Query().Where("email = ?", entity.Email).First(&existing)
    if err != nil {
        // Not found, create it
        facades.Orm().Query().Create(&entity)
    }
    // Already exists, skip
}
```

**Prefer Pattern A** for simplicity. Use Pattern B only if you need to update existing seed data.

## Data Checklist

Before finishing, verify:

- [ ] At least 25 records
- [ ] Every enum/status value has 3+ records
- [ ] Nullable fields mix `nil` and populated values
- [ ] Text fields have diverse values (for sorting/search)
- [ ] Date fields span multiple years (for sorting)
- [ ] Numeric fields have varied values (for sorting)
- [ ] FK fields mix populated and nil (if nullable)
- [ ] Seeder is idempotent (check before insert)
- [ ] Registered in `database_seeder.go` and `kernel.go`
- [ ] `go build ./...` passes
- [ ] `go run . artisan db:seed --seeder=EntitySeeder` runs successfully

## Reference

- Book seeder: `database/seeders/book_seeder.go` (50 records, grouped by genre)
- Config seeder: `database/seeders/config_seeder.go` (idempotent with count check)
- Database seeder: `database/seeders/database_seeder.go` (orchestrator)
- Kernel registration: `database/kernel.go`
