# --- Stage 1: Build ---
FROM golang:1.23-alpine AS builder
ENV GOTOOLCHAIN=auto

WORKDIR /app

RUN apk add --no-cache git

# Chỉ copy go.mod và go.sum trước để cache được layer go mod tidy
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o SangXanh ./cmd/api/main.go

# --- Stage 2: Run ---
FROM alpine:latest

WORKDIR /root/

# Copy binary and .env file
COPY --from=builder /app/SangXanh .

EXPOSE 8080

CMD ["./SangXanh"]
