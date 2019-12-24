DOCKER_CONTEXT = ../..
DOCKER_IMAGE = dealer

GO_BINARIES = ./cmd/dealer/

include ../../rules/rules.mk
include ../../rules/golang.mk
include ../../rules/docker.mk

build: export CGO_ENABLED = 0

clean: ## Clean the build objects
	rm -r dealer
