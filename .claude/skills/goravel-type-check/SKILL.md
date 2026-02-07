---
name: goravel-type-check
description: Cross-reference Go structs (model, requests) with TypeScript interfaces to verify type consistency, field coverage, nullable handling, and JSON tag alignment.
argument-hint: "[EntityName]"
allowed-tools: Read, Grep, Glob
---

# Go ↔ TypeScript Type Consistency Check

Verify type alignment for `$ARGUMENTS`.

## Files to Cross-Reference

### Go (source of truth)
1. `app/models/<entity>.go` — Model struct (JSON response shape)
2. `app/http/requests/<entity>_*request.go` — Create/Update request structs (API input shape)

### TypeScript (must match Go)
3. `resources/js/types/<entity>.ts` — Entity interface, CreateData, UpdateData

### Frontend consumers (must use correct types)
4. `resources/js/pages/<Entity>/sections/<Entity>CreateForm.tsx`
5. `resources/js/pages/<Entity>/sections/<Entity>EditForm.tsx`
6. `resources/js/pages/<Entity>/sections/<Entity>DetailView.tsx`
7. `resources/js/pages/<Entity>/sections/<Entity>Columns.tsx`

## Type Mapping Reference

### Go → TypeScript Field Types

| Go Type | TypeScript Type | Notes |
|---|---|---|
| `string` | `string` | |
| `*string` | `string \| undefined` or `string?` | Optional field |
| `int`, `uint`, `int64` | `number` | |
| `*int`, `*uint` | `number?` | Optional |
| `float64` | `number` | |
| `*float64` | `number?` | Optional |
| `bool` | `boolean` | |
| `*bool` | `boolean?` | Optional |
| `time.Time` | `string` | ISO 8601 / RFC3339 |
| `*time.Time` | `string?` | Optional date |
| `carbon.DateTime` | `string` | Audit fields |
| `[]string` | `string[]` | JSON array |
| `datatypes.JSON` | `Record<string, any>` or typed interface | JSON object |

### JSON Tag → TypeScript Property Name

Go JSON tags determine the API response property names:

```go
// Go model
Title       string     `json:"title"`        // → TS: title: string
PublishedAt *time.Time `json:"publishedAt"`  // → TS: publishedAt?: string
```

**Critical**: The TypeScript interface property name must match the Go JSON tag value exactly.

### Dual-Case Fields

Some entities need both camelCase and snake_case for compatibility:

```typescript
// TS interface
publishedAt?: string;
published_at?: string;  // Backend might send snake_case
```

Check `MarshalJSON()` in the Go model to see which case the API actually sends.

## Checklist

### 1. Model → Entity Interface

For each field in the Go model:

- [ ] Field exists in TypeScript `Entity` interface
- [ ] TypeScript property name matches Go `json:"..."` tag
- [ ] TypeScript type matches Go type (see mapping table)
- [ ] Nullable Go fields (`*string`, `*time.Time`) are optional in TS (`?`)
- [ ] Required Go fields (non-pointer) are required in TS (no `?`)
- [ ] Virtual fields (`gorm:"-"`) are included in TS if exposed via JSON
- [ ] Hidden fields (`json:"-"`) are NOT in TS interface
- [ ] `BaseAuditableModel` fields included: `id`, `createdAt`, `updatedAt`, `createdBy`, `updatedBy`

### 2. Create Request → CreateData Interface

For each field in Go `CreateRequest`:

- [ ] Field exists in TypeScript `CreateData` interface
- [ ] Required fields in Go `Rules()` are required in TS (no `?`)
- [ ] Optional fields are optional in TS (`?`)
- [ ] Type matches (string→string, float64→number, []string→string[])
- [ ] Enum fields use TypeScript union type or enum (not bare `string`)
- [ ] Read-only fields (ID, created_at, etc.) are NOT in CreateData

### 3. Update Request → UpdateData Interface

For each field in Go `UpdateRequest`:

- [ ] Field exists in TypeScript `UpdateData` interface
- [ ] All fields are optional in TS (`?`) — updates are partial
- [ ] Pointer types in Go (`*string`, `*float64`) map to optional TS
- [ ] Same enum types as CreateData
- [ ] Read-only fields are NOT in UpdateData

### 4. API Response Consistency

- [ ] Create form sends data matching Go `CreateRequest` struct fields
- [ ] Edit form sends data matching Go `UpdateRequest` struct fields
- [ ] Create form maps camelCase state → snake_case request body
- [ ] Edit form maps camelCase state → snake_case request body
- [ ] Detail view displays all fields from Entity interface
- [ ] Column definitions reference fields that exist on Entity interface

### 5. Enum Consistency

For each enum/status field:

- [ ] Go validation has `in:VALUE1,VALUE2,...` rule
- [ ] TypeScript has matching union type: `type EntityStatus = 'VALUE1' | 'VALUE2' | ...`
- [ ] Enum values are UPPERCASE in Go and TypeScript (convention)
- [ ] Form dropdowns list all enum values
- [ ] Column status badges handle all enum values
- [ ] i18n `status.*` keys exist for all enum values

### 6. Field Name Mapping (Frontend → API)

Check that forms convert camelCase → snake_case when sending to API:

```typescript
// Form state (camelCase)
const [formData, setFormData] = useState<CreateData>({
    publishedAt: '',
    configType: '',
});

// API request body (snake_case)
const requestBody = {
    published_at: formData.publishedAt,
    config_type: formData.configType,
};
```

- [ ] All camelCase form fields map to snake_case in request body
- [ ] No mismatched field names between form state and API payload
- [ ] Error response field names (`errors.field_name`) map back to form field names

### 7. Nullable Field Handling

- [ ] Go `*string` fields: TS shows empty string or "Not specified" when nil
- [ ] Go `*time.Time` fields: TS handles null dates gracefully
- [ ] Go `*float64` fields: TS handles null numbers (not NaN)
- [ ] Create form initializes optional fields to empty/default values
- [ ] Edit form handles both present and absent optional values from API response

## Common Type Mismatches

### Missing optional marker
```go
// Go: pointer = nullable
Description *string `json:"description"`
```
```typescript
// BAD: required in TS but nullable in Go
description: string;

// GOOD: optional in TS
description?: string;
```

### Wrong enum type
```go
// Go: validated against specific values
rules["status"] = "in:ACTIVE,INACTIVE,PENDING"
```
```typescript
// BAD: accepts any string
status: string;

// GOOD: constrained to valid values
status: 'ACTIVE' | 'INACTIVE' | 'PENDING';
```

### Mismatched field name
```go
// Go: JSON tag is camelCase
ConfigType string `json:"configType"`
```
```typescript
// BAD: doesn't match JSON tag
config_type: string;

// GOOD: matches JSON tag
configType: string;
```

### Missing audit fields
```go
// Go: BaseAuditableModel provides these automatically
type Entity struct {
    BaseAuditableModel  // → id, createdAt, updatedAt, createdBy, updatedBy, deletedAt
    Name string `json:"name"`
}
```
```typescript
// BAD: missing audit fields
interface Entity {
    name: string;
}

// GOOD: extends BaseModel which includes audit fields
interface Entity extends BaseModel {
    name: string;
}
```

## Output

After review, produce a report with:
1. **Field-by-field comparison table**: Go field → TS field → Match/Mismatch
2. **Type mismatches**: Fields where Go and TS types don't align
3. **Missing fields**: Fields in Go but not in TS (or vice versa)
4. **Nullable inconsistencies**: Required in TS but nullable in Go
5. **Enum gaps**: Enum values in Go validation but missing from TS union
6. **Form field mapping issues**: camelCase/snake_case conversion errors

## Reference

- Canonical Go model: `app/models/book.go`
- Canonical TS types: `resources/js/types/book.ts`
- Canonical Go request: `app/http/requests/book_request.go`
- Base model TS: `resources/js/types/crud.ts` (`BaseModel` interface)
- Base Go model: `app/models/base_auditable_model.go` (`BaseAuditableModel`)
