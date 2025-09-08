run:
	go run ./cmd/api/main.go

migrate:
	go run ./cmd/migrate/main.go

docker-build:
	docker build -t job-website-backend .

docker-up:
	docker compose up --build

lint:
	golangci-lint run ./...

snyk-scan:
	snyk test --file=go.mod --json > snyk_go.json || true
	snyk test --docker job-website-backend:latest --json > snyk_docker.json || true

install-hooks:
	pre-commit install