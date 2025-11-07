# Makefile for goagent library
.PHONY: build test run example clean install lint fmt vet deps help build-example

# Variables
BINARY_NAME=example
EXAMPLE_PATH=./example
EXAMPLE_BINARY=$(EXAMPLE_PATH)/$(BINARY_NAME)
MODULE_PATH=github.com/frozenkro/goagent
GOFLAGS=-v

# Default target
help: ## Display this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Library building and testing
build: ## Build the goagent library
	@echo "Building goagent library..."
	go build $(GOFLAGS) ./...

install: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Testing
test: ## Run all tests for the library
	@echo "Running tests..."
	go test $(GOFLAGS) ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Example building and running
build-example: ## Build the example binary
	@echo "Building example..."
	cd $(EXAMPLE_PATH) && go build $(GOFLAGS) -o $(BINARY_NAME) ./example.go

run-example: build-example ## Build and run the example
	@echo "Running example..."
	cd $(EXAMPLE_PATH) && ./$(BINARY_NAME)

example: run-example ## Alias for run-example command

# Code quality
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "Running linter..."
	golangci-lint run

# Utilities
clean: ## Clean build artifacts
	@echo "Cleaning up..."
	rm -f $(EXAMPLE_BINARY)
	rm -f coverage.out coverage.html
	cd $(EXAMPLE_PATH) && rm -rf fib/

deps: ## Show module dependencies
	@echo "Module dependencies:"
	go list -m all

mod-update: ## Update all dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Development helpers
dev-setup: install ## Set up development environment
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi

check: fmt vet test ## Run all checks (format, vet, test)
	@echo "All checks completed successfully!"

# Quick development cycle
dev: clean fmt vet build run-example ## Complete development cycle: clean, format, vet, build, run example

# CI/CD ready targets
ci-test: ## Run tests suitable for CI environment
	go test -race -coverprofile=coverage.out ./...

ci-build: ## Build for CI environment
	CGO_ENABLED=0 go build ./...
	cd $(EXAMPLE_PATH) && CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_NAME) ./example.go

# Release helpers
mod-verify: ## Verify module dependencies
	go mod verify

mod-graph: ## Show module dependency graph
	go mod graph

benchmark: ## Run benchmarks
	go test -bench=. -benchmem ./...