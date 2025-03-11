#!/bin/bash
set -e

echo "Waiting for database to be ready..."
max_retries=30
counter=0
until pg_isready -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER}; do
    >&2 echo "Postgres is unavailable - sleeping"
    counter=$((counter+1))
    if [ $counter -eq $max_retries ]; then
        echo "Failed to connect to database after $max_retries attempts. Exiting."
        exit 1
    fi
    sleep 2
done

echo "Database is ready. Running migrations..."

# Run migrations
/app/migrate -path=/app/database/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

echo "Migrations completed successfully." 