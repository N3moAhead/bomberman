.DEFAULT_GOAL := run

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o bomberman-server ./cmd/bomberman-server/main.go

build-one-shot: vet
	go build -o bomberman-one-shot-server ./cmd/bomberman-one-shot-server/main.go

run: build
	./bomberman-server

run-one-shot: build-one-shot
	./bomberman-one-shot-server

clean:
	rm -f ./bomberman-server
	rm -f ./bomberman-one-shot-server
