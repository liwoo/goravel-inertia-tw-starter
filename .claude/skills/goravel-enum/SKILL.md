---
name: goravel-enum
description: Create a Go enum type and automatically generate the TypeScript equivalent. Produces union types and option arrays for forms.
argument-hint: "[EnumName] [value1,value2,value3]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel Enum Generator

Create Go enum and TypeScript equivalent for `$ARGUMENTS`.

## Step 1: Create Go Enum

Create the enum in `app/http/requests/<enum_name>.go`:

```go
package requests

type EnumName string

const (
    EnumNameValue1 EnumName = "VALUE_1"
    EnumNameValue2 EnumName = "VALUE_2"
    EnumNameValue3 EnumName = "VALUE_3"
)
```

### Naming Conventions

- **Type name**: PascalCase (e.g., `GenderType`, `BookStatus`, `ConfigTypes`)
- **Constant prefix**: Type name + value name in PascalCase
- **Constant values**: UPPER_CASE or Title Case (match your DB values)

### Examples from Codebase

```go
// Simple enum
type GenderType string
const (
    GenderTypeMale   GenderType = "MALE"
    GenderTypeFemale GenderType = "FEMALE"
    GenderTypeOther  GenderType = "OTHER"
)

// Status enum
type BookStatus string
const (
    BookStatusAvailable   BookStatus = "AVAILABLE"
    BookStatusBorrowed    BookStatus = "BORROWED"
    BookStatusMaintenance BookStatus = "MAINTENANCE"
    BookStatusReserved    BookStatus = "RESERVED"
)

// Category enum (Title Case values)
type ConfigTypes string
const (
    ConfigTypeFinancing          ConfigTypes = "Financing"
    ConfigTypeImprovementAspects ConfigTypes = "Improvement Aspects"
    ConfigTypeBusinessCategories ConfigTypes = "Business Categories"
)
```

## Step 2: Generate TypeScript Equivalent (AUTOMATIC)

After creating the Go enum, run the TypeScript generator:

```bash
go run . artisan make:ts-enums --source=app/http/requests --output=resources/js/types
```

This scans `app/http/requests/` for Go enum patterns and generates TypeScript files in `resources/js/types/`.

### Generated Output Format

For each Go enum, it produces a file like `resources/js/types/gender_type.ts`:

```typescript
// Auto-generated from Go enum: GenderType
export type GenderType = 'MALE' | 'FEMALE' | 'OTHER';

export const GENDER_TYPE_OPTIONS: Array<{label: string; value: GenderType}> = [
    { label: 'Male', value: 'MALE' },
    { label: 'Female', value: 'FEMALE' },
    { label: 'Other', value: 'OTHER' },
];
```

## Step 3: Use in Forms

```tsx
import { GenderType, GENDER_TYPE_OPTIONS } from '@/types/gender_type';

<Select
    value={formData.gender}
    onValueChange={(value) => setFormData({ ...formData, gender: value as GenderType })}
>
    <SelectTrigger>
        <SelectValue placeholder="Select gender" />
    </SelectTrigger>
    <SelectContent>
        {GENDER_TYPE_OPTIONS.map((option) => (
            <SelectItem key={option.value} value={option.value}>
                {option.label}
            </SelectItem>
        ))}
    </SelectContent>
</Select>
```

## Step 4: Use in Validation Rules

Reference the enum in your request validation:

```go
func (r *EntityCreateRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "gender": "required|in:MALE,FEMALE,OTHER",
        "status": "required|in:AVAILABLE,BORROWED,MAINTENANCE,RESERVED",
    }
}
```

## Alternative: guts Library

For broader Go-to-TypeScript type conversion (not just enums), consider the `guts` library from Coder:

```bash
go get github.com/coder/guts
```

guts converts full Go type definitions (structs, interfaces, enums) to TypeScript. It's a Go library (not CLI) - you import it into a generator program. Useful when you need to sync entire model definitions, not just enums.

See: https://github.com/coder/guts

## Verification

1. Confirm Go enum file exists in `app/http/requests/`
2. Confirm TypeScript file was generated in `resources/js/types/`
3. Import and use in a React component to verify types work
