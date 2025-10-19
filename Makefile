# Makefile for Phoenix TUI Framework
# Fallback for developers who don't have Task installed
# Recommended: Install Task (https://taskfile.dev/) and use Taskfile.yml instead
#
# Usage: make <target>
# Example: make test, make lint, make bench

.PHONY: help test test-core test-coverage lint lint-fix fmt vet check build build-examples deps deps-update bench bench-core clean dev ci run-basic run-unicode

# Default target
help:
	@echo "Phoenix TUI Framework - Available Commands"
	@echo ""
	@echo "Testing:"
	@echo "  make test           - Run all tests with coverage"
	@echo "  make test-core      - Run core package tests"
	@echo "  make test-coverage  - Generate HTML coverage report"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make lint-fix       - Run golangci-lint with auto-fix"
	@echo "  make fmt            - Format code with gofmt"
	@echo "  make vet            - Run go vet"
	@echo "  make check          - Run all checks (fmt+vet+lint+test)"
	@echo ""
	@echo "Building:"
	@echo "  make build          - Build all packages"
	@echo "  make build-examples - Build example applications"
	@echo ""
	@echo "Dependencies:"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make deps-update    - Update all dependencies"
	@echo ""
	@echo "Benchmarks:"
	@echo "  make bench          - Run all benchmarks"
	@echo "  make bench-core     - Run core benchmarks"
	@echo ""
	@echo "Development:"
	@echo "  make dev            - Pre-commit checks (fmt+vet+lint-fix+test)"
	@echo "  make ci             - CI checks (same as GitHub Actions)"
	@echo "  make clean          - Remove build artifacts"
	@echo ""
	@echo "Examples:"
	@echo "  make run-basic      - Run basic example"
	@echo "  make run-unicode    - Run Unicode example"
	@echo ""
	@echo "Recommended: Install Task (https://taskfile.dev/) for better experience"

# ========================================
# Testing
# ========================================

test:
	@echo "Running all tests..."
	go test -v -race -cover ./...

test-core:
	@echo "Running core tests..."
	cd core && go test -v -race -coverprofile=coverage.out ./...
	cd core && go tool cover -func=coverage.out

test-coverage:
	@echo "Generating coverage report..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# ========================================
# Code Quality
# ========================================

lint:
	@echo "Running linter..."
	golangci-lint run --config .golangci.yml --timeout=5m ./...

lint-fix:
	@echo "Running linter with auto-fix..."
	golangci-lint run --config .golangci.yml --fix ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Running go vet..."
	go vet ./...

check: fmt vet lint test
	@echo "✅ All checks passed!"

# ========================================
# Building
# ========================================

build:
	@echo "Building all packages..."
	go build ./...

build-examples:
	@echo "Building examples..."
	@mkdir -p bin
	go build -o bin/basic.exe ./examples/basic
	go build -o bin/unicode.exe ./examples/unicode
	@echo "Examples built in bin/"

# ========================================
# Dependencies
# ========================================

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# ========================================
# Benchmarks
# ========================================

bench:
	@echo "Running all benchmarks..."
	go test -bench=. -benchmem ./...

bench-core:
	@echo "Running core benchmarks..."
	cd core && go test -bench=. -benchmem ./...

# ========================================
# Development
# ========================================

dev: fmt vet lint-fix test
	@echo "✅ All checks passed! Ready to commit."

ci: fmt vet lint test-coverage build
	@echo "✅ CI checks complete!"

clean:
	@echo "Cleaning build artifacts..."
	@rm -f coverage.out coverage.html
	@find . -name "*.test" -delete 2>/dev/null || true
	@find . -name "*.coverprofile" -delete 2>/dev/null || true
	@rm -f bench-*.txt
	@rm -rf bin/
	@echo "Clean complete!"

# ========================================
# Examples
# ========================================

run-basic:
	@echo "Running basic example..."
	cd examples/basic && go run main.go

run-unicode:
	@echo "Running Unicode example..."
	cd examples/unicode && go run main.go
