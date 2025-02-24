help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run unit tests
	go test -v ./... -count 1

# coverage: ## Run coverage test
# 	go test -v -coverprofile=cover.out -covermode=atomic ./...
# 	go tool cover -html=cover.out -o cover.html

coverage: ## Run coverage test
	go test -coverprofile=coverage.txt

lint: ## Run static tests
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...
