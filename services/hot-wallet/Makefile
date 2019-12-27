DOCKER_IMAGE=hwd
DOCKER_CONTEXT=../..

GO_BINARIES = ./cmd/hwd

include ../../rules/docker.mk
include ../../rules/golang.mk
include ../../rules/rules.mk

clean: ## Clean the build objects
	rm -rf hwd
