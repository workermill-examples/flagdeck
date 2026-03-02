# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/server

# Runtime stage
FROM alpine:3.21

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Set working directory
WORKDIR /

# Copy the binary from builder
COPY --from=builder /app/main .

# Change ownership
RUN chown appuser:appuser /main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Run the application
CMD ["/main"]