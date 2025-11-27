# Golang Profiling and Optimization Worker Project

A mini project demonstrating common performance issues in Go and how to identify and fix them using profiling tools. This project includes a background worker that processes tasks, with both problematic and optimized implementations.

## Project Overview

This project demonstrates three common performance problems in Go:

1. **Goroutine Leaks** - Goroutines that never exit, causing memory leaks
2. **Slice/Map Capacity Issues** - Not pre-allocating capacity, causing multiple reallocations
3. **Excessive Allocations** - Creating unnecessary objects in hot paths

## Project Structure

```
profiling-and-optimize-golang/
├── Makefile                 # Build automation & shortcuts
├── docker-compose.yml       # Full monitoring stack
├── worker/                  # Worker service
│   ├── cmd/
│   │   ├── worker-bad/      # Problematic version
│   │   └── worker-optimized/ # Optimized version
│   ├── internal/
│   │   ├── task/            # Task definitions
│   │   └── worker/          # Worker implementations
│   ├── Dockerfile.bad
│   ├── Dockerfile.optimized
│   └── go.mod
├── prometheus/              # Prometheus configuration
│   └── prometheus.yml
├── grafana/                 # Grafana dashboards & provisioning
│   ├── dashboards/
│   └── provisioning/
├── scripts/
│   ├── start-demo.sh       # Start monitoring stack
│   ├── stop-demo.sh        # Stop monitoring stack
│   ├── profile-cpu.sh      # CPU profiling script
│   ├── profile-mem.sh      # Memory profiling script
│   └── benchmark.sh        # Benchmark comparison script
├── README.md               # This file
└── README-MONITORING.md    # Monitoring setup guide
```

## Problems Demonstrated

### 1. Goroutine Leak

**Problem (bad.go):**
- Spawns goroutines without context cancellation
- Background monitor goroutine runs forever
- Helper goroutines with no exit mechanism
- Unbuffered channels that can block indefinitely

**Solution (optimized.go):**
- Uses `context.Context` for cancellation
- All goroutines check context and exit properly
- Buffered channels to prevent blocking
- Proper cleanup with `WaitGroup`

### 2. Slice/Map Capacity Issues

**Problem (bad.go):**
- Appends to slices without pre-allocating capacity
- Creates maps without pre-sizing
- Multiple reallocations as data grows

**Solution (optimized.go):**
- Pre-allocates slice capacity: `make([]T, 0, expectedSize)`
- Pre-sizes maps: `make(map[K]V, expectedSize)`
- Reuses slices by resetting length: `slice = slice[:0]`

### 3. Excessive Allocations

**Problem (bad.go):**
- Creates new structs in hot paths
- Uses string concatenation (`+=`) instead of `strings.Builder`
- Allocates intermediate variables unnecessarily

**Solution (optimized.go):**
- Reuses structs by updating fields
- Uses `strings.Builder` with object pooling
- Reuses buffers and objects where possible

## Getting Started

### Prerequisites

- Go 1.19 or later
- Docker and Docker Compose (for monitoring stack)
- Basic understanding of Go profiling tools

### Installation

#### Option 1: Run with Docker (Recommended)

Using Make (easiest):
```bash
# Start the entire monitoring stack
make start

# Or use the script directly
./scripts/start-demo.sh
```

This will start:
- Both worker services (bad and optimized)
- Prometheus (metrics collection)
- Grafana (visualization)

Access Grafana at http://localhost:3000 (admin/admin)

**Available Make commands:**
```bash
make help           # Show all available commands
make start          # Start monitoring stack
make stop           # Stop all services
make logs           # View logs
make grafana        # Open Grafana in browser
make profile-cpu    # Run CPU profiling
make benchmark      # Run benchmark comparison
```

See [README-MONITORING.md](README-MONITORING.md) for detailed instructions.

#### Option 2: Run Locally

```bash
# Navigate to worker directory
cd worker

# Build the workers
go build -o bin/worker-bad ./cmd/worker-bad
go build -o bin/worker-optimized ./cmd/worker-optimized
```

## Usage

### Quick Start with Monitoring

The easiest way to see the differences between bad and optimized workers:

```bash
# Start the monitoring stack
make start

# Open Grafana (or visit http://localhost:3000)
make grafana

# View the "Worker Performance Comparison" dashboard
# Watch the goroutine leaks, memory usage, and performance metrics in real-time

# View logs
make logs

# Stop when done
make stop
```

### Running Workers Locally

#### Bad Worker (Problematic Version)

```bash
cd worker
go run ./cmd/worker-bad \
    -tasks=1000 \
    -workers=10 \
    -data-size=100 \
    -duration=30s
```

#### Optimized Worker (Fixed Version)

```bash
cd worker
go run ./cmd/worker-optimized \
    -tasks=1000 \
    -workers=10 \
    -data-size=100 \
    -duration=30s
```

### Command Line Flags

- `-tasks`: Number of tasks to process (default: 1000)
- `-workers`: Number of worker goroutines (default: 10)
- `-data-size`: Size of data array per task (default: 100)
- `-duration`: How long to run the worker (default: 30s)
- `-cpuprofile`: Path to write CPU profile (optional)
- `-memprofile`: Path to write memory profile (optional)

## Monitoring with Grafana

The project includes a complete monitoring stack with Grafana dashboards:

```bash
# Start the monitoring stack
./scripts/start-demo.sh
```

**What you'll see:**
- Real-time goroutine count (showing the leak in bad worker)
- Memory usage comparison
- Task processing rates
- Garbage collection activity
- Performance metrics

For detailed monitoring documentation, see [README-MONITORING.md](README-MONITORING.md)

## Profiling

### CPU Profiling

#### Using the Script

```bash
# Profile the bad worker
./scripts/profile-cpu.sh bad 30s

# Profile the optimized worker
./scripts/profile-cpu.sh optimized 30s
```

#### Manual CPU Profiling

```bash
cd worker

# Bad worker
go run ./cmd/worker-bad -cpuprofile=../pprof/cpu-bad.prof -duration=30s

# Optimized worker
go run ./cmd/worker-optimized -cpuprofile=../pprof/cpu-optimized.prof -duration=30s
```

#### Viewing CPU Profiles

```bash
# Interactive terminal view
go tool pprof pprof/cpu-bad.prof

# Web UI (opens in browser)
go tool pprof -http=:8080 pprof/cpu-bad.prof
```

### Memory Profiling

#### Using the Script

```bash
# Profile the bad worker
./scripts/profile-mem.sh bad 30s

# Profile the optimized worker
./scripts/profile-mem.sh optimized 30s
```

#### Manual Memory Profiling

```bash
cd worker

# Bad worker
go run ./cmd/worker-bad -memprofile=../pprof/mem-bad.prof -duration=30s

# Optimized worker
go run ./cmd/worker-optimized -memprofile=../pprof/mem-optimized.prof -duration=30s
```

### Live Profiling with pprof

When workers are running (locally or in Docker), access pprof endpoints:

```bash
# Goroutine profile (shows leaks)
go tool pprof http://localhost:9090/debug/pprof/goroutine

# Memory profile
go tool pprof http://localhost:9090/debug/pprof/heap

# CPU profile (30 seconds)
go tool pprof http://localhost:9090/debug/pprof/profile?seconds=30
```

#### Viewing Memory Profiles

```bash
# Interactive terminal view
go tool pprof pprof/mem-bad.prof

# Top memory allocations
go tool pprof -top pprof/mem-bad.prof

# Web UI
go tool pprof -http=:8080 pprof/mem-bad.prof
```

### Benchmark Comparison

Run the benchmark script to compare both versions:

```bash
./scripts/benchmark.sh
```

This will:
- Run both workers with the same parameters
- Generate CPU and memory profiles
- Save benchmark results to `pprof/benchmark.txt`

## Analyzing Profiles

### Common pprof Commands

```bash
# Start interactive session
go tool pprof profile.prof

# Inside pprof:
(pprof) top          # Show top functions by CPU/memory
(pprof) top10        # Show top 10 functions
(pprof) list <func>  # Show annotated source for function
(pprof) web          # Generate SVG and open in browser
(pprof) png          # Generate PNG graph
(pprof) help         # Show all commands
```

### Key Metrics to Look For

1. **Goroutine Leaks:**
   - Check `runtime.NumGoroutine()` in output
   - Look for goroutines that never exit in the profile
   - Compare goroutine counts between bad and optimized versions

2. **Memory Allocations:**
   - Use `go tool pprof -alloc_space` for allocation profiles
   - Look for functions with high allocation counts
   - Compare `TotalAlloc` between versions

3. **CPU Usage:**
   - Identify hot paths in CPU profiles
   - Look for functions spending excessive time
   - Check for unnecessary work in loops

## Example Output

### Bad Worker Output

```
=== Bad Worker (Problematic Version) ===
Tasks: 1000, Workers: 10, Data Size: 100
...
Memory Stats:
  Alloc: 15234 KB
  TotalAlloc: 45678 KB
  Sys: 23456 KB
  NumGC: 45
  NumGoroutine: 23

Final goroutine count: 23 (should be ~1, but will be higher due to leaks)
Note: Some goroutines are still running (demonstrating the leak)
```

### Optimized Worker Output

```
=== Optimized Worker (Fixed Version) ===
Tasks: 1000, Workers: 10, Data Size: 100
...
Memory Stats:
  Alloc: 8234 KB
  TotalAlloc: 23456 KB
  Sys: 12345 KB
  NumGC: 12
  NumGoroutine: 1

Final goroutine count: 1 (should be ~1)
All goroutines properly cleaned up!
```

## Learning Objectives

After working with this project, you should understand:

1. How to identify goroutine leaks using profiling and monitoring
2. The importance of pre-allocating slice and map capacity
3. How to reduce allocations in hot paths
4. How to use `go tool pprof` for performance analysis
5. Best practices for goroutine lifecycle management
6. When and how to use object pooling
7. How to instrument Go applications with Prometheus metrics
8. How to set up monitoring with Grafana dashboards
9. How to analyze performance differences in real-time

## Additional Resources

### Documentation
- [Monitoring Setup Guide](README-MONITORING.md) - Detailed monitoring documentation
- [Go pprof Documentation](https://golang.org/pkg/runtime/pprof/)
- [Go Performance Best Practices](https://golang.org/doc/effective_go#performance)
- [Diagnosing Go Programs](https://golang.org/doc/diagnostics.html)
- [Dave Cheney's Go Performance Tips](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)

### Tools
- [Prometheus](https://prometheus.io/) - Metrics collection
- [Grafana](https://grafana.com/) - Metrics visualization
- [pprof](https://github.com/google/pprof) - Profiling tool

## License

This is an educational project for learning Go profiling and optimization techniques.

