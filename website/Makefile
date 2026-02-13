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

dev:
	@if ! command -v air > /dev/null 2>&1; then \
		echo "air is not installed. Installing..."; \
		go install github.com/air-verse/air@latest; \
	fi
	air
