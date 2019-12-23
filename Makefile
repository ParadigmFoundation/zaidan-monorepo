TARGETS = ./services/* ./lib/go

.PHONY: build
build:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) build &&) true

.PHONY: test
test:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) test &&) true

.PHONY:lint
lint:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) lint &&) true

.PHONY: ci
ci:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) ci &&) true

.PHONY: gen
gen:
	$(foreach var,$(wildcard $(TARGETS)), $(MAKE) -C $(var) gen &&) true
