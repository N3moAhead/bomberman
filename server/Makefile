.DEFAULT_GOAL := run

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o bomberman-server ./cmd/bomberman-server/main.go

run: build
	./bomberman-server

clean:
	rm -f ./bomberman-server
