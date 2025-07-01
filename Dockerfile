# 1. Build Stage
FROM golang:1.24 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first (dependency caching)
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of your project files
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# 2. Final Minimal Runtime Stage
FROM debian:bookworm-slim

# Install CA certificates
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the built binary from builder stage
COPY --from=builder /app/main .

# Copy the config folder too (for config files at runtime)
COPY --from=builder /app/config ./config

# Expose application port
EXPOSE 80

# Run the app
CMD ["./main"]
