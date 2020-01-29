all:
	$(MAKE) build

GO_SERVICES = ./services/dealer \
			  ./services/exchange-manager \
			  ./services/hot-wallet \
			  ./services/order-book-manager \
			  ./services/watcher

PY_SERVICES = ./services/hedger \
			  ./services/maker

SERVICES = $(GO_SERVICES) $(PY_SERVICES)
LIBS     = ./lib/go
TARGETS = $(SERVICES) $(LIBS)

include ./rules/rules.mk

.PHONY: docker push start stop start-foreground
docker:
	$(foreach var,$(wildcard $(SERVICES)), $(MAKE) -C $(var) $@ &&) true

push:
	$(foreach var,$(wildcard $(SERVICES)), $(MAKE) -C $(var) $@ &&) true

start: ## Start will start the zaidan stack locally in development mode
	docker-compose \
		-f ./deploy/docker-compose.yml \
		-f ./deploy/docker-compose.dev.yml \
		up -d

start-foreground: ## Start will start the zaidan stack locally in development mode
	docker-compose \
		-f ./deploy/docker-compose.yml \
		-f ./deploy/docker-compose.dev.yml \
		up

stop: ## Stop will will the zaidan stack locally in development mode
	docker-compose \
		-f ./deploy/docker-compose.yml \
		-f ./deploy/docker-compose.dev.yml \
		down -v

goget:
	$(foreach var,$(wildcard $(GO_SERVICES)), $(MAKE) -C $(var) $@ &&) true

%:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) $@ &&) true
