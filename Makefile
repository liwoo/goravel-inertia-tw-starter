# =============================================================================
# Makefile for Books Database Application
# =============================================================================

.PHONY: help build test lint dev clean docker helm deploy

# Default target
.DEFAULT_GOAL := help

# =============================================================================
# Variables
# =============================================================================
DOCKER_REGISTRY ?= docker.io
DOCKER_IMAGE_NAME ?= books-database
APP_NAME := $(DOCKER_IMAGE_NAME)
GO_VERSION := 1.24
NODE_VERSION := 20
IMAGE_NAME := $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME)
HELM_RELEASE := goravel-blog
HELM_CHART := ./helm/goravel-blog

# =============================================================================
# Help
# =============================================================================
help: ## Show this help message
	@echo "Books Database - Development Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"; } /^[a-zA-Z_-]+:.*?##/ { printf "  %-25s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# =============================================================================
# Development
# =============================================================================
dev: ## Start development environment with hot reload
	docker-compose -f docker-compose.dev.yml up -d
	@echo "Development environment started!"
	@echo "  Backend: http://localhost:3000"
	@echo "  Frontend: http://localhost:5173"

dev-stop: ## Stop development environment
	docker-compose -f docker-compose.dev.yml down

dev-logs: ## Follow development logs
	docker-compose -f docker-compose.dev.yml logs -f

dev-tools: ## Start dev environment with tools (pgAdmin, Redis Commander)
	docker-compose -f docker-compose.dev.yml --profile tools up -d

# =============================================================================
# Build
# =============================================================================
build: build-backend build-frontend ## Build both backend and frontend

build-backend: ## Build Go backend
	@echo "Building Go backend..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -a -trimpath \
		-ldflags="-w -s -extldflags '-static'" \
		-o bin/$(APP_NAME) .
	@echo "Backend built: bin/$(APP_NAME)"

build-frontend: ## Build React frontend
	@echo "Building React frontend..."
	npm ci
	npm run build
	@echo "Frontend built: public/"

# =============================================================================
# Testing
# =============================================================================
prepare-test:
	@mkdir -p resources/views
	@touch resources/views/dummy.tmpl

test: test-backend test-frontend ## Run all tests

test-backend: prepare-test ## Run Go tests
	@echo "Running Go tests..."
	APP_ENV=testing go test -v -race -coverprofile=coverage.out ./tests/...

test-frontend: ## Run React tests
	@echo "Running frontend tests..."
	npm run test -- --run

test-feature: prepare-test ## Run feature tests
	@echo "Running feature tests..."
	APP_ENV=testing go test -v ./tests/feature/... -timeout=10m

test-integration: prepare-test ## Run integration tests
	@echo "Running integration tests..."
	APP_ENV=testing go test -v ./tests/integration/... -timeout=10m

test-unit: prepare-test ## Run unit tests
	@echo "Running unit tests..."
	APP_ENV=testing go test -v ./tests/unit/... -timeout=10m

test-coverage: prepare-test ## Run tests with coverage report
	@echo "Running tests with coverage..."
	APP_ENV=testing go test -v -cover -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	npm run test:coverage -- --run
	@echo "Coverage reports generated"

test-docker: ## Run integration tests with Docker
	@echo "Running integration tests in Docker..."
	docker-compose -f docker-compose.test.yml run --rm test

test-e2e: ## Run end-to-end tests
	@echo "Running E2E tests..."
	docker-compose -f docker-compose.test.yml --profile e2e run --rm e2e-test

test-watch: prepare-test ## Watch mode for tests (requires entr)
	@echo "Watching for changes..."
	@find . -name '*.go' | entr -c sh -c 'APP_ENV=testing go test -v ./tests/...'

# =============================================================================
# Linting
# =============================================================================
lint: lint-backend lint-frontend ## Run all linters

lint-backend: ## Lint Go code
	@echo "Linting Go code..."
	golangci-lint run ./...

lint-frontend: ## Lint frontend code
	@echo "Linting frontend code..."
	npm run lint --if-present || true
	npx tsc --noEmit || true

lint-fix: ## Fix linting issues automatically
	@echo "Fixing lint issues..."
	gofmt -w -s .
	npm run lint --fix --if-present || true

# =============================================================================
# Docker
# =============================================================================
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME):latest \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg NODE_VERSION=$(NODE_VERSION) \
		.
	@echo "Docker image built: $(IMAGE_NAME):latest"

docker-build-dev: ## Build Docker image for development
	docker build -t $(IMAGE_NAME):dev \
		--target go-builder \
		.

docker-push: ## Push Docker image to registry
	@echo "Pushing Docker image..."
	docker push $(IMAGE_NAME):latest

docker-scan: ## Scan Docker image for vulnerabilities
	@echo "Scanning Docker image..."
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
		aquasec/trivy:latest image $(IMAGE_NAME):latest

docker-up: ## Start production Docker environment
	docker-compose up -d

docker-down: ## Stop Docker environment
	docker-compose down

docker-logs: ## Follow Docker logs
	docker-compose logs -f

# =============================================================================
# Helm
# =============================================================================
helm-lint: ## Lint Helm chart
	@echo "Linting Helm chart..."
	helm lint $(HELM_CHART)

helm-template: ## Generate Helm templates
	helm template $(HELM_RELEASE) $(HELM_CHART)

helm-template-staging: ## Generate staging Helm templates
	helm template $(HELM_RELEASE) $(HELM_CHART) \
		--values $(HELM_CHART)/values.yaml \
		--values $(HELM_CHART)/values.staging.yaml

helm-template-production: ## Generate production Helm templates
	helm template $(HELM_RELEASE) $(HELM_CHART) \
		--values $(HELM_CHART)/values.yaml \
		--values $(HELM_CHART)/values.production.yaml

helm-diff-staging: ## Show diff for staging deployment
	helm diff upgrade $(HELM_RELEASE) $(HELM_CHART) \
		--namespace staging \
		--values $(HELM_CHART)/values.yaml \
		--values $(HELM_CHART)/values.staging.yaml

helm-diff-production: ## Show diff for production deployment
	helm diff upgrade $(HELM_RELEASE) $(HELM_CHART) \
		--namespace production \
		--values $(HELM_CHART)/values.yaml \
		--values $(HELM_CHART)/values.production.yaml

# =============================================================================
# Database
# =============================================================================
db-migrate: ## Run database migrations
	@echo "Running migrations..."
	go run . artisan migrate

db-seed: ## Seed the database
	@echo "Seeding database..."
	go run . artisan db:seed

db-fresh: ## Drop all tables and re-run migrations
	@echo "Warning: This will drop all tables!"
	go run . artisan migrate:fresh

# =============================================================================
# Clean
# =============================================================================
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf tmp/
	rm -rf coverage/
	rm -f coverage.out coverage.html
	rm -rf public/js public/css public/assets
	@echo "Cleaned"

clean-docker: ## Clean Docker resources
	@echo "Cleaning Docker resources..."
	docker-compose down -v --remove-orphans
	docker-compose -f docker-compose.dev.yml down -v --remove-orphans
	docker-compose -f docker-compose.test.yml down -v --remove-orphans

clean-all: clean clean-docker ## Clean everything
	rm -rf node_modules/
	@echo "All cleaned"

# =============================================================================
# Dependencies
# =============================================================================
deps: deps-backend deps-frontend ## Install all dependencies

deps-backend: ## Install Go dependencies
	@echo "Installing Go dependencies..."
	go mod download
	go mod verify

deps-frontend: ## Install frontend dependencies
	@echo "Installing frontend dependencies..."
	npm ci

deps-update: ## Update all dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	npm update

# =============================================================================
# Security
# =============================================================================
security-scan: ## Run security scans
	@echo "Running security scans..."
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./... || true
	npm audit || true

# =============================================================================
# Misc
# =============================================================================
generate: ## Run code generators
	@echo "Running code generators..."
	go generate ./...

fmt: ## Format code
	@echo "Formatting code..."
	gofmt -w -s .

pre-commit: lint test ## Run pre-commit checks
	@echo "Pre-commit checks passed!"

ci: lint test-backend build ## Run CI pipeline locally
	@echo "CI checks passed!"
