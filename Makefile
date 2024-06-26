.PHONY: all run docker-run

all: docker-run run

# Run the Go application
run:
	go run ./cmd/server/main.go

# Run the Docker container
docker-run:
	docker-compose up -d
