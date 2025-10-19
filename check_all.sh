#!/bin/bash
set -e

echo "Running full local CI checks..."

for dir in clipboard components core layout mouse render style tea terminal testing; do
  echo "==> Checking $dir"
  cd $dir

  # Disable workspace for isolated module checks
  export GOWORK=off

  # go vet
  if [ "$dir" = "clipboard" ]; then
    go vet -unsafeptr=false $(go list ./... | grep -v "/examples" | grep -v "/benchmarks")
  else
    go vet $(go list ./... | grep -v "/examples" | grep -v "/benchmarks")
  fi

  # go build
  BUILDABLE=$(go list -f '{{if (and (or .GoFiles .CgoFiles) (not .Module))}}{{.ImportPath}}{{end}}' ./... 2>/dev/null | grep -v "/examples" | grep -v "/benchmarks" || true)
  if [ -n "$BUILDABLE" ]; then
    go build $BUILDABLE
  fi

  # go test (skip examples and benchmarks)
  TESTABLE=$(go list ./... 2>/dev/null | grep -v "/examples" | grep -v "/benchmarks" || true)
  if [ -n "$TESTABLE" ]; then
    go test -short $TESTABLE
  fi

  cd ..
  echo "✓ $dir passed"
done

echo ""
echo "✓✓✓ All local checks PASSED ✓✓✓"
