build:
	$(MAKE) -C ./services/order-book-manager build

.PHONY: test
test:
	$(MAKE) -C ./services/order-book-manager test
	$(MAKE) -C ./common test

.PHONY: ci
ci:
	$(MAKE) -C ./services/order-book-manager ci
	$(MAKE) -C ./common ci