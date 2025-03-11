#!/bin/bash
# Generate self-signed SSL certificates for development

set -e

# Create directories if they don't exist
mkdir -p nginx/ssl

# Generate a private key
openssl genrsa -out nginx/ssl/server.key 2048

# Generate a CSR (Certificate Signing Request)
openssl req -new -key nginx/ssl/server.key -out nginx/ssl/server.csr -subj "/C=NL/ST=Noord-Holland/L=Amsterdam/O=De Koninklijke Loop/OU=IT/CN=localhost"

# Generate a self-signed certificate
openssl x509 -req -days 365 -in nginx/ssl/server.csr -signkey nginx/ssl/server.key -out nginx/ssl/server.crt

# Set permissions
chmod 600 nginx/ssl/server.key
chmod 600 nginx/ssl/server.crt

echo "Self-signed SSL certificates generated successfully!"
echo "Location: nginx/ssl/"
echo "Files: server.key, server.crt" 