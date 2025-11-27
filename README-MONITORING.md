# Monitoring Setup with Grafana

This guide explains how to set up and use the monitoring stack for the Go Profiling Worker project.

## Overview

The monitoring stack consists of:
- **Prometheus**: Collects metrics from both workers
- **Grafana**: Visualizes metrics in dashboards
- **Worker Services**: Bad and optimized worker versions with metrics endpoints

## Architecture

```
┌─────────────────┐      ┌─────────────────┐
│   Bad Worker    │      │ Optimized Worker│
│   :9090/metrics │      │   :9091/metrics │
└────────┬────────┘      └────────┬────────┘
         │                        │
         │    Scrape Metrics      │
         └────────┬───────────────┘
                  │
          ┌───────▼────────┐
          │   Prometheus   │
          │     :9092      │
          └───────┬────────┘
                  │
          ┌───────▼────────┐
          │    Grafana     │
          │     :3000      │
          └────────────────┘
```

## Quick Start

### 1. Start the Monitoring Stack

Using Make (recommended):

```bash
make start
```

Or run the demo script directly:

```bash
./scripts/start-demo.sh
```

This will:
- Build Docker images for both workers
- Start Prometheus, Grafana, and both workers
- Wait for all services to be healthy
- Display access URLs

**Quick commands:**
```bash
make help      # Show all available commands
make start     # Start the stack
make stop      # Stop the stack
make logs      # View logs
make grafana   # Open Grafana in browser
```

### 2. Access the Services

Once started, access the following URLs:

- **Grafana**: http://localhost:3000
  - Username: `admin`
  - Password: `admin`
  
- **Prometheus**: http://localhost:9092

- **Bad Worker Metrics**: http://localhost:9090/metrics

- **Optimized Worker Metrics**: http://localhost:9091/metrics

### 3. View the Dashboard

1. Open Grafana at http://localhost:3000
2. Login with `admin` / `admin`
3. Navigate to **Dashboards** → **Worker Performance Comparison**

The dashboard shows:
- Active Goroutines (detecting leaks)
- Allocated Memory comparison
- Tasks Processing Rate
- Garbage Collection activity
- Task Processing Duration percentiles
- Real-time gauges for key metrics

### 4. Stop the Stack

```bash
make stop
```

Or use the script:
```bash
./scripts/stop-demo.sh
```

To remove all data (volumes):

```bash
make clean
```

## Metrics Exposed

### Runtime Metrics

| Metric | Description | Type |
|--------|-------------|------|
| `worker_active_goroutines` | Number of active goroutines | Gauge |
| `worker_allocated_memory_bytes` | Current allocated memory | Gauge |
| `worker_total_allocations_bytes` | Total memory allocated | Counter |
| `worker_gc_runs_total` | Total garbage collection runs | Counter |

### Custom Metrics

| Metric | Description | Type |
|--------|-------------|------|
| `worker_tasks_processed_total` | Total tasks processed | Counter |
| `worker_tasks_in_queue` | Tasks currently in queue | Gauge |
| `worker_task_processing_duration_seconds` | Task processing duration | Histogram |
| `worker_task_errors_total` | Total task errors | Counter |

All metrics include the label `worker_type` with values:
- `bad`: Problematic worker
- `optimized`: Fixed worker

## Understanding the Dashboard

### 1. Active Goroutines Panel

**What it shows**: Number of goroutines over time

**What to look for**:
- **Bad worker**: Goroutine count increases continuously (leak)
- **Optimized worker**: Goroutine count stays stable

This demonstrates the goroutine leak in the bad worker where helper goroutines never exit.

### 2. Allocated Memory Panel

**What it shows**: Current memory allocation

**What to look for**:
- **Bad worker**: Higher and more volatile memory usage
- **Optimized worker**: Lower and more stable memory usage

This shows the impact of excessive allocations and poor capacity planning.

### 3. Tasks Processing Rate

**What it shows**: Tasks processed per second

**What to look for**:
- Both workers should process tasks at similar rates
- Sustained processing confirms workers are active

### 4. Garbage Collection Rate

**What it shows**: GC runs per second

**What to look for**:
- **Bad worker**: Higher GC frequency due to excessive allocations
- **Optimized worker**: Lower GC frequency with better memory management

### 5. Task Processing Duration

**What it shows**: p50 and p95 latencies

**What to look for**:
- Duration differences between workers
- Latency spikes indicating performance issues

## Manual Profiling with pprof

Both workers expose pprof endpoints for detailed profiling:

### Bad Worker

```bash
# CPU Profile
go tool pprof http://localhost:9090/debug/pprof/profile?seconds=30

# Memory Profile
go tool pprof http://localhost:9090/debug/pprof/heap

# Goroutine Profile (to see leaks)
go tool pprof http://localhost:9090/debug/pprof/goroutine
```

### Optimized Worker

```bash
# CPU Profile
go tool pprof http://localhost:9091/debug/pprof/profile?seconds=30

# Memory Profile
go tool pprof http://localhost:9091/debug/pprof/heap

# Goroutine Profile
go tool pprof http://localhost:9091/debug/pprof/goroutine
```

### Interactive Analysis

```bash
# Start interactive session
go tool pprof http://localhost:9090/debug/pprof/heap

# Inside pprof:
(pprof) top           # Show top memory consumers
(pprof) list <func>   # Show source code for function
(pprof) web           # Generate graph (requires graphviz)
```

## Prometheus Queries

### Useful PromQL Queries

#### Goroutine Leak Detection

```promql
# Goroutine count increase over 1 hour
increase(worker_active_goroutines[1h])
```

#### Memory Growth Rate

```promql
# Memory increase per minute
rate(worker_total_allocations_bytes[1m])
```

#### GC Pressure

```promql
# GC runs per second
rate(worker_gc_runs_total[1s])
```

#### Task Throughput

```promql
# Tasks per second by operation type
rate(worker_tasks_processed_total[1m])
```

#### Performance Comparison

```promql
# Goroutine count difference
worker_active_goroutines{worker_type="bad"} - worker_active_goroutines{worker_type="optimized"}
```

## Troubleshooting

### Services not starting

Check Docker logs:

```bash
# View all logs
docker-compose logs

# View specific service
docker-compose logs worker-bad
docker-compose logs prometheus
docker-compose logs grafana
```

### Grafana dashboard not showing data

1. Check if Prometheus is scraping targets:
   - Visit http://localhost:9092/targets
   - All targets should be "UP"

2. Check if workers are exposing metrics:
   - Visit http://localhost:9090/metrics
   - Visit http://localhost:9091/metrics
   - Should see metrics output

3. Check Grafana datasource:
   - Go to Configuration → Data Sources
   - Test the Prometheus datasource

### Workers crashing or restarting

```bash
# Check worker logs
docker-compose logs -f worker-bad
docker-compose logs -f worker-optimized

# Check resource usage
docker stats
```

### Port conflicts

If ports are already in use, modify `docker-compose.yml`:

```yaml
ports:
  - "3001:3000"  # Change Grafana to 3001
  - "9093:9090"  # Change Prometheus to 9093
```

## Advanced Configuration

### Adjusting Worker Parameters

Edit `docker-compose.yml` to change worker behavior:

```yaml
worker-bad:
  command: ["./worker-bad", "-tasks=100000", "-workers=50", "-duration=48h"]
```

### Retention Period

Edit `prometheus/prometheus.yml`:

```yaml
global:
  scrape_interval: 10s      # More frequent scraping
  
# In docker-compose.yml, add:
command:
  - '--storage.tsdb.retention.time=30d'  # Keep 30 days
```

### Custom Dashboards

1. Create your dashboard in Grafana
2. Export as JSON
3. Save to `grafana/dashboards/`
4. Restart Grafana to load it

## Educational Use

### Demonstrating Goroutine Leaks

1. Watch the "Active Goroutines" panel
2. Observe bad worker's count increasing
3. Run goroutine profile:
   ```bash
   go tool pprof http://localhost:9090/debug/pprof/goroutine
   (pprof) top
   ```

### Showing Memory Issues

1. Compare "Allocated Memory" between workers
2. Check GC activity differences
3. Profile allocations:
   ```bash
   go tool pprof http://localhost:9090/debug/pprof/allocs
   ```

### Performance Impact

1. Compare task processing rates
2. Analyze duration percentiles
3. Calculate resource efficiency

## References

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Go pprof](https://golang.org/pkg/net/http/pprof/)
- [PromQL Tutorial](https://prometheus.io/docs/prometheus/latest/querying/basics/)

