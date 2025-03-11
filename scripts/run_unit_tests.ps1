#!/usr/bin/env pwsh
# Dit script voert alleen de unit tests uit die geen database verbinding nodig hebben

Write-Host "Uitvoeren van unit tests voor services..."
go test ./services -v

Write-Host "Uitvoeren van unit tests voor handlers..."
go test ./handlers -v

Write-Host "Uitvoeren van unit tests voor auth/service..."
go test ./auth/service -v

Write-Host "Uitvoeren van unit tests voor auth/handlers..."
go test ./auth/handlers -v

Write-Host "Unit tests voltooid." 