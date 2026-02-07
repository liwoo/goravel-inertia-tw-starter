#!/bin/bash
# post-bash-check.sh — PostToolUse hook for Bash
# Runs quality checks after specific shell commands:
#
# After artisan make:* generators → go vet on generated code
# After make:ts-enums             → tsc --noEmit on types
# After go run . artisan migrate  → go build to verify schema
#
# This hook runs on PostToolUse (non-blocking). Output is fed back to Claude.

set -euo pipefail

INPUT=$(cat)
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // empty')
PROJECT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"

if [ -z "$COMMAND" ]; then
    exit 0
fi

# After artisan code generators, run go vet
if echo "$COMMAND" | grep -qE 'artisan make:(model|svc|ctrl|req|page-ctrl|crud-test|audit)'; then
    VET_OUTPUT=$(cd "$PROJECT_DIR" && go vet ./... 2>&1) || true
    if [ -n "$VET_OUTPUT" ]; then
        echo "go vet after code generation:"
        echo "$VET_OUTPUT" | head -30
    fi
fi

# After TS enum generation, run type check
if echo "$COMMAND" | grep -q 'make:ts-enums'; then
    TSC_OUTPUT=$(cd "$PROJECT_DIR" && npx tsc --noEmit 2>&1 | head -20) || true
    if [ -n "$TSC_OUTPUT" ]; then
        echo "TypeScript check after enum generation:"
        echo "$TSC_OUTPUT"
    fi
fi

# After migration commands, verify Go builds
if echo "$COMMAND" | grep -qE 'artisan migrate($|[^:]| )'; then
    BUILD_OUTPUT=$(cd "$PROJECT_DIR" && go build ./... 2>&1) || true
    if [ -n "$BUILD_OUTPUT" ]; then
        echo "Build check after migration:"
        echo "$BUILD_OUTPUT" | head -20
    fi
fi

# After npm/pnpm install, check for audit issues
if echo "$COMMAND" | grep -qE '(npm|pnpm|yarn) install'; then
    AUDIT_OUTPUT=$(cd "$PROJECT_DIR" && npm audit --audit-level=high 2>&1 | tail -5) || true
    if echo "$AUDIT_OUTPUT" | grep -qE 'found [1-9]'; then
        echo "npm audit:"
        echo "$AUDIT_OUTPUT"
    fi
fi

exit 0
