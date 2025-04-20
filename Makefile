.PHONY: build run clean docker docker-run test lint

# Go parameters
BINARY_NAME=server
MAIN_PACKAGE=./cmd/server

# Docker parameters
DOCKER_IMAGE=questionnaire-app
DOCKER_TAG=latest

# Default target
all: test build

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Run the application
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -rf ./data/

# Test the application
test:
	go test -v ./...

# Lint the code
lint:
	go vet ./...
	# Add golangci-lint if available
	which golangci-lint && golangci-lint run || echo "golangci-lint not installed"

# Build Docker image
docker:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run with Docker
docker-run: docker
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Run with Docker Compose
docker-compose-up:
	docker-compose up -d

# Stop Docker Compose
docker-compose-down:
	docker-compose down

# Setup development environment
setup-dev:
	go mod download
	mkdir -p ./data

# Generate sample data
setup-data: build
	mkdir -p ./data
	./$(BINARY_NAME) --port 8080 --data ./data

# Check health
check-health:
	curl -s http://localhost:8080/api/health | grep -q "up" && echo "Service is up" || echo "Service is down"

# Help target
help:
	@echo "Available targets:"
	@echo "  build                - Build the application"
	@echo "  run                  - Run the application"
	@echo "  clean                - Clean build artifacts"
	@echo "  test                 - Run tests"
	@echo "  lint                 - Lint the code"
	@echo "  docker               - Build Docker image"
	@echo "  docker-run           - Run with Docker"
	@echo "  docker-compose-up    - Start with Docker Compose"
	@echo "  docker-compose-down  - Stop Docker Compose"
	@echo "  setup-dev            - Setup development environment"
	@echo "  setup-data           - Generate sample data"
	@echo "  check-health         - Check if service is running"
