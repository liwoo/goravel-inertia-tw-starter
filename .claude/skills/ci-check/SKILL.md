---
name: ci-check
description: Run CI checks locally before pushing. Runs linting, type checks, tests, build verification, and Docker build in sequence - matching the GitHub Actions CI pipeline.
argument-hint: [--fast] [--fix]
disable-model-invocation: true
allowed-tools: Bash, Read, Grep
---

# Local CI Check

Run the CI pipeline locally to catch issues before pushing. Options: **$ARGUMENTS**

Parse arguments:
- `--fast`: Skip Docker build and security scan (quick feedback loop)
- `--fix`: Auto-fix linting issues where possible
- No arguments: run the full pipeline

## Pipeline Steps

Run each step sequentially. Stop on first failure unless noted. Track pass/fail for the summary.

### 1. Frontend Lint
```bash
cd /Users/liwu/Projects/go/blog && npx eslint . --max-warnings=0
```
If `--fix`: `npx eslint . --fix`
**Continue on failure** (matches CI behavior).

### 2. TypeScript Check
```bash
npx tsc --noEmit
```
**Continue on failure** (matches CI behavior).

### 3. Frontend Build
```bash
npm run build
```
**Stop on failure** - if Vite can't build, nothing else matters.

### 4. Go Lint
```bash
golangci-lint run --timeout 5m ./...
```
If golangci-lint is not installed, fall back to:
```bash
go vet ./...
```
If `--fix`: `golangci-lint run --fix ./...`
**Continue on failure**.

### 5. Go Build
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/ci-check-binary .
rm -f /tmp/ci-check-binary
```
**Stop on failure**.

### 6. Go Tests (skip if --fast)
```bash
./scripts/run_tests.sh -v -count=1 ./tests/...
```
If `run_tests.sh` is not available or fails to start the test container, fall back to:
```bash
go test -v -count=1 -timeout 600s ./tests/...
```
**Stop on failure**.

### 7. Frontend Tests (skip if --fast)
```bash
npx vitest run
```
**Stop on failure**.

### 8. Helm Lint (skip if --fast)
```bash
helm lint helm/goravel-blog
helm lint helm/goravel-blog -f helm/goravel-blog/values.staging.yaml
helm lint helm/goravel-blog -f helm/goravel-blog/values.production.yaml
```
**Continue on failure**.

### 9. Security Scan (skip if --fast)
```bash
govulncheck ./...
npm audit --audit-level=high
```
**Continue on failure** (informational).

### 10. Docker Build (skip if --fast)
```bash
docker build -t ci-check-local:test --target go-builder .
```
Only build to go-builder stage for speed. Clean up after.
**Continue on failure**.

## Summary

Output a results table:

| Step | Status | Duration |
|------|--------|----------|
| Frontend Lint | PASS/FAIL/SKIP | Xs |
| TypeScript Check | PASS/FAIL/SKIP | Xs |
| Frontend Build | PASS/FAIL | Xs |
| Go Lint | PASS/FAIL/SKIP | Xs |
| Go Build | PASS/FAIL | Xs |
| Go Tests | PASS/FAIL/SKIP | Xs |
| Frontend Tests | PASS/FAIL/SKIP | Xs |
| Helm Lint | PASS/FAIL/SKIP | Xs |
| Security Scan | PASS/FAIL/SKIP | Xs |
| Docker Build | PASS/FAIL/SKIP | Xs |

Then state:
- **Ready to push** if all critical steps (builds + tests) passed
- **Not ready** if any critical step failed, with the specific failures listed
