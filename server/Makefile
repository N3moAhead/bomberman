DEFAULT_CMD=cd devel && ./podman-compose.sh

.DEFAULT_GOAL := dev

tidy:
	cd src && go fmt ./... 

audit: tidy
	cd src && go vet ./...

dev: audit
	@$(DEFAULT_CMD) up --build;

dev-no-cache: audit
	@$(DEFAULT_CMD) up --build --no-cache;
