DOCKER_IMAGE=hwd
DOCKER_CONTEXT=../..

GO_BINARIES = ./cmd/hwd

include ../../rules/docker.mk
include ../../rules/golang.mk
include ../../rules/rules.mk

ci: ETHEREUM_JSONRPC_URL=http://0x-snapshot:8545

clean: ## Clean the build objects
	rm -rf hwd
