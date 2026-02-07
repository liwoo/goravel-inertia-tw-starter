---
name: playwright-ui-test
description: Use the Playwright MCP to verify UI integrity and functionality of the Goravel Inertia app. Tests page rendering, navigation, forms, CRUD operations, global search, and responsive layout via browser automation.
argument-hint: "[PageName or check-type]"
allowed-tools: Read, Grep, Glob, mcp__playwright__*
---

# Playwright UI Testing (MCP)

Verify UI for `$ARGUMENTS`.

## Prerequisites

- Dev server running: `go run .` (backend on :3000) + `pnpm dev` (Vite on :5173)
- Playwright MCP configured in `.claude.json`:
  ```json
  "playwright": {
    "type": "stdio",
    "command": "npx",
    "args": ["@playwright/mcp@latest"]
  }
  ```

## Available MCP Tools

### Navigation
| Tool | Purpose |
|---|---|
| `browser_navigate` | Go to a URL |
| `browser_navigate_back` | Browser back button |
| `browser_tabs` | List, create, close, select tabs |

### Page Inspection
| Tool | Purpose |
|---|---|
| `browser_snapshot` | Accessibility tree of current page (primary inspection tool) |
| `browser_take_screenshot` | Visual screenshot |
| `browser_console_messages` | JavaScript console output |
| `browser_network_requests` | All network requests since page load |

### Interaction
| Tool | Purpose |
|---|---|
| `browser_click` | Click an element (by accessibility ref or selector) |
| `browser_type` | Type text into an input |
| `browser_fill_form` | Fill multiple form fields at once |
| `browser_select_option` | Select dropdown option |
| `browser_press_key` | Press keyboard key (Enter, Escape, Tab, etc.) |
| `browser_hover` | Hover over element |
| `browser_drag` | Drag and drop between elements |
| `browser_file_upload` | Upload files |
| `browser_handle_dialog` | Accept/dismiss dialogs |

### Assertions
| Tool | Purpose |
|---|---|
| `browser_verify_text_visible` | Verify text appears on page |
| `browser_verify_element_visible` | Verify element is visible |
| `browser_verify_list_visible` | Verify a list of items is visible |
| `browser_verify_value` | Verify input/element value |

### Utilities
| Tool | Purpose |
|---|---|
| `browser_wait_for` | Wait for text or timeout |
| `browser_evaluate` | Run JavaScript in page context |
| `browser_resize` | Resize viewport (responsive testing) |
| `browser_close` | Close the browser |

## Standard Workflow

### Step 1: Navigate to App

```
browser_navigate → http://localhost:5173/login
```

### Step 2: Login

```
browser_snapshot   → Find email/password fields
browser_fill_form  → { email: "admin@example.com", password: "password" }
browser_click      → "Sign In" button
browser_wait_for   → Wait for dashboard text or redirect
```

### Step 3: Navigate to Target Page

```
browser_navigate → http://localhost:5173/admin/books
browser_snapshot → Inspect page structure
```

## UI Integrity Checks

### 1. Page Structure Check

Verify the admin layout renders correctly:

```
browser_navigate → /admin/{page-name}
browser_snapshot → Check for:
```

**Expected elements on every admin page:**
- Sidebar navigation (AppSidebar)
- Site header with page title
- Footer with copyright and version
- Main content area

**Verify with:**
```
browser_verify_text_visible → page title (from i18n)
browser_verify_element_visible → sidebar
browser_verify_text_visible → copyright text in footer
```

### 2. Navigation Check

Verify sidebar links work:

```
browser_snapshot → Find nav items in sidebar
browser_click → target nav item
browser_wait_for → new page title
browser_verify_text_visible → expected page heading
```

### 3. CRUD Page Check

For entity pages using `CrudPage`:

**List view:**
```
browser_navigate → /admin/entity-names
browser_snapshot → Verify:
  - Page title
  - "Add" button (if canCreate)
  - Data table with columns
  - Pagination controls
  - Search input
browser_verify_text_visible → column headers
```

**Create form:**
```
browser_click → "Add" / "Create" button
browser_snapshot → Verify form opened with:
  - All required field labels
  - Input fields, dropdowns, textareas
  - Submit button
  - Cancel button
```

**Fill and submit:**
```
browser_fill_form → { field1: "value1", field2: "value2" }
browser_select_option → status dropdown
browser_click → Submit button
browser_wait_for → success toast message
```

**Detail view:**
```
browser_click → row in table (or eye icon)
browser_snapshot → Verify:
  - All field labels and values displayed
  - Status badge
  - Metadata (ID, created, updated)
  - Edit/Delete action buttons
```

**Edit form:**
```
browser_click → Edit button
browser_snapshot → Verify:
  - Form pre-populated with current values
  - Same fields as create form
  - Metadata section (ID, created, updated)
browser_fill_form → { field1: "updated value" }
browser_click → Save/Update button
browser_wait_for → success toast
```

**Delete:**
```
browser_click → Delete button
browser_handle_dialog → accept confirmation
browser_wait_for → deletion toast or row disappears
```

### 4. Global Search (CMD+K)

```
browser_press_key → Meta+k (Mac) or Control+k
browser_snapshot → Verify search dialog opened:
  - Search input with placeholder
  - Quick access section with entity links
browser_type → search query text
browser_wait_for → search results
browser_snapshot → Verify:
  - Results with entity type badges
  - Result titles and subtitles
  - Keyboard navigation hints (↑↓ Enter ESC)
browser_click → a search result
browser_wait_for → navigation to result page
```

### 5. Responsive Layout Check

```
# Desktop
browser_resize → { width: 1280, height: 720 }
browser_snapshot → Verify full sidebar, desktop columns

# Tablet
browser_resize → { width: 768, height: 1024 }
browser_snapshot → Verify sidebar collapsed, responsive table

# Mobile
browser_resize → { width: 375, height: 812 }
browser_snapshot → Verify mobile layout, hamburger menu, mobile columns
```

### 6. Form Validation Check

```
# Submit empty form
browser_click → Submit button (without filling required fields)
browser_snapshot → Verify:
  - Validation error messages appear
  - Error borders on required fields
  - Form does NOT submit

# Submit invalid data
browser_fill_form → { email: "not-an-email", price: "abc" }
browser_click → Submit button
browser_snapshot → Verify field-level error messages
```

### 7. Permission Gating Check

```
# Login as user without permissions
browser_navigate → /login
browser_fill_form → { email: "limited@example.com", password: "password" }
browser_click → Sign In

# Verify restricted items
browser_navigate → /admin/entity-names
browser_snapshot → Verify:
  - No "Add" button (if no create permission)
  - No edit/delete actions (if no update/delete permission)
  - Or: 403 page if no read permission

# Verify nav items hidden
browser_snapshot → Check sidebar doesn't show restricted nav items
```

### 8. i18n Verification

Check that no raw translation keys appear on the page:

```
browser_snapshot → Scan for untranslated keys like:
  - "page.title" (literal key instead of value)
  - "form.name" or "columns.status"
  - Any text starting with common i18n prefixes
```

## Entity-Specific Checks

### For a New Entity Page

Run these checks after scaffolding with `/inertia-scaffold`:

1. **Page loads**: Navigate to `/admin/entity-names`, verify no errors
2. **Title correct**: Page title from i18n renders
3. **Columns display**: Table columns match expected headers
4. **Create form**: All fields present, dropdowns populated, validation works
5. **Edit form**: Pre-populated with data, same fields as create
6. **Detail view**: All fields displayed with labels
7. **Search**: Entity appears in CMD+K results
8. **Nav entry**: Sidebar shows entity link
9. **Responsive**: Mobile columns render on small viewport
10. **No hardcoded strings**: All text from i18n (no raw keys visible)

## Debugging

### Check Console Errors

```
browser_console_messages → Look for:
  - React errors (hooks, key warnings)
  - Network errors (404, 500)
  - TypeScript runtime errors
  - i18n missing key warnings
```

### Check Network Requests

```
browser_network_requests → Look for:
  - Failed API calls (4xx, 5xx)
  - Slow requests (> 2s)
  - Missing endpoints
```

### Evaluate Page State

```
browser_evaluate → document.title
browser_evaluate → JSON.stringify(window.__INERTIA_PAGE__)
```

## Common Patterns

### Login Helper Flow

Every test session starts with login:

```
1. browser_navigate → /login
2. browser_snapshot → find form
3. browser_fill_form → credentials
4. browser_click → submit
5. browser_wait_for → dashboard or redirect
```

### Assert Toast Message

After CRUD operations:

```
browser_wait_for → "created successfully" or "updated successfully"
# Or use snapshot to find Sonner toast
browser_snapshot → look for toast/notification element
```

### Wait for Data Load

After navigation:

```
browser_wait_for → table row text or "No results" message
# Or wait for specific content
browser_verify_text_visible → expected content
```

## Verification Report

After completing checks, produce:

1. **Pages tested**: List of URLs visited
2. **Pass/Fail**: For each check category
3. **Console errors**: Any JS errors found
4. **Missing elements**: Expected UI not found
5. **i18n issues**: Untranslated keys visible
6. **Responsive issues**: Layout problems at different viewports
7. **Screenshots**: Key states captured

## Reference

- Admin layout: `resources/js/layouts/Admin.tsx`
- Sidebar: `resources/js/components/app-sidebar.tsx`
- CrudPage: `resources/js/components/Crud/CrudPage.tsx`
- Global search: `resources/js/components/GlobalSearch.tsx`
- Login form: `resources/js/components/login-form.tsx`
- Books (canonical page): `resources/js/pages/Books/Index.tsx`
