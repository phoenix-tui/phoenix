# Phoenix TUI - Performance Benchmarks

This directory contains performance benchmarks for Phoenix TUI Framework.

---

## 🎯 Why Benchmarks?

Phoenix TUI is built for **high performance**. We maintain strict performance standards:

- 🚀 **60 FPS target**: Phoenix renders at **37,818 FPS** (630x faster)
- ⚡ **Zero allocations**: Critical paths have 0 memory allocations
- 🧠 **Minimal memory**: <100 bytes per operation on hot paths
- 🌍 **Real-world tested**: Benchmarks include scrolling terminals and code editors

---

## 📊 Quick Start

### Run Benchmarks

```bash
# All critical benchmarks (render, Unicode, etc.)
bash benchmarks/scripts/run_benchmarks.sh

# Individual modules
cd render && go test -bench=. -benchmem ./benchmarks
cd core && go test -bench=. -benchmem ./domain/service
```

### View Results

```bash
# Latest results
cat benchmarks/results/current/render.txt

# Summary
cat benchmarks/results/current/README.md
```

---

## 📂 Directory Structure

```
benchmarks/
├── README.md                 # This file (overview)
├── comparison/               # Comparison with other libraries
│   ├── go.mod               # Separate module for comparison tests
│   └── *_test.go            # Tests comparing Phoenix vs Lipgloss/Bubbletea
├── results/                  # Benchmark results storage
│   ├── README.md            # Detailed results documentation
│   ├── current/             # Latest benchmark runs
│   ├── baseline/            # Stable baseline for comparisons
│   └── history/             # Release milestones
└── scripts/                  # Automation scripts
    ├── run_benchmarks.sh    # Run all benchmarks
    ├── compare.sh           # Compare current vs baseline
    └── set_baseline.sh      # Update baseline
```

---

## 🚀 Current Performance

### Render Performance

| Metric | Result | vs Target |
|--------|--------|-----------|
| **Full Screen Rendering** | 37,818 FPS (26.4 µs) | **630x faster** than 60 FPS |
| **Memory** | 4 B/op | Minimal |
| **Allocations** | 0 allocs/op | Perfect |

### Unicode Performance

| Operation | Result | Status |
|-----------|--------|--------|
| ASCII | 64 ns/op | ✅ 0 allocs |
| Emoji | 110 ns/op | ✅ 0 allocs |
| CJK | 160 ns/op | ✅ 0 allocs |

### Real-World Scenarios

| Scenario | Performance |
|----------|-------------|
| Scrolling Terminal | 88 µs/op |
| Code Editor | 117 µs/op |
| Small UI Change | 28 µs/op |

Full results: see `results/history/` for per-release summaries.

---

## 🔬 Benchmark Categories

### 1. Render Benchmarks (`render/benchmarks/`)

Core rendering performance:
- Full screen rendering (60 FPS target validation)
- Differential rendering (typical case)
- Unicode rendering (emoji, CJK)
- Real-world scenarios (terminals, editors)
- Memory and allocation tracking

**Key Benchmarks**:
- `BenchmarkFullScreen_60FPS` - Must be >24,000 FPS
- `BenchmarkDifferential_SmallChange` - Typical UI update
- `BenchmarkUnicode_Emoji` - Emoji rendering performance

### 2. Unicode Benchmarks (`core/domain/service/`)

Unicode string width calculation:
- ASCII performance
- Emoji and complex emoji
- CJK characters (Chinese, Japanese, Korean)
- Mixed content (ASCII + Unicode)
- Grapheme cluster operations

**Key Benchmarks**:
- `BenchmarkStringWidth_ASCII_Short` - Hot path
- `BenchmarkStringWidth_Emoji_Short` - Emoji width
- `BenchmarkClusterWidth_*` - Individual cluster performance

### 3. Comparison Benchmarks (`comparison/`)

Phoenix vs other libraries:
- Correctness tests (Unicode handling)
- Performance comparisons
- Real-world scenario tests

**Note**: These tests require external dependencies (Lipgloss) and are in a separate module.

---

## 📈 Performance Tracking

### For Users

To verify Phoenix performance on your machine:

```bash
# Clone repo
git clone https://github.com/phoenix-tui/phoenix
cd phoenix

# Run benchmarks
bash benchmarks/scripts/run_benchmarks.sh

# View results
cat benchmarks/results/current/render.txt
```

### For Contributors

See [`results/README.md`](results/README.md) for:
- How to run benchmarks during development
- How to compare with baseline
- How to save results for releases
- Performance regression guidelines

---

## 🎯 Performance Standards

Phoenix maintains these standards:

| Metric | Minimum | Target | Current |
|--------|---------|--------|---------|
| **Render FPS** | >6,000 (100x target) | >24,000 (400x) | 37,818 ✅ |
| **Memory (hot path)** | <500 B/op | <100 B/op | 4 B/op ✅ |
| **Allocations (critical)** | <10 allocs/op | 0 allocs/op | 0 ✅ |
| **Unicode ASCII** | <200 ns/op | <100 ns/op | 64 ns/op ✅ |
| **Unicode Emoji** | <500 ns/op | <200 ns/op | 110 ns/op ✅ |

**Policy**:
- ⚠️ Changes causing **+10% regression** require justification
- ❌ Changes causing **+20% regression** are not accepted without fixes
- ✅ Improvements of any % are welcomed

---

## 🛠️ Tools

### benchstat (recommended)

Statistical comparison of benchmark results:

```bash
# Install
go install golang.org/x/perf/cmd/benchstat@latest

# Compare two runs
benchstat before.txt after.txt

# Example output:
# name                 old time/op  new time/op  delta
# BenchmarkRender-12     34.3µs ± 2%  26.4µs ± 3%  -23.03%  (p=0.000 n=10+10)
```

### Go built-in

```bash
# Basic benchmark
go test -bench=.

# With memory stats
go test -bench=. -benchmem

# Longer runs for stability
go test -bench=. -benchtime=3s

# Profile CPU
go test -bench=. -cpuprofile=cpu.prof

# Profile memory
go test -bench=. -memprofile=mem.prof
```

---

## 📚 Resources

- [Go Benchmark Documentation](https://pkg.go.dev/testing#hdr-Benchmarks)
- [benchstat tool](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- [Phoenix Performance Reports](results/history/)

---

## ❓ FAQ

**Q: Why are allocations important?**
A: Zero allocations = zero GC pressure = consistent latency. Critical for smooth 60 FPS rendering.

**Q: Why 630x faster than 60 FPS?**
A: Performance headroom ensures smooth rendering even on slower machines or under heavy load.

**Q: Can I compare Phoenix with Bubbletea/Lipgloss?**
A: Yes! See `comparison/` directory for correctness and performance tests.

**Q: How do I report performance issues?**
A: Run benchmarks, save results, and open issue with before/after comparison.

---

*Phoenix TUI - Performance Benchmarks*
