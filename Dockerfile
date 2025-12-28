# --- Stage 1: Build ---
FROM golang:1.24-alpine AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o SangXanh ./cmd/api/main.go
ENTRYPOINT ["./SangXanh"]

