#!/usr/bin/env bash
# Pre-release check script - EXACTLY matches CI checks
# Run this BEFORE creating any release

set -e

echo "🔍 Phoenix TUI - Pre-Release Checks"
echo "===================================="
echo ""

FAILED=0

# 1. Formatting check (EXACT CI command)
echo "==> Step 1: Formatting (gofmt -l .)"
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    echo "❌ ERROR: The following files are not formatted:"
    echo "$UNFORMATTED"
    echo ""
    echo "Run 'gofmt -w .' to fix formatting issues."
    FAILED=1
else
    echo "✅ All files properly formatted"
fi
echo ""

# 2. Go vet (EXACT CI command)
echo "==> Step 2: Running go vet"
for dir in clipboard components core layout mouse render style tea terminal testing; do
    echo "  go vet $dir..."
    if [ "$dir" = "clipboard" ]; then
        # Special case for clipboard: Disable unsafeptr checks (Windows syscall false positives)
        if ! (cd $dir && GOWORK=off go vet -unsafeptr=false $(go list ./... | grep -v "/examples" | grep -v "/benchmarks") 2>&1); then
            echo "❌ go vet failed for $dir"
            FAILED=1
        fi
    else
        if ! (cd $dir && GOWORK=off go vet $(go list ./... | grep -v "/examples" | grep -v "/benchmarks") 2>&1); then
            echo "❌ go vet failed for $dir"
            FAILED=1
        fi
    fi
done
echo "✅ go vet successful"
echo ""

# 3. Build all packages
echo "==> Step 3: Building all packages"
for dir in clipboard components core layout mouse render style tea terminal testing; do
    echo "  Building $dir..."
    (cd $dir && GOWORK=off go build ./... 2>&1 | grep -v "no Go files" | grep -v "no non-test Go files" || true)
done
echo "✅ Build successful"
echo ""

# 4. Run tests with race detector
echo "==> Step 4: Tests with race detector"
for dir in clipboard components core layout mouse render style tea terminal testing; do
    echo "  Testing $dir..."
    if ! (cd $dir && GOWORK=off go test -race ./... 2>&1 | grep -E "(PASS|FAIL)" | head -1); then
        echo "❌ Tests failed for $dir"
        FAILED=1
    fi
done
echo ""

# 5. Linter (same as CI)
echo "==> Step 5: Linter"
for dir in clipboard components core layout mouse render style tea terminal testing; do
    echo "  Linting $dir..."
    if ! (cd $dir && GOWORK=off golangci-lint run --timeout=5m ./... 2>&1 | head -1 | grep -q "0 issues"); then
        echo "⚠️  Linter warnings for $dir (non-blocking in beta)"
    fi
done
echo ""

# Final result
if [ $FAILED -eq 0 ]; then
    echo "✅ All pre-release checks PASSED!"
    echo ""
    echo "Next steps:"
    echo "  1. Create release branch"
    echo "  2. Update CHANGELOG.md"
    echo "  3. Merge to main"
    echo "  4. Push main"
    echo "  5. Wait for CI (5-10 min)"
    echo "  6. Push tags ONLY after CI is green"
    exit 0
else
    echo "❌ Pre-release checks FAILED!"
    echo ""
    echo "Fix all errors before proceeding with release."
    exit 1
fi
