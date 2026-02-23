# Phoenix TUI - Benchmark Results

This directory stores benchmark results for performance tracking and comparison.

---

## 📁 Directory Structure

```
results/
├── current/              # Latest benchmark runs (updated frequently)
│   ├── render.txt        # Render performance (benchstat format)
│   ├── core-unicode.txt  # Unicode operations performance
│   └── README.md         # Metadata (date, commit, branch)
│
├── baseline/             # Stable baseline for comparisons
│   ├── render.txt        # Last stable release benchmarks
│   ├── core-unicode.txt
│   └── version.txt       # Baseline version info
│
└── history/              # Milestone results (releases only)
    └── <version>/        # Each release contains:
        ├── render.txt         # Full results
        ├── core-unicode.txt
        └── summary.md         # Human-readable summary
```

---

## 🚀 Quick Start

### Run Benchmarks

```bash
# Run all critical benchmarks
bash benchmarks/scripts/run_benchmarks.sh

# Results saved to benchmarks/results/current/
```

### Compare with Baseline

```bash
# Compare current vs baseline (requires benchstat)
bash benchmarks/scripts/compare.sh

# Install benchstat if needed:
go install golang.org/x/perf/cmd/benchstat@latest
```

### Set New Baseline

```bash
# Update baseline to current results (after release)
bash benchmarks/scripts/set_baseline.sh
```

---

## 📊 Workflow

### During Development

1. **Before making changes**:
   ```bash
   bash benchmarks/scripts/run_benchmarks.sh
   ```
   This captures baseline for your work.

2. **After making changes**:
   ```bash
   bash benchmarks/scripts/run_benchmarks.sh
   bash benchmarks/scripts/compare.sh
   ```
   Check for regressions.

3. **Acceptable Changes**:
   - ✅ **Improvements**: Any % faster is good!
   - ✅ **Stable**: ±5% is noise (acceptable)
   - ⚠️ **Minor regression**: +5-10% requires justification
   - ❌ **Regression**: +10%+ requires fix or explanation

### Before Release

1. **Run full benchmark suite**:
   ```bash
   bash benchmarks/scripts/run_benchmarks.sh
   ```

2. **Compare with baseline**:
   ```bash
   bash benchmarks/scripts/compare.sh
   ```

3. **Save to history** (for releases):
   ```bash
   # Copy results to history
   VERSION="vX.Y.Z"  # Replace with actual version
   mkdir -p benchmarks/results/history/$VERSION
   cp benchmarks/results/current/*.txt benchmarks/results/history/$VERSION/

   # Create summary
   vi benchmarks/results/history/$VERSION/summary.md

   # Commit to git
   git add benchmarks/results/history/$VERSION/
   git commit -m "docs(benchmarks): save $VERSION results"
   ```

4. **Update baseline**:
   ```bash
   bash benchmarks/scripts/set_baseline.sh
   ```

---

## 🎯 Performance Targets

### Render Performance
- **Target**: 60 FPS (16.67 ms/frame)
- **Achieved**: 37,818 FPS (26.4 µs/frame) = **630x faster** ✅
- **Memory**: <100 B/op on hot paths ✅
- **Allocations**: 0 allocs/op on critical paths ✅

### Unicode Performance
- **ASCII**: <100 ns/op ✅
- **Emoji**: <200 ns/op ✅
- **CJK**: <200 ns/op ✅
- **Allocations**: 0 on hot paths ✅

### Real-World Scenarios
- **Scrolling Terminal**: <100 µs/op ✅
- **Code Editor**: <200 µs/op ✅
- **Small Changes**: <50 µs/op ✅

---

## 📝 Best Practices

### What to Store in Git

✅ **DO commit**:
- History milestone results (`history/*/`)
- Baseline results (`baseline/`)
- Scripts and documentation
- Summary files (`.md`)

❌ **DON'T commit frequently**:
- Current results (updated often, causes noise)
- Store current/ only when setting baseline or before major changes

### Benchmark Result Format

Results are stored in **benchstat format** (Go standard):
```
BenchmarkFullScreen_60FPS-12    131121    26442 ns/op    37818 fps    4 B/op    0 allocs/op
```

This format allows:
- Machine comparison with `benchstat`
- Statistical analysis
- Historical tracking
- Easy diff in git

---

## 🔧 Tools

### benchstat

Install:
```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

Usage:
```bash
# Compare two benchmark runs
benchstat baseline/render.txt current/render.txt

# Show only regressions
benchstat -delta-test=ttest baseline/render.txt current/render.txt
```

### Continuous Benchmarking

For CI/CD integration, see `.github/workflows/` (when added).

---

## 📚 References

- [Go Benchmarking Best Practices](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
- [benchstat documentation](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- [Continuous Benchmarking](https://dev.to/vearutop/continuous-benchmarking-with-go-and-github-actions-41ok)

---

*Phoenix TUI - Benchmark Results*
