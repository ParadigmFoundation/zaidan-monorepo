SNAPSHOT_IMAGE_URL=gcr.io/zaidan-eth-net/0x-snapshot:v3

GO_BEFORE_TEST = ganache-start
GO_AFTER_TEST  = ganache-stop

include ../../rules/rules.mk
include ../../rules/golang.mk

.PHONY: ganache-start ganache-stop
ci: export ETHEREUM_JSONRPC_URL = http://0x-snapshot:8545

ganache-start:
	docker run -d --rm -p 8545:8545 -p 8546:8546 --name zrx $(SNAPSHOT_IMAGE_URL)
	sleep 1

ganache-stop:
	docker kill zrx
