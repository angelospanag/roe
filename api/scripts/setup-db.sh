#!/bin/bash

# Database setup script for ROE

set -e

# Default values
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-roe_backend}"

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "Setting up database..."
echo "Database URL: ${DATABASE_URL}"

# Check if migrate is installed
if ! command -v migrate &> /dev/null; then
    echo "Error: golang-migrate is not installed."
    echo "Install it with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Run migrations
echo "Running migrations..."
migrate -database "${DATABASE_URL}" -path db/migrations up

echo "Database setup complete!"
