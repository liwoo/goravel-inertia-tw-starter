# Technical Specification: PR Cleanup & Polish

## Overview

This spec covers the changes needed to clean up the current PR based on the review feedback. The main focus areas are:

- Procurement Management forms
- Events calendar improvements
- BDSP fixes
- Configuration-driven dropdowns across all entities

**Target Deadline:** Tonight/Tomorrow morning for deployment

---

## 1. Procurement Management

### 1.1 Form Field Changes - Convert to Dropdowns

**Affected files:** `resources/js/pages/procurements/` (create/edit forms)

| Field | Current | Required | Data Source |
|-------|---------|----------|-------------|
| `procured_by` | Text input | Dropdown | Config: create new "procured_by" config type |
| `organisation` | Text input | Dropdown | Config: create new "organisations" config type |
| `procurement_type` | Text input | Dropdown | Config: create new "procurement_types" config type |
| `classification` | Text input | Dropdown | Config: create new "classifications" config type |
| `localisation` | Text input | Dropdown | Config: create new "localisations" config type |
| `location` | Text input | Dropdown | Config: create new "locations" config type |
| `partners` | Text input | Dropdown (multi-select) | Config: existing or new "partners" config type |

**Implementation Notes:**
- Fetch config values from `/api/configurations?type={config_type}`
- Use existing dropdown component pattern from the codebase
- Backend still receives strings - no backend changes needed for these fields

---

### 1.2 Reference Number Auto-Generation

**Affected files:**
- Backend: `app/services/procurement_service.go`
- Frontend: Remove reference number input field from create form

**Implementation:**
- Create a reference number generation function
- Suggested format: `PROC-{YEAR}-{SEQUENTIAL_NUMBER}` (e.g., `PROC-2025-00042`)
- Generate on record creation in backend
- Remove the text input from the create form
- Display as read-only in edit/view forms

**Example Backend Logic:**
```go
func GenerateProcurementReference() string {
    year := time.Now().Year()
    // Get next sequential number from DB
    seq := getNextSequence("procurement_references")
    return fmt.Sprintf("PROC-%d-%05d", year, seq)
}
```

---

### 1.3 Remove "Is Published" Toggle - Replace with Publish Button

**Affected files:** Frontend create/edit forms

**Implementation:**
- Default all new procurements to `draft` status
- Remove the "Is Published" toggle from create form
- Add a "Publish" action button (separate from save)
- Button changes status from `draft` to `published`
- Consider placement in detail view header or as action dropdown

**UI Flow:**
1. User creates procurement → saved as draft
2. User reviews procurement → clicks "Publish" button
3. Status changes to published

---

### 1.4 Remove "Interested SMEs" from Create/Edit Forms

**Affected files:** Frontend create/edit forms

**Implementation:**
- Remove `interested_smes` field from create and edit forms entirely
- This field is read-only and populated when SMEs express interest
- Keep it visible in the **view/detail** page only (read-only display)

**Rationale:** SMEs will click "I'm interested" from their own portal view, which populates this field automatically.

---

### 1.5 Split Form into Two Steps (Multi-step Form)

**Affected files:** Frontend create/edit forms and view pages

#### Step 1 - Basic Information
- All the dropdown fields mentioned above
- Date fields (open date, close date)
- Minimum qualifying score
- Reference number (read-only/auto-generated)

#### Step 2 - Details
- `details` field - **Rich text editor** with formatting (bold, italic, lists, etc.)
- `how_to_apply` field - Standard textarea (no rich text needed)

**Implementation:**
- Reference SME forms for multi-step form pattern
- Use a rich text editor component for the `details` field
  - Options: Tiptap, React-Quill, or similar shadcn-compatible editor
  - Must support: bold, italic, bullet lists, numbered lists
- The view page should also be split into tabs matching this structure

**View Page Tabs:**
1. **Overview** - Basic information fields
2. **Details** - Full details and how to apply (rendered HTML)
3. **Interested SMEs** - Read-only list (future implementation)

---

### 1.6 Simplify View UI

**Affected files:** Frontend view/detail pages

**Implementation:**
- Remove unnecessary icons from each field row
- Follow the simpler SME detail view pattern
- Use badges only where they add value (e.g., status badges for draft/published)
- Remove multi-column layouts where unnecessary
- Keep it clean and scannable

**Before (avoid):**
```
📋 Procured By: [icon] Ministry of Trade
📅 Open Date: [icon] 2025-01-15
```

**After (preferred):**
```
Procured By: Ministry of Trade
Open Date: 2025-01-15
Status: [Draft badge]
```

---

### 1.7 Fix Quick Filters

**Investigation needed:** Determine why filters are not working

**Checklist:**
- [ ] Verify filter configuration matches the books example pattern
- [ ] Check that filter keys match backend field names
- [ ] Verify the filter component is properly integrated with the data table
- [ ] Test filter API endpoints return expected results

---

## 2. Events Management

### 2.1 Calendar Styling

**Affected files:** `resources/js/pages/events/` (calendar component)

**Implementation:**
- Wrap calendar in a Card component
- Set card background to white or light shade (not transparent)
- Match internal padding with surrounding UI elements
- Ensure consistent spacing with other page components

**CSS/Styling:**
```tsx
<Card className="bg-white p-4">
  <Calendar ... />
</Card>
```

---

### 2.2 Date Range Support (From/To Dates)

**Affected files:**
- Backend: Migration to add `end_date` field
- Backend: Model and service updates
- Frontend: Update forms with date range picker
- Frontend: Update calendar view for multi-day events

**Database Migration:**
```go
// Add to events table
Schema.Table("events").AddColumn("end_date", "date").Nullable()
```

**Implementation:**
- Add `end_date` field to events table (nullable for single-day events)
- Update create/edit forms:
  - "Start Date" field (required)
  - "End Date" field (optional, defaults to start date if empty)
- Update calendar view to display multi-day events spanning the date range
- Validate that end_date >= start_date

**Example Use Case:** Trade fair spanning 3 days (Nov 27-29) should show on all 3 calendar days.

---

### 2.3 Partner Dropdown from Config

**Affected files:** Frontend forms

**Implementation:**
- Same pattern as procurement
- Fetch from `/api/configurations?type=partners`
- Use existing or create new "partners" config type

---

### 2.4 Remove "Attending SMEs" from Forms

**Affected files:** Frontend create/edit forms

**Implementation:**
- Remove from create/edit forms completely
- Keep as read-only list in view page
- SMEs will register attendance through their own interface ("I'm attending" button)

---

## 3. BDSP Management

### 3.1 Fix "Not Found" Issue

**Priority:** HIGH - Blocking

**Investigation Steps:**
1. Check browser console for errors
2. Verify API endpoint routing
3. Check if BDSP routes are properly registered
4. Test API endpoint directly: `GET /api/bdsps`
5. Check for any merge conflicts that may have broken routes
6. Verify database seeds/migrations ran correctly

**Likely Causes:**
- Route registration issue after merge
- Missing API endpoint
- Database query error

---

### 3.2 Convert to Config Dropdowns

**Affected files:** Frontend forms

| Field | Data Source |
|-------|-------------|
| `partners` | Config: "partners" |
| `product_types` | Config: create new "product_types" config type |

---

## 4. New Configuration Types Required

Add the following configuration categories to the configurations system:

| Config Type | Description | Example Values |
|-------------|-------------|----------------|
| `procured_by` | Entities that can procure | "Ministry of Trade", "World Bank", "UNDP" |
| `organisations` | Organisations list | "Tiyeni", "MoIT", "Books" |
| `procurement_types` | Types of procurement | "Goods", "Services", "Works", "Consulting" |
| `classifications` | Procurement classifications | "Open", "Restricted", "Direct" |
| `localisations` | Localisation options | "Local", "International", "Regional" |
| `locations` | Location options | "Lilongwe", "Blantyre", "Mzuzu", "National" |
| `product_types` | BDSP product types | "Training", "Consulting", "Mentorship" |
| `partners` | Partner organisations | (if not already existing) |

**Implementation:**
- Add these types to the configurations seeder/admin panel
- Ensure admin can manage these values through the configurations UI

---

## 5. Priority Order

### P0 - Critical (Fix First)
- [ ] Fix BDSP "not found" issue
- [ ] Fix filters not working for procurements and events

### P1 - High Priority
- [ ] Convert text inputs to config-driven dropdowns (all entities)
- [ ] Remove interested SMEs / attending SMEs from edit forms
- [ ] Auto-generate reference numbers for procurements

### P2 - Medium Priority
- [ ] Split procurement form into two steps
- [ ] Add rich text editor for procurement details field
- [ ] Calendar styling fixes (padding, background)
- [ ] Simplify view UI (remove unnecessary icons)
- [ ] Add end_date to events table
- [ ] Update event forms and calendar view for date ranges

### P3 - Lower Priority
- [ ] Replace "Is Published" toggle with publish button workflow

---

## 6. Reference Files

Look at these existing implementations for patterns:

| Pattern | Reference Location |
|---------|-------------------|
| Multi-step forms | SME create/edit forms |
| Quick filters | Books implementation |
| Config dropdowns | Existing config usage in codebase |
| Data tables with filters | Books list page |
| Simple detail views | SME detail page |

---

## 7. Testing Checklist

### Procurement
- [ ] All dropdown fields populate from configurations
- [ ] Reference numbers auto-generate correctly (format: `PROC-YYYY-NNNNN`)
- [ ] Multi-step form navigates correctly between steps
- [ ] Rich text editor saves and displays formatted content
- [ ] Quick filters work on procurement list
- [ ] Interested SMEs only visible in view (not create/edit)
- [ ] Draft status is default for new procurements

### Events
- [ ] Calendar displays with proper styling (white card, correct padding)
- [ ] Events can be created with date range (start + end date)
- [ ] Multi-day events display correctly on calendar
- [ ] Partner dropdown populates from config
- [ ] Attending SMEs only visible in view (not create/edit)
- [ ] Quick filters work on events list

### BDSP
- [ ] BDSP pages load without "not found" error
- [ ] Partners dropdown populates from config
- [ ] Product types dropdown populates from config

### General
- [ ] All new config types are available in configurations admin
- [ ] No console errors on any page
- [ ] Forms save correctly with dropdown values
- [ ] View pages display all data correctly

---

## 8. Notes

- **DO NOT** touch configurations UI for now - park that work
- Focus on the form improvements and fixes listed above
- Export functionality will be handled separately
- SME view/portal and Gawani tasks are next phase (after this deployment)

---

## 9. Deployment Checklist

Before merging:
- [ ] All P0 and P1 items completed
- [ ] Run database migrations (especially events.end_date)
- [ ] Seed new configuration types
- [ ] Test all forms end-to-end
- [ ] Verify no regression in existing functionality
