# Go-SecureVibes Makefile
# Provides convenient commands for development, testing, and deployment

.PHONY: help build test clean install fmt lint security docker run coverage deps setup-hooks

# Variables
BINARY_NAME=gosecvibes
VERSION?=dev
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "Go-SecureVibes - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

## build: Build the binary
build:
	@echo "🔨 Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o ${BINARY_NAME} cmd/gosecvibes/main.go
	@echo "✅ Build complete: ${BINARY_NAME}"

## build-all: Build binaries for all platforms
build-all:
	@echo "🔨 Building for all platforms..."
	@mkdir -p dist
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ] && [ "$$arch" = "arm64" ]; then continue; fi; \
			output="dist/${BINARY_NAME}-${VERSION}-$$os-$$arch"; \
			if [ "$$os" = "windows" ]; then output="$$output.exe"; fi; \
			echo "Building $$output..."; \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build ${LDFLAGS} -o $$output cmd/gosecvibes/main.go; \
		done \
	done
	@echo "✅ All builds complete in dist/"

## install: Install the binary to $GOPATH/bin
install:
	@echo "📦 Installing ${BINARY_NAME}..."
	@go install ${LDFLAGS} ./cmd/gosecvibes
	@echo "✅ Installed to $(shell go env GOPATH)/bin/${BINARY_NAME}"

## test: Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out
	@echo ""
	@echo "📊 Coverage report generated: coverage.out"
	@echo "   View HTML report: make coverage-html"

## coverage-html: Generate HTML coverage report
coverage-html: test-coverage
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 HTML coverage report: coverage.html"
	@which open > /dev/null && open coverage.html || echo "Open coverage.html in your browser"

## fmt: Format Go code
fmt:
	@echo "🎨 Formatting code..."
	@gofmt -s -w ${GO_FILES}
	@echo "✅ Code formatted"

## lint: Run linters
lint:
	@echo "🔍 Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "⚠️  golangci-lint not installed"; \
		echo "   Install: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin"; \
		go vet ./...; \
	fi

## vet: Run go vet
vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

## security: Run security scan on itself
security: build
	@echo "🔒 Running security scan..."
	@./${BINARY_NAME} scan . --format both --output self_scan

## deps: Download and verify dependencies
deps:
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod verify
	@go mod tidy
	@echo "✅ Dependencies updated"

## clean: Remove build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f ${BINARY_NAME}
	@rm -rf dist/
	@rm -f coverage.out coverage.html
	@rm -f *.json *.md
	@rm -f self_scan.*
	@echo "✅ Clean complete"

## docker-build: Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t ${BINARY_NAME}:${VERSION} -t ${BINARY_NAME}:latest .
	@echo "✅ Docker image built: ${BINARY_NAME}:${VERSION}"

## docker-run: Run Docker container
docker-run:
	@echo "🐳 Running Docker container..."
	@docker run --rm -v $(PWD):/workspace ${BINARY_NAME}:latest scan /workspace

## docker-compose-up: Start services with docker-compose
docker-compose-up:
	@echo "🐳 Starting services..."
	@docker-compose up --build

## docker-compose-down: Stop services
docker-compose-down:
	@echo "🐳 Stopping services..."
	@docker-compose down

## run: Build and run the scanner on current directory
run: build
	@echo "🚀 Running scanner..."
	@./${BINARY_NAME} scan . --verbose

## run-example: Run scanner on example directory
run-example: build
	@if [ -d "examples" ]; then \
		echo "🚀 Running scanner on examples..."; \
		./${BINARY_NAME} scan examples --verbose; \
	else \
		echo "⚠️  examples directory not found"; \
	fi

## setup-hooks: Install git hooks
setup-hooks:
	@echo "🪝 Setting up git hooks..."
	@git config core.hooksPath .githooks
	@chmod +x .githooks/*
	@echo "✅ Git hooks installed"

## ci: Run CI checks locally
ci: deps fmt vet lint test security
	@echo "✅ All CI checks passed!"

## release: Create a new release (requires VERSION)
release:
	@if [ "${VERSION}" = "dev" ]; then \
		echo "❌ Please specify VERSION (e.g., make release VERSION=v1.0.0)"; \
		exit 1; \
	fi
	@echo "🚀 Creating release ${VERSION}..."
	@git tag -a ${VERSION} -m "Release ${VERSION}"
	@echo "✅ Tag created: ${VERSION}"
	@echo "   Push with: git push origin ${VERSION}"

## benchmark: Run benchmarks
benchmark:
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...

## mod-update: Update all dependencies
mod-update:
	@echo "📦 Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Dependencies updated"

## vulncheck: Check for known vulnerabilities
vulncheck:
	@echo "🔍 Checking for vulnerabilities..."
	@if command -v govulncheck > /dev/null; then \
		govulncheck ./...; \
	else \
		echo "⚠️  govulncheck not installed"; \
		echo "   Install: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi

## all: Build, test, and create release artifacts
all: clean deps fmt vet lint test build-all
	@echo "✅ All tasks complete!"
