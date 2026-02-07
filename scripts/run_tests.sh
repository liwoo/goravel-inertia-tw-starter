#!/bin/bash
# Script to run tests with an isolated PostgreSQL testcontainer
# This ensures environment variables are set BEFORE Go compiles the test binary
# and isolates tests completely from the main application

# Don't use set -e as it can interfere with cleanup on test failures
# Instead, we capture the exit code and handle cleanup explicitly

echo "Starting PostgreSQL testcontainer..."

# Generate unique identifier for this test run
TEST_RUN_ID=$$

# Cleanup function - called on exit, interrupt, or termination
cleanup() {
    local exit_code=$?
    echo ""
    echo "Cleaning up test container and temp files..."
    if [ -n "$CONTAINER_ID" ]; then
        docker rm -f "$CONTAINER_ID" > /dev/null 2>&1 || true
    fi
    if [ -n "$TEST_STORAGE_DIR" ] && [ -d "$TEST_STORAGE_DIR" ]; then
        rm -rf "$TEST_STORAGE_DIR" 2>/dev/null || true
    fi
    exit $exit_code
}

# Set up trap for multiple signals to ensure cleanup always runs
trap cleanup EXIT INT TERM

# Start a PostgreSQL container
CONTAINER_ID=$(docker run -d \
    --name books-test-db-$TEST_RUN_ID \
    -e POSTGRES_USER=testuser \
    -e POSTGRES_PASSWORD=testpassword123 \
    -e POSTGRES_DB=books_test \
    -p 0:5432 \
    postgres:16-alpine)

if [ -z "$CONTAINER_ID" ]; then
    echo "Failed to start PostgreSQL container"
    exit 1
fi

echo "Container started: $CONTAINER_ID"

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
    if docker exec "$CONTAINER_ID" pg_isready -U testuser -d books_test > /dev/null 2>&1; then
        echo "PostgreSQL is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "PostgreSQL failed to become ready in time"
        exit 1
    fi
    echo "Waiting... ($i/30)"
    sleep 1
done

# Get the mapped port
MAPPED_PORT=$(docker port "$CONTAINER_ID" 5432 | cut -d: -f2)
if [ -z "$MAPPED_PORT" ]; then
    echo "Failed to get mapped port"
    exit 1
fi
echo "PostgreSQL available on port: $MAPPED_PORT"

# Create isolated storage directory for tests
TEST_STORAGE_DIR="/tmp/books-test-storage-$TEST_RUN_ID"
mkdir -p "$TEST_STORAGE_DIR/framework/sessions"
mkdir -p "$TEST_STORAGE_DIR/framework/cache"
mkdir -p "$TEST_STORAGE_DIR/logs"
echo "Test storage directory: $TEST_STORAGE_DIR"

# Set environment variables for tests - completely isolated from main app
export APP_ENV=testing
export APP_DEBUG=true
export APP_KEY=testkeyfortestingonlyabc12345678
export APP_NAME=books_test

# Disable 2FA requirement for tests (allows simple JWT auth without TOTP)
export AUTH_REQUIRE_2FA=false

# Database - use the test container
export DB_CONNECTION=postgres
export DB_HOST=127.0.0.1
export DB_PORT=$MAPPED_PORT
export DB_DATABASE=books_test
export DB_USERNAME=testuser
export DB_PASSWORD=testpassword123
export DB_SSLMODE=disable

# Use isolated storage paths
export STORAGE_PATH="$TEST_STORAGE_DIR"

# Session - use isolated file storage
export SESSION_DRIVER=file
export SESSION_FILES_PATH="$TEST_STORAGE_DIR/framework/sessions"
export SESSION_COOKIE="books_test_session_$TEST_RUN_ID"

# JWT - use a completely different secret for tests
export JWT_SECRET="test-jwt-secret-$TEST_RUN_ID-isolated-from-prod"

# HTTP request timeout - increase for CI where bcrypt is slower
export HTTP_REQUEST_TIMEOUT=30

# Disable Redis for tests (use memory/file instead)
export REDIS_HOST=""
export CACHE_DRIVER=file
export CACHE_PATH="$TEST_STORAGE_DIR/framework/cache"

# Logging
export LOG_CHANNEL=single
export LOG_PATH="$TEST_STORAGE_DIR/logs/test.log"

echo "Running tests with isolated environment..."
echo "DB_HOST=$DB_HOST DB_PORT=$DB_PORT DB_DATABASE=$DB_DATABASE"

# Run the tests with the -p=1 flag to run packages sequentially
# This prevents race conditions when multiple packages try to migrate/seed simultaneously
# Store the exit code so cleanup runs but we still return the correct status
go test -p=1 "$@"
TEST_EXIT_CODE=$?

# Exit with the test exit code (cleanup will run via trap)
exit $TEST_EXIT_CODE
