DOCKER_IMAGE=hwd
DOCKER_CONTEXT=../..

GO_BINARIES = ./cmd/hwd

include ../../rules/docker.mk
include ../../rules/golang.mk
include ../../rules/rules.mk

build: export CGO_ENABLED = 0

clean: ## Clean the build objects
	rm -rf hwd
