#!/bin/bash

# Start Demo Script
# Starts the entire monitoring stack with Docker Compose

set -e

echo "======================================"
echo "Starting Go Profiling Demo"
echo "======================================"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker first."
    exit 1
fi

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Error: docker-compose is not installed."
    echo "Please install docker-compose: https://docs.docker.com/compose/install/"
    exit 1
fi

# Navigate to project root
cd "$(dirname "$0")/.."

echo "Building and starting services..."
echo ""

# Build and start all services
docker-compose up -d --build

echo ""
echo "Waiting for services to be ready..."
echo ""

# Wait for Prometheus
echo -n "Waiting for Prometheus..."
until curl -s http://localhost:9092/-/healthy > /dev/null 2>&1; do
    echo -n "."
    sleep 2
done
echo " ✓"

# Wait for Grafana
echo -n "Waiting for Grafana..."
until curl -s http://localhost:3000/api/health > /dev/null 2>&1; do
    echo -n "."
    sleep 2
done
echo " ✓"

# Wait for Workers
echo -n "Waiting for Bad Worker..."
until curl -s http://localhost:9090/metrics > /dev/null 2>&1; do
    echo -n "."
    sleep 2
done
echo " ✓"

echo -n "Waiting for Optimized Worker..."
until curl -s http://localhost:9091/metrics > /dev/null 2>&1; do
    echo -n "."
    sleep 2
done
echo " ✓"

echo ""
echo "======================================"
echo "All services are ready!"
echo "======================================"
echo ""
echo "Access URLs:"
echo "  - Grafana:         http://localhost:3000"
echo "    (login: admin / admin)"
echo "  - Prometheus:      http://localhost:9092"
echo "  - Bad Worker:      http://localhost:9090/metrics"
echo "  - Optimized Worker: http://localhost:9091/metrics"
echo ""
echo "To view logs:"
echo "  docker-compose logs -f [service-name]"
echo ""
echo "To stop all services:"
echo "  ./scripts/stop-demo.sh"
echo ""
echo "Press Ctrl+C to view logs, or run 'docker-compose logs -f' to follow all logs"
echo ""

# Show container status
docker-compose ps

