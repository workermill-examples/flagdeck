# Multi-stage build for FlagDeck API
FROM golang:1.24-alpine AS builder

# Install git for go modules that might need it
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed ./cmd/seed

# Production stage
FROM alpine:3.21

# Install ca-certificates for HTTPS requests and wget for healthcheck
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN addgroup -g 1001 appgroup && adduser -S appuser -u 1001 -G appgroup

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/seed .

# Change ownership to appuser
RUN chown appuser:appgroup main seed

# Switch to non-root user
USER appuser

# Expose port (will be overridden by docker-compose)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT:-8080}/health || exit 1

# Command runs seed first, then main server
CMD ["/bin/sh", "-c", "./seed && ./main"]