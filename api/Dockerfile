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
FROM alpine:3.21
RUN apk --no-cache add ca-certificates wget
COPY --from=builder /app/main /main
COPY --from=builder /app/seed /seed
CMD ["/bin/sh", "-c", "./seed && ./main"]