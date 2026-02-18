#!/bin/bash

# Database setup script for MIC Website Backend
# This script creates the database and runs migrations

set -e

echo "=== MIC Website Backend - Database Setup ==="
echo ""

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL environment variable is not set."
    echo "Please set it in your .env file or export it:"
    echo "  export DATABASE_URL='postgres://user:pass@localhost:5432/startups?sslmode=disable'"
    exit 1
fi

echo "Using DATABASE_URL: ${DATABASE_URL%%@*}"
echo ""

# Extract database name from DATABASE_URL
DB_NAME=$(echo $DATABASE_URL | sed -n 's/.*\/\([^?]*\).*/\1/p')

if [ -z "$DB_NAME" ]; then
    echo "Error: Could not extract database name from DATABASE_URL"
    exit 1
fi

# Extract connection info without database name
DB_CONN=$(echo $DATABASE_URL | sed 's/\/[^?]*/\/postgres/')

echo "Creating database: $DB_NAME"
echo ""

# Create database (connect to default 'postgres' database to create new one)
psql "$DB_CONN" -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || echo "Database $DB_NAME already exists or connection failed (you may need to create it manually)"
echo ""

echo "Running migrations..."
echo ""

# Run migrations
for migration in migrations/*.up.sql; do
    if [ -f "$migration" ]; then
        echo "Running migration: $(basename $migration)"
        psql "$DATABASE_URL" -f "$migration" || {
            echo "Error running migration: $migration"
            exit 1
        }
    fi
done

echo ""
echo "=== Database setup complete! ==="
echo ""
echo "You can now start the server with:"
echo "  go run cmd/server/main.go"

