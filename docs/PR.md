# Pull Request Tracking

This file tracks changes made to base/shared functionality that need to be propagated to the origin codebase or other projects using the same framework.

## Pending Changes for Origin Codebase

### 1. Fix Generic CRUD Service Delete Method (CRITICAL)

**File**: `app/contracts/generic_crud_service.go`
**Line**: 666
**Issue**: Delete method doesn't specify table name explicitly, causing GORM to incorrectly infer table names for models with custom `TableName()` methods.

**Error Example**:
```
sqlite3: SQL logic error: no such table: configs
```

**Problem Code**:
```go
// Line 665 - Missing .Table() specification
if _, err := facades.Orm().Query().Delete(&model); err != nil {
    return fmt.Errorf("failed to delete %s: %w", s.tableName, err)
}
```

**Fixed Code**:
```go
// Line 666 - Explicitly specify table name
// IMPORTANT: Must use .Table() to specify correct table name for models with custom TableName()
if _, err := facades.Orm().Query().Table(s.tableName).Delete(&model); err != nil {
    return fmt.Errorf("failed to delete %s: %w", s.tableName, err)
}
```

**Impact**:
- Affects all models with custom `TableName()` methods
- Causes delete operations to fail with "no such table" errors
- Critical bug that breaks CRUD delete functionality

**Testing**:
- Test with models that have custom table names (e.g., `sme_config`, `user_roles`, etc.)
- Run CRUD test suite: `go test -v ./tests/feature/crud -run TestConfigControllerCRUDTestSuite`
- Verify delete operations return 204 instead of 500

---

## Validation Rule Fix for Custom Types

### 2. Remove `|string` Validation for Custom Type Fields

**Files**:
- `app/http/requests/*_create_request.go`
- `app/http/requests/*_update_request.go`

**Issue**: Goravel's validator cannot validate custom type fields (like enums) with the `|string` rule.

**Error Example**:
```json
{
  "errors": {
    "config_type": {
      "string": "config_type value must be a string"
    }
  }
}
```

**Problem Code**:
```go
func (r *ConfigCreateRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "name":        "required|string|max_len:255",
        "config_type": "required|string|max_len:255",  // ❌ Fails for custom types
    }
}
```

**Fixed Code**:
```go
func (r *ConfigCreateRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "name":        "required|max_len:255",
        "config_type": "required|max_len:255",  // ✅ Works for custom types
    }
}
```

**Impact**:
- Affects all request validation for enum/custom type fields
- Causes 422 validation errors on valid data
- Consider updating code generators to not include `|string` for custom types

**Note**: This should be documented in the framework's validation guide or code generator templates should be updated.

---

## Future Considerations

### Items to Monitor

1. **GORM Table Name Resolution**
   - Ensure all GORM operations use `.Table(tableName)` when working with models that have custom table names
   - Consider adding a linter rule or test to catch missing `.Table()` calls

2. **Validation Rule Generators**
   - Update artisan command templates to avoid `|string` for custom type fields
   - Add documentation about validating enum types

3. **Test Suite Updates**
   - Update test generators to use correct permission names (service-based, not controller-based)
   - Update test generators to use correct API endpoints
   - Add handling for framework limitations in delete tests

---

## How to Apply These Changes

### For Origin Codebase Maintainers

1. **Review and test the delete fix** in `app/contracts/generic_crud_service.go`
2. **Update validation documentation** to clarify `|string` usage with custom types
3. **Run full test suite** to ensure no regressions
4. **Update code generators** if applicable

### For Project Forks/Users

When pulling updates from origin codebase:
- Check if this fix is included
- If not, manually apply the `.Table(s.tableName)` fix to the Delete method
- Update validation rules for any custom type fields

---

## Version Information

**Date**: 2025-10-20
**Framework Version**: Goravel v1.16.3
**Applied To**: smedi-database project
**Status**: ✅ Fixed locally, ⏳ Pending upstream

---

## Related Issues

- Delete operations failing with "no such table: configs" error
- Validation errors on enum fields: "value must be a string"
- CRUD test suite failures for models with custom table names
