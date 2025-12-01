#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
TEST_PATH="./tests/..."
VERBOSE="-v"
TIMEOUT="30m"
COVERAGE=""

# Function to print colored output
print_color() {
    color=$1
    message=$2
    echo -e "${color}${message}${NC}"
}

# Function to show usage
show_usage() {
    echo "Usage: ./run_tests.sh [options] [test-path]"
    echo ""
    echo "Options:"
    echo "  all         Run all tests (default)"
    echo "  unit        Run unit tests only"
    echo "  integration Run integration tests only"
    echo "  feature     Run feature tests only"
    echo "  coverage    Run with coverage report"
    echo "  quick       Run without verbose output"
    echo "  clean       Clean test artifacts before running"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./run_tests.sh                    # Run all tests"
    echo "  ./run_tests.sh unit               # Run unit tests"
    echo "  ./run_tests.sh coverage           # Run all tests with coverage"
    echo "  ./run_tests.sh clean feature      # Clean and run feature tests"
    echo "  ./run_tests.sh TestBookCRUD       # Run specific test"
}

# Function to clean test artifacts
clean_artifacts() {
    print_color "$YELLOW" "🧹 Cleaning test artifacts..."
    rm -rf tests/*/database/
    rm -rf tests/*/storage/
    find . -name "*.sqlite" -path "*/tests/*" -delete
    print_color "$GREEN" "✅ Test artifacts cleaned"
}

# Function to setup test environment
setup_test_env() {
    # Create dummy view file to prevent template errors
    mkdir -p resources/views
    touch resources/views/dummy.tmpl
    
    # Ensure test database directory exists
    mkdir -p database
    
    # Set test environment
    export APP_ENV=testing
}

# Function to run tests
run_tests() {
    local test_type=$1
    local test_path=$2
    
    print_color "$YELLOW" "🧪 Running $test_type tests..."
    print_color "$YELLOW" "📍 Path: $test_path"
    
    # Build test command with APP_ENV explicitly set
    cmd="APP_ENV=testing go test $VERBOSE $test_path -timeout=$TIMEOUT"
    if [ -n "$COVERAGE" ]; then
        cmd="$cmd -cover -coverprofile=coverage.out"
    fi

    # Run tests
    if eval $cmd; then
        print_color "$GREEN" "✅ $test_type tests passed!"
        return 0
    else
        print_color "$RED" "❌ $test_type tests failed!"
        return 1
    fi
}

# Parse arguments
while [ $# -gt 0 ]; do
    case "$1" in
        help|--help|-h)
            show_usage
            exit 0
            ;;
        all)
            TEST_PATH="./tests/..."
            ;;
        unit)
            TEST_PATH="./tests/unit/..."
            ;;
        integration)
            TEST_PATH="./tests/integration/..."
            ;;
        feature)
            TEST_PATH="./tests/feature/..."
            ;;
        coverage)
            COVERAGE="1"
            ;;
        quick)
            VERBOSE=""
            ;;
        clean)
            clean_artifacts
            ;;
        Test*)
            # Specific test name
            TEST_PATH="./tests/... -run $1"
            ;;
        *)
            # Custom path
            if [[ $1 == ./* ]]; then
                TEST_PATH="$1"
            else
                print_color "$RED" "Unknown option: $1"
                show_usage
                exit 1
            fi
            ;;
    esac
    shift
done

# Main execution
print_color "$GREEN" "🚀 Goravel Blog Test Runner"
print_color "$GREEN" "=========================="

# Setup environment
setup_test_env

# Determine test type for display
TEST_TYPE="all"
if [[ $TEST_PATH == *"unit"* ]]; then
    TEST_TYPE="unit"
elif [[ $TEST_PATH == *"integration"* ]]; then
    TEST_TYPE="integration"
elif [[ $TEST_PATH == *"feature"* ]]; then
    TEST_TYPE="feature"
fi

# Run tests
if run_tests "$TEST_TYPE" "$TEST_PATH"; then
    EXIT_CODE=0
    
    # Show coverage report if requested
    if [ -n "$COVERAGE" ]; then
        print_color "$YELLOW" "📊 Coverage Report:"
        go tool cover -func=coverage.out | tail -10
        print_color "$YELLOW" "💡 To view detailed coverage: go tool cover -html=coverage.out"
    fi
else
    EXIT_CODE=1
fi

print_color "$GREEN" "=========================="
print_color "$GREEN" "🏁 Test run completed"

exit $EXIT_CODE