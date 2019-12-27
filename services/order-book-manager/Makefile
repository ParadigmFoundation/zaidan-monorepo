DOCKER_CONTEXT = ../..
DOCKER_IMAGE = obm

GO_BINARIES = ./cmd/obm/ ./cmd/obm-ui

include ../../rules/docker.mk
include ../../rules/golang.mk

clean: ## Clean the build objects
	rm -r obm obm-ui
