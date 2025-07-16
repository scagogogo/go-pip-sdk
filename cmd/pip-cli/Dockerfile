# Build stage
FROM golang:1.21-alpine AS builder

# Set build arguments
ARG VERSION=dev

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o pip-cli ./cmd/pip-cli

# Runtime stage
FROM python:3.13-alpine

# Install system dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    bash \
    curl \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1001 -S pipuser && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G pipuser -g pipuser pipuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/pip-cli /usr/local/bin/pip-cli

# Make binary executable
RUN chmod +x /usr/local/bin/pip-cli

# Create directories for user
RUN mkdir -p /app/workspace && \
    chown -R pipuser:pipuser /app

# Switch to non-root user
USER pipuser

# Set environment variables
ENV PYTHONUNBUFFERED=1
ENV PIP_NO_CACHE_DIR=1
ENV PIP_DISABLE_PIP_VERSION_CHECK=1

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pip-cli version || exit 1

# Default command
ENTRYPOINT ["pip-cli"]
CMD ["help"]

# Labels
LABEL maintainer="scagogogo"
LABEL description="Go Pip SDK CLI - Command-line interface for Python package management"
LABEL version="${VERSION}"
LABEL org.opencontainers.image.source="https://github.com/scagogogo/go-pip-sdk"
LABEL org.opencontainers.image.documentation="https://scagogogo.github.io/go-pip-sdk/"
LABEL org.opencontainers.image.licenses="MIT"
