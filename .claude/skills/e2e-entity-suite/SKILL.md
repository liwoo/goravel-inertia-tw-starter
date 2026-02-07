---
name: e2e-entity-suite
description: Run a comprehensive E2E browser test suite for a newly scaffolded Goravel entity. Tests navigation, CRUD operations, search, filters, row actions, form validation, FK dropdowns, and cross-entity integration using Playwright MCP.
argument-hint: "[EntityName] [base-url]"
allowed-tools: Read, Grep, Glob, Bash, mcp__playwright__*
---

# E2E Entity Test Suite (Playwright MCP)

**Agent**: Blessings Phiri — Senior QA & Reliability Engineer

Run the complete E2E test suite for `$ARGUMENTS`.

## Prerequisites

- Dev server running: `go run .` (backend) + `pnpm dev` (Vite frontend)
- Dev database migrated: `go run . artisan migrate`
- **Seed data exists**: Run `/fake-data` skill first, then `go run . artisan db:seed --seeder=EntitySeeder` to populate 25+ records. Without seed data, sorting, pagination, and filter tests are meaningless.
- Playwright MCP configured in project `.claude.json`
- Backend CRUD tests passing before UI testing

## Pre-Flight: Verify Seed Data

Before starting browser tests, verify the entity has sufficient data:

```bash
# Check record count via API (should return 25+ records)
# Or use the browser to navigate to the page and check the stats card / total count
```

If the page shows 0 records or the table is empty:
1. Run `/fake-data EntityName` to create the seeder
2. Run `go run . artisan db:seed --seeder=EntitySeeder`
3. Refresh the page

## Test Protocol

For each test: **execute the action**, **inspect with `browser_snapshot`**, **record PASS/FAIL**, and **note any issues**.

Use `browser_snapshot` (accessibility tree) for assertions — it's faster and more reliable than screenshots.

---

## Phase 1: Login & Setup

### Test 1.1: Login

```
1. browser_navigate → http://localhost:5173/login
2. browser_snapshot → find email/password fields
3. browser_type → email field with admin credentials
4. browser_type → password field
5. browser_click → Sign In button
6. browser_wait_for → dashboard or admin page loads
7. browser_snapshot → verify logged in (sidebar visible, user name in header)
```

**PASS criteria**: Redirected to admin area, sidebar navigation visible.

---

## Phase 2: Navigation

### Test 2.1: Sidebar Navigation Entry

```
1. browser_snapshot → find entity nav item in sidebar
2. Verify: entity appears in correct position with icon
3. browser_click → entity nav item
4. browser_wait_for → page title text
5. browser_snapshot → verify URL is /admin/entity-names
```

**PASS criteria**: Nav item visible, click navigates to correct URL.

### Test 2.2: Page Load & Structure

```
1. browser_snapshot → verify page structure:
   - Page title (from i18n, NOT raw key like "page.title")
   - "Add entity" button
   - Stats cards (if enabled)
   - Filter tabs (All, Active, Inactive, etc.)
   - Data table with column headers
   - Pagination or record count
```

**PASS criteria**: All structural elements present, no raw i18n keys visible.

**Watch for**:
- i18n keys rendered as literals (e.g., `page.title` instead of "Authors")
- Missing stats cards (backend StatsBuilder not configured)
- Empty table (dev DB not seeded or migrations not run)

---

## Phase 3: Stats & Filters

### Test 3.1: Stats Cards

```
1. browser_snapshot → find stats cards
2. Verify: total count matches table row count
3. Verify: status-specific counts (active, inactive) are numbers
```

**PASS criteria**: Stats cards display with numeric values.

### Test 3.2: Filter Tabs

```
1. browser_snapshot → find filter tab buttons (All, Active, Inactive, etc.)
2. browser_click → "Inactive" tab (or a non-default filter)
3. browser_wait_for → table to update
4. browser_snapshot → verify:
   - URL updated with ?status=INACTIVE (or relevant filter param)
   - Table shows only matching records (or "No results" message)
   - Active tab is visually highlighted
5. browser_click → "All" tab to reset
6. browser_snapshot → verify all records shown again
```

**PASS criteria**: Filter tabs update URL, table contents change accordingly.

**Watch for**:
- Filter not appending query param to URL
- Table not refreshing after filter change
- "All" tab not clearing the filter

---

## Phase 4: Create Operation

### Test 4.1: Open Create Form

```
1. browser_click → "Add entity" button
2. browser_snapshot → verify create form dialog/drawer:
   - Section headers (e.g., "Entity Information", "Contact Details")
   - All field labels with correct text (from i18n)
   - Required field indicators (*)
   - Input types match field types (text, email, date, number, textarea, select)
   - Submit/Save button
```

**PASS criteria**: Form opens with all expected fields and labels.

### Test 4.2: Fill and Submit Create Form

```
1. Fill all required fields with valid test data
2. Fill optional fields (at least one date, one nullable text)
3. For enum/status fields: browser_click → Select, choose an option
4. browser_click → Submit/Save button
5. browser_wait_for → success toast message (e.g., "created successfully")
6. browser_snapshot → verify:
   - Form closed/dismissed
   - New record appears in table
   - Stats updated (total count incremented)
```

**PASS criteria**: Record created, toast shown, table updated.

**Watch for**:
- **Empty string date error**: if submit fails with 500, check if nullable date/numeric fields send `""` instead of `null`
- **snake_case binding failure**: if all fields are empty after creation, check if API request uses camelCase instead of snake_case keys
- **Missing dev DB migration**: if submit returns 500 "relation does not exist", run `go run . artisan migrate`

### Test 4.3: Verify Created Record

```
1. browser_snapshot → find the newly created record in table
2. Verify: all visible columns show correct data
3. Verify: status badge displays correctly
```

**PASS criteria**: Record visible with correct data in all columns.

---

## Phase 5: Detail View

### Test 5.1: Open Detail View

```
1. browser_click → "Open menu" button on the created record's row
2. browser_snapshot → verify action menu: View, Edit, Delete
3. browser_click → "View" option
4. browser_snapshot → verify detail view:
   - All fields displayed with labels
   - Status shown as badge
   - Nullable empty fields show "Not specified" (or equivalent)
   - Metadata section (ID, Created, Last Updated)
   - Dates formatted correctly
```

**PASS criteria**: Detail view shows all fields with proper formatting.

**Watch for**:
- **Float precision**: price/decimal fields showing `23.989999771118164` instead of `23.99`
- Raw i18n keys instead of translated labels
- Missing metadata section

### Test 5.2: Close Detail View

```
1. browser_press_key → Escape (or click close button)
2. browser_snapshot → verify detail view closed, table visible again
```

---

## Phase 6: Edit Operation

### Test 6.1: Open Edit Form

```
1. browser_click → "Open menu" button on the record's row
2. browser_click → "Edit" option
3. browser_snapshot → verify edit form:
   - All fields pre-populated with current values
   - Same fields as create form
   - Metadata section present (ID, Created, Updated)
   - Status dropdown shows current value
```

**PASS criteria**: Form opens pre-populated with correct data.

**Watch for**:
- **Float precision in inputs**: price field showing `23.989999771118164`
- Date fields not formatted as YYYY-MM-DD
- FK dropdown showing "Select..." instead of current value (expected for legacy records without FK set)

### Test 6.2: Modify and Save

```
1. Clear and modify at least one required field
2. Change status dropdown to a different value
3. browser_click → Save/Update button
4. browser_wait_for → success toast (e.g., "updated successfully")
5. browser_snapshot → verify:
   - Form closed
   - Table shows updated values
```

**PASS criteria**: Record updated, toast shown, table reflects changes.

---

## Phase 7: Table Search

### Test 7.1: Search — Match

```
1. browser_click → search input in the table toolbar
2. browser_type → a known value (e.g., part of the created record's name)
3. browser_wait_for → table to filter
4. browser_snapshot → verify:
   - Table shows only matching records
   - Result count updated
```

**PASS criteria**: Search returns matching records.

### Test 7.2: Search — No Match

```
1. Clear search field
2. browser_type → "xyznonexistent12345"
3. browser_wait_for → table to update
4. browser_snapshot → verify:
   - "No results" or empty table message displayed
   - Record count shows 0
5. Clear search field to restore full list
```

**PASS criteria**: No results message shown for non-matching query.

---

## Phase 8: Column Sorting

**Requires**: 25+ seeded records (run `/fake-data` first)

### Test 8.1: Sort by Primary Column (Ascending)

```
1. browser_snapshot → note the current first row's primary field value (e.g., name)
2. browser_click → sortable column header (e.g., "Name" or primary field)
3. browser_wait_for → table to refresh
4. browser_snapshot → verify:
   - Sort indicator (arrow) appears on the clicked column
   - First row shows the alphabetically/numerically lowest value
   - Rows are in ascending order (compare first 3-5 rows)
   - URL may update with ?sort=field&direction=ASC
```

**PASS criteria**: Column sorts ascending, sort indicator visible.

### Test 8.2: Sort by Primary Column (Descending)

```
1. browser_click → same column header again (toggle to descending)
2. browser_wait_for → table to refresh
3. browser_snapshot → verify:
   - Sort indicator shows descending direction
   - First row shows the alphabetically/numerically highest value
   - Order is reversed from ascending
```

**PASS criteria**: Column sorts descending on second click.

### Test 8.3: Sort by Date Column

```
1. browser_click → "Created" column header (or another date column)
2. browser_wait_for → table to refresh
3. browser_snapshot → verify:
   - Dates are in chronological order
   - Sort indicator on date column
```

**PASS criteria**: Date column sorts correctly.

**Watch for**:
- Sort not working → check column definition has `sortable: true`
- Sort indicator missing → check CrudPage column rendering
- Records don't change order → check API handles `sort` and `direction` query params

---

## Phase 9: Pagination

**Requires**: 25+ seeded records (more than one page worth)

### Test 9.1: Verify Pagination Controls

```
1. browser_navigate → /admin/entity-names (reset to first page)
2. browser_snapshot → verify:
   - Pagination controls visible (page numbers or Next/Previous buttons)
   - Current page highlighted (page 1)
   - Total record count shown (e.g., "Showing 1-10 of 30")
   - Table shows correct number of rows per page (typically 10)
```

**PASS criteria**: Pagination controls present with correct counts.

### Test 9.2: Navigate to Next Page

```
1. browser_click → "Next" button or page 2
2. browser_wait_for → table to refresh
3. browser_snapshot → verify:
   - Table shows different records than page 1
   - Current page indicator updated to page 2
   - URL updated with ?page=2
   - "Previous" button now enabled
```

**PASS criteria**: Second page shows different records, page indicator updated.

### Test 9.3: Navigate Back to First Page

```
1. browser_click → "Previous" button or page 1
2. browser_wait_for → table to refresh
3. browser_snapshot → verify:
   - Original first-page records displayed
   - Page indicator back to 1
```

**PASS criteria**: First page restored correctly.

**Watch for**:
- No pagination shown → table has fewer records than page size (need more seed data)
- Page 2 shows same records → API not handling `page` query param
- Total count wrong → stats not refreshing with pagination

---

## Phase 10: Global Search (CMD+K)

### Test 10.1: Open Global Search

```
1. browser_press_key → Meta+k
2. browser_snapshot → verify search dialog:
   - Search input with placeholder
   - Quick access section showing entity type (if user has permission)
```

### Test 10.2: Search for Entity

```
1. browser_type → search query matching the created record
2. browser_wait_for → search results (or short delay)
3. browser_snapshot → verify:
   - Results include the entity with correct type badge
   - Title and subtitle display correctly
   - Entity type badge has correct color
4. browser_press_key → Escape to close
```

**PASS criteria**: Entity appears in global search with correct badge and data.

**Watch for**:
- Entity not appearing → check `search_controller.go` has the search method and permission check
- Missing type badge → check `search_config.tsx` has the entity type in `SEARCH_ENTITIES`
- Empty results → check entity service has `WithSearchFields` configured

---

## Phase 11: Row Actions

### Test 11.1: Action Menu

```
1. browser_click → "Open menu" (three dots) on a table row
2. browser_snapshot → verify menu items:
   - "View" option
   - "Edit" option
   - Separator
   - "Delete" option (in destructive/red style)
3. browser_press_key → Escape to close menu
```

**PASS criteria**: All action menu items present.

### Test 11.2: Delete — Cancel

```
1. browser_click → "Open menu" on a row
2. browser_click → "Delete" option
3. browser_handle_dialog → accept: false (Cancel the confirmation)
4. browser_snapshot → verify record still in table (not deleted)
```

**PASS criteria**: Cancel prevents deletion, record remains.

### Test 11.3: Delete — Confirm (optional, use test record)

Only run this on a test record you intend to discard:

```
1. browser_click → "Open menu" on the test row
2. browser_click → "Delete" option
3. browser_handle_dialog → accept: true (Confirm deletion)
4. browser_wait_for → table to update
5. browser_snapshot → verify:
   - Record removed from table
   - Stats count decremented
   - Success toast or table refresh
```

**PASS criteria**: Record deleted, table and stats updated.

---

## Phase 12: Cross-Entity Integration (Conditional)

**Skip this phase** if the entity is NOT referenced as a foreign key by another entity. Only run when another entity has an FK dropdown pointing to this one (e.g., Author → Book's author dropdown).

### Test 12.1: FK Dropdown in Related Entity

```
1. browser_navigate → /admin/related-entity (e.g., /admin/books)
2. browser_click → "Add" button to open create form
3. browser_snapshot → verify FK field is a dropdown (Select), NOT a text input
4. browser_click → the FK dropdown
5. browser_snapshot → verify:
   - Options load from API (entity names appear in dropdown)
   - The entity created in Phase 4 appears as an option
6. browser_press_key → Escape to close
```

**PASS criteria**: FK dropdown loads related entity records.

### Test 12.2: FK Dropdown in Edit Form

```
1. Navigate to related entity list, open edit form for an existing record
2. browser_snapshot → verify FK dropdown:
   - Present as Select (not text input)
   - Pre-selected with current FK value (if set)
   - Or shows "Select..." placeholder (if FK is null for legacy records)
```

**PASS criteria**: FK dropdown present and functional in edit form.

**Watch for**:
- Dropdown shows "Select..." for all records → legacy records have null FK (expected)
- Dropdown falls back to text input → API returned empty array (check `/api/entity-names` endpoint)
- API returns `first_name` but JS reads `firstName` → check model JSON tags match

---

## Phase 13: Responsive Layout (Optional)

### Test 13.1: Mobile View

```
1. browser_resize → { width: 375, height: 812 }
2. browser_snapshot → verify:
   - Mobile columns render (compact layout)
   - Sidebar collapsed to hamburger menu
   - Table is scrollable or stacked
3. browser_resize → { width: 1280, height: 720 } (restore desktop)
```

---

## Phase 14: Console & Network Check

### Test 14.1: Console Errors

```
1. browser_console_messages → level: "error"
2. Verify: no React errors, no 404/500 errors, no i18n missing key warnings
```

### Test 14.2: Network Failures

```
1. browser_network_requests → includeStatic: false
2. Verify: no failed API calls (4xx, 5xx status codes)
```

---

## Bug Watch Checklist

These are known issues to specifically check for during testing:

| Bug | Symptom | Root Cause | Fix |
|-----|---------|------------|-----|
| Nullable field `""` | 500 on create/edit with empty date/numeric | Empty string sent instead of `null` | Use `\|\| null` in form requestData |
| Float precision | Price shows `23.989999` | Raw float64 from JSON | Use `.toFixed(2)` in display |
| Missing dev migration | 500 "relation does not exist" | Dev DB not migrated | Run `go run . artisan migrate` |
| Snake_case binding | All fields empty after create | camelCase keys in request body | Convert to snake_case in requestData |
| Search not working | Entity missing from CMD+K | Missing search method/permission | Add to `search_controller.go` + `search_config.tsx` |
| FK dropdown empty | Text input instead of Select | API returns empty array | Check `/api/entity-names` returns data |
| i18n raw keys | "page.title" shown literally | Namespace not registered | Add to `locales/index.ts` |
| Filter not working | Tab click does nothing | Missing URL query param handling | Check filter config and API param support |

---

## Test Report Format

After completing all phases, produce a summary:

```
## E2E Test Report: [EntityName]

### Environment
- URL: http://localhost:5173
- Date: YYYY-MM-DD
- Browser: Chromium (Playwright MCP)

### Results

| # | Test | Result | Notes |
|---|------|--------|-------|
| 1.1 | Login | PASS/FAIL | |
| 2.1 | Sidebar Navigation | PASS/FAIL | |
| 2.2 | Page Load & Structure | PASS/FAIL | |
| 3.1 | Stats Cards | PASS/FAIL | |
| 3.2 | Filter Tabs | PASS/FAIL | |
| 4.1 | Open Create Form | PASS/FAIL | |
| 4.2 | Create Record | PASS/FAIL | |
| 4.3 | Verify Created Record | PASS/FAIL | |
| 5.1 | Detail View | PASS/FAIL | |
| 5.2 | Close Detail View | PASS/FAIL | |
| 6.1 | Open Edit Form | PASS/FAIL | |
| 6.2 | Edit Record | PASS/FAIL | |
| 7.1 | Table Search (match) | PASS/FAIL | |
| 7.2 | Table Search (no match) | PASS/FAIL | |
| 8.1 | Sort Ascending | PASS/FAIL | |
| 8.2 | Sort Descending | PASS/FAIL | |
| 8.3 | Sort by Date | PASS/FAIL | |
| 9.1 | Pagination Controls | PASS/FAIL | |
| 9.2 | Next Page | PASS/FAIL | |
| 9.3 | Previous Page | PASS/FAIL | |
| 10.1 | Global Search Open | PASS/FAIL | |
| 10.2 | Global Search Results | PASS/FAIL | |
| 11.1 | Row Action Menu | PASS/FAIL | |
| 11.2 | Delete Cancel | PASS/FAIL | |
| 11.3 | Delete Confirm | PASS/FAIL/SKIP | |
| 12.1 | FK Dropdown (create) | PASS/FAIL/N/A | |
| 12.2 | FK Dropdown (edit) | PASS/FAIL/N/A | |
| 13.1 | Mobile Layout | PASS/FAIL/SKIP | |
| 14.1 | Console Errors | PASS/FAIL | |
| 14.2 | Network Failures | PASS/FAIL | |

### Issues Found
1. [Description, severity, suggested fix]

### Bugs Fixed During Testing
1. [What was wrong → what was fixed]
```

---

## Cleanup

After testing:

```
1. browser_close → Close the browser
2. If test records were created, note whether to keep or remove them
```

## Reference

- Playwright MCP tools guide: `/playwright-ui-test` skill
- CrudPage component: `resources/js/components/Crud/CrudPage.tsx`
- Global search: `resources/js/components/GlobalSearch.tsx` + `resources/js/config/search_config.tsx`
- Admin layout: `resources/js/layouts/Admin.tsx`
- Login: `resources/js/components/login-form.tsx`
