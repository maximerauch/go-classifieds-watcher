.PHONY: run clean build test fmt

# Run via Docker Compose
run:
	docker-compose up --build

# Run locally (without Docker, needs Go installed)
run-local:
	go run cmd/watcher/main.go

# Clean artifacts
clean:
	docker-compose down --rmi local
	rm -f data/seen.json
	rm -f main

# Format code
fmt:
	go fmt ./...

# Test
test:
	go test ./...