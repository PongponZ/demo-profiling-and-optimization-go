# Demo Profiling and Optimization in Go

โปรเจกต์ demo สำหรับการเรียนรู้เทคนิคการ profiling และ optimization ใน Go พร้อมตัวอย่างการใช้งาน pprof, benchmarking, และ monitoring ด้วย Prometheus + Grafana

## 📋 Table of Contents

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Services](#services)
- [Makefile Commands](#makefile-commands)
- [Profiling](#profiling)
- [Benchmarking](#benchmarking)
- [Monitoring](#monitoring)
- [API Endpoints](#api-endpoints)

## 🏗️ Architecture

```
┌─────────────┐      ┌──────────────┐      ┌──────────────┐
│ basic-setup │─────>│   RabbitMQ   │─────>│ super-worker │
│  (Web API)  │      │  (Queue)     │      │  (Consumer)  │
└─────────────┘      └──────────────┘      └──────────────┘
      │                                            │
      │                                            │
      v                                            v
┌──────────────────────────────────────────────────────────┐
│              Prometheus (Metrics Collector)              │
└──────────────────────────────────────────────────────────┘
      │
      v
┌──────────────────────────────────────────────────────────┐
│              Grafana (Monitoring Dashboard)              │
└──────────────────────────────────────────────────────────┘
```

### Components:
- **basic-setup**: Web server ที่รับ HTTP requests และส่ง jobs ไปยัง RabbitMQ
- **super-worker**: Worker service ที่ consume jobs จาก RabbitMQ
- **RabbitMQ**: Message broker สำหรับ job queue
- **Prometheus**: Metrics collection และ storage
- **Grafana**: Dashboard สำหรับแสดงผล metrics
- **Node Exporter**: System-level metrics collector

## 📦 Prerequisites

- Docker & Docker Compose
- Go 1.21+ (สำหรับ local development)
- Make (optional but recommended)

## 🚀 Quick Start

### ใช้งานด้วย Make (แนะนำ)

```bash
# Build และ start services
make build
make up

# หรือ build และ start พร้อมกัน
make build && make up
```

### ใช้งานด้วย Docker Compose

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# View logs
docker-compose logs -f
```

## 🌐 Services

เมื่อ start services แล้ว จะสามารถเข้าถึงได้ที่:

| Service | URL | Credentials |
|---------|-----|-------------|
| Web Server | http://localhost:3010 | - |
| Grafana | http://localhost:3000 | admin/admin |
| Prometheus | http://localhost:9090 | - |
| RabbitMQ UI | http://localhost:15672 | guest/guest |
| pprof (basic-setup) | http://localhost:6060/debug/pprof/ | - |
| pprof (super-worker) | http://localhost:6061/debug/pprof/ | - |
| Metrics (basic-setup) | http://localhost:2112/metrics | - |
| Metrics (super-worker) | http://localhost:2113/metrics | - |

## 📖 Makefile Commands

### Docker Management

```bash
make build          # Build Docker images
make up             # Start all services
make down           # Stop all services
make restart        # Restart all services
make logs           # Show logs from all services
make logs-basic     # Show logs from basic-setup only
make logs-worker    # Show logs from super-worker only
make clean          # Remove containers, volumes, and images
```

### Testing & Benchmarking

```bash
make test           # Run Go tests
make benchmark      # Run benchmarks
make benchmark-cpu  # Run benchmarks with CPU profiling
make benchmark-mem  # Run benchmarks with memory profiling
```

### Profiling

```bash
make profile-cpu        # Generate CPU profile (30s)
make profile-mem        # Generate memory profile
make profile-allocs     # Generate allocations profile
make profile-goroutine  # Generate goroutine profile
make profile-trace      # Generate execution trace (5s)

# View profiles
make view-cpu       # View CPU profile in browser
make view-mem       # View memory profile in browser
make view-trace     # View execution trace
```

### Application Testing

```bash
make publish        # Publish 100 jobs to RabbitMQ
make publish-1000   # Publish 1000 jobs
make publish-10000  # Publish 10000 jobs

make test-goleak    # Test goroutine leak endpoint
make test-block     # Test mutex blocking endpoint
make test-alloc     # Test heavy allocation endpoint
make test-cpu       # Test CPU intensive endpoint
```

### Open Dashboards

```bash
make grafana        # Open Grafana dashboard
make prometheus     # Open Prometheus UI
make rabbitmq       # Open RabbitMQ management UI
```

## 🔍 Profiling

### pprof Web Interface

```bash
# CPU profiling
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profiling
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap

# Goroutine profiling
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/goroutine

# Trace profiling
curl http://localhost:6060/debug/pprof/trace?seconds=5 -o trace.out
go tool trace trace.out
```

### Available pprof Endpoints

- `/debug/pprof/` - Index page
- `/debug/pprof/profile` - CPU profile
- `/debug/pprof/heap` - Memory profile
- `/debug/pprof/goroutine` - Goroutine profile
- `/debug/pprof/allocs` - Allocation profile
- `/debug/pprof/block` - Block profile
- `/debug/pprof/mutex` - Mutex profile
- `/debug/pprof/trace` - Execution trace

## 📊 Benchmarking

### Run Benchmarks

```bash
# Run all benchmarks
cd basic-setup/benchmark
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkStringConcat -benchmem

# With CPU profiling
go test -bench=. -benchmem -cpuprofile=cpu.prof

# With memory profiling
go test -bench=. -benchmem -memprofile=mem.prof
```

### Available Benchmarks

- **String Concatenation**: `string_concat_test.go`
  - Plus operator
  - fmt.Sprintf
  - strings.Builder
  - bytes.Buffer

- **Slice Capacity**: `capacity_test.go`
  - With/without pre-allocation
  
- **Type Conversion**: `strconv_test.go`
  - fmt.Sprintf vs strconv.Itoa

## 📈 Monitoring

### Grafana Dashboard

1. เข้า Grafana: http://localhost:3000
2. Login ด้วย `admin/admin`
3. ไปที่ **Dashboards** → **Go Apps Monitoring**

Dashboard จะแสดง:
- **CPU Usage**: CPU usage percentage ของแต่ละ service
- **Memory Usage**: Allocated memory
- **Goroutines Count**: จำนวน goroutines ที่กำลังทำงาน
- **Heap Memory**: Heap memory in use
- **GC Duration**: Garbage collection duration
- **CPU Gauges**: Real-time CPU usage

### Prometheus Metrics

ตัวอย่าง metrics ที่เก็บ:

```promql
# CPU Usage
rate(process_cpu_seconds_total{job="basic-setup"}[1m])

# Memory Usage
go_memstats_alloc_bytes{job="basic-setup"}

# Goroutines
go_goroutines{job="basic-setup"}

# GC Duration
rate(go_gc_duration_seconds_sum{job="basic-setup"}[1m])
```

### Custom Metrics

สามารถเพิ่ม custom metrics ได้โดยใช้ Prometheus client library:

```go
import "github.com/prometheus/client_golang/prometheus"

var requestCount = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests",
    },
)

func init() {
    prometheus.MustRegister(requestCount)
}
```

## 🌐 API Endpoints

### basic-setup Service

```bash
# Health check
GET http://localhost:3010/

# Publish jobs to RabbitMQ
GET http://localhost:3010/publish/:number

# Test goroutine leak
GET http://localhost:3010/goleak

# Test mutex blocking
GET http://localhost:3010/block

# Test heavy allocations
GET http://localhost:3010/alloc

# Test CPU intensive operation
GET http://localhost:3010/cpu
```

### Examples

```bash
# Publish 100 jobs
curl http://localhost:3010/publish/100

# Test endpoints
curl http://localhost:3010/goleak
curl http://localhost:3010/block
curl http://localhost:3010/alloc
curl http://localhost:3010/cpu
```

## 📁 Project Structure

```
.
├── basic-setup/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   └── handler/
│   └── benchmark/
│       ├── string_concat_test.go
│       ├── capacity_test.go
│       └── strconv_test.go
├── super-worker/
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── controller/
│       ├── usecase/
│       ├── repo/
│       └── entity/
├── libs/
│   └── rabbitmq.go
├── grafana/
│   └── provisioning/
│       ├── datasources/
│       │   └── prometheus.yml
│       └── dashboards/
│           ├── dashboard.yml
│           └── go-apps-monitoring.json
├── docker-compose.yml
├── Dockerfile
├── prometheus.yml
├── Makefile
└── README.md
```

## 🛠️ Development

### Run Locally (Without Docker)

```bash
# Start RabbitMQ only
docker-compose up rabbitmq -d

# Run basic-setup
make dev
# หรือ
go run ./basic-setup/cmd/main.go

# Run super-worker (terminal อื่น)
make worker-dev
# หรือ
go run ./super-worker/cmd/main.go
```

### Install Dependencies

```bash
make deps
# หรือ
go mod download
go mod tidy
```

## 🐛 Troubleshooting

### Services ไม่ start

```bash
# ดู logs
make logs

# หรือดู logs แยกตาม service
make logs-basic
make logs-worker
```

### RabbitMQ connection error

รอให้ RabbitMQ health check ผ่านก่อน (ประมาณ 10-15 วินาที)

```bash
# Check RabbitMQ status
docker-compose ps rabbitmq
```

### Port conflicts

ตรวจสอบว่า ports ไม่ซ้ำกับ services อื่นที่กำลังรันอยู่:
- 3010 (web server)
- 3000 (Grafana)
- 5672, 15672 (RabbitMQ)
- 6060, 6061 (pprof)
- 9090 (Prometheus)
- 2112, 2113 (metrics)

## 📚 Resources

- [Go pprof Documentation](https://pkg.go.dev/net/http/pprof)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Grafana Documentation](https://grafana.com/docs/)
- [Go Benchmarking](https://pkg.go.dev/testing#hdr-Benchmarks)

## 📝 License

MIT License

## 👥 Contributing

PRs welcome! Feel free to contribute to this demo project.
