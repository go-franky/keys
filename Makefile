GO=go
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOLANGCILINT?=golangci-lint

.PHONY: default
default: help

.PHONY: coverage
coverage: ## Run coverage analysis with go cover
	$(GO) test -race -v -count=1 -coverprofile=cover.out ./...
	$(GO) tool cover -html=cover.out

.PHONY: fmt
fmt: ## Run fmt on go files
	@test -z $(shell $(GO) fmt $$($(GO) list ./... | grep -v /vendor/)) || (echo "Unsucsessfull format - files changed" && exit 1) # This will return non-0 if unsuccessful  run `go fmt ./...` to fix

.PHONY: help
help: ## Show this help
	@echo "Usage: make <target>"
	@echo
	@echo "Targets:"
	@grep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## Run the linter
	$(GOLANGCILINT) run --config .github/golangci.yml -v

.PHONY: test
test: ## Run the test suite
	$(GO) test -v --race -count=1 ./... 2>&1 | tee ${TEST_REPORT};

.PHONY: test_with_localstack
test_with_localstack: ## Run the test suite and spin up localstack in docker to stub AWS
	$(GO) test -v -count=1 -race -cover -tags=localstack \
		github.com/go-franky/keys/aws

.PHONY: vet
vet: ## Run vet on go files
	$(GO) vet ./... > ${VET_REPORT} 2>&1 ;

