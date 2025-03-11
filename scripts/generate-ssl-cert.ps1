#!/usr/bin/env pwsh
# Generate self-signed SSL certificates for development

$ErrorActionPreference = "Stop"

# Check if OpenSSL is installed
try {
    openssl version | Out-Null
} catch {
    Write-Error "OpenSSL is not installed or not in PATH. Please install OpenSSL."
    exit 1
}

# Create directories if they don't exist
New-Item -Path "nginx/ssl" -ItemType Directory -Force | Out-Null

# Generate a private key
Write-Host "Generating private key..." -ForegroundColor Green
openssl genrsa -out nginx/ssl/server.key 2048

# Generate a CSR (Certificate Signing Request)
Write-Host "Generating CSR..." -ForegroundColor Green
openssl req -new -key nginx/ssl/server.key -out nginx/ssl/server.csr -subj "/C=NL/ST=Noord-Holland/L=Amsterdam/O=De Koninklijke Loop/OU=IT/CN=localhost"

# Generate a self-signed certificate
Write-Host "Generating self-signed certificate..." -ForegroundColor Green
openssl x509 -req -days 365 -in nginx/ssl/server.csr -signkey nginx/ssl/server.key -out nginx/ssl/server.crt

Write-Host "Self-signed SSL certificates generated successfully!" -ForegroundColor Green
Write-Host "Location: nginx/ssl/" -ForegroundColor Cyan
Write-Host "Files: server.key, server.crt" -ForegroundColor Cyan 