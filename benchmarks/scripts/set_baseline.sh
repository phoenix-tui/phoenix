#!/usr/bin/env bash
# Phoenix TUI - Set Baseline
# Copies current results to baseline for future comparisons

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
RESULTS_DIR="$PROJECT_ROOT/benchmarks/results"
CURRENT_DIR="$RESULTS_DIR/current"
BASELINE_DIR="$RESULTS_DIR/baseline"

# Check if current results exist
if [ ! -d "$CURRENT_DIR" ] || [ ! -f "$CURRENT_DIR/render.txt" ]; then
    echo "Error: No current results found"
    echo "Run benchmarks first:"
    echo "  bash benchmarks/scripts/run_benchmarks.sh"
    exit 1
fi

# Get version from git tag or commit
VERSION=$(git describe --tags --always 2>/dev/null || git rev-parse --short HEAD)

echo "==> Setting baseline to current results"
echo ""
echo "Version: $VERSION"
echo "Branch:  $(git rev-parse --abbrev-ref HEAD)"
echo "Commit:  $(git rev-parse --short HEAD)"
echo ""

# Confirm action
read -p "Set this as new baseline? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

# Create baseline directory
mkdir -p "$BASELINE_DIR"

# Copy current to baseline
cp "$CURRENT_DIR/render.txt" "$BASELINE_DIR/render.txt"
cp "$CURRENT_DIR/core-unicode.txt" "$BASELINE_DIR/core-unicode.txt"

# Save version info
echo "$VERSION" > "$BASELINE_DIR/version.txt"
echo "Branch: $(git rev-parse --abbrev-ref HEAD)" >> "$BASELINE_DIR/version.txt"
echo "Commit: $(git rev-parse HEAD)" >> "$BASELINE_DIR/version.txt"
echo "Date:   $(date '+%Y-%m-%d %H:%M:%S')" >> "$BASELINE_DIR/version.txt"

echo ""
echo "✓ Baseline updated!"
echo ""
echo "Baseline set to: $VERSION"
echo "Location: $BASELINE_DIR"
echo ""
echo "To compare future changes:"
echo "  bash benchmarks/scripts/compare.sh"
