## Description

**What does this PR do?**

<!-- Provide a clear and concise description of the changes -->

Fixes # (issue)

## Type of Change

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📚 Documentation update
- [ ] 🔧 Refactoring (no functional changes)
- [ ] ⚡ Performance improvement
- [ ] ✅ Test coverage improvement

## Which Libraries Are Affected?

- [ ] phoenix/core
- [ ] phoenix/style
- [ ] phoenix/tea
- [ ] phoenix/layout
- [ ] phoenix/render
- [ ] phoenix/components
- [ ] phoenix/mouse
- [ ] phoenix/clipboard
- [ ] phoenix/testing
- [ ] phoenix/terminal
- [ ] Documentation only
- [ ] Infrastructure (CI/CD, scripts, etc.)

## Checklist

### Code Quality

- [ ] My code follows the Phoenix coding style (gofmt, golangci-lint)
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] My changes generate no new warnings or errors

### Testing

- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Test coverage has not decreased (run `go test -cover ./...`)
- [ ] I have tested on multiple platforms (if applicable):
  - [ ] Linux
  - [ ] macOS
  - [ ] Windows

### Documentation

- [ ] I have updated the documentation accordingly
- [ ] I have updated the CHANGELOG.md (if user-facing change)
- [ ] I have added/updated code examples (if new feature)
- [ ] I have updated the README.md (if API change)

### DDD Architecture (if applicable)

- [ ] Domain layer: Pure business logic, no dependencies on other layers
- [ ] Application layer: Use cases, orchestration only
- [ ] Infrastructure layer: Technical details, I/O operations
- [ ] API layer: Public interface, follows fluent builder pattern

## Testing Evidence

**How has this been tested?**

<!-- Describe the tests you ran to verify your changes -->

```bash
# Example: Test commands you ran
go test -v ./components/...
go test -cover ./...
```

**Test coverage:**

```
# Paste test coverage output here
```

## Screenshots (if applicable)

<!-- Add screenshots for UI changes or terminal output -->

## Performance Impact

**Does this change affect performance?**

- [ ] No performance impact
- [ ] Performance improvement (include benchmarks below)
- [ ] Performance regression (justify why necessary)

```bash
# Benchmark results (if applicable)
go test -bench=. -benchmem ./...
```

## Breaking Changes

**Does this PR introduce breaking changes?**

- [ ] No breaking changes
- [ ] Yes (describe migration path below)

**Migration guide (if breaking):**

<!-- Explain how users should update their code -->

## Additional Context

<!-- Add any other context about the PR here -->

## Related Issues/PRs

<!-- Link to related issues or pull requests -->

- Closes #
- Related to #
- Depends on #

---

## For Maintainers

### Review Checklist

- [ ] Code quality meets Phoenix standards
- [ ] Tests are comprehensive and pass
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated (if needed)
- [ ] Breaking changes are documented (if any)
- [ ] Performance impact is acceptable
- [ ] DDD architecture is maintained
- [ ] Multi-module versioning is considered (if library change)
