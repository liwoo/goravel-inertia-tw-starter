#!/bin/bash

# Script to clean the test database

echo "Cleaning test database..."

# Remove the test database file
if [ -f "database/test.sqlite" ]; then
    rm database/test.sqlite
    echo "✓ Removed test.sqlite"
fi

# Create a new empty database
touch database/test.sqlite
echo "✓ Created new test.sqlite"

# Run migrations
echo "Running migrations..."
APP_ENV=testing go run . artisan migrate

echo "✓ Test database cleaned and ready!"