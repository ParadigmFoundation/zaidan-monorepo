GO_BUILD ?= go build
GO_TEST  ?= go test ./... -race
CGO_ENABLED ?= 0

build: ## Build the binaries
ifeq ($(GO_BINARIES),)
	@echo "GO_BINARIES not defined in Makefile"
else
	@go env -w CGO_ENABLED=$(CGO_ENABLED)
	$(foreach var,$(wildcard $(GO_BINARIES)), $(GO_BUILD) $(var);)
endif

test: $(GO_BEFORE_TEST) run-tests $(GO_AFTER_TEST) ## Run the tests

run-tests:
	CGO_ENABLED=1 $(GO_TEST)

gen: ## Generate required files
	go generate ./...

ci: ## Run the CI
	$(MAKE) run-tests

lint:
	golangci-lint run --new --fast --timeout=5m

goget:
	go get ./...
