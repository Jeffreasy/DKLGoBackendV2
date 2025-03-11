#!/bin/bash
set -e

if ! pgrep -f "/app/dklautomationgo" > /dev/null; then
    echo "Application is not running"
    exit 1
fi

# Check if the application is responding
if ! wget -q --spider http://localhost:${PORT}/health; then
    echo "Health check failed: application is not responding"
    exit 1
fi

echo "Health check passed"
exit 0
