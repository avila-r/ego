.PHONY: help test test-verbose test-race coverage lint fmt vet clean install-tools bench check-deps tidy

# Default target
.DEFAULT_GOAL := help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Coverage parameters
COVERAGE_DIR=coverage
COVERAGE_FILE=$(COVERAGE_DIR)/coverage.out
COVERAGE_HTML=$(COVERAGE_DIR)/coverage.html

# Linter
GOLANGCI_LINT=golangci-lint

help: ## Display this help message
	@echo "Extended Go (ego) - Makefile commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: ## Run tests for all packages
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	$(GOTEST) -v -count=1 ./...

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	$(GOTEST) -race -short ./...

coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	$(GOCMD) tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print "Total coverage: " $$3}'

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

lint: ## Run golangci-lint
	@echo "Running linter..."
	@which $(GOLANGCI_LINT) > /dev/null || (echo "golangci-lint not found. Run 'make install-tools' first" && exit 1)
	$(GOLANGCI_LINT) run ./...

fmt: ## Format code using gofmt
	@echo "Formatting code..."
	$(GOFMT) ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "All checks passed!"

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	$(GOMOD) tidy

check-deps: ## Check for outdated dependencies
	@echo "Checking for outdated dependencies..."
	$(GOCMD) list -u -m all

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@which $(GOLANGCI_LINT) > /dev/null || \
		(echo "Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin)
	@echo "Tools installed successfully!"

clean: ## Clean build artifacts and coverage reports
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(COVERAGE_DIR)
	rm -f *.out *.test
	@echo "Clean complete!"

ci: check coverage ## Run CI pipeline (all checks + coverage)
	@echo "CI pipeline complete!"

# Module-specific tests (examples)
test-box: ## Run tests for box package only
	$(GOTEST) -v ./box/...

test-collection: ## Run tests for collection package only
	$(GOTEST) -v ./collection/...

test-stream: ## Run tests for stream package only
	$(GOTEST) -v ./stream/...

# Quick targets
quick-test: ## Run tests without cache
	$(GOTEST) -count=1 ./...

watch-test: ## Watch and run tests on file changes (requires entr)
	@which entr > /dev/null || (echo "entr not found. Install it with: brew install entr (macOS) or apt-get install entr (Linux)" && exit 1)
	@echo "Watching for changes... (press Ctrl+C to stop)"
	@find . -name '*.go' | entr -c make quick-test