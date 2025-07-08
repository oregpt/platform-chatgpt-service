# ChatGPT Service Dockerfile
# Multi-stage build for a smaller final image

# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o chatgpt-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

# Add labels for better maintainability
LABEL maintainer="OreGPT Team" \
      description="ChatGPT Integration Service for OreGPT Agent Platform" \
      version="1.0.0"

# Install CA certificates for HTTPS and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user to run the application
RUN adduser -D -H -h /app appuser

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/chatgpt-service .

# Set ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose the port the service runs on
EXPOSE 8080

# Set health check
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -qO- http://localhost:8080/health || exit 1

# Run the service
CMD ["./chatgpt-service"]
