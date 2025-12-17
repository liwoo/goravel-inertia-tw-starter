# TypeScript Type Generation from Go

This project includes enhanced tooling to automatically generate TypeScript types and enums from Go request structs and enum definitions.

## Features

### 1. Enhanced `make:ui` Command

The `make:ui` command now uses AST-based parsing to generate accurate TypeScript types from Go request structs, with automatic enum detection.

**Usage:**
```bash
go run . artisan make:ui --page=EntityName --request=EntityName
```

**What it generates:**

1. **TypeScript Type Definitions** (`resources/js/types/{entity}.ts`):
   - Parses Go request struct fields using Go AST
   - Automatically detects and includes all enums from `app/http/requests`
   - Generates union types and option arrays for enums
   - Creates interfaces for:
     - Main entity interface (extends BaseModel)
     - CreateData interface (matches CreateRequest)
     - UpdateData interface (all fields optional)
     - ListResponse interface
     - ListRequest interface
     - FormErrors interface
     - Stats interface

2. **UI Components** (forms, tables, detail views, etc.)

**Example:**

Given a Go enum:
```go
// app/http/requests/gender_type.go
type GenderType string

const (
    GenderMale   GenderType = "MALE"
    GenderFemale GenderType = "FEMALE"
)
```

And a Go request struct:
```go
// app/http/requests/sme_create_request.go
type SmeCreateRequest struct {
    Name   string      `json:"name"`
    Gender GenderType  `json:"gender"`
    Age    *int        `json:"age"`
}
```

The generated TypeScript will include:
```typescript
// Enum types
export type GenderType = 'MALE' | 'FEMALE';

export const GENDER_TYPE_OPTIONS: { value: GenderType; label: string }[] = [
  { value: 'MALE', label: 'Male' },
  { value: 'FEMALE', label: 'Female' }
];

// Interfaces
export interface Sme extends BaseModel {
  name: string;
  gender: GenderType;
  age?: number;
}

export interface SmeCreateData {
  name: string;
  gender: GenderType;
  age?: number;
}
```

### 2. Standalone Enum Generator

Generate TypeScript enum types from all Go enums in a directory.

**Usage:**
```bash
go run . artisan make:ts-enums [--source=app/http/requests] [--output=resources/js/types]
```

**Flags:**
- `--source` or `-s`: Source directory to scan for Go enums (default: `app/http/requests`)
- `--output` or `-o`: Output directory for TypeScript files (default: `resources/js/types`)

**What it generates:**

For each Go enum type found, it creates a separate TypeScript file with:
- Union type definition
- Options array for forms/selects

**Example:**

```bash
go run . artisan make:ts-enums
```

Output:
```
✓ Generated gender_type.ts (2 values)
✓ Generated education_type.ts (5 values)
✓ Generated nationality.ts (225 values)
```

Each file contains:
```typescript
export type GenderType = 'MALE' | 'FEMALE';

export const GENDER_TYPE_OPTIONS: { value: GenderType; label: string }[] = [
  { value: 'MALE', label: 'Male' },
  { value: 'FEMALE', label: 'Female' }
];
```

## How It Works

### Type Mapping

The type generator uses Go's AST parser to accurately extract type information:

| Go Type | TypeScript Type |
|---------|----------------|
| `string` | `string` |
| `int`, `int64`, `uint`, `float64` | `number` |
| `bool` | `boolean` |
| `time.Time` | `string` (ISO date) |
| `*string` (pointer) | `string?` (optional) |
| `[]string` (slice) | `string[]` |
| Custom enum types | Union types |

### Enum Detection

The generator automatically detects Go enums by looking for:
1. Type aliases to `string`: `type GenderType string`
2. Associated constants: `const GenderMale GenderType = "MALE"`

### Label Generation

Human-readable labels are automatically generated from constant names:
- `GenderMale` → `"Male"`
- `EducationPrimary` → `"Primary"`
- `StatusActive` → `"Active"`

The generator removes common prefixes (type name, "Type", "Status", "Kind") and adds spaces before capital letters.

## Best Practices

### 1. Organizing Enums

Place all enum type definitions in `app/http/requests` for automatic detection:

```
app/http/requests/
├── gender_type.go
├── education_type.go
├── status_type.go
└── sme_create_request.go
```

### 2. Enum Definition Pattern

Follow this pattern for consistent enum generation:

```go
package requests

type StatusType string

const (
    StatusActive   StatusType = "ACTIVE"
    StatusInactive StatusType = "INACTIVE"
    StatusPending  StatusType = "PENDING"
)
```

### 3. Request Struct Pattern

Use JSON tags and pointers for optional fields:

```go
type EntityCreateRequest struct {
    Name        string      `json:"name"`           // Required
    Description *string     `json:"description"`    // Optional
    Status      StatusType  `json:"status"`         // Enum type
    Tags        []string    `json:"tags"`           // Array
}
```

### 4. Using Generated Types

Import and use the generated types in your TypeScript/React code:

```typescript
import { GenderType, GENDER_TYPE_OPTIONS } from '@/types/gender_type';
import { SmeCreateData } from '@/types/sme';

// Use the union type
const gender: GenderType = 'MALE';

// Use in a form select
<Select options={GENDER_TYPE_OPTIONS} />

// Use the interface
const createSme = (data: SmeCreateData) => {
  // TypeScript will validate the shape
};
```

## Workflow

### Initial Setup

1. Generate standalone enum files:
   ```bash
   go run . artisan make:ts-enums
   ```

2. Import enums in your shared types file if needed.

### Creating New Entities

1. Create the Go request struct in `app/http/requests/`:
   ```go
   // app/http/requests/product_create_request.go
   type ProductCreateRequest struct {
       Name  string `json:"name"`
       Price int    `json:"price"`
   }
   ```

2. Generate the full UI:
   ```bash
   go run . artisan make:ui --page=Product --request=Product
   ```

3. The TypeScript types will automatically include any new enums found in the requests directory.

### Adding New Enums

1. Create the enum in `app/http/requests/`:
   ```go
   // app/http/requests/product_status.go
   type ProductStatus string

   const (
       ProductStatusAvailable ProductStatus = "AVAILABLE"
       ProductStatusSoldOut   ProductStatus = "SOLD_OUT"
   )
   ```

2. Regenerate enum files:
   ```bash
   go run . artisan make:ts-enums
   ```

3. Import where needed:
   ```typescript
   import { ProductStatus, PRODUCT_STATUS_OPTIONS } from '@/types/product_status';
   ```

## Technical Details

### Implementation

- **Parser**: Uses Go's `go/ast` and `go/parser` packages for accurate code analysis
- **No runtime dependencies**: All type generation happens at build time
- **Fallback**: If enhanced parsing fails, falls back to the original string-based parser

### Files

- `app/console/commands/type_generator_enhanced.go` - AST-based type generator
- `app/console/commands/ui_maker.go` - Enhanced UI maker with type generation
- `app/console/commands/enum_generator.go` - Standalone enum generator

## Troubleshooting

### Enums not detected

Make sure:
1. Enum type is defined as `type EnumName string`
2. Constants are properly declared with the type
3. Files are in the source directory (default: `app/http/requests`)

### Incorrect TypeScript types

The generator will fall back to basic parsing if AST parsing fails. Check:
1. Go files are syntactically correct
2. JSON tags are properly formatted
3. Source files are accessible

### Label formatting issues

Labels are auto-generated. To customize:
1. Modify the `generateEnumLabel()` function in `type_generator_enhanced.go`
2. Or manually edit the generated TypeScript files

## Future Enhancements

Potential improvements:
- [ ] Support for nested struct types
- [ ] Custom type mappings configuration
- [ ] Validation rule extraction from Go struct tags
- [ ] Documentation comment preservation
- [ ] Integration with `guts` library for more advanced scenarios
