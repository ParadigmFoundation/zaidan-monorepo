build:
	$(MAKE) -C ./services/dealer build
	$(MAKE) -C ./services/order-book-manager build

.PHONY: test
test:
	$(MAKE) -C ./services/dealer test
	$(MAKE) -C ./services/hot-wallet test
	$(MAKE) -C ./services/order-book-manager test
	$(MAKE) -C ./lib/go test

.PHONY: ci
ci:
	$(MAKE) -C ./services/dealer ci
	$(MAKE) -C ./services/hot-wallet ci
	$(MAKE) -C ./services/order-book-manager ci
	$(MAKE) -C ./lib/go ci
