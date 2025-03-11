# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o dklautomationgo .

# Copy migrate binary
RUN cp $(go env GOPATH)/bin/migrate .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata bash postgresql-client

# Copy binary and other necessary files from builder
COPY --from=builder /app/dklautomationgo /app/
COPY --from=builder /app/migrate /app/
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/database/migrations /app/database/migrations

# Create scripts directory
RUN mkdir -p /app/scripts

# Create migrate script
RUN echo '#!/bin/bash\n\
set -e\n\
\n\
echo "Waiting for database to be ready..."\n\
max_retries=30\n\
counter=0\n\
until pg_isready -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER}; do\n\
    >&2 echo "Postgres is unavailable - sleeping"\n\
    counter=$((counter+1))\n\
    if [ $counter -eq $max_retries ]; then\n\
        echo "Failed to connect to database after $max_retries attempts. Exiting."\n\
        exit 1\n\
    fi\n\
    sleep 2\n\
done\n\
\n\
echo "Database is ready. Running migrations..."\n\
\n\
# Run migrations\n\
/app/migrate -path=/app/database/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up\n\
\n\
echo "Migrations completed successfully."\n\
' > /app/scripts/migrate.sh

# Create healthcheck script
RUN echo '#!/bin/bash\n\
set -e\n\
\n\
if ! pgrep -f "/app/dklautomationgo" > /dev/null; then\n\
    echo "Application is not running"\n\
    exit 1\n\
fi\n\
\n\
# Check if the application is responding\n\
if ! wget -q --spider http://localhost:${PORT}/health; then\n\
    echo "Health check failed: application is not responding"\n\
    exit 1\n\
fi\n\
\n\
echo "Health check passed"\n\
exit 0\n\
' > /app/scripts/healthcheck.sh

# Make scripts executable
RUN chmod +x /app/scripts/migrate.sh /app/scripts/healthcheck.sh

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 CMD ["/app/scripts/healthcheck.sh"]

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/dklautomationgo"] 