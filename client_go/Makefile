.DEFAULT_GOAL := run

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o client_go ./cmd/client_go/main.go

run: build
	./client_go

clean:
	rm -f client_go
