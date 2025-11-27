#!/bin/bash

# Benchmark Script
# Compares performance between bad and optimized workers

set -e

OUTPUT_DIR="pprof"
BENCHMARK_FILE="$OUTPUT_DIR/benchmark.txt"

mkdir -p "$OUTPUT_DIR"
mkdir -p bin

echo "=== Benchmark Comparison ===" > "$BENCHMARK_FILE"
echo "Date: $(date)" >> "$BENCHMARK_FILE"
echo "" >> "$BENCHMARK_FILE"

# Build both versions
echo "Building workers..."
go build -o "./bin/worker-bad" "./cmd/worker-bad"
go build -o "./bin/worker-optimized" "./cmd/worker-optimized"

# Test parameters
TASKS=10000
WORKERS=20
DATA_SIZE=200
DURATION=10s

echo "Running benchmarks..."
echo "Parameters: Tasks=$TASKS, Workers=$WORKERS, DataSize=$DATA_SIZE, Duration=$DURATION"
echo ""

# Benchmark bad worker
echo "=== Bad Worker ===" | tee -a "$BENCHMARK_FILE"
echo "Running bad worker..." | tee -a "$BENCHMARK_FILE"
time ./bin/worker-bad \
    -tasks=$TASKS \
    -workers=$WORKERS \
    -data-size=$DATA_SIZE \
    -duration=$DURATION \
    -cpuprofile="$OUTPUT_DIR/bench-cpu-bad.prof" \
    -memprofile="$OUTPUT_DIR/bench-mem-bad.prof" 2>&1 | tee -a "$BENCHMARK_FILE" || true

echo "" >> "$BENCHMARK_FILE"
echo "---" >> "$BENCHMARK_FILE"
echo "" >> "$BENCHMARK_FILE"

# Wait a bit between runs
sleep 2

# Benchmark optimized worker
echo "=== Optimized Worker ===" | tee -a "$BENCHMARK_FILE"
echo "Running optimized worker..." | tee -a "$BENCHMARK_FILE"
time ./bin/worker-optimized \
    -tasks=$TASKS \
    -workers=$WORKERS \
    -data-size=$DATA_SIZE \
    -duration=$DURATION \
    -cpuprofile="$OUTPUT_DIR/bench-cpu-optimized.prof" \
    -memprofile="$OUTPUT_DIR/bench-mem-optimized.prof" 2>&1 | tee -a "$BENCHMARK_FILE" || true

echo ""
echo "Benchmark results saved to: $BENCHMARK_FILE"
echo ""
echo "To compare profiles:"
echo "  go tool pprof -http=:8080 $OUTPUT_DIR/bench-cpu-bad.prof"
echo "  go tool pprof -http=:8081 $OUTPUT_DIR/bench-cpu-optimized.prof"

