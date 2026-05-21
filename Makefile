.PHONY: test build up down logs

test:
	cd backend && go test ./...
	cd web && npm test -- --run

build:
	docker compose build

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f
