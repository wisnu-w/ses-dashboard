#!/bin/bash

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-ses_monitoring}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}

echo "Running database optimization migration..."

# Run the migration
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f internal/infrastructure/database/migration/0008_add_search_indexes.sql

if [ $? -eq 0 ]; then
    echo "✅ Database indexes created successfully!"
    echo "Search and filter performance should be significantly improved."
else
    echo "❌ Failed to create database indexes."
    exit 1
fi