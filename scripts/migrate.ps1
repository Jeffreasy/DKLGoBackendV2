# scripts/migrate.ps1

Write-Host "Running database migrations..." -ForegroundColor Cyan

# Zorg ervoor dat je in de juiste directory bent
Set-Location $PSScriptRoot\..

# Controleer of de .env bestand bestaat en laad het
$envPath = Join-Path -Path (Get-Location) -ChildPath ".env"
if (Test-Path $envPath) {
    Write-Host "Loading database settings from .env file..." -ForegroundColor Cyan
    Get-Content $envPath | ForEach-Object {
        if ($_ -match "^\s*([^=]+)=(.*)$") {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
        }
    }
}

# Controleer of de database omgevingsvariabelen zijn ingesteld
if (-not $env:DB_HOST -or -not $env:DB_PORT -or -not $env:DB_USER -or -not $env:DB_PASSWORD -or -not $env:DB_NAME) {
    Write-Host "Database environment variables not set. Using default values..." -ForegroundColor Yellow
    $env:DB_HOST = "localhost"
    $env:DB_PORT = "5432"
    $env:DB_USER = "postgres"
    $env:DB_PASSWORD = "postgres"
    $env:DB_NAME = "dklautomationgo"
    $env:DB_SSLMODE = "disable"
}

# Controleer of migrate.exe bestaat
if (-not (Test-Path ".\migrate.exe")) {
    Write-Host "migrate.exe not found. Installing golang-migrate..." -ForegroundColor Yellow
    
    # Probeer golang-migrate te installeren
    try {
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
        
        # Zoek het pad naar de migrate executable
        $migratePath = "$env:GOPATH\bin\migrate.exe"
        if (-not (Test-Path $migratePath)) {
            $migratePath = "$env:USERPROFILE\go\bin\migrate.exe"
        }
        
        if (Test-Path $migratePath) {
            Copy-Item $migratePath .
            Write-Host "Copied migrate.exe from $migratePath" -ForegroundColor Green
        } else {
            Write-Host "Could not find migrate.exe. Please install it manually." -ForegroundColor Red
            Write-Host "Run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest" -ForegroundColor Yellow
            exit 1
        }
    } catch {
        Write-Host "Error installing golang-migrate: $_" -ForegroundColor Red
        Write-Host "Please install golang-migrate manually:" -ForegroundColor Yellow
        Write-Host "1. Run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest" -ForegroundColor Yellow
        Write-Host "2. Copy migrate.exe from %USERPROFILE%\go\bin to this directory" -ForegroundColor Yellow
        exit 1
    }
}

# Maak de database URL
$dbUrl = "postgres://$($env:DB_USER):$($env:DB_PASSWORD)@$($env:DB_HOST):$($env:DB_PORT)/$($env:DB_NAME)?sslmode=$($env:DB_SSLMODE)"

# Controleer of de database URL correct is
if ($dbUrl -eq "postgres://:@:/?sslmode=") {
    Write-Host "Invalid database URL. Please check your environment variables." -ForegroundColor Red
    exit 1
}

# Toon de database URL (met verborgen wachtwoord)
$displayUrl = "postgres://$($env:DB_USER):***@$($env:DB_HOST):$($env:DB_PORT)/$($env:DB_NAME)?sslmode=$($env:DB_SSLMODE)"
Write-Host "Using database URL: $displayUrl" -ForegroundColor Cyan

# Voer de migraties uit
Write-Host "Running migrations..." -ForegroundColor Cyan
try {
    # Gebruik eerst de eenvoudige relatieve pad methode (werkt meestal het beste)
    $migrateCommand = ".\migrate.exe -path database/migrations -database $dbUrl up"
    Write-Host "Executing: $migrateCommand" -ForegroundColor Gray
    
    $result = Invoke-Expression $migrateCommand
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Migrations completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "Simple migration method failed with exit code $LASTEXITCODE" -ForegroundColor Red
        Write-Host $result -ForegroundColor Red
        
        # Probeer een alternatieve methode met absoluut pad
        Write-Host "Trying alternative migration method with absolute path..." -ForegroundColor Yellow
        
        # Zorg voor het juiste pad naar de migraties (gebruik forward slashes)
        $migrationsPath = (Resolve-Path ".\database\migrations").Path -replace '\\', '/'
        
        # Zorg ervoor dat we niet dubbel 'file://' hebben
        $migrateCommand = ".\migrate.exe -path file://$migrationsPath -database $dbUrl up"
        Write-Host "Executing: $migrateCommand" -ForegroundColor Gray
        
        $altResult = Invoke-Expression $migrateCommand
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Migrations completed successfully with alternative method!" -ForegroundColor Green
        } else {
            Write-Host "Alternative migration method failed with exit code $LASTEXITCODE" -ForegroundColor Red
            Write-Host $altResult -ForegroundColor Red
            
            # Laatste poging met een andere pad notatie
            Write-Host "Trying final migration method..." -ForegroundColor Yellow
            $finalMigrateCommand = ".\migrate.exe -path ./database/migrations -database $dbUrl up"
            Write-Host "Executing: $finalMigrateCommand" -ForegroundColor Gray
            
            $finalResult = Invoke-Expression $finalMigrateCommand
            
            if ($LASTEXITCODE -eq 0) {
                Write-Host "Migrations completed successfully with final method!" -ForegroundColor Green
            } else {
                Write-Host "All migration methods failed. Please check your migration files and database connection." -ForegroundColor Red
                exit 1
            }
        }
    }
} catch {
    Write-Host "Error running migrations: $_" -ForegroundColor Red
    exit 1
}

# Controleer of de migraties succesvol zijn uitgevoerd door te kijken of de tabellen bestaan
Write-Host "Verifying migrations..." -ForegroundColor Cyan
try {
    # Stel de PGPASSWORD omgevingsvariabele in voor psql
    $env:PGPASSWORD = $env:DB_PASSWORD
    
    # Zoek psql
    $psqlPath = $null
    $possiblePaths = @(
        "psql",
        "psql.exe",
        "C:\Program Files\PostgreSQL\*\bin\psql.exe",
        "${env:ProgramFiles}\PostgreSQL\*\bin\psql.exe"
    )

    foreach ($path in $possiblePaths) {
        if ($path -match "\*") {
            # Path contains wildcard, use Get-ChildItem
            $foundPaths = Get-ChildItem -Path $path -ErrorAction SilentlyContinue
            if ($foundPaths) {
                $psqlPath = $foundPaths[0].FullName
                break
            }
        } else {
            # Try to find command in PATH
            try {
                $command = Get-Command $path -ErrorAction SilentlyContinue
                if ($command) {
                    $psqlPath = $command.Source
                    break
                }
            } catch {
                # Command not found, continue to next path
            }
        }
    }
    
    if ($psqlPath) {
        # Controleer of de tabellen bestaan
        $tablesExist = & $psqlPath -h $env:DB_HOST -p $env:DB_PORT -U $env:DB_USER -d $env:DB_NAME -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'contact_formulieren') AND EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'aanmeldingen');" -t 2>&1
        
        if ($tablesExist -match "t") {
            Write-Host "Migration verification successful: Tables exist in the database." -ForegroundColor Green
        } else {
            Write-Host "Migration verification warning: Tables may not exist in the database." -ForegroundColor Yellow
        }
    }
} catch {
    Write-Host "Error verifying migrations: $_" -ForegroundColor Yellow
    Write-Host "Please verify manually that the migrations were successful." -ForegroundColor Yellow
}

Write-Host "Migration process completed. You can now start the application with 'go run main.go'" -ForegroundColor Cyan