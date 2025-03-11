#!/usr/bin/env pwsh
# Script to run integration tests with the correct environment variables

# Start the test database container
Write-Host "Starting test database container..." -ForegroundColor Cyan
docker-compose -f docker-compose.test.yml up -d test-db

# Wait for the database to be ready
Write-Host "Waiting for test database to be ready..." -ForegroundColor Cyan
Start-Sleep -Seconds 5

# Create the test database if it doesn't exist
Write-Host "Creating test database if it doesn't exist..." -ForegroundColor Cyan
docker exec dklautomationgo-test-db psql -U postgres -c "DROP DATABASE IF EXISTS dklautomationgo_test;" 2>&1 | Out-Null
docker exec dklautomationgo-test-db psql -U postgres -c "CREATE DATABASE dklautomationgo_test WITH ENCODING 'UTF8' OWNER postgres;" 2>&1 | Out-Null

# Set environment variables for the test database
$env:TEST_DB_HOST = "localhost"
$env:TEST_DB_PORT = "5433"
$env:TEST_DB_USER = "postgres"
$env:TEST_DB_PASSWORD = "Bootje@12"
$env:TEST_DB_NAME = "dklautomationgo_test"

try {
    # Run the integration tests
    Write-Host "Running integration tests..." -ForegroundColor Green
    go test -v ./tests/integration/...
}
finally {
    # Stop the test database container
    Write-Host "Stopping test database container..." -ForegroundColor Cyan
    docker-compose -f docker-compose.test.yml down
} 