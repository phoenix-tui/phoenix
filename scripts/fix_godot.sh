#!/bin/bash
# Fix godot linter issues - add periods to end of comments
# Usage: ./scripts/fix_godot.sh

# Find all Go files and fix comments that end without period
find clipboard -name "*.go" -type f | while read -r file; do
    # Backup file
    cp "$file" "$file.bak"

    # Fix single-line comments that don't end with period, URL, or special chars
    # Pattern: // Text$ → // Text.
    # Exclude: URLs (http://), already has period, ends with special chars (:,?,!,),})
    sed -i 's|^\([[:space:]]*//[[:space:]][^/].*[^.?!:)}\n]\)$|\1.|g' "$file"

    # Check if file changed
    if ! diff -q "$file" "$file.bak" > /dev/null 2>&1; then
        echo "Fixed: $file"
        rm "$file.bak"
    else
        # No changes - restore backup
        mv "$file.bak" "$file"
    fi
done

echo ""
echo "Godot fixes applied. Run 'golangci-lint run ./clipboard/...' to verify."
