#!/bin/bash
# Phoenix TUI Framework - Module Version Bumper
# Updates all cross-module require versions in go.mod files.
#
# Usage: bash scripts/bump-versions.sh v0.2.4
#
# This script updates require versions but does NOT run go mod tidy,
# because the new version doesn't exist on proxy yet (tags not pushed).
# go.work handles local resolution during development.

set -e

VERSION="$1"

if [ -z "$VERSION" ]; then
  echo "Error: Version required!"
  echo ""
  echo "Usage: $0 <version>"
  echo "Example: $0 v0.2.4"
  exit 1
fi

# Validate version format
if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-z]+\.[0-9]+)?$ ]]; then
  echo "Error: Invalid version format!"
  echo "Expected: vX.Y.Z or vX.Y.Z-beta.N or vX.Y.Z-rc.N"
  echo "Got: $VERSION"
  exit 1
fi

echo "Phoenix TUI - Version Bumper"
echo "============================"
echo "Target version: $VERSION"
echo ""

# All Phoenix module import paths
PHOENIX_MODULES=(
  "github.com/phoenix-tui/phoenix/clipboard"
  "github.com/phoenix-tui/phoenix/components"
  "github.com/phoenix-tui/phoenix/core"
  "github.com/phoenix-tui/phoenix/layout"
  "github.com/phoenix-tui/phoenix/mouse"
  "github.com/phoenix-tui/phoenix/render"
  "github.com/phoenix-tui/phoenix/style"
  "github.com/phoenix-tui/phoenix/tea"
  "github.com/phoenix-tui/phoenix/terminal"
  "github.com/phoenix-tui/phoenix/testing"
)

# Find all go.mod files
GOMOD_FILES=$(find . -name "go.mod" -not -path "./.git/*" | sort)

UPDATED=0
for gomod in $GOMOD_FILES; do
  CHANGED=false
  for mod in "${PHOENIX_MODULES[@]}"; do
    # Match require lines like: github.com/phoenix-tui/phoenix/core v0.2.0
    # Also matches with // indirect suffix
    if grep -q "${mod} v" "$gomod" 2>/dev/null; then
      # Replace any version with target version (preserve // indirect if present)
      if sed -i "s|${mod} v[^ ]*|${mod} ${VERSION}|g" "$gomod" 2>/dev/null; then
        CHANGED=true
      fi
    fi
  done
  if [ "$CHANGED" = true ]; then
    echo "  Updated: $gomod"
    UPDATED=$((UPDATED + 1))
  fi
done

echo ""
echo "Updated $UPDATED go.mod files"
echo ""

# Verify changes
echo "Verifying..."
MISMATCHED=0
for gomod in $GOMOD_FILES; do
  for mod in "${PHOENIX_MODULES[@]}"; do
    # Check for any Phoenix module with wrong version
    BAD=$(grep "${mod} v" "$gomod" 2>/dev/null | grep -v "${VERSION}" | grep -v "replace" | grep -v "module" || true)
    if [ -n "$BAD" ]; then
      echo "  MISMATCH in $gomod: $BAD"
      MISMATCHED=$((MISMATCHED + 1))
    fi
  done
done

if [ "$MISMATCHED" -gt 0 ]; then
  echo ""
  echo "ERROR: $MISMATCHED version mismatches found!"
  exit 1
fi

echo "All Phoenix cross-references point to $VERSION"
echo ""
echo "Next steps:"
echo "  1. git add -A && git commit -m \"chore: bump module versions to $VERSION\""
echo "  2. Push, wait for CI"
echo "  3. bash scripts/create-release-tags.sh $VERSION"
