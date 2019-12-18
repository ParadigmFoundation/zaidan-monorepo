TARGETS = ./services/* ./lib/go

.PHONY: build
build:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) build;)

.PHONY: test
test:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) test;)

.PHONY: ci
ci:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) ci;)

.PHONY: gen
gen:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) gen;)
