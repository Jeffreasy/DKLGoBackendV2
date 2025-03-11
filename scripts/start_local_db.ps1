# scripts/start_local_db.ps1

Write-Host "Checking PostgreSQL database..." -ForegroundColor Cyan

# Zorg ervoor dat je in de juiste directory bent
Set-Location $PSScriptRoot\..

# Controleer of de database omgevingsvariabelen zijn ingesteld
if (-not $env:DB_HOST -or -not $env:DB_PORT -or -not $env:DB_USER -or -not $env:DB_PASSWORD -or -not $env:DB_NAME) {
    Write-Host "Database environment variables not set. Using default values..." -ForegroundColor Yellow
    $env:DB_HOST = "localhost"
    $env:DB_PORT = "5432"
    $env:DB_USER = "postgres"
    $env:DB_PASSWORD = "Bootje@12"
    $env:DB_NAME = "dklautomationgo"
    $env:DB_SSLMODE = "disable"
}

# Controleer of psql is geïnstalleerd
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

if (-not $psqlPath) {
    Write-Host "psql not found. Please make sure PostgreSQL is installed and in your PATH." -ForegroundColor Red
    Write-Host "You can download PostgreSQL from https://www.postgresql.org/download/windows/" -ForegroundColor Yellow
    
    # Voeg PostgreSQL bin directory toe aan PATH als we het kunnen vinden
    $pgBinPath = "C:\Program Files\PostgreSQL"
    if (Test-Path $pgBinPath) {
        $pgVersions = Get-ChildItem -Path $pgBinPath -Directory
        if ($pgVersions) {
            $latestVersion = $pgVersions | Sort-Object -Property Name -Descending | Select-Object -First 1
            $pgBinPath = Join-Path -Path $latestVersion.FullName -ChildPath "bin"
            if (Test-Path $pgBinPath) {
                Write-Host "Adding PostgreSQL bin directory to PATH: $pgBinPath" -ForegroundColor Yellow
                $env:PATH += ";$pgBinPath"
                $psqlPath = Join-Path -Path $pgBinPath -ChildPath "psql.exe"
                if (Test-Path $psqlPath) {
                    Write-Host "Found psql at: $psqlPath" -ForegroundColor Green
                } else {
                    exit 1
                }
            } else {
                exit 1
            }
        } else {
            exit 1
        }
    } else {
        exit 1
    }
}

Write-Host "Found psql at: $psqlPath" -ForegroundColor Green

# Stel de PGPASSWORD omgevingsvariabele in voor psql
$env:PGPASSWORD = $env:DB_PASSWORD

# Test de verbinding met PostgreSQL
Write-Host "Testing connection to PostgreSQL server..." -ForegroundColor Cyan
$connectionSuccess = $false
$maxAttempts = 3
$attempts = 0

while (-not $connectionSuccess -and $attempts -lt $maxAttempts) {
    try {
        $testConnection = & $psqlPath -h $env:DB_HOST -p $env:DB_PORT -U $env:DB_USER -c "SELECT version();" -d "postgres" 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            $connectionSuccess = $true
            Write-Host "Successfully connected to PostgreSQL server." -ForegroundColor Green
            Write-Host $testConnection -ForegroundColor Gray
        } else {
            $attempts++
            Write-Host "Failed to connect to PostgreSQL server (Attempt $attempts of $maxAttempts):" -ForegroundColor Red
            Write-Host $testConnection -ForegroundColor Red
            
            if ($attempts -lt $maxAttempts) {
                Write-Host "Please enter the correct PostgreSQL password for user '$env:DB_USER':" -ForegroundColor Yellow
                $securePassword = Read-Host -AsSecureString
                $bstr = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($securePassword)
                $env:DB_PASSWORD = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($bstr)
                $env:PGPASSWORD = $env:DB_PASSWORD
            } else {
                Write-Host "Maximum connection attempts reached. Please check your PostgreSQL installation and credentials." -ForegroundColor Red
                exit 1
            }
        }
    } catch {
        $attempts++
        Write-Host "Error connecting to PostgreSQL: $_" -ForegroundColor Red
        
        if ($attempts -lt $maxAttempts) {
            Write-Host "Please enter the correct PostgreSQL password for user '$env:DB_USER':" -ForegroundColor Yellow
            $securePassword = Read-Host -AsSecureString
            $bstr = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($securePassword)
            $env:DB_PASSWORD = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($bstr)
            $env:PGPASSWORD = $env:DB_PASSWORD
        } else {
            Write-Host "Maximum connection attempts reached. Please check your PostgreSQL installation and credentials." -ForegroundColor Red
            exit 1
        }
    }
}

# Controleer of de database bestaat
Write-Host "Checking if database '$env:DB_NAME' exists..." -ForegroundColor Cyan
$dbExists = & $psqlPath -h $env:DB_HOST -p $env:DB_PORT -U $env:DB_USER -c "SELECT 1 FROM pg_database WHERE datname = '$env:DB_NAME';" -d "postgres" -t 2>&1

if ($dbExists -notmatch "1") {
    Write-Host "Database '$env:DB_NAME' does not exist. Creating..." -ForegroundColor Yellow
    
    try {
        $createResult = & $psqlPath -h $env:DB_HOST -p $env:DB_PORT -U $env:DB_USER -c "CREATE DATABASE $env:DB_NAME;" -d "postgres" 2>&1
        
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Failed to create database:" -ForegroundColor Red
            Write-Host $createResult -ForegroundColor Red
            exit 1
        }
        
        Write-Host "Database '$env:DB_NAME' created successfully." -ForegroundColor Green
    } catch {
        Write-Host "Error creating database: $_" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "Database '$env:DB_NAME' already exists." -ForegroundColor Green
}

# Controleer of de uuid-ossp extensie beschikbaar is
Write-Host "Checking if uuid-ossp extension is available..." -ForegroundColor Cyan
$extensionExists = & $psqlPath -h $env:DB_HOST -p $env:DB_PORT -U $env:DB_USER -c "SELECT 1 FROM pg_available_extensions WHERE name = 'uuid-ossp';" -d "$env:DB_NAME" -t 2>&1

if ($extensionExists -notmatch "1") {
    Write-Host "Extension 'uuid-ossp' is not available. Please install the PostgreSQL contrib package." -ForegroundColor Red
    exit 1
} else {
    Write-Host "Extension 'uuid-ossp' is available." -ForegroundColor Green
}

Write-Host "PostgreSQL database check completed successfully!" -ForegroundColor Green
Write-Host "You can now run the migrations with .\scripts\migrate.ps1" -ForegroundColor Cyan

# Toon de verbindingsgegevens
Write-Host "`nDatabase connection details:" -ForegroundColor Cyan
Write-Host "Host: $env:DB_HOST" -ForegroundColor White
Write-Host "Port: $env:DB_PORT" -ForegroundColor White
Write-Host "User: $env:DB_USER" -ForegroundColor White
Write-Host "Password: ********" -ForegroundColor White
Write-Host "Database: $env:DB_NAME" -ForegroundColor White
Write-Host "SSL Mode: $env:DB_SSLMODE" -ForegroundColor White

# Sla de database instellingen op in een .env bestand voor gebruik door andere scripts
$envContent = @"
DB_HOST=$env:DB_HOST
DB_PORT=$env:DB_PORT
DB_USER=$env:DB_USER
DB_PASSWORD=$env:DB_PASSWORD
DB_NAME=$env:DB_NAME
DB_SSLMODE=$env:DB_SSLMODE
"@

$envPath = Join-Path -Path (Get-Location) -ChildPath ".env"
$envContent | Out-File -FilePath $envPath -Encoding utf8 -Force
Write-Host "Database settings saved to .env file" -ForegroundColor Green