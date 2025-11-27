#!/bin/bash

# Stop Demo Script
# Stops all services and cleans up

set -e

echo "======================================"
echo "Stopping Go Profiling Demo"
echo "======================================"
echo ""

# Navigate to project root
cd "$(dirname "$0")/.."

# Stop all services
echo "Stopping services..."
docker-compose down

echo ""
echo "Services stopped successfully!"
echo ""
echo "To remove volumes (data will be lost):"
echo "  docker-compose down -v"
echo ""

