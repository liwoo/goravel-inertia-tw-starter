---
name: rebrand
description: Rebrand the application with a new name, updating i18n translations, env files, favicon, and all frontend/backend references
argument-hint: <new-brand-name>
disable-model-invocation: true
allowed-tools: Read, Edit, Write, Bash, Grep, Glob
---

# Rebrand Application

Rebrand this application to **$ARGUMENTS**.

Follow every step below. Do NOT skip any. After all steps, run the verification build.

## Step 1: Derive branding values

From the brand name `$ARGUMENTS`, compute:
- **fullName**: The full brand name exactly as provided (e.g., "Toyota Portal")
- **initials**: First letter of each word, max 2 characters, uppercase (e.g., "TP")
- **copyright**: `{{year}} $ARGUMENTS. All rights reserved.`

If the user provided just a company name (e.g., "Toyota"), the portal name is "$ARGUMENTS Portal" unless they specified otherwise.

## Step 2: Update i18n translation files

Update these JSON files under `resources/js/locales/en/`:

### auth.json
- `branding.portalName` → the full portal name
- `branding.copyright` → the copyright string with `{{year}}` interpolation
- `login.title` → the full portal name

### nav.json
- `sidebar.portalName` → the full portal name

### settings.json
- `backupCodes.downloadHeader` → `"<fullName> - Two-Factor Authentication Backup Codes"`

## Step 3: Update .env files

Update `APP_NAME` in `.env`, `.env.example`, and `.env.testing` (if they exist):
```
APP_NAME=<fullName>
```

Also update `VITE_APP_NAME` if it exists in any .env file.

## Step 4: Update chart-actions default

In `resources/js/components/ui/chart-actions.tsx`, update the default value for `appName`:
- Change `appName = "Admin Portal"` → `appName = "<fullName>"`
- Update the JSDoc comment above it too

## Step 5: Update favicon

Generate a new SVG favicon at `public/images/favicon.svg` with the computed initials:

```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">
  <rect width="32" height="32" rx="6" fill="#18181b"/>
  <text x="16" y="22" font-family="system-ui, -apple-system, sans-serif" font-size="16" font-weight="700" fill="#ffffff" text-anchor="middle"><INITIALS></text>
</svg>
```

Verify `resources/views/app.tmpl` points to `/images/favicon.svg`.

## Step 6: Verify no stale references

Run a grep for the OLD brand name (currently "Admin Portal") across the project (excluding `node_modules/`). If any references remain in source files, update them.

## Step 7: Build verification

Run:
1. `npx tsc --noEmit` - must pass with zero errors
2. `npm run build` - must produce a successful Vite build

Report the results. If either fails, fix the errors before finishing.

## Summary

After completion, output a table showing:
| File | What changed |
|------|-------------|
| ... | ... |

And confirm the build status.
