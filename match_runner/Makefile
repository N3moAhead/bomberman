# Go parameters
BINARY_NAME=match-runner
CMD_PATH=./cmd/match_runner
VERSION ?= $(shell git describe --tags --always --dirty || echo "dev")

# Docker parameters
DOCKER_IMAGE_NAME=bomberman-match-runner
DOCKER_TAG ?= latest

.PHONY: all build run clean test lint docker-build docker-run help

all: build

# Build the application
build:
	@echo "==> Building $(BINARY_NAME) v$(VERSION)..."
	@go build -ldflags="-X 'main.version=$(VERSION)'" -o $(BINARY_NAME) $(CMD_PATH)

# Run the application locally (directly on my host)
run:
	@echo "==> Running $(BINARY_NAME)..."
	@go run $(CMD_PATH)

# Clean the built artifacts
clean:
	@echo "==> Cleaning..."
	@rm -f $(BINARY_NAME)

# Run tests
test:
	@echo "==> Running tests..."
	@go test ./...

# Run the linter
lint:
	@echo "==> Running linter..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run ./...

# Build the Docker image
docker-build:
	@echo "==> Building Docker image $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)..."
	@docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

# Run the application inside a Docker container
docker-run:
	@echo "==> Running Docker container $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)..."
	@docker run --rm -it \
		--network=host \
		-v /var/run/podman/podman.sock:/var/run/podman/podman.sock \
		-e RABBITMQ_URL=${RABBITMQ_URL} \
		$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# Display help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
