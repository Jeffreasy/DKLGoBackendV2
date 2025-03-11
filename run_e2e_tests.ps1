#!/usr/bin/env pwsh
# Script om E2E tests uit te voeren

Write-Host "🧪 E2E tests starten..." -ForegroundColor Cyan

# Zorg ervoor dat de test omgeving schoon is
Write-Host "🧹 Opruimen van vorige test omgeving..." -ForegroundColor Yellow
docker-compose -f tests/e2e/docker-compose.e2e.yml down -v --remove-orphans

# Start de test omgeving
Write-Host "🚀 Starten van de test omgeving..." -ForegroundColor Yellow
docker-compose -f tests/e2e/docker-compose.e2e.yml up -d

# Wacht tot de services klaar zijn
Write-Host "⏳ Wachten tot services gereed zijn..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Voer de tests uit
Write-Host "🧪 E2E tests uitvoeren..." -ForegroundColor Green
go test -v ./tests/e2e/scenarios/...

# Sla de exit code op
$testResult = $LASTEXITCODE

# Ruim de test omgeving op
Write-Host "🧹 Opruimen van test omgeving..." -ForegroundColor Yellow
docker-compose -f tests/e2e/docker-compose.e2e.yml down -v --remove-orphans

# Geef de juiste exit code terug
if ($testResult -ne 0) {
    Write-Host "❌ E2E tests gefaald!" -ForegroundColor Red
    exit $testResult
} else {
    Write-Host "✅ E2E tests succesvol afgerond!" -ForegroundColor Green
    exit 0
} 