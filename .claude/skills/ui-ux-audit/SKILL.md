---
name: ui-ux-audit
description: Audit and refactor an entity's UI for table column density, icon overuse, form dropdown opportunities, validation patterns, status color consistency, loading/empty states, and shared config reusability.
argument-hint: "[EntityName]"
allowed-tools: Read, Write, Edit, Grep, Glob
---

# UI/UX Audit & Refactor

**Agent**: Thoko Nkhoma — Senior Frontend Engineer

Audit and fix the UI for `$ARGUMENTS`.

## Files to Review

1. `resources/js/pages/<Entity>/sections/<Entity>Columns.tsx`
2. `resources/js/pages/<Entity>/sections/<Entity>CreateForm.tsx`
3. `resources/js/pages/<Entity>/sections/<Entity>EditForm.tsx`
4. `resources/js/pages/<Entity>/sections/<Entity>DetailView.tsx`
5. `resources/js/pages/<Entity>/sections/<Entity>PageConfig.tsx`
6. `resources/js/types/<entity_name>.ts`
7. `app/http/requests/<entity>_request.go` (for validation rules & enum values)

---

## Part 1: Table Column Standards

### Rule 1.1: Max 3-4 Visible Desktop Columns

| Column | Purpose | Example |
|--------|---------|---------|
| **Primary (composite)** | Entity identity — name + subtitle | Name + email |
| **Key attribute** | Most important data point | Status badge, price |
| **Secondary attribute** | Supporting context | Category, date |
| **Metadata (optional)** | System info | Created date |

**Audit**: Count columns in `getEntityColumns()`. If > 4, merge related data into composite columns or move low-value columns to detail view only.

### Rule 1.2: Primary Column — Max 2 Lines

```tsx
// GOOD — 2 lines
<div className="flex items-center gap-3">
  <div className="p-2 rounded-lg bg-muted">
    <User className="h-4 w-4 text-muted-foreground" />
  </div>
  <div>
    <div className="font-medium">{entity.name}</div>
    <div className="text-sm text-muted-foreground">{entity.email}</div>
  </div>
</div>

// BAD — 3+ lines stacked
<div>
  <div className="font-medium">{entity.name}</div>
  <div className="text-sm">{entity.email}</div>
  <div className="text-xs">{entity.website}</div>   // move to detail view
</div>
```

### Rule 1.3: No Data Duplication

If a field appears in the composite primary column, don't also give it a standalone column.

### Rule 1.4: Minimize Icons in Data Columns

| Usage | Verdict |
|-------|---------|
| Primary column avatar icon (in muted box) | OK |
| Status badge icon (inside Badge) | OK |
| Calendar icon next to date text | NO — date format is self-evident |
| DollarSign icon next to `$24.99` | NO — the `$` prefix is enough |
| Mail icon next to `user@example.com` | NO — the `@` makes it obvious |
| Hash icon next to ISBN | NO — redundant decoration |

**Rule**: If the data format already communicates its type, the icon adds visual noise.

```tsx
// BAD — icon pollution
<div className="flex items-center">
  <Calendar className="w-3 h-3 mr-1" />
  {date.toLocaleDateString()}
</div>

// GOOD — clean data
<div className="text-sm text-muted-foreground">
  {date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
</div>
```

### Rule 1.5: Mobile — Single Composite Column

One column with entity name, key attribute, and status badge arranged compactly:

```tsx
{
  key: 'name',
  label: t('columns.entity'),
  render: (item) => (
    <div className="flex items-center justify-between">
      <div>
        <div className="font-medium">{item.name}</div>
        <div className="text-sm text-muted-foreground">{item.subtitle}</div>
      </div>
      <StatusBadge status={item.status} />
    </div>
  ),
}
```

---

## Part 2: Form UX Standards

### Rule 2.1: Dropdown-First Approach

**If a field has a finite set of known values, it MUST be a `<Select>` — never a plain text `<Input>`.**

Prevents: typos, inconsistent data, poor filtering, bad UX.

### Dropdown-Worthy Fields Reference

| Field Pattern | Component | Options Source |
|---------------|-----------|---------------|
| status, type, category | `<Select>` | Go enum (`\|in:...` validation) |
| priority | `<Select>` | Low, Medium, High, Critical |
| gender | `<Select>` | Male, Female, Other, Prefer not to say |
| nationality, country | `<Combobox>` | Searchable predefined list |
| language | `<Combobox>` | Searchable predefined list |
| currency | `<Select>` | USD, EUR, GBP, MWK, ZAR |
| education_level | `<Select>` | None → Postgraduate |
| marital_status | `<Select>` | Single, Married, Divorced, Widowed |
| employment_status | `<Select>` | Employed, Self-Employed, Unemployed, Student, Retired |
| sector, industry | `<Combobox>` | Searchable predefined list |
| FK relationship | `<Select>` | Fetched from related entity API |

### Detection Heuristic

A field should be a dropdown if:
1. Go request has `|in:VALUE1,VALUE2,...` validation
2. Field name contains: status, type, category, gender, nationality, country, language, currency, level, role, priority, sector
3. Field has fewer than ~30 possible values
4. Multiple records would share the same value

### Select vs Combobox

- **`<Select>`**: < 15 options (status, gender, priority)
- **`<Combobox>`** (searchable): 15-200 options (nationality, country)
- **`<Input>`**: Only for truly free-form text (name, bio, description, email, URL)

### Rule 2.2: Form Validation Pattern

**Tier 1 — Client-side** (before submit):
```tsx
const newErrors: Record<string, string> = {};
if (!formData.name.trim()) newErrors.name = t('validation.nameRequired');
if (formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
  newErrors.email = t('validation.emailInvalid');
}
if (Object.keys(newErrors).length > 0) { setErrors(newErrors); return; }
```

**Tier 2 — Server-side** (from API response):
```tsx
if (!response.ok) {
  const errorData = await response.json().catch(() => ({}));
  if (errorData.errors) {
    const serverErrors: Record<string, string> = {};
    Object.entries(errorData.errors).forEach(([key, messages]) => {
      const camelKey = key.replace(/_([a-z])/g, (_, c) => c.toUpperCase());
      serverErrors[camelKey] = Array.isArray(messages) ? messages[0] : String(messages);
    });
    setErrors(serverErrors);
  }
  return;
}
```

### Rule 2.3: Consistent Field Layout

- Pairs of short fields: `grid grid-cols-1 md:grid-cols-2 gap-4`
- Full-width: `<Textarea>`, long URL inputs, rich content
- Icon style: pick ONE approach (icons on all fields OR no icons) — don't mix

---

## Part 3: Loading, Empty, and Error States

### Rule 3.1: Loading State

CrudDataTable already accepts a `loading` prop. Ensure the page passes it:

```tsx
<CrudDataTable
  loading={isLoading}
  // ...
/>
```

The table shows a spinner overlay during data fetch. No custom skeleton needed.

### Rule 3.2: Empty State

CrudDataTable accepts `emptyMessage`. Use translated text:

```tsx
<CrudDataTable
  emptyMessage={t('page.noResults')}
  // ...
/>
```

Add to the entity's i18n namespace:
```json
{
  "page": {
    "noResults": "No authors found. Create one to get started."
  }
}
```

**Audit**: Check that `emptyMessage` is set and uses `t()`, not a hardcoded English string.

### Rule 3.3: Form Error States

Every form must handle:

1. **Validation errors** — red border on field + error text below
2. **Network errors** — toast or inline message when fetch fails
3. **Submit button state** — disabled/loading while saving (`setIsSaving`)

```tsx
// Error display per field
{errors.fieldName && (
  <p className="text-sm text-destructive">{errors.fieldName}</p>
)}

// Network error catch
} catch (error) {
  onError?.(error);
  // Optionally show a generic toast
}
```

### Rule 3.4: Stats Cards Fallback

If `StatsEnabled` is true in the page controller but the API returns no stats, the cards should show `0` — not crash or show `undefined`.

---

## Part 4: Status Color Standardization

### The Problem

Every entity defines its own `getStatusConfig()` with ad-hoc color strings. This leads to:
- Inconsistent dark mode patterns (`dark:bg-blue-900` vs `dark:bg-blue-900/30`)
- "ACTIVE" being different greens across entities
- Copy-paste drift between Columns and DetailView for the same entity

### Standard Status Color Palette

All entities MUST use these exact color classes:

| Status Meaning | Light | Dark |
|---------------|-------|------|
| **Positive** (Active, Available, Approved, Enabled) | `bg-green-100 text-green-800` | `dark:bg-green-900/30 dark:text-green-400` |
| **Negative** (Inactive, Disabled, Rejected, Deleted) | `bg-gray-100 text-gray-800` | `dark:bg-gray-900/30 dark:text-gray-400` |
| **Info** (Borrowed, Pending, Processing, In Review) | `bg-blue-100 text-blue-800` | `dark:bg-blue-900/30 dark:text-blue-400` |
| **Warning** (Maintenance, Expiring, Attention) | `bg-orange-100 text-orange-800` | `dark:bg-orange-900/30 dark:text-orange-400` |
| **Accent** (Reserved, Special, Featured) | `bg-purple-100 text-purple-800` | `dark:bg-purple-900/30 dark:text-purple-400` |
| **Danger** (Error, Failed, Overdue, Blocked) | `bg-red-100 text-red-800` | `dark:bg-red-900/30 dark:text-red-400` |
| **Unknown / Fallback** | `bg-gray-100 text-gray-800` | `dark:bg-gray-900/30 dark:text-gray-400` |

### Key Rule: `/30` opacity on dark backgrounds

ALWAYS use `dark:bg-{color}-900/30` (with 30% opacity), NOT `dark:bg-{color}-900` (full opacity). The `/30` keeps dark mode badges subtle rather than glaring.

### Audit Action

1. Grep for `bg-green-100|bg-blue-100|bg-orange-100|bg-purple-100|bg-red-100|bg-gray-100` in entity files
2. Verify dark mode classes use `/30` opacity pattern
3. Verify the same status value gets the same color in Columns AND DetailView
4. Verify fallback/unknown status uses gray

### Future: Shared Status Utility

Consider extracting to a shared utility:

```tsx
// resources/js/config/status-colors.ts
export const STATUS_COLORS = {
  positive: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  negative: 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400',
  info: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  warning: 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400',
  accent: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400',
  danger: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
} as const;
```

This is recommended but not required yet — the audit should flag inconsistencies first.

---

## Part 5: Shared Config & Reusability

### Rule 5.1: Shared Option Lists

Dropdown options that appear across multiple entities should live in a shared file:

```
resources/js/config/options.ts
```

```typescript
// resources/js/config/options.ts

export const NATIONALITIES = [
  { value: 'Malawian', label: 'Malawian' },
  { value: 'Zambian', label: 'Zambian' },
  { value: 'Mozambican', label: 'Mozambican' },
  { value: 'Tanzanian', label: 'Tanzanian' },
  { value: 'South African', label: 'South African' },
  { value: 'Zimbabwean', label: 'Zimbabwean' },
  { value: 'Kenyan', label: 'Kenyan' },
  { value: 'Ugandan', label: 'Ugandan' },
  { value: 'Nigerian', label: 'Nigerian' },
  { value: 'Ghanaian', label: 'Ghanaian' },
  { value: 'British', label: 'British' },
  { value: 'American', label: 'American' },
  { value: 'Indian', label: 'Indian' },
  { value: 'Chinese', label: 'Chinese' },
  { value: 'Other', label: 'Other' },
] as const;

export const GENDERS = [
  { value: 'MALE', label: 'Male' },
  { value: 'FEMALE', label: 'Female' },
  { value: 'OTHER', label: 'Other' },
  { value: 'PREFER_NOT_TO_SAY', label: 'Prefer not to say' },
] as const;

export const EDUCATION_LEVELS = [
  { value: 'NONE', label: 'None' },
  { value: 'PRIMARY', label: 'Primary' },
  { value: 'SECONDARY', label: 'Secondary' },
  { value: 'TERTIARY', label: 'Tertiary' },
  { value: 'POSTGRADUATE', label: 'Postgraduate' },
] as const;

export const MARITAL_STATUSES = [
  { value: 'SINGLE', label: 'Single' },
  { value: 'MARRIED', label: 'Married' },
  { value: 'DIVORCED', label: 'Divorced' },
  { value: 'WIDOWED', label: 'Widowed' },
] as const;

export const EMPLOYMENT_STATUSES = [
  { value: 'EMPLOYED', label: 'Employed' },
  { value: 'SELF_EMPLOYED', label: 'Self-Employed' },
  { value: 'UNEMPLOYED', label: 'Unemployed' },
  { value: 'STUDENT', label: 'Student' },
  { value: 'RETIRED', label: 'Retired' },
] as const;
```

### Rule 5.2: Don't Duplicate Status Configs

Within a single entity, the status config must be defined ONCE and shared between Columns and DetailView:

```tsx
// BAD — duplicated in AuthorColumns.tsx AND AuthorDetailView.tsx
// Each file has its own getAuthorStatusConfig() with potentially different colors

// GOOD — define once, import in both
// resources/js/pages/Authors/sections/authorStatusConfig.ts
export function getAuthorStatusConfig(t: TFunction) {
  return {
    ACTIVE: { label: t('status.active'), icon: <CheckCircle />, color: STATUS_COLORS.positive },
    INACTIVE: { label: t('status.inactive'), icon: <XCircle />, color: STATUS_COLORS.negative },
  };
}
```

### Rule 5.3: When to Extract

Extract to shared config when:
- Same options appear in 2+ entities
- Same color mapping is copy-pasted across files
- Same formatting logic (date, currency) is repeated

---

## Audit Checklist

### Table Columns
- [ ] Desktop columns <= 4
- [ ] Primary column max 2 lines
- [ ] No data duplicated across columns
- [ ] Icons only on: primary avatar, status badges
- [ ] No icons on: dates, prices, emails, plain text
- [ ] Mobile column is single composite
- [ ] Sortable columns marked `sortable: true`

### Form Dropdowns
- [ ] Status/type/category use `<Select>`
- [ ] Nationality/country use `<Combobox>` or `<Select>`
- [ ] Gender/education/marital use `<Select>`
- [ ] FK fields use `<Select>` with API fetch
- [ ] All Go `|in:...` validated fields use `<Select>`
- [ ] No free-text for fields with known values

### Form Validation
- [ ] Required fields validated client-side
- [ ] Email/URL format checked
- [ ] Server errors mapped to form fields
- [ ] Error state: `border-destructive` + text below
- [ ] `setIsSaving` called before/after fetch

### Loading / Empty / Error
- [ ] `loading` prop passed to CrudDataTable
- [ ] `emptyMessage` uses `t()` not hardcoded string
- [ ] Form handles network errors (catch block)
- [ ] Stats cards handle missing/zero data

### Status Colors
- [ ] All dark mode classes use `/30` opacity
- [ ] Same status = same color across Columns and DetailView
- [ ] Fallback uses gray for unknown status
- [ ] Positive = green, Negative = gray, Info = blue, Warning = orange

### Shared Config
- [ ] Dropdown options that repeat across entities live in `config/options.ts`
- [ ] Status config not duplicated between Columns and DetailView
- [ ] Date/currency formatting consistent

---

## Report Format

```
## UI/UX Audit Report: [EntityName]

### Table Columns
| Issue | File:Line | Current | Fix |
|-------|-----------|---------|-----|
| 5 columns visible | Columns:30 | 5 cols | Remove email (in composite) |
| Calendar icon on date | Columns:112 | <Calendar> icon | Remove, date is self-evident |

### Form Dropdowns
| Field | File:Line | Current | Fix |
|-------|-----------|---------|-----|
| nationality | CreateForm:245 | <Input> | <Select> with NATIONALITIES |

### Status Colors
| Issue | File:Line | Current | Fix |
|-------|-----------|---------|-----|
| Missing /30 opacity | DetailView:28 | dark:bg-purple-900 | dark:bg-purple-900/30 |

### Loading/Empty States
| Issue | File:Line | Fix |
|-------|-----------|-----|
| Hardcoded empty msg | Index:45 | Use t('page.noResults') |

### Fixes Applied
1. [What was changed]
```

## Verify

After fixes:
```bash
npx tsc --noEmit
npx eslint "resources/js/pages/<Entity>/**/*.tsx" --max-warnings=0
```

## Reference

- CrudDataTable (loading/empty): `resources/js/components/Crud/CrudDataTable.tsx`
- Book columns (icon overuse example): `resources/js/pages/Books/sections/BookColumns.tsx`
- Author form (text input should be dropdown): `resources/js/pages/Authors/sections/AuthorCreateForm.tsx`
- Shared config location: `resources/js/config/options.ts` (create if missing)
