.PHONY: help build up down restart logs clean test benchmark profile-cpu profile-mem profile-trace

# Default target
help:
	@echo "Available targets:"
	@echo "  make build          - Build all Docker images"
	@echo "  make up             - Start all services"
	@echo "  make down           - Stop all services"
	@echo "  make restart        - Restart all services"
	@echo "  make logs           - Show logs from all services"
	@echo "  make logs-basic     - Show logs from basic-setup"
	@echo "  make logs-worker    - Show logs from super-worker"
	@echo "  make clean          - Remove all containers, volumes, and images"
	@echo "  make test           - Run Go tests"
	@echo "  make benchmark      - Run Go benchmarks"
	@echo "  make profile-cpu    - Generate CPU profile"
	@echo "  make profile-mem    - Generate memory profile"
	@echo "  make profile-trace  - Generate execution trace"
	@echo "  make publish        - Publish jobs (default 100)"
	@echo "  make grafana        - Open Grafana dashboard"
	@echo "  make prometheus     - Open Prometheus UI"
	@echo "  make rabbitmq       - Open RabbitMQ management"

# Docker commands
build:
	docker-compose build

up:
	docker-compose up -d
	@echo "Services started!"
	@echo "- Web server:      http://localhost:3010"
	@echo "- Grafana:         http://localhost:3000 (admin/admin)"
	@echo "- Prometheus:      http://localhost:9090"
	@echo "- RabbitMQ UI:     http://localhost:15672 (guest/guest)"
	@echo "- pprof (basic):   http://localhost:6060/debug/pprof/"
	@echo "- pprof (worker):  http://localhost:6061/debug/pprof/"
	@echo "- metrics (basic): http://localhost:2112/metrics"
	@echo "- metrics (worker):http://localhost:2113/metrics"

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

logs-basic:
	docker-compose logs -f basic-setup

logs-worker:
	docker-compose logs -f super-worker

clean:
	docker-compose down -v --rmi all --remove-orphans
	@echo "Cleaned up all containers, volumes, and images"

# Go commands
test:
	go test -v ./...

benchmark:
	go test -bench=. -benchmem ./basic-setup/benchmark/...

benchmark-cpu:
	go test -bench=. -benchmem -cpuprofile=cpu.prof ./basic-setup/benchmark/...

benchmark-mem:
	go test -bench=. -benchmem -memprofile=mem.prof ./basic-setup/benchmark/...

# Profiling commands
profile-cpu:
	@echo "Generating CPU profile for 30 seconds..."
	curl http://localhost:6060/debug/pprof/profile?seconds=30 -o cpu.prof
	@echo "View with: go tool pprof -http=:8080 cpu.prof"

profile-mem:
	@echo "Generating memory profile..."
	curl http://localhost:6060/debug/pprof/heap -o mem.prof
	@echo "View with: go tool pprof -http=:8080 mem.prof"

profile-allocs:
	@echo "Generating allocations profile..."
	curl http://localhost:6060/debug/pprof/allocs -o allocs.prof
	@echo "View with: go tool pprof -http=:8080 allocs.prof"

profile-goroutine:
	@echo "Generating goroutine profile..."
	curl http://localhost:6060/debug/pprof/goroutine -o goroutine.prof
	@echo "View with: go tool pprof -http=:8080 goroutine.prof"

profile-trace:
	@echo "Generating execution trace for 5 seconds..."
	curl http://localhost:6060/debug/pprof/trace?seconds=5 -o trace.out
	@echo "View with: go tool trace trace.out"

# Application commands
publish:
	@echo "Publishing 100 jobs to RabbitMQ..."
	curl http://localhost:3010/publish/100

publish-1000:
	@echo "Publishing 1000 jobs to RabbitMQ..."
	curl http://localhost:3010/publish/1000

publish-10000:
	@echo "Publishing 10000 jobs to RabbitMQ..."
	curl http://localhost:3010/publish/10000

# Test endpoints
test-goleak:
	@echo "Testing goroutine leak endpoint..."
	curl http://localhost:3010/goleak

test-block:
	@echo "Testing blocking endpoint..."
	curl http://localhost:3010/block

test-alloc:
	@echo "Testing allocation endpoint..."
	curl http://localhost:3010/alloc

test-cpu:
	@echo "Testing CPU intensive endpoint..."
	curl http://localhost:3010/cpu

# Open dashboards (requires xdg-open or open command)
grafana:
	@which xdg-open > /dev/null && xdg-open http://localhost:3000 || open http://localhost:3000 || echo "Open http://localhost:3000 in your browser"

prometheus:
	@which xdg-open > /dev/null && xdg-open http://localhost:9090 || open http://localhost:9090 || echo "Open http://localhost:9090 in your browser"

rabbitmq:
	@which xdg-open > /dev/null && xdg-open http://localhost:15672 || open http://localhost:15672 || echo "Open http://localhost:15672 in your browser"

# Development
dev:
	go run ./basic-setup/cmd/main.go

worker-dev:
	go run ./super-worker/cmd/main.go

# Dependencies
deps:
	go mod download
	go mod tidy

# View profiles with pprof web interface
view-cpu:
	go tool pprof -http=:8080 cpu.prof

view-mem:
	go tool pprof -http=:8080 mem.prof

view-trace:
	go tool trace trace.out

