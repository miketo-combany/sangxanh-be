# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary for Linux AMD64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api

# Stage 2: Run
FROM alpine:latest

# Set working directory in container
WORKDIR /root/

# Copy the compiled binary from builder
COPY --from=builder /app/app .
COPY .env .

# Expose port if your app listens on one (e.g., 8080)
EXPOSE 8080

# Run the binary
CMD ["./app"]
