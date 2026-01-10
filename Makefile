.PHONY: all build run test clean fmt swag swagfmt coverage coverage-html coverage-summary

APP_NAME=api
BUILD_DIR=bin
CMD_DIR=cmd/api

# Testable packages (excluding api, docs, mocks, seeder, store, dto)
TESTABLE_PKGS := \
	./internal/auth/... \
	./internal/cache/... \
	./internal/config/... \
	./internal/env/... \
	./internal/gateway/... \
	./internal/handler/... \
	./internal/middleware/... \
	./internal/packer/... \
	./internal/response/... \
	./internal/service/... \
	./internal/types/...

all: fmt swagfmt swag test build

fmt:
	go fmt ./...

swagfmt:
	swag fmt -g cmd/api/main.go

swag:
	swag init -g cmd/api/main.go -o internal/docs

test:
	go test -v ./...

# Run tests with coverage for testable packages only
coverage:
	@echo "Running tests with coverage for testable packages..."
	@go test -coverprofile=coverage.out -coverpkg=$(shell echo $(TESTABLE_PKGS) | tr ' ' ',' | sed 's/\.\.\.//g') $(TESTABLE_PKGS)
	@echo ""
	@echo "=== Coverage Summary ==="
	@go tool cover -func=coverage.out | tail -1

# Generate HTML coverage report
coverage-html: coverage
	@go tool cover -html=coverage.out -o coverage.html
	@echo ""
	@echo "Coverage report generated: coverage.html"
	@echo "Open in browser to view detailed coverage"

# Show detailed per-package coverage
coverage-summary:
	@echo "=== Detailed Package Coverage ==="
	@echo ""
	@for pkg in auth cache config env gateway handler middleware packer response service types; do \
		coverage=$$(go test -coverprofile=/tmp/$$pkg.out -coverpkg=./internal/$$pkg ./internal/$$pkg 2>&1 | grep "coverage:" | awk '{print $$5}'); \
		if [ -n "$$coverage" ]; then \
			printf "%-20s %s\n" "$$pkg:" "$$coverage"; \
		fi; \
	done
	@echo ""
	@echo "=== Overall Coverage (testable packages) ==="
	@$(MAKE) -s coverage | tail -1 | awk '{print $$3}'

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)

run: build
	./$(BUILD_DIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
