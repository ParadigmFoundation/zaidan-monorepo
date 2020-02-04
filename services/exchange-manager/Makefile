DOCKER_CONTEXT = ../..
DOCKER_IMAGE = em

GO_BINARIES = ./cmd/em
GO_TEST = go test -v ./... -race

include ../../rules/docker.mk
include ../../rules/golang.mk

clean: ## Clean the build objects
	rm -r em
