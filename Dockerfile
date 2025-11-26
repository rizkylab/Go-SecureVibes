# Multi-stage build for optimal image size
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.Version=${VERSION:-dev} -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o gosecvibes \
    cmd/gosecvibes/main.go

# Final stage - minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates git

# Create non-root user
RUN addgroup -g 1000 scanner && \
    adduser -D -u 1000 -G scanner scanner

# Copy binary from builder
COPY --from=builder /build/gosecvibes /usr/local/bin/gosecvibes

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set working directory
WORKDIR /workspace

# Change ownership
RUN chown -R scanner:scanner /workspace

# Switch to non-root user
USER scanner

# Set entrypoint
ENTRYPOINT ["gosecvibes"]

# Default command
CMD ["--help"]

# Labels
LABEL maintainer="rizkylab"
LABEL description="Go-SecureVibes - Comprehensive Security Scanner"
LABEL version="${VERSION:-dev}"
LABEL org.opencontainers.image.source="https://github.com/rizkylab/Go-SecureVibes"
LABEL org.opencontainers.image.description="A comprehensive security scanner with multi-agent architecture"
LABEL org.opencontainers.image.licenses="MIT"
