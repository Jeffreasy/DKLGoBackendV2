#!/usr/bin/env pwsh
# Docker Helper Script for DKL Automationgo

param (
    [Parameter(Position=0, Mandatory=$true)]
    [ValidateSet("start", "stop", "restart", "logs", "build", "shell", "db", "backup", "restore")]
    [string]$Command,
    
    [Parameter(Position=1, Mandatory=$false)]
    [string]$Arg
)

$ErrorActionPreference = "Stop"

# Check if Docker is installed
try {
    docker --version | Out-Null
} catch {
    Write-Error "Docker is not installed or not in PATH. Please install Docker Desktop."
    exit 1
}

# Check if Docker Compose is installed
try {
    docker-compose --version | Out-Null
} catch {
    Write-Error "Docker Compose is not installed or not in PATH. Please install Docker Desktop."
    exit 1
}

# Function to check if containers are running
function Test-ContainersRunning {
    $containers = docker ps --filter "name=dklautomationgo" --format "{{.Names}}"
    return $containers.Count -gt 0
}

# Execute the requested command
switch ($Command) {
    "start" {
        Write-Host "Starting containers..." -ForegroundColor Green
        docker-compose up -d
        Write-Host "Containers started. Application is available at http://localhost:8080" -ForegroundColor Green
    }
    "stop" {
        Write-Host "Stopping containers..." -ForegroundColor Yellow
        docker-compose down
        Write-Host "Containers stopped." -ForegroundColor Yellow
    }
    "restart" {
        Write-Host "Restarting containers..." -ForegroundColor Yellow
        docker-compose down
        docker-compose up -d
        Write-Host "Containers restarted. Application is available at http://localhost:8080" -ForegroundColor Green
    }
    "logs" {
        if ($Arg -eq "db") {
            Write-Host "Showing database logs..." -ForegroundColor Cyan
            docker-compose logs -f db
        } else {
            Write-Host "Showing application logs..." -ForegroundColor Cyan
            docker-compose logs -f app
        }
    }
    "build" {
        Write-Host "Building and starting containers..." -ForegroundColor Green
        docker-compose up -d --build
        Write-Host "Containers built and started. Application is available at http://localhost:8080" -ForegroundColor Green
    }
    "shell" {
        if (-not (Test-ContainersRunning)) {
            Write-Error "Containers are not running. Start them first with: ./scripts/docker-helper.ps1 start"
            exit 1
        }
        Write-Host "Opening shell in app container..." -ForegroundColor Cyan
        docker-compose exec app sh
    }
    "db" {
        if (-not (Test-ContainersRunning)) {
            Write-Error "Containers are not running. Start them first with: ./scripts/docker-helper.ps1 start"
            exit 1
        }
        Write-Host "Opening PostgreSQL shell..." -ForegroundColor Cyan
        docker-compose exec db psql -U postgres -d dklautomationgo
    }
    "backup" {
        if (-not (Test-ContainersRunning)) {
            Write-Error "Containers are not running. Start them first with: ./scripts/docker-helper.ps1 start"
            exit 1
        }
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $backupFile = "backup_$timestamp.sql"
        Write-Host "Creating database backup to $backupFile..." -ForegroundColor Green
        docker-compose exec -T db pg_dump -U postgres dklautomationgo > $backupFile
        Write-Host "Backup created: $backupFile" -ForegroundColor Green
    }
    "restore" {
        if (-not (Test-ContainersRunning)) {
            Write-Error "Containers are not running. Start them first with: ./scripts/docker-helper.ps1 start"
            exit 1
        }
        if (-not $Arg) {
            Write-Error "Please specify a backup file to restore. Usage: ./scripts/docker-helper.ps1 restore backup_file.sql"
            exit 1
        }
        if (-not (Test-Path $Arg)) {
            Write-Error "Backup file not found: $Arg"
            exit 1
        }
        Write-Host "Restoring database from $Arg..." -ForegroundColor Yellow
        Get-Content $Arg | docker-compose exec -T db psql -U postgres -d dklautomationgo
        Write-Host "Database restored from $Arg" -ForegroundColor Green
    }
} 