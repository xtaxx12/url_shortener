# ============================================
# STAGE 1: Build Stage
# ============================================
# Using specific Go version for reproducibility
FROM golang:1.21-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first (better layer caching)
COPY go.mod ./

# Copy source code
COPY . .

# Download and verify dependencies (generates go.sum if missing)
RUN go mod tidy && go mod download && go mod verify

# Build the binary with optimizations
# CGO_ENABLED=0: Build static binary without C dependencies
# -ldflags="-w -s": Strip debug info and symbol table (reduces size)
# -trimpath: Remove file system paths from binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=1.0.0" \
    -trimpath \
    -o /app/server \
    ./cmd/api

# ============================================
# STAGE 2: Production Stage (Minimal Image)
# ============================================
# Using scratch for absolute minimal size
# Alternative: FROM alpine:3.19 (adds ~5MB but includes shell for debugging)
FROM scratch

# Import CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary from builder
COPY --from=builder /app/server /server

# Expose port
EXPOSE 8080

# Set user as non-root (security best practice)
# Note: scratch image doesn't have users, but we document the intention
# In Alpine variant, use: USER 1000:1000

# Health check metadata (used by orchestrators)
# Note: scratch doesn't have curl, health checks done via HTTP probe
LABEL org.opencontainers.image.title="URL Shortener API" \
      org.opencontainers.image.description="High-performance URL shortener service" \
      org.opencontainers.image.version="1.0.0" \
      org.opencontainers.image.vendor="DevOps Portfolio"

# Run the binary
ENTRYPOINT ["/server"]
