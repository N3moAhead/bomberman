.PHONY: all build build-os clean fmt image image-os image-os-public image-public run run-os vet

.DEFAULT_GOAL := run

# Go parameters
BINARY_NAME := bomberman-server
BINARY_NAME_OS := bomberman-one-shot-server

# Docker parameters
IMAGE_NAME := docker.io/nemoahead/bomberman-server
IMAGE_NAME_OS := docker.io/nemoahead/bomberman-os-server

# Build rules
all: build build-os

build: vet
	go build -o $(BINARY_NAME) ./cmd/bomberman-server/main.go

build-os: vet
	go build -o $(BINARY_NAME_OS) ./cmd/bomberman-one-shot-server/main.go

# Run rules
run: build
	./$(BINARY_NAME)

run-os: build-os
	./$(BINARY_NAME_OS)

# Docker rules
image: vet
	podman build -f Dockerfile.server -t $(IMAGE_NAME) .

image-public: image
	podman push $(IMAGE_NAME):latest

image-os: vet
	podman build -f Dockerfile.os-server -t $(IMAGE_NAME_OS) .

image-os-public: image-os
	podman push $(IMAGE_NAME_OS):latest

# Code quality rules
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

# Cleanup
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME_OS)
