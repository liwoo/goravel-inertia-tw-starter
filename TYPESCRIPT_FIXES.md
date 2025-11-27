# TypeScript Linting Issues - Fix Details

## Summary
- **Initial errors:** ~90 across 25+ files
- **After fixes:** 73 errors remaining
- **Fixed:** ~17 errors

## Fixed Issues

The following have been fixed:
1. SharedData index signature for usePage constraint
2. Missing @/components/ui/form module (removed unused import)
3. CrudPage.tsx preserveState/preserveScroll options removed
4. CrudPage.tsx implicit any in filter callbacks
5. FilterConditionBuilder.tsx operator comparison cast
6. GlobalSearch.tsx permissionAction type (added PermissionAction type)
7. NotificationDrawer.tsx null checks for counts
8. PermissionMatrix.tsx implicit any types (added MatrixRole, MatrixService, MatrixAction interfaces)
9. site-header.tsx User null checks (conditional rendering)
10. exportUtils.ts jsPDF getNumberOfPages method (type assertion)
11. PermissionsContext.tsx null checks for auth.user

---

## 1. SharedData Index Signature Issue (5 occurrences)

**Files affected:**
- `resources/js/types/app.d.ts`
- `resources/js/layouts/Admin.tsx`
- `resources/js/components/site-header.tsx`
- `resources/js/pages/Bdsp/Index.tsx`
- `resources/js/pages/Event/Index.tsx`
- `resources/js/pages/ProcurementNotice/Index.tsx`
- `resources/js/pages/Sme/Index.tsx`

**Error:**
```
Type 'SharedData' does not satisfy the constraint 'PageProps'.
Index signature for type 'string' is missing in type 'SharedData'.
```

**Fix:**
Add index signature to `SharedData` interface in `resources/js/types/app.d.ts`:

```typescript
export interface SharedData {
    pageTitle: string;
    auth: {
        user: User | null;
    };
    [key: string]: any; // Add this line
}
```

---

## 2. Missing @/components/ui/form Module

**File affected:**
- `resources/js/components/Forms/CrudForm.tsx`

**Error:**
```
Cannot find module '@/components/ui/form' or its corresponding type declarations.
```

**Fix Options:**
1. **Option A** - Create the missing `form.tsx` component (shadcn/ui Form component)
2. **Option B** - Remove the unused import if `Form` is not actually used (it's imported but the component uses a native `<form>` element)

Looking at the file, `Form` is imported but never used in the JSX. **Remove the unused import:**

```typescript
// Remove this line:
import { Form } from '@/components/ui/form';
```

---

## 3. CrudPage.tsx Type Issues (6 errors)

**File:** `resources/js/components/Crud/CrudPage.tsx`

### 3.1 Index Signature Issue (lines 200, 227)
**Error:**
```
Element implicitly has an 'any' type because expression of type 'string' can't be used to index type 'ListRequest | Record<string, any>'.
```

**Fix:** Cast filters to a type with index signature:
```typescript
// Line 200: Change
const currentValue = actualFilters[key] !== undefined ? actualFilters[key] : (filters as any)[key];
// To (use Record<string, any> cast)
const currentValue = actualFilters[key] !== undefined
  ? actualFilters[key]
  : (filters as Record<string, any>)[key];
```

### 3.2 Implicit Any in Filter (line 485)
**Error:**
```
Parameter '_' implicitly has an 'any' type.
Parameter 'i' implicitly has an 'any' type.
```

**Fix:**
```typescript
// Change
const newConditions = appliedDynamicFilter.conditions.filter((_, i) => i !== index);
// To
const newConditions = appliedDynamicFilter.conditions.filter((_: FilterCondition | CompoundFilter, i: number) => i !== index);
```

### 3.3 preserveState Option (lines 621, 671)
**Error:**
```
Object literal may only specify known properties, and 'preserveState' does not exist in type 'ReloadOptions<RequestPayload>'.
```

**Fix:** Inertia.js v1.x renamed `preserveState` to `preserveState` in `router.visit()` but `router.reload()` doesn't have this option. Remove it:
```typescript
// Change
router.reload({
    only: ['data', 'filters', 'stats'],
    preserveState: false,
    preserveScroll: true
});
// To
router.reload({
    only: ['data', 'filters', 'stats'],
    preserveScroll: true
});
```

### 3.4 Implicit Any in Line 799
**Error:**
```
Parameter 'c' implicitly has an 'any' type.
```

**Fix:**
```typescript
// Change
return appliedDynamicFilter.conditions.filter(c => !('logic' in c)).length;
// To
return appliedDynamicFilter.conditions.filter((c: FilterCondition | CompoundFilter) => !('logic' in c)).length;
```

---

## 4. CrudPage.test.tsx Issues (4 errors)

**File:** `resources/js/components/Crud/CrudPage.test.tsx`

### 4.1 Element vs HTMLElement (line 111)
**Fix:**
```typescript
// Change
const selectTrigger = container.querySelector('[data-slot="select-trigger"]');
fireEvent.click(selectTrigger);
// To
const selectTrigger = container.querySelector('[data-slot="select-trigger"]') as HTMLElement;
if (selectTrigger) fireEvent.click(selectTrigger);
```

### 4.2 ForwardRefExoticComponent Issues (lines 357, 380, 403)
The test mock forms need to be proper forwardRef components:

```typescript
// Change
createForm: () => <div>Create Form</div>,
// To
createForm: React.forwardRef<any, any>(() => <div>Create Form</div>),
```

---

## 5. FilterConditionBuilder.tsx Operator Comparison (2 errors)

**File:** `resources/js/components/Filters/FilterConditionBuilder.tsx`

**Error:**
```
This comparison appears to be unintentional because the types 'FilterOperator' and '""' have no overlap.
```

**Fix:** The `FilterOperator` type doesn't include empty string. Update the type or change the comparison:
```typescript
// Option 1: Add empty string to type
type FilterOperator = '' | 'eq' | 'neq' | 'gt' | ...;

// Option 2: Cast the comparison
if ((condition.operator as string) === '') ...
```

---

## 6. GlobalSearch.tsx Permission Type Mismatch

**File:** `resources/js/components/GlobalSearch.tsx`

**Error at line 53:**
```
Argument of type 'string' is not assignable to parameter of type '"delete" | "export" | "create" | "read" | ...'.
```

**Fix:** The `permissionAction` in `search_config.tsx` is typed as `string` but `canPerformAction` expects a union type. Update the config type:

```typescript
// In search_config.tsx, change:
permissionAction: string;
// To:
permissionAction: 'create' | 'read' | 'update' | 'delete' | 'export' | 'bulk_update' | 'bulk_delete' | 'write' | 'manage';
```

---

## 7. MessageSidebar.tsx User Type Incompatibility (4 errors)

**File:** `resources/js/components/Messages/MessageSidebar.tsx`

**Errors:**
- User type missing `role` property
- Conversation `latest_message` type mismatch

**Fix:**
1. Ensure the `User` interface in MessageSidebar matches the app's User type
2. Update the `Conversation` interface to include all required `Message` fields

---

## 8. NotificationDrawer.tsx Null Checks (6 errors)

**File:** `resources/js/components/Notifications/NotificationDrawer.tsx`

**Errors:**
```
'counts.total' is possibly 'undefined'.
'counts' is possibly 'null'.
```

**Fix:** Add null checks before accessing:
```typescript
// Change
if (counts.total > 0) ...
// To
if (counts?.total && counts.total > 0) ...
```

---

## 9. PermissionMatrix.tsx Implicit Any Types (20 errors)

**File:** `resources/js/components/Permissions/PermissionMatrix.tsx`

**Fix:** Add type annotations to all callback parameters:
```typescript
// Example fixes:
.map((role) => ...)  // Change to: .map((role: Role) => ...)
.reduce((total, role) => ...) // Change to: .reduce((total: number, role: Role) => ...)
```

---

## 10. site-header.tsx User Null Checks (2 errors)

**File:** `resources/js/components/site-header.tsx`

**Errors:**
```
Type 'User | null' is not assignable to type 'User'.
```

**Fix:** Add null checks or pass undefined:
```typescript
// Change
<UserMenu user={auth.user} />
// To
{auth.user && <UserMenu user={auth.user} />}
```

---

## 11. SME-Related Type Issues (~25 errors)

### 11.1 SmeEditFormSimple.tsx - Snake_case vs CamelCase (~20 errors)
**File:** `resources/js/pages/Sme/sections/SmeEditFormSimple.tsx`

The code uses snake_case property names (`registration_number`, `tax_identification_number`) but the `Sme` type uses camelCase (`registrationNumber`, `taxIdentificationNumber`).

**Fix:** Update all property references to use camelCase:
```typescript
// Change
item.registration_number
// To
item.registrationNumber
```

### 11.2 SmeEditBusinessInfoTab.tsx - Possibly Undefined (~12 errors)
**File:** `resources/js/pages/Sme/sections/edit-tabs/SmeEditBusinessInfoTab.tsx`

**Fix:** Add null checks or use optional chaining:
```typescript
// Change
businessData.businessImprovementAspects.length
// To
businessData.businessImprovementAspects?.length ?? 0
```

### 11.3 SmeEditFormalizationTab.tsx - Missing Properties
**Fix:** The default object is missing `smeId` and `id`. Add them or make them optional in the type.

### 11.4 SmeDetailView.tsx - Missing 'region' Property
**Line 144:** Property 'region' doesn't exist on Sme type. Either add it to the type or remove the reference.

---

## 12. exportUtils.ts jsPDF Method

**File:** `resources/js/utils/exportUtils.ts`

**Error at line 353:**
```
Property 'getNumberOfPages' does not exist on type 'jsPDF'.
```

**Fix:** The jsPDF type definitions may be outdated. Use type assertion:
```typescript
// Change
const pageCount = doc.getNumberOfPages();
// To
const pageCount = (doc as any).getNumberOfPages();
```

Or install updated types: `npm install @types/jspdf@latest`

---

## 13. Books-Related Date Parsing Issues (5 errors)

**Files:**
- `resources/js/components/Books/BookForms.tsx`
- `resources/js/pages/Books/sections/BookColumns.tsx`
- `resources/js/pages/Books/sections/BookDetailView.tsx`
- `resources/js/pages/Books/sections/BookEditForm.tsx`

**Error:**
```
Argument of type 'string | undefined' is not assignable to parameter of type 'string | number | Date'.
```

**Fix:** Add null check before creating Date:
```typescript
// Change
new Date(item.publishedDate)
// To
item.publishedDate ? new Date(item.publishedDate) : null
```

---

## 14. RolesIndex.tsx Type Mismatches (3 errors)

**File:** `resources/js/pages/Permissions/RolesIndex.tsx`

**Errors:**
- `RoleListResponse` not assignable to `PaginatedResult<Role>` (from/null issue)
- ForwardRefExoticComponent issues with form components

**Fix:**
1. Update the `PaginatedResult` type to allow `null` for `from`
2. Ensure form components are properly wrapped with `forwardRef`

---

## 15. Test Setup IntersectionObserver

**File:** `resources/js/test/setup.ts`

**Fix:** Update the mock to include all required properties:
```typescript
global.IntersectionObserver = class IntersectionObserver {
  root = null;
  rootMargin = '';
  thresholds = [];
  // ... other methods
}
```

---

## Priority Order for Fixes

1. **High Priority (Blocks compilation):**
   - SharedData index signature (affects multiple files)
   - Missing ui/form module

2. **Medium Priority (Type safety):**
   - CrudPage.tsx type issues
   - SME snake_case vs camelCase
   - Null checks in NotificationDrawer

3. **Lower Priority (Tests/Edge cases):**
   - Test file fixes
   - PermissionMatrix implicit any
   - IntersectionObserver mock
