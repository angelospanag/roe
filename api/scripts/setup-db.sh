#!/bin/bash

# Database setup script for ROE

set -e

DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-roe_backend}"

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "Running migrations against ${DATABASE_URL}..."
go tool goose -dir db/migrations postgres "${DATABASE_URL}" up

echo "Database setup complete!"
