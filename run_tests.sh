#!/bin/bash

# Test runner script for organized test structure

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
TEST_TYPE=""
VERBOSE=""
SPECIFIC_TEST=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--type)
            TEST_TYPE="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE="-v"
            shift
            ;;
        -s|--specific)
            SPECIFIC_TEST="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [options]"
            echo "Options:"
            echo "  -t, --type <type>      Run specific test type (unit|integration|feature|permissions|auth|crud|all)"
            echo "  -v, --verbose          Run tests with verbose output"
            echo "  -s, --specific <test>  Run specific test by name"
            echo "  -h, --help             Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Ensure test environment
export APP_ENV=testing
mkdir -p resources/views && touch resources/views/dummy.tmpl

# Function to run tests and report results
run_tests() {
    local test_path=$1
    local test_name=$2
    
    echo -e "${YELLOW}Running $test_name tests...${NC}"
    
    if [[ -n "$SPECIFIC_TEST" ]]; then
        APP_ENV=testing go test $VERBOSE $test_path -run "$SPECIFIC_TEST"
    else
        APP_ENV=testing go test $VERBOSE $test_path
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $test_name tests passed${NC}"
    else
        echo -e "${RED}✗ $test_name tests failed${NC}"
        exit 1
    fi
}

# Main test execution
case $TEST_TYPE in
    unit)
        run_tests "./tests/unit/..." "Unit"
        ;;
    integration)
        run_tests "./tests/integration/..." "Integration"
        ;;
    feature)
        run_tests "./tests/feature/..." "Feature"
        ;;
    permissions)
        run_tests "./tests/feature/permissions/..." "Permissions"
        ;;
    auth)
        run_tests "./tests/feature/auth/..." "Authentication"
        ;;
    crud)
        run_tests "./tests/feature/crud/..." "CRUD"
        ;;
    all|"")
        echo -e "${YELLOW}Running all tests...${NC}"
        run_tests "./tests/unit/..." "Unit"
        run_tests "./tests/integration/..." "Integration"
        run_tests "./tests/feature/..." "Feature"
        echo -e "${GREEN}✓ All tests passed!${NC}"
        ;;
    *)
        echo -e "${RED}Unknown test type: $TEST_TYPE${NC}"
        echo "Valid types: unit, integration, feature, permissions, auth, crud, all"
        exit 1
        ;;
esac

# Clean up test data if all tests passed
if [ $? -eq 0 ]; then
    echo -e "${YELLOW}Cleaning up test data...${NC}"
    rm -rf tests/feature/storage/framework/sessions/*
    touch tests/feature/storage/framework/sessions/.gitkeep
    echo -e "${GREEN}✓ Cleanup complete${NC}"
fi