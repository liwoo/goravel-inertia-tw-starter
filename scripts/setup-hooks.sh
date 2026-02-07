#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up Git hooks for Goravel Blog...${NC}"

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo -e "${YELLOW}pre-commit is not installed. Installing...${NC}"
    
    # Try to install pre-commit based on the OS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install pre-commit
        else
            echo -e "${RED}Homebrew not found. Please install pre-commit manually:${NC}"
            echo "pip install pre-commit"
            exit 1
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        if command -v pip &> /dev/null; then
            pip install pre-commit
        else
            echo -e "${RED}pip not found. Please install pre-commit manually:${NC}"
            echo "https://pre-commit.com/#install"
            exit 1
        fi
    else
        echo -e "${RED}Unsupported OS. Please install pre-commit manually:${NC}"
        echo "https://pre-commit.com/#install"
        exit 1
    fi
fi

# Install the git hook scripts
echo -e "${YELLOW}Installing pre-commit hooks...${NC}"
pre-commit install

# Install golangci-lint if not present
if ! command -v golangci-lint &> /dev/null; then
    echo -e "${YELLOW}golangci-lint is not installed. Installing...${NC}"
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
fi

# Run pre-commit on all files to check current state
echo -e "${YELLOW}Running pre-commit on all files...${NC}"
pre-commit run --all-files || true

echo -e "${GREEN}✅ Git hooks setup complete!${NC}"
echo -e "${GREEN}Pre-commit will now run automatically before each commit.${NC}"
echo ""
echo -e "${YELLOW}To manually run hooks on all files:${NC}"
echo "pre-commit run --all-files"
echo ""
echo -e "${YELLOW}To skip hooks for a single commit:${NC}"
echo "git commit --no-verify"