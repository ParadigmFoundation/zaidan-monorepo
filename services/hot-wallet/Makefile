DOCKER_IMAGE=hwd
DOCKER_CONTEXT=../..

GO_BINARIES = ./cmd/hwd

ETHEREUM_JSONRPC_URL=http://0x-snapshot:8545

include ../../rules/docker.mk
include ../../rules/golang.mk
include ../../rules/rules.mk

clean: ## Clean the build objects
	rm -rf hwd
