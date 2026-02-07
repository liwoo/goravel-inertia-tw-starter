---
name: goravel-crud-nav
description: Add i18n-aware navigation entry for a Goravel entity. Updates sidebar menu with icon, permission gating, and nav i18n keys.
argument-hint: "[EntityName] [icon-name]"
allowed-tools: Read, Write, Edit, Grep, Glob
---

# Goravel CRUD Navigation (i18n-aware)

Add navigation for `$ARGUMENTS`.

> **Note**: For global search (CMD+K) integration, use `/goravel-crud-search` after this skill.

## File 1: Navigation — `resources/js/config/navigation.ts`

### Step 1: Add Icon Import

```typescript
import {
    // ... existing imports
    YourIcon,  // from lucide-react
} from "lucide-react"
```

Browse icons at https://lucide.dev/icons

### Step 2: Add Navigation Entry

Navigation items use **i18n keys** (not hardcoded strings). Keys are resolved at render time via `useTranslation('nav')`.

**For main navigation** (`navMain` array):
```typescript
{
    title: "main.entityNames",        // i18n key in nav namespace
    url: "/admin/entity-names",
    icon: YourIcon,
    requiredService: "entities",       // Must match ServiceRegistry value
    requiredAction: "read" as const,
},
```

**For document/secondary section** (`documents` array):
```typescript
{
    name: "documents.entityNames",     // i18n key in nav namespace
    url: "/admin/entity-names",
    icon: YourIcon,
    requiredService: "entities",
    requiredAction: "read" as const,
},
```

**For admin-only section** (`navSecondary` array):
```typescript
{
    title: "secondary.entityNames",    // i18n key in nav namespace
    url: "/admin/entity-names",
    icon: YourIcon,
    requireSuperAdmin: true,
},
```

### Step 3: Add i18n Key to Nav Namespace

Edit `resources/js/locales/en/nav.json` — add the key matching what you used above:

```json
{
  "main": {
    "entityNames": "Entity Names"
  }
}
```

Or for documents section:
```json
{
  "documents": {
    "entityNames": "Entity Names"
  }
}
```

## Navigation Permission Fields

| Field | Purpose |
|---|---|
| `title` / `name` | **i18n key** in `nav` namespace (resolved at render time) |
| `requiredService` | Service name from `permission_constants.go` (e.g., "books") |
| `requiredAction` | Permission action: "read", "create", "update", "delete", "manage" |
| `requiredRole` | Specific role slug required |
| `requireSuperAdmin` | Only visible to super admins |

## Choosing the Right Section

| Section | Array | Use for |
|---|---|---|
| Main nav | `navMain` | Primary entity pages (Books, Users, Applications) |
| Documents | `documents` | Document-type entities (Reports, Files) |
| Secondary | `navSecondary` | Admin-only pages (Settings, Configs, Audit Logs) |

## Verification

1. Check sidebar shows the new navigation item (with correct permission gating)
2. Verify the i18n key resolves correctly (check `nav.json`)
3. Log in as a user without the required permission — item should be hidden

## Verify

After adding the navigation entry and i18n key:

```bash
# TypeScript compiles (catches wrong icon imports, missing fields)
npx tsc --noEmit

# Lint the navigation config
npx eslint resources/js/config/navigation.ts --max-warnings=0
```

## Reference

See `resources/js/config/navigation.ts` for existing entries.
See `resources/js/locales/en/nav.json` for existing i18n keys.
