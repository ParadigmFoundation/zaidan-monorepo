DOCKER_CONTEXT = ../..
DOCKER_IMAGE = dealer
CGO_ENABLED=1

GO_BINARIES = ./cmd/dealer/

include ../../rules/rules.mk
include ../../rules/golang.mk
include ../../rules/docker.mk

clean: ## Clean the build objects
	rm -r dealer
