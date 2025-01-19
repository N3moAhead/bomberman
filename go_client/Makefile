DEFAULT_CMD=cd devel && ./podman-compose.sh

.DEFAULT_GOAL := run

tidy:
	cd src && go fmt ./...

audit: tidy
	cd src && go vet ./...

run: audit
	go run ./src/cmd/client/main.go

# Does sadly not work correctly at the moment
dev: audit
	@$(DEFAULT_CMD) up --build;

dev-no-cache: audit
	@$(DEFAULT_CMD) up --build --no-cache;
