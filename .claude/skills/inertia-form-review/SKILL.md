---
name: inertia-form-review
description: Audit and review existing form components for i18n compliance, type safety, missing validations, dropdown opportunities, accessibility, and consistency between create/edit forms.
argument-hint: "[EntityName]"
allowed-tools: Read, Grep, Glob
---

# Inertia Form Review & Audit (i18n-aware)

Review forms for `$ARGUMENTS`.

## Files to Review

1. `resources/js/pages/<EntityName>/sections/<EntityName>CreateForm.tsx`
2. `resources/js/pages/<EntityName>/sections/<EntityName>EditForm.tsx`
3. `resources/js/pages/<EntityName>/sections/<EntityName>DetailView.tsx`
4. `resources/js/pages/<EntityName>/sections/<EntityName>Columns.tsx`
5. `resources/js/pages/<EntityName>/sections/<EntityName>PageConfig.tsx`
6. `resources/js/types/<entity_name>.ts` (TypeScript types)
7. `resources/js/locales/en/<entities>.json` (translation file)
8. `resources/js/locales/index.ts` (namespace registration)
9. Corresponding Go request files in `app/http/requests/`

## Checklist

### 1. i18n Compliance

- [ ] Translation namespace JSON file exists at `resources/js/locales/en/<entities>.json`
- [ ] Namespace registered in `resources/js/locales/index.ts` (import + ns array + resources)
- [ ] **Forms**: Use `const { t } = useTranslation('entities')` hook (React components)
- [ ] **Detail view**: Uses `const { t } = useTranslation('entities')` hook
- [ ] **Columns/Filters/PageConfig**: Use `t: TFunction` parameter (plain functions)
- [ ] **No hardcoded user-visible strings** — all go through `t()`
- [ ] All field labels use `t('form.*')` keys
- [ ] All placeholders use `t('form.enter*')` keys
- [ ] All validation errors use `t('validation.*')` keys
- [ ] Toast messages use `t('toast.created')` / `t('toast.updated')`
- [ ] Status/enum display labels use `t('status.*')` keys
- [ ] Column headers use `t('columns.*')` keys
- [ ] Filter labels use `t('filters.*')` keys
- [ ] Section headers use `t('form.entityInfo')`, `t('form.metadata')`
- [ ] Page title uses `t('page.title')` / `t('page.headTitle')`
- [ ] Confirm dialogs use `t('confirm.*')` with interpolation

### 2. i18n Key Completeness

Check that the JSON file has all required sections:
- [ ] `page.*` — title, headTitle, myEntities
- [ ] `stats.*` — if stats enabled
- [ ] `status.*` — for each enum value + `unknown`
- [ ] `columns.*` — for each visible column
- [ ] `filters.*` — for each filter + `allStatus`/`allEntities`
- [ ] `actions.*` — for page actions
- [ ] `form.*` — for each field label, placeholder, section header, metadata labels
- [ ] `validation.*` — for each required field
- [ ] `toast.*` — created, updated
- [ ] `confirm.*` — bulk delete, with `{{count}}` interpolation

### 3. Structure & Pattern Compliance

- [ ] Uses `forwardRef` pattern
- [ ] Has `useImperativeHandle(ref, () => ({ handleSubmit }))`
- [ ] Has `displayName` set: `EntityCreateForm.displayName = 'EntityCreateForm'`
- [ ] Extends `CrudFormProps` (create) or `CrudEditFormProps<Entity>` (edit)
- [ ] Has `setIsSaving` prop passed through

### 4. Type Safety

- [ ] Form state uses typed interface (`EntityCreateData` / `EntityUpdateData`)
- [ ] All `useState` calls have proper generic types
- [ ] Select `onValueChange` casts to correct enum type: `value as EnumType`
- [ ] Number fields use `parseFloat(e.target.value)` or `parseInt()`
- [ ] No `any` types (except `errors: Record<string, string>`)

### 5. Create <> Edit Consistency

- [ ] Both forms have the **same fields** (unless intentionally different)
- [ ] Field order matches between create and edit
- [ ] Same validation rules in both forms
- [ ] Both use `useTranslation` with the **same namespace**
- [ ] Both use the same `t()` keys for the same field labels
- [ ] Edit form initializes state from `item` prop (handle dual-case)
- [ ] Edit form has Metadata section (ID, Created, Updated) — create does NOT
- [ ] Edit form uses `PUT` method, create uses `POST`
- [ ] Edit form URL includes `/${item.id}`, create does not

### 6. API Integration

- [ ] Correct API endpoint: `/api/entity-names` (hyphenated, plural)
- [ ] Headers include all three: `Content-Type`, `Accept`, `X-Requested-With`
- [ ] Request body converts camelCase -> snake_case for backend
- [ ] **Nullable fields use `|| null`** (NOT empty string `""`) — PostgreSQL rejects `""` for date, numeric columns
- [ ] Required fields keep their value as-is (no `|| null`)
- [ ] Error response maps `errorData.errors` to field-level errors
- [ ] `setIsSaving?.(true)` before fetch, `setIsSaving?.(false)` in finally

### 7. Field Rendering

- [ ] Every field uses icon layout: `flex items-start gap-3` with icon box
- [ ] Required fields have `*` in label (via `t('form.name')` = `"Name *"`)
- [ ] Error state: `border-destructive` class on input + error text below
- [ ] Placeholder text on all inputs (via `t('form.enter*')`)
- [ ] Enum/category fields use `<Select>` dropdown (NOT free-text input)
- [ ] Select options use `t('status.*')` for display labels
- [ ] Boolean fields use `<Switch>` component
- [ ] Long text fields use `<Textarea>` with appropriate `rows`
- [ ] Date fields use `<Input type="date">`
- [ ] Number fields use `<Input type="number">`

### 8. Dropdown Opportunity Audit

Check for fields that should be dropdowns but are free-text:

- [ ] Status fields -> should use enum Select with `t('status.*')`
- [ ] Type/category fields -> should use enum Select
- [ ] Country/region fields -> should use predefined options
- [ ] Boolean-like fields (yes/no, active/inactive) -> should use Switch or Select
- [ ] **FK relationship fields** -> should use Select fetching related entity (e.g., `authorId` -> Author dropdown)

If free-text fields should be dropdowns:
1. Check if Go enum exists in `app/http/requests/`
2. Check if TS enum exists in `resources/js/types/`
3. If FK field, check if related entity API exists (`/api/related-entities`)
4. If neither exists, recommend creating with `/goravel-enum`

### 9. Decimal/Float Field Display

- [ ] Price/decimal fields in edit forms use `parseFloat(value.toFixed(2))` for initialization
- [ ] Price/decimal columns use `.toFixed(2)` or `toLocaleString()` for display
- [ ] Never display raw float64 values (may show precision artifacts like `23.989999771118164`)

### 10. FK Dropdown Integration

For entities with foreign key relationships:

- [ ] Form fetches related records via `useEffect` with appropriate API call
- [ ] Uses `<Select>` when records are available, `<Input>` fallback when not
- [ ] Sets BOTH FK ID (`relatedId`) and display name on selection
- [ ] Handles `res?.data?.data || res?.data || []` for nested paginated response
- [ ] Edit form pre-selects the current FK value via `value={formData.relatedId?.toString() || ''}`

### 9. Accessibility

- [ ] All inputs have matching `htmlFor` / `id` pairs
- [ ] Labels are descriptive (not just "Field 1")
- [ ] Error messages are associated with fields
- [ ] Form prevents default submit: `onSubmit={(e) => e.preventDefault()}`

## Common i18n Issues

### Missing namespace registration
```typescript
// locales/index.ts must have:
import entities from './en/entities.json';
ns: [..., 'entities'],
resources: { en: { ..., entities } },
```

### Wrong i18n pattern for file type
```tsx
// BAD: using TFunction in a React component
export function EntityDetailView({ t }: { t: TFunction }) { ... }

// GOOD: using useTranslation hook in React components
export function EntityDetailView() {
    const { t } = useTranslation('entities');
}

// BAD: using useTranslation in a plain function
export function getEntityColumns(): CrudColumn[] {
    const { t } = useTranslation('entities');  // Hook rules violation!
}

// GOOD: accepting TFunction as parameter in plain functions
export function getEntityColumns(t: TFunction): CrudColumn[] { ... }
```

### Hardcoded strings that should be translated
```tsx
// BAD
<Label>Name *</Label>
onSuccess('Entity created successfully');
if (confirm('Are you sure?'))

// GOOD
<Label>{t('form.name')}</Label>
onSuccess(t('toast.created'));
if (confirm(t('confirm.bulkDelete', { count: ids.length })))
```

### Missing interpolation variables
```json
// BAD: hardcoded count
"bulkDelete": "Are you sure you want to delete 5 items?"

// GOOD: uses interpolation
"bulkDelete": "Are you sure you want to delete {{count}} item(s)?"
```

## Output

After review, produce a report with:
1. **Pass/Fail** for each checklist item
2. **i18n gaps** — missing translation keys or unregistered namespace
3. **Hardcoded strings** — any user-visible text not going through `t()`
4. **Pattern violations** — wrong i18n usage (hook vs TFunction)
5. **Issues found** with specific line numbers
6. **Recommended fixes** with code snippets
7. **Dropdown opportunities** if any free-text fields should be constrained

## Reference

See `resources/js/pages/Books/` for the canonical i18n-aware pattern.
See `resources/js/locales/en/books.json` for complete translation key structure.
