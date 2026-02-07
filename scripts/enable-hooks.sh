#!/bin/bash

# Enable git hooks for the project

echo "Configuring Git to use .githooks directory..."
git config core.hooksPath .githooks

echo "✅ Git hooks enabled!"
echo "The pre-commit hook will now run before each commit."
echo ""
echo "To disable hooks temporarily, use: git commit --no-verify"
echo "To disable hooks permanently, use: git config --unset core.hooksPath"