#!/bin/bash
# Phoenix TUI Framework - Release Tag Creator
# Creates all 11 tags (10 module tags + 1 root tag) for a release

set -e  # Exit on error

VERSION="$1"

if [ -z "$VERSION" ]; then
  echo "❌ Error: Version required!"
  echo ""
  echo "Usage: $0 <version>"
  echo ""
  echo "Examples:"
  echo "  $0 v0.1.0-beta.2"
  echo "  $0 v0.1.0-rc.1"
  echo "  $0 v0.1.0"
  echo "  $0 v1.0.0"
  exit 1
fi

# Validate version format
if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-z]+\.[0-9]+)?$ ]]; then
  echo "❌ Error: Invalid version format!"
  echo "Expected: vX.Y.Z or vX.Y.Z-beta.N or vX.Y.Z-rc.N"
  echo "Got: $VERSION"
  exit 1
fi

# Module list (10 Phoenix libraries)
MODULES="clipboard components core layout mouse render style tea terminal testing"

echo "🚀 Phoenix TUI Framework - Release Tag Creator"
echo "================================================"
echo ""
echo "Version: $VERSION"
echo "Modules: 10 libraries + 1 root"
echo ""

# Check if tags already exist
echo "🔍 Checking for existing tags..."
TAG_EXISTS=false
for module in $MODULES; do
  if git rev-parse "$module/$VERSION" >/dev/null 2>&1; then
    echo "⚠️  Tag $module/$VERSION already exists!"
    TAG_EXISTS=true
  fi
done
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "⚠️  Tag $VERSION already exists!"
  TAG_EXISTS=true
fi

if [ "$TAG_EXISTS" = true ]; then
  echo ""
  echo "❌ Error: Some tags already exist!"
  echo "Delete existing tags first with:"
  echo "  git tag -d <tag>"
  echo "  git push origin --delete <tag>"
  exit 1
fi

echo "✅ No existing tags found"
echo ""

# Create module-specific tags
echo "📦 Creating module-specific tags..."
for module in $MODULES; do
  TAG="$module/$VERSION"
  echo "  Creating: $TAG"
  git tag -a "$TAG" -m "Phoenix TUI Framework $VERSION - $module module

Complete $module implementation with:
- DDD architecture (domain/application/infrastructure/api)
- High test coverage (90%+ target achieved)
- Production-ready code quality
- Cross-platform support (Linux/macOS/Windows)

Part of unified Phoenix TUI Framework release $VERSION.
All modules released together for version consistency.

See CHANGELOG.md for full release notes."
done

echo "✅ Module tags created: 10"
echo ""

# Create root tag
echo "📦 Creating root tag..."
git tag -a "$VERSION" -m "Phoenix TUI Framework $VERSION (root umbrella module)

See CHANGELOG.md for full release notes.

This is the root umbrella module containing replace directives for all 10 Phoenix libraries:
- clipboard
- components
- core
- layout
- mouse
- render
- style
- tea
- terminal
- testing

Import individual modules directly for best experience:
  go get github.com/phoenix-tui/phoenix/tea@tea/$VERSION
  go get github.com/phoenix-tui/phoenix/components@components/$VERSION"

echo "✅ Root tag created: 1"
echo ""

# List all created tags
echo "📋 All created tags (11 total):"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
for module in $MODULES; do
  echo "  $module/$VERSION"
done
echo "  $VERSION (root)"
echo ""

# Show current commit
COMMIT=$(git rev-parse HEAD)
COMMIT_SHORT=$(git rev-parse --short HEAD)
echo "📌 Tagged commit: $COMMIT_SHORT"
echo ""

# Instructions for pushing
echo "✅ Tags created successfully!"
echo ""
echo "Next steps:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1. Review tags:"
echo "   git tag | grep $VERSION"
echo ""
echo "2. Push all tags to origin:"
echo "   git push origin --tags"
echo ""
echo "3. Verify on pkg.go.dev after ~10 minutes:"
echo "   https://pkg.go.dev/github.com/phoenix-tui/phoenix@$VERSION"
echo ""
echo "4. Create GitHub Release:"
echo "   gh release create $VERSION --title \"Phoenix TUI $VERSION\" --notes-file CHANGELOG.md"
echo ""
echo "⚠️  REMINDER: Only push tags AFTER CI is green!"
echo "   https://github.com/phoenix-tui/phoenix/actions"
