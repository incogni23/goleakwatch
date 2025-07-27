.PHONY: help test test-race test-coverage bench lint clean build examples release

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Testing
test: ## Run tests
	go test -v ./...

test-race: ## Run tests with race detection
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

# Code quality
lint: ## Run linters
	golangci-lint run

lint-fix: ## Run linters with auto-fix
	golangci-lint run --fix

# Building
build: ## Build the library
	go build -v ./...

examples: ## Build and run examples
	cd examples && go run .

# Cleaning
clean: ## Clean build artifacts
	go clean
	rm -f coverage.out coverage.html
	find . -name "*.test" -delete

# Development
deps: ## Download dependencies
	go mod download
	go mod tidy

# Release
release: ## Create a new release (usage: make release VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "Please specify VERSION=v1.0.0"; exit 1; fi
	git tag $(VERSION)
	git push origin $(VERSION)

# Install tools
install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

# Security
security: ## Run security checks
	govulncheck ./...

# All checks
check: lint test-race security ## Run all checks (lint, test, security) 