---
name: inertia-scaffold
description: Full UI scaffold orchestrator - generates all frontend components for a Goravel entity with i18n (types, translations, page, columns, forms, detail view, page config, page controller).
argument-hint: "[EntityName]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
disable-model-invocation: true
---

# Full Inertia UI Scaffold (i18n-aware)

Generate complete frontend UI for `$ARGUMENTS`.

## Prerequisites

Before running this scaffold, ensure the backend is ready AND tested:
- [ ] Migration created and run (`/goravel-crud-migration`)
- [ ] Model generated and fixed (`/goravel-crud-model`)
- [ ] Service layer created (`/goravel-crud-service`)
- [ ] Permissions registered (`/goravel-crud-permissions`)
- [ ] Request validators created (`/goravel-crud-request`)
- [ ] API controller created (`/goravel-crud-controller`)
- [ ] API routes registered (`/goravel-crud-routes`)
- [ ] **CRUD tests generated, fixed, and ALL passing** (`/goravel-crud-test`)

**CRITICAL**: Do NOT start UI work until all CRUD tests pass. Backend bugs (Bind, validation key mismatches, GORM mapping) are much harder to debug through the UI.

## Scaffold Sequence

### Step 1: TypeScript Types & i18n Namespace
**Skill**: `/inertia-types`

Create `resources/js/types/<entity_name>.ts` with:
- Entity interface (matching Go model, dual-case fields)
- CreateData interface
- UpdateData interface
- ListResponse type alias
- Stats interface

Create `resources/js/locales/en/<entities>.json` with all sections:
- `page`, `stats`, `status`, `columns`, `filters`, `actions`, `form`, `validation`, `toast`, `confirm`

Register in `resources/js/locales/index.ts`:
- Import the JSON
- Add to `ns` array
- Add to `resources.en`

**Verify**: Namespace imports without errors.

---

### Step 2: Go Page Controller
**Skill**: `/inertia-page-ctrl`

Run: `go run . artisan make:page-ctrl --controller=$ARGUMENTS`

Configure `app/http/controllers/<entities>/<entities>_page_controller.go`:
- Set `ResourceType`, `PageComponent`, `Service`, `ServiceIdentifier`
- Optionally enable `StatsEnabled` + `StatsBuilder`

Register web route in `routes/web.go`:
```go
router.Get("/admin/entity-names", entityPageController.Index)
```

**Verify**: `go build ./...` compiles.

---

### Step 3: Generate UI Files
**Skill**: `/inertia-page` (uses `make:ui` command)

Run: `go run . artisan make:ui --page=$ARGUMENTS --request=$ARGUMENTS`

This generates the initial file structure under `resources/js/pages/<Entity>/`.

**Verify**: Directory exists with Index.tsx and sections/.

---

### Step 4: Column Definitions
**Skill**: `/inertia-columns`

Create/edit `resources/js/pages/<Entity>/sections/<Entity>Columns.tsx`:
- `getEntityColumns(t: TFunction)` — Desktop columns with translated headers
- `getEntityColumnsMobile(t: TFunction)` — Compact mobile columns
- `getEntityFilters(t: TFunction)` — Filter definitions with translated labels
- Status config function with `t('status.*')` labels

**Verify**: All `t()` keys exist in the JSON namespace.

---

### Step 5: Create Form
**Skill**: `/inertia-form`

Create `resources/js/pages/<Entity>/sections/<Entity>CreateForm.tsx`:
- `useTranslation('entities')` hook
- `t('form.*')` for labels/placeholders
- `t('validation.*')` for error messages
- `t('toast.created')` for success message
- `t('status.*')` for Select dropdown options
- forwardRef + useImperativeHandle pattern

**Verify**: All `t()` keys exist in the JSON namespace.

---

### Step 6: Edit Form
**Skill**: `/inertia-form`

Create `resources/js/pages/<Entity>/sections/<Entity>EditForm.tsx`:
- Same i18n pattern as create form
- `t('toast.updated')` for success message
- Metadata section with `t('form.metadata')`, `t('form.entityId')`, `t('form.created')`, `t('form.lastUpdated')`

**Verify**: Same fields and `t()` keys as create form.

---

### Step 7: Detail View
**Skill**: `/inertia-detail`

Create `resources/js/pages/<Entity>/sections/<Entity>DetailView.tsx`:
- `useTranslation('entities')` hook
- Reuses `t('form.*')` keys (strip `*` suffix with `.replace(' *', '')`)
- `t('status.*')` for status badges
- `t('form.notSpecified')` for empty values

**Verify**: All fields from entity displayed with translated labels.

---

### Step 8: Page Config
**Skill**: `/inertia-page-config`

Create `resources/js/pages/<Entity>/sections/<Entity>PageConfig.tsx`:
- `getEntityStatsConfigs(t: TFunction)` — Stats card titles
- `getEntitySimpleFilters(t: TFunction, stats)` — Filter labels with badges
- `getEntityPageActions(t: TFunction, permissions, handlers)` — Action labels
- `getEntityBulkActions(t: TFunction)` — Confirm dialogs with `{{count}}`

**Verify**: All `t()` keys exist in the JSON namespace.

---

### Step 9: Barrel Export
Create `resources/js/pages/<Entity>/sections/index.ts`:
```typescript
export { EntityDetailView } from './EntityDetailView';
export { EntityCreateForm } from './EntityCreateForm';
export { EntityEditForm } from './EntityEditForm';
export { getEntityColumns, getEntityColumnsMobile, getEntityFilters } from './EntityColumns';
export { getEntityStatsConfigs, getEntitySimpleFilters, getEntityPageActions, getEntityBulkActions } from './EntityPageConfig';
```

**Verify**: All exports resolve without errors.

---

### Step 10: Wire Up Index Page
**Skill**: `/inertia-page`

Update `resources/js/pages/<Entity>/Index.tsx`:
- `useTranslation('entities')` hook
- `t('page.title')`, `t('page.headTitle')`, `t('page.myEntities')`
- Pass `t` to all config functions: `getEntityColumns(t)`, `getEntityFilters(t)`, etc.
- `renderStatsCards(stats, getEntityStatsConfigs(t))`
- `createSimpleFilters(getEntitySimpleFilters(t, stats))`
- `createPageActions(getEntityPageActions(t, permissions, handlers))`

**Verify**: Page renders without TypeScript errors.

---

### Step 11: Navigation
**Skill**: `/goravel-crud-nav`

Add to `resources/js/config/navigation.ts`:
```typescript
{
    title: "main.entities",  // i18n key for nav namespace
    url: "/admin/entity-names",
    icon: EntityIcon,
    requiredService: "entities",
    requiredAction: "read" as const,
}
```

Add to `resources/js/locales/en/nav.json`:
```json
{
  "main": {
    "entities": "Entities"
  }
}
```

**Verify**: Entity appears in sidebar with correct permission gating.

---

### Step 12: Global Search (CMD+K)
**Skill**: `/goravel-crud-search`

Add to backend `search_controller.go`:
- Permission check in `GlobalSearch()` using `auth.ServiceEntity`
- `searchEntities()` method using the entity service

Add to frontend `search_config.tsx`:
- Entity type in `SearchEntityType` union
- Entity config in `SEARCH_ENTITIES` array (icon, colors, permissions)

**Verify**: CMD+K search returns results for the new entity.

---

### Step 13: Form Review
**Skill**: `/inertia-form-review`

Run audit on all frontend files to verify:
- i18n compliance (no hardcoded strings)
- Correct pattern usage (hook vs TFunction)
- Translation key completeness
- Create <> Edit consistency
- Type safety

---

### Step 14: File Downloads & Export
**Skill**: `/file-downloads`

Add export functionality to the entity page:
- Define export fields with formatters (dates, currency, arrays)
- Wire ExportDialog with format picker and field selection
- Add export button to page actions
- Add i18n keys for export labels

---

### Step 15: UX Audit
**Skill**: `/ui-ux-audit`

Review all generated UI for quality:
- Table columns: max 3-4, no icon overuse, no data duplication
- Forms: dropdowns for constrained fields (nationality, gender, etc.), proper validation
- Status colors: consistent palette, dark mode `/30` opacity
- Loading/empty states: translated messages, proper error handling
- Shared config: extract reusable option lists to `config/options.ts`

**Fix issues before proceeding to seeding and E2E testing.**

---

### Step 16: Seed Data
**Skill**: `/fake-data`

Create a database seeder with 25+ realistic records:
- Run `/fake-data EntityName` to create the seeder file
- Register in `database_seeder.go` and `kernel.go`
- Run `go run . artisan db:seed --seeder=EntitySeeder`

**Required for E2E testing** — sorting, pagination, and filter tests need sufficient data.

---

### Step 17: E2E Browser Testing
**Skill**: `/e2e-entity-suite`

Run the full Playwright MCP test suite to verify all UI functionality end-to-end:
- Login, navigation, page structure
- CRUD operations (create, detail view, edit, delete)
- Table search, sorting, pagination
- Global search (CMD+K), filter tabs
- Row action menus, FK dropdowns
- Console errors, network failures

**This step catches bugs that static checks miss** (nullable field errors, float precision, missing migrations on dev DB, search integration gaps).

---

## Post-Scaffold Verification

```bash
# 1. TypeScript compiles
npx tsc --noEmit

# 2. Go compiles
go build ./...

# 3. Dev DB migrated
go run . artisan migrate

# 4. i18n keys complete (manual check)
# Open the translation JSON and verify all sections exist

# 5. Dev server runs
# pnpm dev (manual check)

# 6. E2E browser tests
# /e2e-entity-suite EntityName
```

## i18n Quick Reference

| Component Type | i18n Pattern | Example |
|---------------|-------------|---------|
| React component (forms, detail, Index) | `useTranslation('entities')` hook | `const { t } = useTranslation('entities')` |
| Plain function (columns, filters, config) | `t: TFunction` parameter | `function getColumns(t: TFunction)` |
| Navigation items | i18n key string in config | `title: "main.entities"` |
| Status/enum labels | `t('status.*')` | `t('status.active')` |
| Form labels | `t('form.*')` | `t('form.name')` |
| Validation errors | `t('validation.*')` | `t('validation.nameRequired')` |
| Toast messages | `t('toast.*')` | `t('toast.created')` |
| Interpolation | `{{variable}}` | `t('confirm.bulkDelete', { count })` |

## Quick Reference (File Mapping)

| Step | Skill | File(s) |
|------|-------|---------|
| 1 | `/inertia-types` | `types/<entity>.ts` + `locales/en/<entities>.json` + `locales/index.ts` |
| 2 | `/inertia-page-ctrl` | `controllers/<entities>/<entities>_page_controller.go` |
| 3 | `/inertia-page` | `pages/<Entity>/Index.tsx` |
| 4 | `/inertia-columns` | `sections/<Entity>Columns.tsx` |
| 5-6 | `/inertia-form` | `sections/<Entity>CreateForm.tsx`, `<Entity>EditForm.tsx` |
| 7 | `/inertia-detail` | `sections/<Entity>DetailView.tsx` |
| 8 | `/inertia-page-config` | `sections/<Entity>PageConfig.tsx` |
| 9 | -- | `sections/index.ts` |
| 10 | `/inertia-page` | `pages/<Entity>/Index.tsx` |
| 11 | `/goravel-crud-nav` | `navigation.ts` + `nav.json` |
| 12 | `/goravel-crud-search` | `search_controller.go` + `search_config.tsx` |
| 13 | `/inertia-form-review` | Review only |
| 14 | `/file-downloads` | Export fields + ExportDialog in `Index.tsx` |
| 15 | `/ui-ux-audit` | Columns, forms, colors, loading states audit |
| 16 | `/fake-data` | `database/seeders/<entity>_seeder.go` |
| 17 | `/e2e-entity-suite` | Playwright browser testing (all CRUD + search + filters) |
