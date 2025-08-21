#!/bin/bash

# Simple linting script for Goravel Blog

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Running Go linting checks...${NC}"

# 1. gofmt
echo -e "${YELLOW}Checking code formatting with gofmt...${NC}"
GOFMT_FILES=$(gofmt -l .)
if [ -n "$GOFMT_FILES" ]; then
    echo -e "${RED}The following files need formatting:${NC}"
    echo "$GOFMT_FILES"
    echo -e "${YELLOW}Run 'gofmt -w -s .' to fix${NC}"
    exit 1
else
    echo -e "${GREEN}✅ Code formatting OK${NC}"
fi

# 2. go vet
echo -e "${YELLOW}Running go vet...${NC}"
if go vet ./...; then
    echo -e "${GREEN}✅ go vet passed${NC}"
else
    echo -e "${RED}❌ go vet found issues${NC}"
    exit 1
fi

# 3. staticcheck (if available)
if command -v staticcheck &> /dev/null; then
    echo -e "${YELLOW}Running staticcheck...${NC}"
    if staticcheck ./...; then
        echo -e "${GREEN}✅ staticcheck passed${NC}"
    else
        echo -e "${RED}❌ staticcheck found issues${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}staticcheck not installed, skipping...${NC}"
    echo -e "${YELLOW}Install with: go install honnef.co/go/tools/cmd/staticcheck@latest${NC}"
fi

# 4. Check for unused dependencies
echo -e "${YELLOW}Checking for unused dependencies...${NC}"
go mod tidy
if [ -n "$(git status --porcelain go.mod go.sum)" ]; then
    echo -e "${RED}go.mod or go.sum was modified by go mod tidy${NC}"
    echo -e "${YELLOW}Please run 'go mod tidy' and commit the changes${NC}"
    exit 1
else
    echo -e "${GREEN}✅ Dependencies OK${NC}"
fi

echo -e "${GREEN}✨ All linting checks passed!${NC}"