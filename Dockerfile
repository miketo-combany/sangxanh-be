# --- Stage 1: Build ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cài git sớm để go mod tidy không lỗi
RUN apk add --no-cache git

# Chỉ copy go.mod và go.sum trước để cache được layer go mod tidy
COPY go.mod go.sum ./
RUN go mod tidy

# Sau đó mới copy toàn bộ source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o SangXanh ./cmd/api/main.go
