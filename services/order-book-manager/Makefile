.PHONY = help buil gen test ci

build: ## Build the binaries
	CGO_ENABLED=0 go build -o obm    ./cmd/obm/main.go 
	CGO_ENABLED=0 go build -o obm-ui ./cmd/obm-ui/main.go

clean: ## Clean the build objects
	rm -r obm obm-ui

gen: ## Generate required files
	go generate ./...

test: ## Run the tests
	go test ./... -race

container: ## Build a docker container
	docker build . -t zaidan-obm

ci: ## Run the CI
	go test ./... -race

help: ## Show this help.
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"
