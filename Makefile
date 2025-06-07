help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

lint: ## Run static tests
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

test: ## Run unit tests
	go run github.com/rakyll/gotest@latest -v $(go list ./... | grep -v /examples/)

race: ## Run race tests
	go run github.com/rakyll/gotest@latest -race -v ./...

coverage: ## Run coverage test
	go run github.com/rakyll/gotest@latest -coverprofile=coverage.txt
	go run github.com/boumenot/gocover-cobertura@latest < coverage.txt > coverage.xml

audit: ## Find vulnerabilities
	go list -json -m all | go run github.com/sonatype-nexus-community/nancy@latest sleuth

outdated: ## Find outdated packages
#	go list -u -m -json -mod=readonly all | go run github.com/psampaz/go-mod-outdated@latest -update -direct
	go list -u -m -json -mod=readonly all | go run github.com/psampaz/go-mod-outdated@latest -direct

weight: ## Analyze weight binary size
	go run github.com/jondot/goweight@latest

# tools: ## Install develop tools
# 	go install github.com/boumenot/gocover-cobertura@latest
# 	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# 	go install github.com/jondot/goweight@latest
#   go install github.com/psampaz/go-mod-outdated@latest
# 	go install github.com/rakyll/gotest@latest
# 	go install github.com/sonatype-nexus-community/nancy@latest