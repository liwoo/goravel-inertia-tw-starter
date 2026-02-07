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

### ToCreateData() Method

Maps request fields to database columns:

```go
func (r *EntityCreateRequest) ToCreateData() map[string]interface{} {
    data := map[string]interface{}{
        "title":       r.Title,
        "description": r.Description,
    }

    // Handle optional fields
    if r.Notes != nil {
        data["notes"] = *r.Notes
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

### ToUpdateData() - Only Include Provided Fields

```go
func (r *EntityUpdateRequest) ToUpdateData() map[string]interface{} {
    data := map[string]interface{}{}

    if r.Title != nil {
        data["title"] = *r.Title
    }
    if r.Description != nil {
        data["description"] = *r.Description
    }
    if r.Tags != nil {
        data["tags"] = *r.Tags
    }

    return data
}
```

## Step 4: Exclude Read-Only Fields

Fields that are calculated, aggregated, or system-managed should NOT appear in either request struct. Examples:
- Scores, totals, aggregates
- `created_by`, `updated_by` (set by framework)
- `created_at`, `updated_at` (set by framework)

## Next Step

Run `/goravel-crud-controller` to generate the API controller.
