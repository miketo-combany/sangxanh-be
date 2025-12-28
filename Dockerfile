# --- Stage 1: Build ---
FROM golang:1.24-alpine AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o SangXanh ./cmd/api/main.go

# --- Stage 2: Runtime ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/SangXanh .

# Cloud Run will set PORT environment variable
ENV SERVER_HOST=0.0.0.0
EXPOSE 8080

CMD ["./SangXanh"]

