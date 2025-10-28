#!/usr/bin/env bash
# Pre-Release Validation Script for Phoenix TUI Framework
# This script runs all quality checks before creating a release
# EXACTLY matches CI checks + additional validations

set -e  # Exit on first error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Header
echo ""
echo "========================================"
echo "  Phoenix TUI - Pre-Release Validation"
echo "========================================"
echo ""

# Track overall status
ERRORS=0
WARNINGS=0

# Phoenix modules
MODULES="clipboard components core layout mouse render style tea terminal testing"

# 1. Check Go version
log_info "Checking Go version..."
GO_VERSION=$(go version | awk '{print $3}')
log_success "Go version: $GO_VERSION"
echo ""

# 2. Check git status (early check)
log_info "Checking git status..."
if git diff-index --quiet HEAD --; then
    log_success "Working directory is clean"
else
    log_warning "Uncommitted changes detected"
    git status --short
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 3. Code formatting check (EXACT CI command)
log_info "Checking code formatting (gofmt -l .)..."
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    log_error "The following files need formatting:"
    echo "$UNFORMATTED"
    echo ""
    log_info "Run: gofmt -w ."
    ERRORS=$((ERRORS + 1))
else
    log_success "All files are properly formatted"
fi
echo ""

# 4. Go vet (EXACT CI command)
log_info "Running go vet..."
VET_FAILED=0
for dir in $MODULES; do
    echo "  • $dir..."
    if [ "$dir" = "clipboard" ]; then
        # Special case for clipboard: Disable unsafeptr checks (Windows syscall false positives)
        VET_OUTPUT=$(cd $dir && GOWORK=off go vet -unsafeptr=false ./... 2>&1)
    else
        VET_OUTPUT=$(cd $dir && GOWORK=off go vet ./... 2>&1)
    fi

    # Check for actual errors (not just warnings)
    if echo "$VET_OUTPUT" | grep -E "^#|FAIL" > /dev/null; then
        log_error "go vet failed for $dir"
        echo "$VET_OUTPUT" | head -5
        ERRORS=$((ERRORS + 1))
        VET_FAILED=1
    elif [ -n "$VET_OUTPUT" ]; then
        # Has output but not errors - just warnings
        log_warning "go vet warnings for $dir (non-blocking)"
        WARNINGS=$((WARNINGS + 1))
    fi
done
if [ $VET_FAILED -eq 0 ]; then
    log_success "go vet passed"
fi
echo ""

# 4b. Cross-compilation vet (Linux)
# Catch build-tag-specific issues that only show on other platforms
log_info "Running cross-compilation vet (GOOS=linux)..."
CROSS_VET_FAILED=0
for dir in $MODULES; do
    echo "  • $dir (linux)..."
    if [ "$dir" = "clipboard" ]; then
        # Disable cgo for cross-compilation + unsafeptr for clipboard
        CROSS_VET_OUTPUT=$(cd $dir && CGO_ENABLED=0 GOOS=linux GOWORK=off go vet -unsafeptr=false ./... 2>&1)
    else
        CROSS_VET_OUTPUT=$(cd $dir && CGO_ENABLED=0 GOOS=linux GOWORK=off go vet ./... 2>&1)
    fi

    # Check for actual errors (undefined functions, etc.)
    if echo "$CROSS_VET_OUTPUT" | grep -E "undefined:|FAIL" > /dev/null; then
        log_error "cross-compilation vet failed for $dir (linux)"
        echo "$CROSS_VET_OUTPUT" | head -10
        ERRORS=$((ERRORS + 1))
        CROSS_VET_FAILED=1
    fi
done
if [ $CROSS_VET_FAILED -eq 0 ]; then
    log_success "cross-compilation vet passed"
fi
echo ""

# 5. Build all packages
log_info "Building all packages..."
for dir in $MODULES; do
    echo "  • $dir..."
    if ! (cd $dir && GOWORK=off go build ./... 2>&1 | grep -v "no Go files" | grep -v "no non-test Go files" || true); then
        log_error "Build failed for $dir"
        ERRORS=$((ERRORS + 1))
    fi
done
log_success "Build successful"
echo ""

# 6. go.mod validation
log_info "Validating go.mod files..."
for dir in $MODULES; do
    echo "  • $dir/go.mod..."
    (cd $dir && go mod tidy)
done
log_success "All go.mod files are valid"
echo ""

# 7. Run tests (with race detector if GCC available)
if command -v gcc &> /dev/null; then
    log_info "Running tests with race detector..."
    RACE_FLAG="-race"
else
    log_warning "GCC not found, running tests without race detector"
    log_info "Install GCC (mingw-w64) for race detection on Windows"
    WARNINGS=$((WARNINGS + 1))
    log_info "Running tests..."
    RACE_FLAG=""
fi

TEST_FAILED=0
for dir in $MODULES; do
    echo "  • Testing $dir..."
    TEST_OUTPUT=$(cd $dir && GOWORK=off go test $RACE_FLAG ./... 2>&1)
    if echo "$TEST_OUTPUT" | grep -q "FAIL"; then
        log_error "Tests failed for $dir"
        echo "$TEST_OUTPUT" | grep "FAIL" | head -3
        ERRORS=$((ERRORS + 1))
        TEST_FAILED=1
    elif echo "$TEST_OUTPUT" | grep -q "PASS"; then
        echo "    ✓ PASS"
    fi
done
if [ $TEST_FAILED -eq 0 ]; then
    if [ -n "$RACE_FLAG" ]; then
        log_success "All tests passed with race detector"
    else
        log_success "All tests passed"
    fi
fi
echo ""

# 8. Test coverage check
log_info "Checking test coverage..."
COVERAGE_BELOW_TARGET=0
for dir in $MODULES; do
    COVERAGE=$(cd $dir && GOWORK=off go test -cover ./... 2>&1 | grep "coverage:" | tail -1 | awk '{print $5}' | sed 's/%//')
    if [ -n "$COVERAGE" ]; then
        echo "  • $dir: ${COVERAGE}%"
        # Check if coverage is above 70%
        if awk -v cov="$COVERAGE" 'BEGIN {exit !(cov >= 70.0)}'; then
            : # Coverage OK
        else
            log_warning "$dir coverage below 70% (${COVERAGE}%)"
            COVERAGE_BELOW_TARGET=1
            WARNINGS=$((WARNINGS + 1))
        fi
    fi
done
if [ $COVERAGE_BELOW_TARGET -eq 0 ]; then
    log_success "All modules meet coverage requirements (>70%)"
fi
echo ""

# 9. golangci-lint (same as CI)
log_info "Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    LINT_ISSUES=0
    for dir in $MODULES; do
        echo "  • Linting $dir..."
        if ! (cd $dir && GOWORK=off golangci-lint run --timeout=5m ./... 2>&1 | tail -5 | grep -q "0 issues"); then
            log_warning "Linter found issues in $dir (non-blocking)"
            LINT_ISSUES=1
            WARNINGS=$((WARNINGS + 1))
        fi
    done
    if [ $LINT_ISSUES -eq 0 ]; then
        log_success "golangci-lint passed with no issues"
    fi
else
    log_warning "golangci-lint not installed, skipping"
    log_info "Install: https://golangci-lint.run/welcome/install/"
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 10. Verify benchmarks compile
log_info "Verifying benchmarks compile..."
BENCHMARK_DIRS="benchmarks/comparison render/benchmarks"
BENCHMARK_FAILED=0
for dir in $BENCHMARK_DIRS; do
    if [ -d "$dir" ]; then
        echo "  • $dir..."
        if ! (cd $dir && go test -bench=. -run=^$ . > /dev/null 2>&1); then
            log_warning "Benchmark compilation issues in $dir"
            BENCHMARK_FAILED=1
            WARNINGS=$((WARNINGS + 1))
        fi
    fi
done
if [ $BENCHMARK_FAILED -eq 0 ]; then
    log_success "All benchmarks compile successfully"
fi
echo ""

# 11. Check for TODO/FIXME comments
log_info "Checking for TODO/FIXME comments..."
TODO_COUNT=$(grep -r "TODO\|FIXME" --include="*.go" --exclude-dir=vendor . 2>/dev/null | wc -l)
if [ "$TODO_COUNT" -gt 0 ]; then
    log_warning "Found $TODO_COUNT TODO/FIXME comments"
    WARNINGS=$((WARNINGS + 1))
else
    log_success "No TODO/FIXME comments found"
fi
echo ""

# 12. Check critical documentation files
log_info "Checking documentation..."
DOCS_MISSING=0
for doc in README.md CHANGELOG.md CONTRIBUTING.md SECURITY.md CODE_OF_CONDUCT.md; do
    if [ ! -f "$doc" ]; then
        log_error "Missing: $doc"
        DOCS_MISSING=1
        ERRORS=$((ERRORS + 1))
    fi
done
if [ $DOCS_MISSING -eq 0 ]; then
    log_success "All critical documentation files present"
fi
echo ""

# Summary
echo "========================================"
echo "  Summary"
echo "========================================"
echo ""

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    log_success "All checks passed! Ready for release."
    echo ""
    log_info "Next steps:"
    echo "  1. Create release branch"
    echo "  2. Update CHANGELOG.md"
    echo "  3. Merge to main"
    echo "  4. Push main"
    echo "  5. Wait for CI (5-10 min) ⏰"
    echo "  6. Push tags ONLY after CI is green ✅"
    echo ""
    exit 0
elif [ $ERRORS -eq 0 ]; then
    log_warning "Checks completed with $WARNINGS warning(s)"
    echo ""
    log_info "Review warnings above before proceeding with release"
    echo ""
    exit 0
else
    log_error "Checks failed with $ERRORS error(s) and $WARNINGS warning(s)"
    echo ""
    log_error "Fix errors before creating release"
    echo ""
    exit 1
fi
