VERSION=latest

# ============================================================================ #
# HELPER
# ============================================================================ #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ============================================================================ #
# DEVELOPMENT
# ============================================================================ #

## build: build the application
# TODO: convert to podman
.PHONY: build
build:
	@go build -o=bin/gecho main.go

## clean: remove build artifacts and temporary files
# TODO: convert to podman
clean:
	@rm -r bin/
	
## run: run the application
# TODO: convert to podman
.PHONY: run
run: build
	@./bin/gecho

## up: TODO
.PHONY: up
up:

## down: TODO
.PHONY: down
down:

# ============================================================================ #
# QUALITY CONTROL
# ============================================================================ #

## ci: perform all quality control checks and tests
.PHONY: ci
ci: tidy lint test vuln

## tidy: verify the correct dependencies are used
.PHONY: tidy
tidy:
	@go mod tidy --diff
	@go mod verify

## lint: analyze for stylistic and logic errors
.PHONY: lint
lint: lint-container lint-go

## lint-container: analyze containers for stylistic and logic errors
.PHONY: lint-container
lint-container:
	@if ! command -v podman > /dev/null; then \
		sudo apt-get install -y podman; \
	fi
	@podman run -e 'HOME=/' --rm -i ghcr.io/hadolint/hadolint:latest < Containerfile

## lint-go: analyze code for stylistic and logic errors
.PHONY: lint-go
lint-go:
	@if ! command -v golangci-lint > /dev/null; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0; \
	fi
	@golangci-lint run

## test: perform tests against the application
.PHONY: test
test:
	@go test -vet=all -shuffle=on -race ./...

## vuln: check for known vulnerable dependencies
.PHONY: vuln
vuln:
	@if ! command -v govulncheck > /dev/null; then \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi
	@govulncheck ./...
