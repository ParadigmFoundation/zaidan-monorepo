.PHONY = help build gen test ci lint

build: ## Build the binaries
	go build

clean: ## Clean the build objects
	@echo "not implemented"

gen: ## Generate required files
	go generate ./...

test: ## Run the tests
	go test ./... -race

container: ## Build a docker container
	@echo "not implemented"

ci: ## Run the CI
	$(MAKE) test

lint:
	golangci-lint run

help: ## Show this help.
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"
