.PHONY: help start stop restart build clean logs ps profile-cpu profile-mem benchmark test grafana prometheus

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

## help: Display this help message
help:
	@echo "$(BLUE)Go Profiling Worker - Available Commands$(NC)"
	@echo ""
	@echo "$(GREEN)Main Commands:$(NC)"
	@echo "  $(YELLOW)make start$(NC)         - Start the monitoring stack (Grafana, Prometheus, Workers)"
	@echo "  $(YELLOW)make stop$(NC)          - Stop all services"
	@echo "  $(YELLOW)make restart$(NC)       - Restart all services"
	@echo "  $(YELLOW)make logs$(NC)          - View logs from all services"
	@echo "  $(YELLOW)make ps$(NC)            - Show running containers"
	@echo ""
	@echo "$(GREEN)Build Commands:$(NC)"
	@echo "  $(YELLOW)make build$(NC)         - Build Docker images"
	@echo "  $(YELLOW)make build-local$(NC)   - Build workers locally"
	@echo ""
	@echo "$(GREEN)Profiling Commands:$(NC)"
	@echo "  $(YELLOW)make profile-cpu$(NC)   - Run CPU profiling on bad worker"
	@echo "  $(YELLOW)make profile-mem$(NC)   - Run memory profiling on bad worker"
	@echo "  $(YELLOW)make benchmark$(NC)     - Run benchmark comparison"
	@echo ""
	@echo "$(GREEN)Monitoring Access:$(NC)"
	@echo "  $(YELLOW)make grafana$(NC)       - Open Grafana dashboard (http://localhost:3000)"
	@echo "  $(YELLOW)make prometheus$(NC)    - Open Prometheus (http://localhost:9092)"
	@echo ""
	@echo "$(GREEN)Cleanup Commands:$(NC)"
	@echo "  $(YELLOW)make clean$(NC)         - Stop services and remove volumes (deletes data)"
	@echo "  $(YELLOW)make clean-build$(NC)   - Remove built binaries"
	@echo ""
	@echo "$(GREEN)Development:$(NC)"
	@echo "  $(YELLOW)make test$(NC)          - Run tests"
	@echo "  $(YELLOW)make tidy$(NC)          - Run go mod tidy"
	@echo ""

## start: Start the monitoring stack
start:
	@echo "$(GREEN)Starting monitoring stack...$(NC)"
	@./scripts/start-demo.sh

## stop: Stop all services
stop:
	@echo "$(YELLOW)Stopping all services...$(NC)"
	@./scripts/stop-demo.sh

## restart: Restart all services
restart: stop start
	@echo "$(GREEN)Services restarted!$(NC)"

## build: Build Docker images
build:
	@echo "$(GREEN)Building Docker images...$(NC)"
	@docker-compose build

## build-local: Build workers locally
build-local:
	@echo "$(GREEN)Building workers locally...$(NC)"
	@cd worker && go build -o ../bin/worker-bad ./cmd/worker-bad
	@cd worker && go build -o ../bin/worker-optimized ./cmd/worker-optimized
	@echo "$(GREEN)Binaries created in bin/$(NC)"

## logs: View logs from all services
logs:
	@docker-compose logs -f

## logs-bad: View bad worker logs
logs-bad:
	@docker-compose logs -f worker-bad

## logs-optimized: View optimized worker logs
logs-optimized:
	@docker-compose logs -f worker-optimized

## logs-prometheus: View Prometheus logs
logs-prometheus:
	@docker-compose logs -f prometheus

## logs-grafana: View Grafana logs
logs-grafana:
	@docker-compose logs -f grafana

## ps: Show running containers
ps:
	@docker-compose ps

## profile-cpu: Run CPU profiling on bad worker
profile-cpu:
	@echo "$(GREEN)Running CPU profiling on bad worker (30s)...$(NC)"
	@./scripts/profile-cpu.sh bad 30s

## profile-cpu-optimized: Run CPU profiling on optimized worker
profile-cpu-optimized:
	@echo "$(GREEN)Running CPU profiling on optimized worker (30s)...$(NC)"
	@./scripts/profile-cpu.sh optimized 30s

## profile-mem: Run memory profiling on bad worker
profile-mem:
	@echo "$(GREEN)Running memory profiling on bad worker (30s)...$(NC)"
	@./scripts/profile-mem.sh bad 30s

## profile-mem-optimized: Run memory profiling on optimized worker
profile-mem-optimized:
	@echo "$(GREEN)Running memory profiling on optimized worker (30s)...$(NC)"
	@./scripts/profile-mem.sh optimized 30s

## benchmark: Run benchmark comparison
benchmark:
	@echo "$(GREEN)Running benchmark...$(NC)"
	@./scripts/benchmark.sh

## grafana: Open Grafana in browser
grafana:
	@echo "$(GREEN)Opening Grafana...$(NC)"
	@open http://localhost:3000 || xdg-open http://localhost:3000 || echo "Please open http://localhost:3000 in your browser"

## prometheus: Open Prometheus in browser
prometheus:
	@echo "$(GREEN)Opening Prometheus...$(NC)"
	@open http://localhost:9092 || xdg-open http://localhost:9092 || echo "Please open http://localhost:9092 in your browser"

## clean: Stop services and remove volumes
clean:
	@echo "$(YELLOW)Stopping services and removing volumes...$(NC)"
	@docker-compose down -v
	@echo "$(GREEN)Cleanup complete!$(NC)"

## clean-build: Remove built binaries
clean-build:
	@echo "$(YELLOW)Removing built binaries...$(NC)"
	@rm -rf bin/
	@rm -f worker-bad worker-optimized
	@echo "$(GREEN)Build artifacts removed!$(NC)"

## test: Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	@cd worker && go test -v ./...

## tidy: Run go mod tidy
tidy:
	@echo "$(GREEN)Running go mod tidy...$(NC)"
	@cd worker && go mod tidy
	@echo "$(GREEN)Dependencies updated!$(NC)"

## up: Alias for start
up: start

## down: Alias for stop
down: stop

## status: Alias for ps
status: ps

