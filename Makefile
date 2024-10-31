.POSIX:

export PATH := $(abspath bin/):${PATH}

.PHONY: all
all: help ## Default target (help)

.PHONY: check-mkcert
check-mkcert: ## Check if mkcert is installed
	@which mkcert > /dev/null || (echo "mkcert is not installed. Please install it from https://github.com/FiloSottile/mkcert" && exit 1)

.PHONY: create-certificate
create-certificate: check-mkcert ## Create a new local certificate if it doesn't exist
# if the certificate does not exist, create a new one
	@if [ ! -f src/cmd/run-server/localhost.pem ]; then \
		echo "* Creating a new certificate"; \
		cd src/cmd/run-server && mkcert -install && mkcert localhost; \
	else \
		echo "* Certificate already exists"; \
	fi

.PHONY: mod-init
mod-init: create-certificate ## Initialize go mod
# if go.mod does not exist, initialize it
	@if [ ! -f src/go.mod ]; then \
		echo "* Initializing go mod"; \
		cd src && go mod init .; \
	else \
		echo "* go.mod already exists"; \
	fi

.PHONY: populate-generated-folder
populate-generated-folder: mod-init ## Populate the generated folder
# if the generated folder exist, remove it
	@if [ -d generated ]; then \
		rm -rf generated; \
	fi
# create the generated folder
	@cd src && go run ./cmd/generate-v1 --destination ../generated

.PHONY: run-server
run-server: populate-generated-folder ## Run the server
	cd src && go run ./cmd/run-server -certificate cmd/run-server/localhost.pem -key cmd/run-server/localhost-key.pem

.PHONY: build-all
build-all: mod-init ## Build all the binaries
	@echo "* Building all binaries"
	@cd src && go build -o ../bin/add-module ./cmd/add-module
	@cd src && go build -o ../bin/add-provider ./cmd/add-provider
	@cd src && go build -o ../bin/bump-versions ./cmd/bump-versions
	@cd src && go build -o ../bin/generate-v0 ./cmd/generate-v1
	@cd src && go build -o ../bin/run-server ./cmd/run-server
	@cd src && go build -o ../bin/validate ./cmd/validate
	@cd src && go build -o ../bin/verify-gpg-key ./cmd/verify-gpg-key

.PHONY: test
test: mod-init ## Run the tests for all the packages
	@cd src && go test ./...

.PHONY: clean
clean: ## Clean generated files and binaries
# remove all bins under src
	@rm -rf bin/
# remove generated folders
	@rm -rf bin generated

.PHONY: fmt
fmt: ## Format the Go code
	@cd src && go fmt ./...

.PHONY: lint
lint: ## Run static analysis on the code
	@cd src && golangci-lint run

.PHONY: vet
vet: ## Examine the code for potential issues
	@cd src && go vet ./...

.PHONY: install-tools
install-tools: ## Install necessary development tools
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: help
help: ## Prints this help message.
	@echo ""
	@echo "Opentofu Registry Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "The available targets for execution are listed below."
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*$$"; OFS = ""} \
    /^# .*$$/ { doc=$$0; sub(/^# /, "", doc); next } \
    /^[a-zA-Z0-9_-]+:.*## .*$$/ { target=$$1; sub(/:$$/, "", target); desc=$$0; sub(/^[^#]*## /, "", desc); if (!seen[target]++) { printf "\033[1m%-30s\033[0m %s\n", target, desc } } \
    /^[a-zA-Z0-9_-]+:.*$$/ { target=$$1; sub(/:$$/, "", target); if (!seen[target]++) { if (doc != "") { printf "\033[1m%-30s\033[0m %s\n", target, doc; doc="" } else { printf "\033[1m%-30s\033[0m\n", target } } }' $(MAKEFILE_LIST)