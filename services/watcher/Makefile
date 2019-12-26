DOCKER_IMAGE = watcher
DOCKER_CONTEXT = ../..

GO_BINARIES  = ./cmd/watcher

build: export CGO_ENABLED = 0

include ../../rules/docker.mk
include ../../rules/golang.mk
include ../../rules/rules.mk
