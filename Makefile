.PHONY: dev build-backend build-frontend test lint docker-build k8s-apply clean

## Start full local environment with Docker Compose
dev:
	docker-compose up --build

## Build the Go backend binary
build-backend:
	cd backend && go build -o bin/server ./cmd/server

## Build the frontend for production
build-frontend:
	cd frontend && npm run build

## Run all tests
test:
	cd backend && go test -v -race ./...
	cd frontend && npm run test

## Lint all code
lint:
	cd backend && go vet ./...
	cd frontend && npm run lint

## Build Docker images
docker-build:
	docker build -t trophy-collector-backend ./backend
	docker build -t trophy-collector-frontend ./frontend

## Apply all Kubernetes manifests
k8s-apply:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/

## Remove build artifacts
clean:
	rm -rf backend/bin
	rm -rf frontend/dist

## Get a Steam API key
steam-key:
	@echo "Get your Steam API key at: https://steamcommunity.com/dev/apikey"
