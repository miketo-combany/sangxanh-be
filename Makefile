.PHONY: run

run:
	GOOS=linux GOARCH=amd64 go run cmd/api/main.go