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

# Build the application binaries - static compilation
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o seed ./cmd/seed

# Production stage - minimal image
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/seed /app/seed

# Expose port
EXPOSE 8080

# Run seed then main server
CMD ["/bin/sh", "-c", "/app/seed && /app/main"]