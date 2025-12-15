.PHONY: deps tidy mocks test coverage-html ci
GOBIN ?= $(HOME)/go/bin
COVERAGE_FILE ?= coverage.out
COVERAGE_TARGET_PACKAGE ?= $(shell go list ./... \
| grep -Ev '/(testhelper|mocks|db|config|main)$$' \
| paste -sd, -)
COVERAGE_THRESHOLD ?= 80

deps:
	go mod download

tidy:
	go mod tidy

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

test-tools:
	go install github.com/vektra/mockery/v2@v2.53.5

test: deps test-tools
	go test -coverprofile=$(COVERAGE_FILE) -coverpkg=$(COVERAGE_TARGET_PACKAGE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html

check-coverage:
	@cov=$$(go tool cover -func=$(COVERAGE_FILE) | awk '/^total:/ {gsub(/%/,"",$$3); print $$3}'); \
	echo "coverage: $${cov}% (threshold $(COVERAGE_THRESHOLD)%)"; \
	awk -v c=$${cov} -v t=$(COVERAGE_THRESHOLD) 'BEGIN { exit(c < t) }'

ci: tidy lint test check-coverage
