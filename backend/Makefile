.PHONY: help generate test-unit test-coverage test-acceptance build clean coverage-html test-verbose ci

# Variables
BINARY_NAME=gtd-todo
BUILD_DIR=bin
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html
COVERAGE_THRESHOLD=80

# Default target
all: ci

# Help target
help:
	@echo "Available targets:"
	@echo "  generate        Generate mock files"
	@echo "  test-unit       Run unit tests with coverage"
	@echo "  test-coverage   Run unit tests and verify 80%+ coverage"
	@echo "  test-acceptance Run acceptance tests"
	@echo "  build          Build the application"
	@echo "  clean          Clean build artifacts"
	@echo "  coverage-html   Generate HTML coverage report"
	@echo "  test-verbose    Run tests with verbose output"
	@echo "  ci             Run full CI pipeline"

# Generate mock files
generate:
	@echo "🔄 Generating mock files..."
	go generate ./...
	@echo "✅ Mock generation completed"

# Run unit tests with coverage (excluding mock files)
test-unit:
	@echo "🧪 Running unit tests with coverage..."
	go test -coverprofile=$(COVERAGE_FILE) ./...
	@echo "🧹 Removing mock files from coverage..."
	@awk '!/mock_.*\.go:/' $(COVERAGE_FILE) > $(COVERAGE_FILE).tmp && mv $(COVERAGE_FILE).tmp $(COVERAGE_FILE)
	@echo "✅ Unit tests completed"

# Verify coverage threshold (excluding mocks)
test-coverage: test-unit
	@echo "📊 Checking coverage threshold ($(COVERAGE_THRESHOLD)%)..."
	@total_coverage=$$(go tool cover -func=$(COVERAGE_FILE) | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
	echo "Overall coverage: $${total_coverage}%"; \
	if [ $$(echo "$${total_coverage} < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage $${total_coverage}% is below $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	else \
		echo "✅ Coverage $${total_coverage}% meets requirement ($(COVERAGE_THRESHOLD)%+)"; \
	fi

# Run acceptance tests
test-acceptance:
	@echo "🎯 Running acceptance tests..."
	go test -tags=acceptance ./tests/acceptance/...
	@echo "✅ Acceptance tests completed"

# Build the application
build:
	@echo "🔨 Building application..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "✅ Clean completed"

# Generate HTML coverage report
coverage-html: test-unit
	@echo "📋 Generating HTML coverage report..."
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "✅ Coverage report generated: $(COVERAGE_HTML)"

# Run tests with verbose output
test-verbose:
	@echo "🔍 Running tests with verbose output..."
	go test -v ./...
	@echo "✅ Verbose tests completed"

# Full CI pipeline
ci: generate test-coverage test-acceptance build
	@echo "🎉 CI pipeline completed successfully!"