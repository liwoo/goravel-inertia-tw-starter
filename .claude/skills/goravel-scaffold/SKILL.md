---
name: goravel-scaffold
description: Full end-to-end CRUD scaffolding for a new Goravel entity. Orchestrates all 19 steps from migration to navigation.
argument-hint: "[EntityName] [table_name]"
disable-model-invocation: true
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel Full CRUD Scaffold

Complete end-to-end scaffolding for `$ARGUMENTS`.

This skill orchestrates the full 19-step scaffolding sequence. Each step references a granular skill for details.

## Pre-Flight

Before starting, gather requirements:
1. **Entity name** (PascalCase): e.g., `Event`, `BusinessFormalisation`
2. **Table name** (snake_case plural): e.g., `events`, `business_formalisations`
3. **Fields**: List all fields with types
4. **Enums**: Any enum types needed (use `/goravel-enum` first)
5. **Relationships**: Foreign keys to other entities
6. **Permission scope**: Which actions (CRUD + extras)?
7. **Read-only fields**: Calculated/aggregated fields?

## Scaffolding Sequence

### Phase 1: Database & Model

- [ ] **Step 1**: Create migration
  ```bash
  go run . artisan make:migration create_<table>_table
  ```
- [ ] **Step 2**: Edit migration with field definitions, run `go run . artisan migrate`
- [ ] **Step 3**: Add audit fields
  ```bash
  go run . artisan make:audit --table=<table>
  ```
- [ ] **Step 4**: Fix audit migration (remove duplicate deleted_at), run `go run . artisan migrate`
- [ ] **Step 5**: Register migrations in `database/kernel.go`
- [ ] **Step 6**: Generate model
  ```bash
  go run . artisan make:model --table=<table> <ModelName>
  ```
- [ ] **Step 7**: Fix model (array fields, audit types, carbon.DateTime, SearchFields, TableName)
- [ ] **Step 8**: Create GORM interaction test, run it
  ```bash
  APP_ENV=testing go test -v ./tests/unit -run Test<Model>ModelTestSuite
  ```

See `/goravel-crud-migration` and `/goravel-crud-model` for details.

### Phase 2: Service & Permissions

- [ ] **Step 9**: Generate service
  ```bash
  go run . artisan make:svc --model=<Model> <entity>
  ```
- [ ] **Step 10**: Configure builder (search, sort, filter, validation, scope)
- [ ] **Step 11**: Register permissions in `app/auth/permission_constants.go`
- [ ] **Step 12**: Sync permissions
  ```bash
  go run . artisan permissions:setup
  ```

See `/goravel-crud-service` and `/goravel-crud-permissions` for details.

### Phase 3: Controller & Routes

- [ ] **Step 13**: Generate requests
  ```bash
  go run . artisan make:req --model=<entity> <entity>
  ```
- [ ] **Step 14**: Fix requests (validation rules, ToCreateData, carbon.DateTime, read-only exclusion)
- [ ] **Step 15**: Generate controller
  ```bash
  go run . artisan make:ctrl --model=<entity> <entity>
  ```
- [ ] **Step 16**: Fix controller (naming, service constant)
- [ ] **Step 17**: Register routes in `routes/api.go`

See `/goravel-crud-request`, `/goravel-crud-controller`, and `/goravel-crud-routes` for details.

### Phase 4: Testing

- [ ] **Step 18**: Generate CRUD tests
  ```bash
  go run . artisan make:crud-test --controller=<Entity>
  ```
- [ ] **Step 19**: Fix tests (permissions, endpoints, test data, arrays)
- [ ] **Step 20**: Run tests until all pass
  ```bash
  APP_ENV=testing go test -v ./tests/feature/crud -run Test<Entity>CRUDTestSuite
  ```

See `/goravel-crud-test` for details.

### Phase 5: UI & Navigation

- [ ] **Step 21**: Generate page controller
  ```bash
  go run . artisan make:page-ctrl --controller=<Entity>
  ```
- [ ] **Step 22**: Generate UI files
  ```bash
  go run . artisan make:ui --page=<Entity> --request=<Entity>
  ```
- [ ] **Step 23**: Register web route in `routes/web.go`
- [ ] **Step 24**: Add navigation entry in `resources/js/config/navigation.ts`
- [ ] **Step 25**: Add global search in `search_controller.go` and `search_config.tsx`

See `/goravel-crud-page`, `/goravel-crud-nav`, and `/goravel-crud-search` for details.

### Phase 6: Optional Enhancements

- [ ] Generate Swagger docs: `go run . artisan make:swagger-docs`
- [ ] Add simple filters (enum-based filtering with badges)
- [ ] Add statistics widget to dashboard
- [ ] Create database seeder

## Post-Scaffold Verification

1. All CRUD tests pass
2. API endpoints respond correctly (test with curl or Swagger)
3. Admin page loads and displays data
4. Navigation shows with correct permission gating
5. Global search includes the new entity
6. Different user roles see appropriate data (scope filtering)

## CRITICAL Rules

- **NEVER** modify core framework files (generic_crud_service, generic_crud_controller, CRUDPage)
- **ALWAYS** use artisan commands for generation, then fix
- **ALWAYS** test after each major phase
- Use hyphenated endpoints: `/entity-names` not `/entity_names`
- Service constants MUST match between controller and permission_constants.go
