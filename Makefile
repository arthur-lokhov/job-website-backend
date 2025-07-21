run:
	go run ./cmd/api/main.go

migrate:
	go run ./cmd/migrate/main.go

docker-build:
	docker build -t job-website-backend .

docker-up:
	docker compose up --build