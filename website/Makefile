CSS_INPUT = ./static/css/app.css
CSS_OUTPUT = ./static/css/dist.css

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
