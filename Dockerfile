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
RUN go build -o main ./cmd/app

# 2. Final Minimal Runtime Stage
FROM debian:bookworm-slim

# Set working directory
WORKDIR /app

# Copy the built binary from builder stage
COPY --from=builder /app/main .

# Copy the config folder too (for config files at runtime)
COPY --from=builder /app/config ./config

# Copy the .env file into the container (make sure the .env file exists in your project directory)
COPY --from=builder /app/.env ./.env

# Expose application port
EXPOSE 8080

# Run the app
CMD ["./main"]
