.PHONY: dev build test lint docker-build k8s-apply

dev:
	docker-compose up --build

build-backend:
	cd backend && go build -o bin/server ./cmd/server

build-frontend:
	cd frontend && npm run build

test:
	cd backend && go test ./...
	cd frontend && npm run test

lint:
	cd backend && go vet ./...
	cd frontend && npm run lint

docker-build:
	docker build -t trophy-collector-backend ./backend
	docker build -t trophy-collector-frontend ./frontend

k8s-apply:
	kubectl apply -f k8s/