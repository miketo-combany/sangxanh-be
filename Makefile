.PHONY: run dev build test clean docker-build docker-up docker-down docker-logs docker-restart docker-dev help

# Local development
run:
	go run cmd/api/main.go

dev:
	air -c .air.toml

build:
	go build -o bin/SangXanh ./cmd/api/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/ tmp/

# Docker commands - Production
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app

docker-restart:
	docker compose restart app

docker-clean:
	docker compose down -v
	docker system prune -f

# Docker commands - Development (with hot reload)
docker-dev:
	docker compose --profile dev up -d app-dev

docker-dev-logs:
	docker compose logs -f app-dev

docker-dev-down:
	docker compose --profile dev down

docker-dev-restart:
	docker compose --profile dev restart app-dev

# Combined commands
docker-rebuild:
	docker compose down
	docker compose build --no-cache
	docker compose up -d

docker-rebuild-dev:
	docker compose --profile dev down
	docker compose --profile dev build --no-cache
	docker compose --profile dev up -d app-dev

# Show container status
docker-ps:
	docker compose ps

# Enter container shell
docker-shell:
	docker compose exec app sh

docker-shell-dev:
	docker compose exec app-dev sh

init: 
	go run cmd/sync/main.go
	

# Help
help:
	@echo "Available commands:"
	@echo "  make run              - Run app locally"
	@echo "  make dev              - Run app locally with hot reload (air)"
	@echo "  make build            - Build binary"
	@echo "  make test             - Run tests"
	@echo "  make clean            - Clean build artifacts"
	@echo ""
	@echo "Docker - Production:"
	@echo "  make docker-build     - Build docker image"
	@echo "  make docker-up        - Start containers"
	@echo "  make docker-down      - Stop containers"
	@echo "  make docker-logs      - View logs"
	@echo "  make docker-restart   - Restart app"
	@echo "  make docker-rebuild   - Rebuild and restart"
	@echo "  make docker-clean     - Remove containers and volumes"
	@echo ""
	@echo "Docker - Development:"
	@echo "  make docker-dev       - Start dev container with hot reload"
	@echo "  make docker-dev-logs  - View dev logs"
	@echo "  make docker-dev-down  - Stop dev container"
	@echo "  make docker-dev-restart - Restart dev container"
	@echo "  make docker-rebuild-dev - Rebuild and restart dev"
	@echo ""
	@echo "Utilities:"
	@echo "  make docker-ps        - Show container status"
	@echo "  make docker-shell     - Enter production container"
	@echo "  make docker-shell-dev - Enter dev container"