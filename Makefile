# Variables
BINARY_NAME=post-comments
DOCKER_COMPOSE=docker-compose.yml

.PHONY: all build run test docker-build docker-run

all: build

# Build the Go binary
build:
	go build -o ./bin/$(BINARY_NAME) ./cmd/server

# Run the Go application
run:
	go run ./cmd/server/main.go

# Run the tests
test:
	go test ./...

# Build the Docker image
docker-build:
	docker build -t $(BINARY_NAME) .

# Run the Docker container
docker-run:
	docker-compose -f $(DOCKER_COMPOSE) up --build

# Clean the build
clean:
	go clean
	rm -f ./bin/$(BINARY_NAME)