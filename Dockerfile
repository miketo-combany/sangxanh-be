# --- Stage 1: Build ---
FROM golang:1.24-alpine AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o SangXanh ./cmd/api/main.go

# --- Stage 2: Run ---
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/SangXanh .

# Ensure the server listens on all interfaces
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
EXPOSE 8080

CMD ["./SangXanh"]

