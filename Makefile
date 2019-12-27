all:
	$(MAKE) build

SERVICES = ./services/*
LIBS     = ./lib/go
TARGETS = $(SERVICES) $(LIBS)

include ./rules/rules.mk

.PHONY: docker push start
docker:
	$(foreach var,$(wildcard $(SERVICES)), $(MAKE) -C $(var) $@ &&) true

push:
	$(foreach var,$(wildcard $(SERVICES)), $(MAKE) -C $(var) $@ &&) true

start: ## Start will start the zaidan stack locally in development mode
	docker-compose \
		-f ./deploy/docker-compose.yml \
		-f ./deploy/docker-compose.dev.yml \
		up

%:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) $@ &&) true
