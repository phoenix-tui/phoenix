#!/bin/bash
# Phoenix TUI Framework - Lint All Modules
# Runs golangci-lint on all 10 Phoenix modules

set -e  # Exit on first error

modules="clipboard components core layout mouse render style tea terminal testing"

echo "Running golangci-lint on all Phoenix modules..."
echo ""

failed_modules=""
passed_count=0

for module in $modules; do
    echo "==> Linting $module"
    cd "$module"

    if golangci-lint run ./...; then
        echo "✓ $module passed"
        ((passed_count++))
    else
        echo "✗ $module has linting issues"
        failed_modules="$failed_modules $module"
    fi

    cd ..
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ -z "$failed_modules" ]; then
    echo "✓✓✓ All modules passed linting ($passed_count/10) ✓✓✓"
    exit 0
else
    echo "✗✗✗ Some modules failed linting ✗✗✗"
    echo "Failed modules:$failed_modules"
    exit 1
fi