#!/usr/bin/env pwsh
# Dit script maakt een testdatabase aan voor de tests

# Configuratie
$DB_HOST = if ($env:TEST_DB_HOST) { $env:TEST_DB_HOST } else { "localhost" }
$DB_PORT = if ($env:TEST_DB_PORT) { $env:TEST_DB_PORT } else { "5432" }
$DB_USER = if ($env:TEST_DB_USER) { $env:TEST_DB_USER } else { "postgres" }
$DB_PASSWORD = if ($env:TEST_DB_PASSWORD) { $env:TEST_DB_PASSWORD } else { "Bootje@12" }
$DB_NAME = if ($env:TEST_DB_NAME) { $env:TEST_DB_NAME } else { "dklautomationgo_test" }

# Controleer of Docker draait
$dockerRunning = docker ps 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "Docker is niet actief. Start Docker en probeer het opnieuw."
    exit 1
}

# Controleer of de database container draait
$dbContainer = docker ps --filter "name=dklautomationgo-db" --format "{{.Names}}"
if (-not $dbContainer) {
    Write-Host "De database container 'dklautomationgo-db' is niet actief. Start de container en probeer het opnieuw."
    exit 1
}

# Maak de testdatabase aan
Write-Host "Aanmaken van testdatabase '$DB_NAME'..."
$createDbCmd = "docker exec dklautomationgo-db psql -U $DB_USER -c 'CREATE DATABASE $DB_NAME;'"
Invoke-Expression $createDbCmd

if ($LASTEXITCODE -eq 0) {
    Write-Host "Testdatabase '$DB_NAME' is succesvol aangemaakt."
} else {
    Write-Host "Er is een fout opgetreden bij het aanmaken van de testdatabase. Mogelijk bestaat deze al."
}

# Stel de omgevingsvariabelen in voor de tests
$env:TEST_DB_HOST = $DB_HOST
$env:TEST_DB_PORT = $DB_PORT
$env:TEST_DB_USER = $DB_USER
$env:TEST_DB_PASSWORD = $DB_PASSWORD
$env:TEST_DB_NAME = $DB_NAME

Write-Host "Omgevingsvariabelen zijn ingesteld voor de tests."
Write-Host "Je kunt nu de tests uitvoeren met: go test ./..." 