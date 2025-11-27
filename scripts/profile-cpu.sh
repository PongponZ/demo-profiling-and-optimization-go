#!/bin/bash

# CPU Profiling Script
# Usage: ./scripts/profile-cpu.sh [bad|optimized] [duration]

set -e

VERSION=${1:-bad}
DURATION=${2:-30s}
OUTPUT_DIR="pprof"
WORKER_BIN=""

mkdir -p "$OUTPUT_DIR"
mkdir -p bin

if [ "$VERSION" = "bad" ]; then
    WORKER_DIR="./cmd/worker-bad"
    WORKER_BIN="./bin/worker-bad"
    OUTPUT_FILE="$OUTPUT_DIR/cpu-bad.prof"
elif [ "$VERSION" = "optimized" ]; then
    WORKER_DIR="./cmd/worker-optimized"
    WORKER_BIN="./bin/worker-optimized"
    OUTPUT_FILE="$OUTPUT_DIR/cpu-optimized.prof"
else
    echo "Usage: $0 [bad|optimized] [duration]"
    exit 1
fi

echo "Building $VERSION worker..."
go build -o "$WORKER_BIN" "$WORKER_DIR" || exit 1

echo "Starting CPU profiling for $VERSION worker..."
echo "Duration: $DURATION"
echo "Output: $OUTPUT_FILE"
echo ""

# Run with CPU profiling
"$WORKER_BIN" \
    -tasks=5000 \
    -workers=20 \
    -data-size=200 \
    -duration="$DURATION" \
    -cpuprofile="$OUTPUT_FILE" 2>&1 | tee "$OUTPUT_DIR/cpu-${VERSION}-output.log" || true

echo ""
echo "CPU profile saved to: $OUTPUT_FILE"
echo ""
echo "To view the profile, run:"
echo "  go tool pprof $OUTPUT_FILE"
echo ""
echo "Or open in web UI:"
echo "  go tool pprof -http=:8080 $OUTPUT_FILE"

