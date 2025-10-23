#!/usr/bin/env bash
# Phoenix TUI - Benchmark Comparison
# Compares current benchmarks with baseline using benchstat

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
RESULTS_DIR="$PROJECT_ROOT/benchmarks/results"
CURRENT_DIR="$RESULTS_DIR/current"
BASELINE_DIR="$RESULTS_DIR/baseline"

# Check if benchstat is installed
if ! command -v benchstat &> /dev/null; then
    echo "Error: benchstat not installed"
    echo "Install with: go install golang.org/x/perf/cmd/benchstat@latest"
    exit 1
fi

# Check if baseline exists
if [ ! -d "$BASELINE_DIR" ] || [ ! -f "$BASELINE_DIR/render.txt" ]; then
    echo "Error: No baseline found"
    echo "Run this first to set baseline:"
    echo "  bash benchmarks/scripts/set_baseline.sh"
    exit 1
fi

# Check if current results exist
if [ ! -d "$CURRENT_DIR" ] || [ ! -f "$CURRENT_DIR/render.txt" ]; then
    echo "Error: No current results found"
    echo "Run benchmarks first:"
    echo "  bash benchmarks/scripts/run_benchmarks.sh"
    exit 1
fi

echo "==> Phoenix TUI Benchmark Comparison"
echo ""

# Show baseline version
if [ -f "$BASELINE_DIR/version.txt" ]; then
    echo "Baseline: $(cat "$BASELINE_DIR/version.txt")"
else
    echo "Baseline: Unknown version"
fi

echo "Current:  $(git rev-parse --short HEAD) ($(git rev-parse --abbrev-ref HEAD))"
echo ""

# Compare render benchmarks
echo "==> Render Performance Comparison"
echo ""
benchstat "$BASELINE_DIR/render.txt" "$CURRENT_DIR/render.txt" | head -30
echo ""

# Compare core benchmarks
echo "==> Unicode Performance Comparison"
echo ""
benchstat "$BASELINE_DIR/core-unicode.txt" "$CURRENT_DIR/core-unicode.txt" | head -30
echo ""

echo "==> Interpretation Guide"
echo ""
echo "  +X%  : X% slower (regression)"
echo "  -X%  : X% faster (improvement)"
echo "  ~    : No significant change"
echo ""
echo "Performance targets:"
echo "  • Render: >24,000 FPS (60 FPS target = 400x faster)"
echo "  • Memory: <100 B/op on hot paths"
echo "  • Allocs: 0 allocs/op on critical paths"
