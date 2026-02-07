---
name: goravel-crud-request
description: Generate create and update request validators for a Goravel entity. Handles validation rules, field mapping, and carbon.DateTime conversion.
argument-hint: "[entity_name]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel CRUD Request Generator

Generate request validators for `$ARGUMENTS`.

## Step 1: Generate Requests

```bash
go run . artisan make:req --model=<entity_name> <entity_name>
```

This creates:
- `app/http/requests/<entity_name>_create_request.go`
- `app/http/requests/<entity_name>_update_request.go`

## Step 2: Review and Fix Create Request

### Struct Tags (CRITICAL)

All `form` and `json` tags MUST use **snake_case** for Goravel's `ctx.Request().Bind()` to work:

```go
// CORRECT — Bind() populates fields from JSON payloads
FirstName string  `form:"first_name" json:"first_name"`
BirthDate *string `form:"birth_date" json:"birth_date"`
PhotoURL  *string `form:"photo_url" json:"photo_url"`

// WRONG — Bind() will NOT populate these fields
FirstName string `form:"firstName" json:"firstName"`
```

### Struct Fields

- Use concrete types for required fields: `string`, `float64`, `bool`
- Use pointer types for optional fields: `*string`, `*float64`, `*bool`
- Date fields should be `string` (not `carbon.DateTime`)
- Array fields should be `[]string` or `[]int`

### Validation Rules

```go
func (r *EntityCreateRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "title":       "required|max_len:255",
        "description": "required",
        "email":       "required|email",
        "amount":      "required",       // NO |numeric for float64
        "date":        "",               // NO |date for *string - validate in PrepareForValidation
        "status":      "required|in:ACTIVE,INACTIVE,PENDING",
    }
}
```

### Validation Gotchas

| Field Type | Wrong | Correct |
|---|---|---|
| `float64` | `"required\|numeric"` | `"required"` (numeric validation breaks Go float64) |
| `*string` date | `"date"` | Custom validation in `PrepareForValidation()` |
| Custom string type | `"required\|string"` | `"required\|max_len:255"` (remove `\|string`) |

### Validation Rules Key Format (CRITICAL)

The generic CrudController validates the **data map** from `ToCreateData()` / `ToUpdateData()`, NOT the struct fields. Therefore:
- Create `Rules()` keys MUST match `ToCreateData()` output keys (**camelCase**)
- Update `Rules()` keys MUST match `ToUpdateData()` output keys (**snake_case**)
- Same applies to `Messages()` and `Attributes()` methods

```go
// CREATE Rules — camelCase keys to match ToCreateData() output
func (r *EntityCreateRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "firstName": "required|max_len:100",   // camelCase
        "lastName":  "required|max_len:100",   // camelCase
        "status":    "in:ACTIVE,INACTIVE",
    }
}

// UPDATE Rules — snake_case keys to match ToUpdateData() output
func (r *EntityUpdateRequest) Rules(ctx http.Context) map[string]string {
    rules := map[string]string{}
    if r.FirstName != nil {
        rules["first_name"] = "required|max_len:100"   // snake_case
    }
    return rules
}
```

### ToCreateData() Method (camelCase Keys)

Keys MUST use **camelCase** to match model `json:"..."` tags. The generic service's `setFieldsRecursively()` matches data map keys against model json tags during creation.

```go
func (r *EntityCreateRequest) ToCreateData() map[string]interface{} {
    data := map[string]interface{}{
        "firstName":   r.FirstName,    // camelCase — matches model json:"firstName"
        "lastName":    r.LastName,     // camelCase
        "description": r.Description,
    }

    // Handle optional fields
    if r.Notes != nil {
        data["notes"] = *r.Notes
    }

    // Handle multi-word optional fields (camelCase)
    if r.BirthDate != nil && *r.BirthDate != "" {
        data["birthDate"] = *r.BirthDate  // camelCase
    }

    // Handle array fields
    if r.Tags != nil {
        data["tags"] = r.Tags
    } else {
        data["tags"] = []string{}  // Default empty array
    }

    return data
}
```

### carbon.DateTime Conversion (CRITICAL)

```go
// If model uses NON-POINTER carbon.DateTime:
//   Date carbon.DateTime `json:"date"`
if r.Date != "" {
    parsedDate := carbon.Parse(r.Date)
    if parsedDate.Error == nil {
        data["date"] = *carbon.NewDateTime(parsedDate)  // DEREFERENCE with *
    }
}

// If model uses POINTER *carbon.DateTime:
//   DateOfBirth *carbon.DateTime `json:"date_of_birth"`
if r.DateOfBirth != nil && *r.DateOfBirth != "" {
    parsedDate := carbon.Parse(*r.DateOfBirth)
    if parsedDate.Error == nil {
        data["date_of_birth"] = carbon.NewDateTime(parsedDate)  // NO dereference
    }
}
```

**Rule**: `carbon.NewDateTime()` returns `*carbon.DateTime`. Dereference with `*` only if model field is non-pointer.

## Step 3: Fix Update Request

Same structure but ALL fields should be pointers (partial updates):

```go
type EntityUpdateRequest struct {
    Title       *string   `json:"title"`
    Description *string   `json:"description"`
    Tags        *[]string `json:"tags"`
}
```

### ToUpdateData() - snake_case Keys (CRITICAL)

Keys MUST use **snake_case** to match DB column names. The generic service's `Update()` passes the data map directly to GORM's `Update()`, which needs DB column names.

```go
func (r *EntityUpdateRequest) ToUpdateData() map[string]interface{} {
    data := map[string]interface{}{}

    if r.FirstName != nil {
        data["first_name"] = *r.FirstName  // snake_case — matches DB column
    }
    if r.LastName != nil {
        data["last_name"] = *r.LastName    // snake_case
    }
    if r.Description != nil {
        data["description"] = *r.Description
    }
    if r.BirthDate != nil {
        data["birth_date"] = *r.BirthDate  // snake_case
    }
    if r.Tags != nil {
        data["tags"] = *r.Tags
    }

    return data
}
```

### Why Create and Update Use Different Key Formats

| Method | Key Format | Reason |
|--------|-----------|--------|
| `ToCreateData()` | camelCase | `setFieldsRecursively()` matches model `json:"..."` tags |
| `ToUpdateData()` | snake_case | GORM's `Update()` matches DB column names directly |
```

## Step 4: Exclude Read-Only Fields

Fields that are calculated, aggregated, or system-managed should NOT appear in either request struct. Examples:
- Scores, totals, aggregates
- `created_by`, `updated_by` (set by framework)
- `created_at`, `updated_at` (set by framework)

## Verify

After fixing both request files:

```bash
# Vet the requests package
go vet ./app/http/requests/...

# Confirm full project compiles
go build ./...
```

## Next Step

Run `/goravel-crud-controller` to generate the API controller.
