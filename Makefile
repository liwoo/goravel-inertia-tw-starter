.PHONY: test test-feature test-integration test-unit test-coverage help

# Default target
help:
	@echo "Available commands:"
	@echo "  make test              - Run all tests"
	@echo "  make test-feature      - Run feature tests only"
	@echo "  make test-integration  - Run integration tests only"
	@echo "  make test-unit         - Run unit tests only"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make test-watch        - Run tests in watch mode (requires entr)"

# Prepare test environment
prepare-test:
	@mkdir -p resources/views
	@touch resources/views/dummy.tmpl

# Run all tests
test: prepare-test
	@echo "🧪 Running all tests..."
	@APP_ENV=testing go test -v ./tests/... -timeout=10m

# Run feature tests
test-feature: prepare-test
	@echo "🧪 Running feature tests..."
	@APP_ENV=testing go test -v ./tests/feature/... -timeout=10m

# Run integration tests
test-integration: prepare-test
	@echo "🧪 Running integration tests..."
	@APP_ENV=testing go test -v ./tests/integration/... -timeout=10m

# Run unit tests
test-unit: prepare-test
	@echo "🧪 Running unit tests..."
	@APP_ENV=testing go test -v ./tests/unit/... -timeout=10m

# Run tests with coverage
test-coverage: prepare-test
	@echo "🧪 Running tests with coverage..."
	@APP_ENV=testing go test -v -cover -coverprofile=coverage.out ./tests/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Watch mode for tests (requires entr to be installed)
test-watch: prepare-test
	@echo "👀 Watching for changes..."
	@find . -name '*.go' | entr -c sh -c 'APP_ENV=testing go test -v ./tests/...'