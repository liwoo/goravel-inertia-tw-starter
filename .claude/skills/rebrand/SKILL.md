---
name: rebrand
description: Rebrand the application with a new name, updating Go module path, Docker/Helm/CI names, DB names, i18n translations, env files, favicon, and all frontend/backend references
argument-hint: <new-brand-name>
disable-model-invocation: true
allowed-tools: Read, Edit, Write, Bash, Grep, Glob
---

# Rebrand Application

Rebrand this application to **$ARGUMENTS**.

Follow every step below. Do NOT skip any. After all steps, run the verification build.

---

## Step 1: Derive branding values

From the brand name `$ARGUMENTS`, compute:
- **fullName**: The full brand name exactly as provided (e.g., "Toyota Portal")
- **shortName**: Single word brand (e.g., "Toyota")
- **kebabName**: Lowercase hyphenated (e.g., "toyota-portal")
- **snakeName**: Lowercase underscored (e.g., "toyota_portal")
- **initials**: First letter of each word, max 2 characters, uppercase (e.g., "TP")
- **copyright**: `{{year}} $ARGUMENTS. All rights reserved.`
- **moduleName**: `<kebabName>` — the Go module path (e.g., "toyota-portal")
- **dbName**: `<snakeName>_db` — the database name (e.g., "toyota_portal_db")
- **dbUser**: `<snakeName>` — the database username
- **imageName**: `<kebabName>` — the Docker image name

If the user provided just a company name (e.g., "Toyota"), the portal name is "$ARGUMENTS Portal" unless they specified otherwise.

---

## Step 2: Discover current brand values

Before replacing, identify the current brand strings by reading:

1. `go.mod` line 1 — current Go module path (e.g., `books-database`)
2. `.env` or `.env.example` — current `APP_NAME`, `DB_DATABASE`, `DB_USERNAME`
3. `docker-compose/docker-compose.yml` — current `name:` and `container_name:` prefix
4. `Dockerfile` — current `org.opencontainers.image.vendor`
5. `helm/goravel-blog/Chart.yaml` — current `description`, `keywords`, `home`

Store these as the **old** values for the replacement mapping.

---

## Step 3: Replace Go module path

This is the most critical step — every Go import depends on the module path.

1. Read current module from `go.mod` line 1 (e.g., `module books-database`)
2. Replace **old module** → **new moduleName** across ALL `.go` files, `.golangci.yml`, and `.gitignore`:

```bash
# Replace in all Go files and config
grep -rl '<OLD_MODULE>' --include='*.go' --include='.golangci.yml' --include='.gitignore' \
  --exclude-dir=node_modules --exclude-dir=vendor --exclude-dir=.git . \
  | xargs sed -i '' 's/<OLD_MODULE>/<NEW_MODULE>/g'

# Fix go.mod separately (no extension match)
sed -i '' 's/<OLD_MODULE>/<NEW_MODULE>/g' go.mod
```

3. Run `go build ./...` to verify all imports resolve.

---

## Step 4: Replace Docker & infrastructure names

Replace in order (most specific first to avoid partial matches). Target files: `*.yml`, `*.yaml`, `Makefile`, `Dockerfile`, `*.sh`, `*.md`, `.env*`.

### 4a. Docker Compose names

| Pattern | Example old | Example new |
|---------|-------------|-------------|
| Project name | `books-database` | `<kebabName>` |
| Container prefix | `books-app`, `books-postgres`, `books-redis`, etc. | `<kebabName>-app`, etc. |
| Network names | `books-network`, `books-network-dev`, `books-network-test` | `<kebabName>-network`, etc. |
| Test containers | `books-test-runner`, `books-frontend-test`, etc. | `<kebabName>-test-runner`, etc. |

### 4b. Database names

| Pattern | Example old | Example new |
|---------|-------------|-------------|
| DB name | `books_db` | `<dbName>` |
| DB user | `books` (as username) | `<dbUser>` |
| DB password | `books_password` | `<snakeName>_password` |
| Test DB | `books_test` | `<snakeName>_test` |
| Test password | `books_test_password` | `<snakeName>_test_password` |
| Production DB | `books_db_production` | `<dbName>_production` |
| Production user | `books_production` | `<snakeName>_production` |
| Cache prefix | `books_cache` | `<snakeName>_cache` |

### 4c. Display names

| Pattern | Example old | Example new |
|---------|-------------|-------------|
| `"Books Database"` | Full display name | `"<fullName>"` |
| `"Books"` (standalone brand) | Short brand | `"<shortName>"` |

### 4d. Helm values

| Pattern | Files | Example |
|---------|-------|---------|
| Ingress hosts | `values.yaml`, `values.staging.yaml`, `values.production.yaml` | `books.example.com` → `<kebabName>.example.com` |
| TLS secret names | same | `books-tls-secret` → `<kebabName>-tls-secret` |
| K8s namespaces | skill files, values | `books-staging` → `<kebabName>-staging` |
| Chart description | `Chart.yaml` | Update description and keywords |
| Chart home/sources | `Chart.yaml` | Update GitHub URLs |
| Image repository | `values.yaml` | `books-database` → `<kebabName>` |

### 4e. CI/CD

| File | What to update |
|------|---------------|
| `.github/workflows/ci.yml` | `IMAGE_NAME` fallback |
| `.github/workflows/cd.yml` | `IMAGE_NAME` fallback |
| `Makefile` | `DOCKER_IMAGE_NAME` |

### Recommended approach

Use targeted `sed` in order from longest/most-specific patterns to shortest to avoid partial match corruption:

```bash
# 1. Compound names first (e.g., books-test-runner before books-test)
# 2. Underscore names (books_test_password before books_test before books)
# 3. Display names (Books Database before Books)
# 4. Single word last (books, Books)
```

**Important**: Run replacements on ALL these file types:
```
--include='*.yml' --include='*.yaml' --include='*.md' --include='*.sh'
--include='*.json' --include='Makefile' --include='Dockerfile'
```

Also handle dotfiles separately since `--include` doesn't always match them:
- `.env`, `.env.example`, `.env.testing`
- `.gitignore`, `.golangci.yml`

---

## Step 5: Replace in skill files

The `.claude/skills/` directory contains hardcoded module paths and brand names in `SKILL.md` files. These MUST also be updated:

```bash
grep -rl '<OLD_MODULE>' .claude/skills/ | xargs sed -i '' 's/<OLD_MODULE>/<NEW_MODULE>/g'
grep -rl '<OLD_BRAND>' .claude/skills/ | xargs sed -i '' 's/<OLD_BRAND>/<NEW_BRAND>/g'
```

---

## Step 6: Update i18n translation files

Update these JSON files under `resources/js/locales/en/`:

### auth.json
- `branding.portalName` → the full portal name
- `branding.copyright` → the copyright string with `{{year}}` interpolation
- `login.title` → the full portal name

### nav.json
- `sidebar.portalName` → the full portal name

### settings.json
- `backupCodes.downloadHeader` → `"<fullName> - Two-Factor Authentication Backup Codes"`

---

## Step 7: Update .env files

Update `APP_NAME` in `.env`, `.env.example`, and `.env.testing` (if they exist):
```
APP_NAME=<fullName>
```

Also update `VITE_APP_NAME` if it exists in any .env file.

**Note**: DB_DATABASE, DB_USERNAME, DB_PASSWORD should already be updated from Step 4b. Verify they are correct.

---

## Step 8: Update chart-actions default

In `resources/js/components/ui/chart-actions.tsx`, update the default value for `appName`:
- Change `appName = "..."` → `appName = "<fullName>"`
- Update the JSDoc comment above it too

---

## Step 9: Update favicon

Generate a new SVG favicon at `public/images/favicon.svg` with the computed initials:

```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">
  <rect width="32" height="32" rx="6" fill="#18181b"/>
  <text x="16" y="22" font-family="system-ui, -apple-system, sans-serif" font-size="16" font-weight="700" fill="#ffffff" text-anchor="middle"><INITIALS></text>
</svg>
```

Verify `resources/views/app.tmpl` points to `/images/favicon.svg`.

---

## Step 10: Update MEMORY.md

Update the module path in `.claude/projects/-Users-liwu-Projects-go-blog/memory/MEMORY.md`:
- `Module path:` line should reflect the new module name

---

## Step 11: Verify no stale references

Run a case-insensitive grep for ALL old brand strings across the project:

```bash
# Check for old module path
grep -rli '<OLD_MODULE>' --exclude-dir=node_modules --exclude-dir=vendor \
  --exclude-dir=.git --exclude-dir=.playwright-mcp --exclude-dir=storage \
  --exclude-dir=tmp --exclude='go.sum' . 2>/dev/null

# Check for old brand name (case-insensitive)
grep -rli '<OLD_BRAND>' --exclude-dir=node_modules --exclude-dir=vendor \
  --exclude-dir=.git --exclude-dir=.playwright-mcp --exclude-dir=storage \
  --exclude-dir=tmp --exclude='go.sum' . 2>/dev/null
```

If any references remain in source files, fix them. Ignore:
- `go.sum` (external dependency checksums)
- `storage/logs/` (runtime logs)
- `tmp/` (compiled binaries)

---

## Step 12: Build verification

Run:
1. `go build ./...` — must pass (verifies Go module rename)
2. `npx tsc --noEmit` — must pass (verifies TypeScript)
3. `npm run build` — must produce a successful Vite build

Report the results. If any fail, fix the errors before finishing.

---

## Summary

After completion, output a table showing:

| Category | Old | New | Files affected |
|----------|-----|-----|---------------|
| Go module | ... | ... | ~N files |
| Docker project | ... | ... | N files |
| DB names | ... | ... | N files |
| Display name | ... | ... | N files |
| Helm/K8s | ... | ... | N files |
| i18n | ... | ... | N files |
| CI/CD | ... | ... | N files |

And confirm all build checks pass.
