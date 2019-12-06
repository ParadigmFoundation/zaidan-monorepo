build:
	$(MAKE) -C ./services/order-book-manager build

test:
	$(MAKE) -C ./services/order-book-manager test
	$(MAKE) -C ./common test

ci:
	$(MAKE) -C ./services/order-book-manager ci
	$(MAKE) -C ./common ci