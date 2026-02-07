#!/bin/bash
# post-edit-lint.sh — PostToolUse hook for Edit|Write
# Auto-formats files and runs targeted lint checks after every file change.
#
# Go files:   goimports (format + organize imports), then go vet on the package
# TS/TSX:     prettier --write
# Helm YAML:  helm lint
#
# This hook runs on PostToolUse (non-blocking). Output is fed back to Claude
# as context so it can fix any issues automatically.

set -euo pipefail

INPUT=$(cat)
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

if [ -z "$FILE_PATH" ] || [ ! -f "$FILE_PATH" ]; then
    exit 0
fi

EXTENSION="${FILE_PATH##*.}"
PROJECT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"

case "$EXTENSION" in
    go)
        # Auto-format with goimports (handles imports + gofmt)
        if command -v goimports &>/dev/null; then
            goimports -w "$FILE_PATH" 2>/dev/null || true
        elif command -v gofmt &>/dev/null; then
            gofmt -w "$FILE_PATH" 2>/dev/null || true
        fi

        # Run go vet on the package containing the file
        PACKAGE_DIR=$(dirname "$FILE_PATH")
        VET_OUTPUT=$(cd "$PROJECT_DIR" && go vet "$PACKAGE_DIR/..." 2>&1) || true
        if [ -n "$VET_OUTPUT" ]; then
            echo "go vet issues in $(basename "$PACKAGE_DIR"):"
            echo "$VET_OUTPUT" | head -20
        fi
        ;;

    ts|tsx)
        # Auto-format with prettier
        if [ -f "$PROJECT_DIR/node_modules/.bin/prettier" ]; then
            cd "$PROJECT_DIR" && npx prettier --write "$FILE_PATH" 2>/dev/null || true
        fi
        ;;

    yaml|yml)
        # Helm lint if it's a helm chart file
        if echo "$FILE_PATH" | grep -q "helm/goravel-blog/"; then
            LINT_OUTPUT=$(helm lint "$PROJECT_DIR/helm/goravel-blog" 2>&1) || true
            if echo "$LINT_OUTPUT" | grep -qE '(ERROR|WARNING)'; then
                echo "Helm lint:"
                echo "$LINT_OUTPUT" | grep -E '(ERROR|WARNING|Chart)' | head -10
            fi
        fi
        ;;
esac

exit 0
