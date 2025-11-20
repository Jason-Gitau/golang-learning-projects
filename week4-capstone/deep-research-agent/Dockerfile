# Multi-stage Dockerfile for Deep Research Agent
# Build stage: ~800MB (includes full Go toolchain)
# Final stage: ~50MB (only runtime + binary)

# ============================================
# Stage 1: Builder
# ============================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /build

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=1 needed for SQLite
# -ldflags reduces binary size
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w" \
    -o research-agent \
    main.go

# ============================================
# Stage 2: Runtime
# ============================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    sqlite \
    sqlite-libs \
    tzdata \
    curl

# Create app user (security: don't run as root)
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/research-agent .

# Create necessary directories with proper permissions
RUN mkdir -p \
    /app/data \
    /app/uploads/files \
    /app/exports \
    /app/config && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose API port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Environment variables with defaults
ENV GIN_MODE=release \
    PORT=8080 \
    DATABASE_PATH=/app/data/research.db \
    UPLOAD_DIR=/app/uploads/files \
    EXPORT_DIR=/app/exports

# Default command: start API server
CMD ["./research-agent", "serve", "--port", "8080", "--workers", "10"]
