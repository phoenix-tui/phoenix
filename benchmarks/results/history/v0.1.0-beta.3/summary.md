# Phoenix TUI v0.1.0-beta.3 - Benchmark Results

**Release Date**: 2025-10-23
**Branch**: develop
**Commit**: 5b403dc

---

## 🎯 Performance Summary

### Render Performance (Critical Path)

| Metric | Result | Status |
|--------|--------|--------|
| **Full Screen 60FPS** | 37,818 FPS (26.4 µs/op) | ✅ 630x faster than 60 FPS target |
| **Memory** | 4 B/op | ✅ Minimal |
| **Allocations** | 0 allocs/op | ✅ Perfect |
| **vs beta.2 baseline** | +30% faster | ✅ Improved |

### Unicode Performance (Core)

| Operation | Result | Allocations |
|-----------|--------|-------------|
| ASCII Short | 64 ns/op | 0 allocs |
| Emoji Short | 110 ns/op | 0 allocs |
| CJK Short | 160 ns/op | 0 allocs |
| Mixed Long | 17.4 µs/op | 1 alloc |

### Real-World Scenarios

| Scenario | Performance |
|----------|-------------|
| Scrolling Terminal | 88 µs/op |
| Code Editor | 117 µs/op |
| Small Change | 28 µs/op |

---

## 📝 Release Notes

**Key Changes**:
- ✅ ExecProcess API for interactive commands (vim, ssh, claude)
- ✅ Terminal raw mode management
- ✅ Race condition fix in inputReader (generation counter)
- ✅ Umbrella module github.com/phoenix-tui/phoenix
- ✅ Keybindings fixes (Ctrl+W, Alt+Backspace)

**Performance Impact**:
- 🚀 30% render performance improvement
- ✅ Zero allocations maintained
- ✅ No overhead from race fix

---

## 🔬 Full Results

See:
- `render.txt` - Full render benchmark results
- `core-unicode.txt` - Full Unicode benchmark results

Compare with baseline:
```bash
benchstat baseline/render.txt history/v0.1.0-beta.3/render.txt
```
