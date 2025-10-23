#!/usr/bin/env bash
# Phoenix TUI - Benchmark Runner
# Runs all critical benchmarks and saves results for comparison

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
RESULTS_DIR="$PROJECT_ROOT/benchmarks/results"
CURRENT_DIR="$RESULTS_DIR/current"

echo "==> Running Phoenix TUI Benchmarks"
echo ""

# Ensure results directory exists
mkdir -p "$CURRENT_DIR"

# Save metadata
echo "# Phoenix TUI Benchmark Results" > "$CURRENT_DIR/README.md"
echo "" >> "$CURRENT_DIR/README.md"
echo "**Date**: $(date '+%Y-%m-%d %H:%M:%S')" >> "$CURRENT_DIR/README.md"
echo "**Branch**: $(git rev-parse --abbrev-ref HEAD)" >> "$CURRENT_DIR/README.md"
echo "**Commit**: $(git rev-parse --short HEAD)" >> "$CURRENT_DIR/README.md"
echo "" >> "$CURRENT_DIR/README.md"

# Run render benchmarks
echo "==> Running render benchmarks..."
cd "$PROJECT_ROOT/render"
export GOWORK=off
go test -bench=. -benchmem -benchtime=2s ./benchmarks > "$CURRENT_DIR/render.txt" 2>&1
echo "✓ Render benchmarks saved to current/render.txt"

# Run core Unicode benchmarks
echo "==> Running core Unicode benchmarks..."
cd "$PROJECT_ROOT/core"
export GOWORK=off
go test -bench=. -benchmem -benchtime=2s ./domain/service > "$CURRENT_DIR/core-unicode.txt" 2>&1
echo "✓ Core benchmarks saved to current/core-unicode.txt"

echo ""
echo "==> Benchmark Summary"
echo ""

# Extract key metrics from render
echo "Render Performance:"
grep "BenchmarkFullScreen_60FPS" "$CURRENT_DIR/render.txt" | tail -1
echo ""

# Extract key metrics from core
echo "Unicode Performance:"
grep "BenchmarkStringWidth_ASCII_Short" "$CURRENT_DIR/core-unicode.txt" | tail -1
grep "BenchmarkStringWidth_Emoji_Short" "$CURRENT_DIR/core-unicode.txt" | tail -1
echo ""

echo "✓ All benchmarks complete!"
echo "Results saved to: $CURRENT_DIR"
echo ""
echo "To compare with baseline:"
echo "  bash benchmarks/scripts/compare.sh"
