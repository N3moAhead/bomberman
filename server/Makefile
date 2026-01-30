.DEFAULT_GOAL := run

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o bomberman-server ./cmd/bomberman-server/main.go

image: vet
	podman build -f Dockerfile.server -t docker.io/nemoahead/bomberman-server .

image-public: image
	podman push docker.io/nemoahead/bomberman-server:latest

build-os: vet
	go build -o bomberman-one-shot-server ./cmd/bomberman-one-shot-server/main.go

image-os: vet
	podman build -f Dockerfile.os-server -t docker.io/nemoahead/bomberman-os-server .

image-os-public: image-os
	podman push docker.io/nemoahead/bomberman-os-server:latest

run: build
	./bomberman-server

run-os: build-one-shot
	./bomberman-one-shot-server

clean:
	rm -f ./bomberman-server
	rm -f ./bomberman-one-shot-server
