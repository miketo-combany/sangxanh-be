# --- Stage 1: Build ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY . .

RUN go build -o SangXanh ./cmd/api/main.go

# --- Stage 2: Run ---
FROM alpine:latest

WORKDIR /root/

# Copy binary and .env file
COPY --from=builder /app/SangXanh .

EXPOSE 8080

CMD ["./SangXanh"]
