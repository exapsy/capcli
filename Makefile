BINARY_NAME_CLI=capcli
BINARY_NAME_CLI_GUI=capgcli

.PHONY: vendor
vendor:
	@echo "Installing dependencies..."
	@go mod tidy && go mod vendor
	@echo "✓ Done installing dependencies"

.PHONY: build
cli-build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME_CLI) -v ./cmd/cli/main.go
	@echo "✓ Done building"

.PHONY: run
cli-run:
	@echo "Running..."
	@go run ./cmd/cli/main.go
	@echo "✓ Done running"

