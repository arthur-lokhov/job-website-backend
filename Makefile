# make commands
migrate-up:
	@echo "Applying migrations..."
	goose up

migrate-down:
	@echo "Reverting migrations..."
	goose down

migrate-status:
	@echo "Migration status:"
	goose status

tidy:
	go mod tidy
