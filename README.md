# Books Database

[![CI](https://github.com/Tiyeni/books-database/actions/workflows/ci.yml/badge.svg)](https://github.com/Tiyeni/books-database/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-ready-blue.svg)](https://hub.docker.com/)

A Goravel-based admin portal with React/Inertia.js frontend, featuring JWT auth, RBAC permissions, CRUD scaffolding, i18n, dark mode, and client-side data exports.

## Prerequisites

- **Go 1.24+** — [golang.org/dl](https://golang.org/dl/)
- **Node.js 20+** — [nodejs.org](https://nodejs.org/)
- **PostgreSQL 16+** — via Docker or local install
- **Docker** — for testcontainers and dev environment

Optional:
- **Air** for Go hot reload: `go install github.com/air-verse/air@latest`
- **Claude Code** for agentic development (see [Agentic Mode](#agentic-mode))

## Quick Start

```bash
# Clone and install
git clone <repository-url> && cd blog
go mod download
npm install

# Setup environment
cp .env.example .env
# Edit .env with your database credentials (see Environment Variables below)

# Initialize database
go run . artisan key:generate
go run . artisan migrate
go run . artisan seed          # Seeds RBAC permissions + sample data

# Create admin user
go run . artisan user:create-admin

# Run (two terminals)
air                             # Terminal 1: Backend (or: go run .)
npm run dev                     # Terminal 2: Frontend

# Open http://localhost:3500
```

## Environment Variables

### Application (.env)

```env
APP_NAME="Books Database"
APP_ENV=local
APP_DEBUG=true
APP_URL=http://localhost
APP_HOST=127.0.0.1
APP_PORT=3000
APP_KEY=                        # Generate with: go run . artisan key:generate

JWT_SECRET=                     # Any random string for local dev

DB_CONNECTION=postgres
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=books_db
DB_USERNAME=postgres
DB_PASSWORD=yourpassword

SESSION_DRIVER=file
SESSION_LIFETIME=120

AUTH_REQUIRE_2FA=true           # Set false for simpler local dev
```

### Test Credentials

For E2E browser testing (Playwright MCP) and local development, set these variables so test scripts know which user to authenticate as:

```env
TEST_USER=admin@example.com     # Email of an admin user created via user:create-admin
TEST_PASSWORD=your-test-pass    # That user's password
```

These are used by:
- `/e2e-entity-suite` — logs into the app via the browser to test CRUD flows
- `/playwright-ui-test` — verifies page rendering, forms, and navigation
- Local smoke testing after seeding

Create the test user once:
```bash
go run . artisan user:create-admin
# Use the email/password you set in TEST_USER/TEST_PASSWORD
```

---

## Manual Mode

Step-by-step commands for building a new CRUD entity without Claude Code.

### 1. Database & Model

```bash
# Create table migration
go run . artisan make:migration create_lenders_table
# Edit the generated file in database/migrations/, then:
go run . artisan migrate

# Add audit fields (created_by, updated_by, deleted_by, ip_address, user_agent)
go run . artisan make:audit --table=lenders
# Remove duplicate deleted_at from generated file if present, then:
go run . artisan migrate

# Register both migrations in database/kernel.go

# Generate model from table
go run . artisan make:model-from-table --table=lenders --model=Lender
# Fix: array fields, carbon.DateTime types, SearchFields, TableName
```

### 2. Service & Permissions

```bash
# Generate service
go run . artisan make:svc --svc=Lender
# Configure builder: search fields, sort, filter, validation, scope

# Register in app/auth/permission_constants.go:
#   ServiceLenders ServiceRegistry = "lenders"

# Sync to database
go run . artisan permissions:setup
```

### 3. Requests & Controller

```bash
# Generate request validators
go run . artisan make:req --model=Lender --resource=Lender
# Fix: validation rules, snake_case form tags, ToCreateData/ToUpdateData

# Generate API controller
go run . artisan make:api-ctrl --controller=Lender
# Fix: service constant must match permission_constants.go

# Register in routes/api.go (search/filter routes BEFORE {id} route)
```

### 4. Testing

```bash
# Generate CRUD tests
go run . artisan make:crud-test --svc=lender

# Run tests (uses PostgreSQL testcontainer)
./scripts/run_tests.sh -v ./tests/feature/crud -run TestLenderCrudSuite
```

**STOP**: All tests must pass before proceeding to UI work.

### 5. Frontend

```bash
# Generate page controller
go run . artisan make:page-ctrl --controller=Lender

# Generate complete UI hierarchy
go run . artisan make:ui --page=Lender --request=Lender
# Creates: pages/Lender/Index.tsx, sections/*.tsx, types/lender.ts

# Register web route in routes/web.go:
#   router.Get("/admin/lenders", lendersPageController.Index)

# Add to resources/js/config/navigation.ts
# Add to search_controller.go + search_config.tsx for CMD+K search
```

### 6. Verify

```bash
go build ./...                  # Go compiles
npx tsc --noEmit                # TypeScript compiles
go run . artisan migrate        # Dev DB migrated
# Visit http://localhost:3500/admin/lenders
```

### Quick Reference (Manual)

| # | Command / Action | Purpose |
|---|-----------------|---------|
| 1 | `make:migration` | Create table schema |
| 2 | `migrate` | Apply migration |
| 3 | `make:audit` | Add audit fields |
| 4 | `migrate` | Apply audit migration |
| 5 | Register in `kernel.go` | Wire migrations |
| 6 | `make:model-from-table` | Generate model |
| 7 | `make:svc` | Generate service |
| 8 | Edit `permission_constants.go` | Register service |
| 9 | `permissions:setup` | Sync permissions |
| 10 | `make:req` | Generate request validators |
| 11 | `make:api-ctrl` | Generate API controller |
| 12 | Register in `routes/api.go` | Wire API endpoints |
| 13 | `make:crud-test` | Generate tests |
| 14 | `run_tests.sh` | Run & fix tests |
| 15 | `make:page-ctrl` | Generate page controller |
| 16 | `make:ui` | Generate UI files |
| 17 | Register in `routes/web.go` | Wire page route |
| 18 | Add to `navigation.ts` | Sidebar entry |
| 19 | Update `search_controller.go` | Global search (optional) |

---

## Agentic Mode

Use [Claude Code](https://claude.ai/claude-code) with the project's skill system to scaffold entities end-to-end. Skills are step-by-step recipes that Claude executes using the project's artisan generators, then fixes the output automatically.

### Agent Personas

| Agent | Persona | Focus |
|-------|---------|-------|
| `goravel-crud-engineer` | Chikondi Banda | Backend: migrations, models, services, controllers, routes, tests |
| `goravel-inertia-ui-engineer` | Thoko Nkhoma | Frontend: types, pages, forms, columns, detail views, exports |
| `goravel-qa-engineer` | Blessings Phiri | QA: CRUD review, type checking, test suites, E2E browser tests |
| `goravel-devops-engineer` | Kondwani Mwale | DevOps: Docker, Helm, CI/CD, deploy, infrastructure debug |

### Backend Skills (Chikondi)

| Skill | Purpose |
|-------|---------|
| `/goravel-crud-migration` | Create table + audit migrations |
| `/goravel-crud-model` | Generate model with post-gen fixes |
| `/goravel-crud-service` | Service layer with builder pattern |
| `/goravel-crud-permissions` | Register in permission system |
| `/goravel-crud-request` | Create/update request validators |
| `/goravel-crud-controller` | API controller with permissions |
| `/goravel-crud-routes` | Register API + web routes |
| `/goravel-crud-test` | Generate & fix CRUD tests |
| `/goravel-crud-page` | Page controller + UI generation |
| `/goravel-crud-nav` | Navigation entry with i18n |
| `/goravel-crud-search` | Global search (CMD+K) |
| `/goravel-enum` | Go enum + TypeScript generation |
| `/goravel-scaffold` | Full backend orchestrator (19 steps) |
| `/fake-data` | Database seeder with 25+ records |

### Frontend Skills (Thoko)

| Skill | Purpose |
|-------|---------|
| `/inertia-types` | TypeScript interfaces + i18n namespace |
| `/inertia-page` | Index.tsx with CrudPage wrapper |
| `/inertia-columns` | Table columns, mobile columns, filters |
| `/inertia-form` | Create/edit forms with validation |
| `/inertia-detail` | Read-only detail view |
| `/inertia-page-config` | Stats, filters, page actions, bulk actions |
| `/inertia-page-ctrl` | Go page controller for Inertia |
| `/inertia-custom-page` | Non-CRUD pages (dashboards, reports) |
| `/inertia-form-review` | Audit forms for i18n, types, consistency |
| `/inertia-scaffold` | Full frontend orchestrator (17 steps) |
| `/multi-step-form` | Wizard-style multi-step forms |
| `/ui-ux-audit` | Table density, icons, dropdowns, colors |
| `/file-downloads` | CSV/Excel/PDF/JSON export |

### QA Skills (Blessings)

| Skill | Purpose |
|-------|---------|
| `/goravel-crud-review` | Full CRUD implementation audit |
| `/goravel-type-check` | Go struct <> TypeScript interface consistency |
| `/goravel-test-suite` | Write comprehensive test suites |
| `/playwright-ui-test` | Browser UI verification |
| `/e2e-entity-suite` | 14-phase, 30-test E2E suite |

### DevOps Skills (Kondwani)

| Skill | Purpose |
|-------|---------|
| `/deploy` | Deploy to staging/production |
| `/docker-dev` | Local Docker development |
| `/infra-debug` | Debug pods, containers, logs |
| `/helm-values` | Update Helm chart values |
| `/ci-check` | Run CI pipeline locally |

### Cross-Cutting Skills

| Skill | Purpose |
|-------|---------|
| `/rebrand` | Rename the app across all files |
| `/broadcast-notification` | Entity-specific notification system |
| `/custom-pages` | Non-CRUD pages (portals, modals, widgets) |
| `/dashboard-visualization` | Charts, KPIs, data dashboards |

### Full Entity Scaffold (Agentic)

To scaffold a complete entity end-to-end:

**Backend first** — run `/goravel-scaffold EntityName table_name`:
```
Phase 1: Database & Model (migration, audit, model)
Phase 2: Service & Permissions
Phase 3: Requests, Controller & Routes
Phase 4: CRUD Tests (must all pass before UI)
Phase 5: Page controller, UI files, navigation, search
```

**Frontend next** — run `/inertia-scaffold EntityName`:
```
Steps 1-9:   Types, page controller, UI files, columns, forms, detail, config
Steps 10-13: Index page, navigation, search, form review
Steps 14-15: File downloads, UX audit
Steps 16-17: Seed data, E2E browser testing
```

### E2E Testing with Playwright MCP

After scaffolding, run `/e2e-entity-suite EntityName` to execute browser tests:

```
Phase 1:  Login & Navigation
Phase 2:  Page Structure (stats, filters, columns)
Phase 3:  Create Entity (form validation, submission)
Phase 4:  Detail View
Phase 5:  Edit Entity
Phase 6:  Table Search
Phase 7:  Filter Tabs
Phase 8:  Sorting (requires 25+ seeded records)
Phase 9:  Pagination
Phase 10: Global Search (CMD+K)
Phase 11: Row Actions (edit, delete from menu)
Phase 12: Cross-Entity Integration (conditional)
Phase 13: Responsive Layout
Phase 14: Console Errors & Network Failures
```

Requires `TEST_USER` and `TEST_PASSWORD` env vars and seeded data (`/fake-data`).

---

## Docker Development

### Local with Docker Compose

```bash
cd docker-compose

# Standard (built image)
docker compose up -d

# With hot reload (source mounted)
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# With admin tools (pgAdmin on :5050, Redis Commander on :8081)
docker compose --profile tools up -d

# Tear down
docker compose down           # Keep data
docker compose down -v        # Remove volumes
```

### Services

| Service | Port | Purpose |
|---------|------|---------|
| App | 3000 | Go backend + frontend |
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache + sessions |
| pgAdmin | 5050 | DB management (tools profile) |
| Redis Commander | 8081 | Cache browser (tools profile) |
| MinIO | 9000/9001 | S3 storage (storage profile) |
| Mailhog | 8025 | Email testing (mail profile) |

### Docker Build

3-stage build: Go builder, Node builder, Alpine runtime.

```bash
docker build -t books-database .
# Binary compressed with UPX, runs as non-root user
```

---

## Testing

### Automated Tests (Go)

Tests use PostgreSQL 16 via testcontainers — no local DB setup needed.

```bash
# Run all tests
./scripts/run_tests.sh -v ./tests/...

# Run specific suite
./scripts/run_tests.sh -v ./tests/feature/crud -run TestBookCrudSuite

# Run with coverage
./scripts/run_tests.sh ./tests/... -v -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -15
```

The test script (`scripts/run_tests.sh`):
- Starts a PostgreSQL 16 container on a random port
- Sets `APP_ENV=testing`, `AUTH_REQUIRE_2FA=false`
- Creates isolated storage directories
- Runs tests with `-p=1` (sequential packages, prevents migration races)
- Auto-cleans up container on exit

### Frontend Tests

```bash
npm run test                    # Vitest in watch mode
npm run test -- --run           # Single run
npm run test:coverage -- --run  # With coverage
```

### Test Helpers

| Helper | Purpose |
|--------|---------|
| `helpers.SetupJWTUser(email, password, role)` | Create test user with role |
| `helpers.AssignPermissionToRole(role, perms)` | Grant permissions |
| `helpers.CleanTestDatabase()` | Reset between tests |
| `s.loginUser(email, password)` | Get auth cookie |

### Writing Tests

```go
type LenderCrudSuite struct {
    suite.Suite
    tests.TestCase
    server     *httptest.Server
    client     *http.Client
    authCookie *http.Cookie
    testUser   *models.User
}

func (s *LenderCrudSuite) SetupSuite() {
    s.InitApp()
    s.server = httptest.NewServer(facades.Route())
    jar, _ := cookiejar.New(nil)
    s.client = &http.Client{Jar: jar}
    role, _ := helpers.AssignPermissionToRole("admin", allPermissions)
    user, _ := helpers.SetupJWTUser("test@example.com", "password", role)
    s.testUser = user
    s.authCookie = s.loginUser("test@example.com", "password")
}
```

---

## CI/CD

### CI Pipeline (GitHub Actions)

8 parallel jobs on push to `main`/`develop`:

| Job | What it does |
|-----|-------------|
| `go-lint` | golangci-lint + go vet |
| `go-test` | Full test suite with testcontainers + coverage |
| `go-build` | CGO_ENABLED=0 static binary |
| `frontend-lint` | ESLint + TypeScript type check |
| `frontend-test` | Vitest with coverage |
| `frontend-build` | Vite production build |
| `dependency-scan` | govulncheck + npm audit |
| `docker-build` | Build, Trivy scan, push (if DEPLOYABLE=true) |
| `helm-lint` | Helm lint + template validation |

### CD Pipeline

Sequential deploy with safety gates:
1. Build, scan, push Docker image
2. `helm upgrade --atomic` (auto-rollback on failure)
3. Smoke test (`/health` endpoint)
4. Auto-rollback job if smoke test fails

### Helm Chart

```bash
# Lint
helm lint helm/goravel-blog

# Template with staging values
helm template goravel-blog helm/goravel-blog \
  --values helm/goravel-blog/values.yaml \
  --values helm/goravel-blog/values.staging.yaml

# Deploy
helm upgrade --install goravel-blog helm/goravel-blog \
  --values helm/goravel-blog/values.staging.yaml \
  --set secrets.jwtSecret=$JWT_SECRET \
  --atomic --timeout 5m
```

3 value files: `values.yaml` (defaults), `values.staging.yaml`, `values.production.yaml`.

---

## Project Structure

```
app/
  auth/                         # Permission constants, helpers, scoped access
  console/commands/             # Artisan commands
  contracts/                    # Interfaces (ServiceBuilder, etc.)
  http/
    controllers/                # API + page controllers (per entity)
    middleware/                  # JWT, CORS, permissions
    requests/                   # Request validation structs
  models/                       # GORM models
  providers/                    # Service providers
  services/                     # Business logic (builder pattern)
database/
  migrations/                   # Schema migrations
  seeders/                      # Data seeders
resources/js/
  components/                   # Reusable UI (CrudPage, ExportDialog, etc.)
  config/                       # navigation.ts, search_config.tsx
  contexts/                     # React contexts (permissions, theme)
  hooks/                        # Custom hooks (useExport, etc.)
  locales/en/                   # i18n translation JSON files
  pages/                        # Entity pages (Index.tsx + sections/)
  types/                        # TypeScript interfaces
  utils/                        # Utilities (exportUtils, etc.)
routes/
  api.go                        # API routes
  web.go                        # Web/Inertia routes
helm/goravel-blog/              # Helm chart
docker-compose/                 # Docker Compose files
scripts/                        # Shell scripts (tests, linting, hooks)
.claude/skills/                 # Claude Code skill definitions
.github/workflows/              # CI/CD pipelines
```

## Key Conventions

- **Endpoints**: hyphenated (`/entity-names` not `/entity_names`)
- **Request struct tags**: snake_case (`form:"first_name" json:"first_name"`)
- **Create data**: camelCase keys (matches model json tags)
- **Update data**: snake_case keys (matches DB columns via GORM)
- **i18n**: `useTranslation('namespace')` in components, `t: TFunction` param in configs
- **Permissions**: `auth.ServiceEntity` + `auth.PermissionAction` pattern
- **Services**: `contracts.NewServiceBuilder[T]` builder pattern

## Troubleshooting

### Test Issues

```bash
# Migration race conditions
# Use -p=1 flag (run_tests.sh does this automatically)
./scripts/run_tests.sh -v ./tests/...

# JWT returns 302 redirect (not 401)
# JWT middleware redirects unauthenticated requests — use cookie jar in tests
# For unauthenticated tests, use a FRESH http.Client (no jar)

# JSON numbers decode as float64
# Assert with float64 type: s.Equal(float64(100), result["price"])

# Timestamps have second precision only
# TimestampsTz() creates timestamp(0) — don't rely on millisecond ordering
```

### Database Issues

```bash
# Fresh start
go run . artisan migrate:fresh
go run . artisan seed

# Foreign key order errors
# Check migration order in database/kernel.go

# Dev DB not updated after test changes
# Tests use testcontainers — run migrate on dev DB separately:
go run . artisan migrate
```

### Frontend Issues

```bash
# TypeScript errors
npx tsc --noEmit

# Clear build cache
rm -rf node_modules/.vite
npm run dev

# Full reinstall
rm -rf node_modules package-lock.json && npm install
```

### Permission Issues

```bash
# Re-sync permissions
go run . artisan permissions:setup

# Check user role
go run . artisan user:show email@example.com

# Debug: set LOG_LEVEL=debug, look for "CheckScopedPermission" in logs
```

## License

MIT
