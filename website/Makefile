CSS_INPUT = ./static/css/app.css
CSS_OUTPUT = ./static/css/dist.css

.PHONY: fmt vet templ tailwind build run dev

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

templ:
	go tool templ generate

tailwind:
	npx @tailwindcss/cli -i $(CSS_INPUT) -o $(CSS_OUTPUT) --minify

build: templ tailwind vet
	go build -o bombahead ./cmd/website/main.go

run: build
	./bombahead

image-website:
	podman build -f Dockerfile.Website -t docker.io/nemoahead/bomberman-website

image-website-public: image-website
	podman push docker.io/nemoahead/bomberman-website:latest

image-matchmaker:
	podman build -f Dockerfile.Matchmaker -t docker.io/nemoahead/bomberman-matchmaker

image-matchmaker-public: image-matchmaker
	podman push docker.io/nemoahead/bomberman-matchmaker:latest

dev:
	@if ! command -v air > /dev/null 2>&1; then \
		echo "air is not installed. Installing..."; \
		go install github.com/air-verse/air@latest; \
	fi
	air
