# --- Stage 1: Build ---
FROM golang:1.23-alpine AS builder
ENV GOTOOLCHAIN=auto

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum for dependency caching
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o SangXanh ./cmd/api/main.go

# --- Stage 2: Run ---
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/SangXanh .

# Copy .env file if it exists (optional)
COPY .env* ./ 

EXPOSE 8080

CMD ["./SangXanh"]
