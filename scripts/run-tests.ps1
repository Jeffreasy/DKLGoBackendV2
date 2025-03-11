param (
    [switch]$Unit,
    [switch]$Integration,
    [switch]$Coverage
)

# Start test database
Write-Host "Starting test database..." -ForegroundColor Cyan
docker-compose -f docker-compose.test.yml up -d test-db

# Wait for database to be ready
Write-Host "Waiting for test database to be ready..." -ForegroundColor Cyan
Start-Sleep -Seconds 5

try {
    if ($Unit -or (!$Unit -and !$Integration)) {
        # Run unit tests
        Write-Host "Running unit tests..." -ForegroundColor Green
        if ($Coverage) {
            go test -v -coverprofile=coverage.out ./auth/... ./database/... ./handlers/... ./models/... ./services/...
            go tool cover -html=coverage.out -o coverage.html
            Start-Process coverage.html
        } else {
            go test -v ./auth/... ./database/... ./handlers/... ./models/... ./services/...
        }
    }

    if ($Integration) {
        # Run integration tests
        Write-Host "Running integration tests..." -ForegroundColor Green
        go test -v ./tests/integration/...
    }
} finally {
    # Stop test database
    Write-Host "Stopping test database..." -ForegroundColor Cyan
    docker-compose -f docker-compose.test.yml down
} 