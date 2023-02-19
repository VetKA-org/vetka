DEFAULT_GOAL := help
APP = cmd/app/app
DATABASE_URI ?= postgres://postgres:postgres@127.0.0.1:5432/vetka?sslmode=disable

help: ## Display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

install-tools: ## Install additional tools
	pre-commit install
.PHONY: install-tools

build: ## Build whole project
	go build -o $(APP) cmd/app/*.go
.PHONY: build

clean: ## Remove build artifacts and downloaded test tools
	rm -f $(APP)
.PHONY: clean

lint: ## Run linters on the source code
	golangci-lint run
.PHONY: lint

unit-tests: ## Run unit tests
	@go test -v -race ./internal/... -coverprofile=coverage.out.tmp -covermode atomic
	@cat coverage.out.tmp | grep -v "_mock.go" > coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
.PHONY: unit-tests

int-tests: build ### Run integration tests
	go clean -testcache && DATABASE_URI=$(DATABASE_URI) SECRET=xxx go test -v ./test/integration/...
.PHONY: int-tests
