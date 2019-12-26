DOCKER_CONTEXT ?= .
DOCKER_FILE ?= Dockerfile
DOCKER_REGISTRY ?= gcr.io/zaidan-io

docker: ## Build a docker container
ifeq ($(DOCKER_IMAGE),)
	@echo "DOCKER_IMAGE not defined in Makefile"
else
	docker build $(DOCKER_CONTEXT) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE) \
		-f $(DOCKER_FILE)
endif

push: ## Pushes a docker image to the registry
ifeq ($(DOCKER_IMAGE),)
	@echo "DOCKER_IMAGE not defined in Makefile"
else
	docker tag $(DOCKER_REGISTRY)/$(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(or $(DOCKER_TAG), latest)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)
	#docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(or $(DOCKER_TAG), latest)
endif
