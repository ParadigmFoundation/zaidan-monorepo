SNAPSHOT_IMAGE_URL=gcr.io/zaidan-eth-net/0x-snapshot:latest

ci: export ETHEREUM_JSONRPC_URL = http://0x-snapshot:8545
ci: run-test

test: ganache-start run-test ganache-stop

run-test:
	go test -v ./...

ganache-start:
	docker run -d --rm -p 8545:8545 -p 8546:8546 --name zrx $(SNAPSHOT_IMAGE_URL)
	sleep 0.5

ganache-stop:
	docker kill zrx

.PHONY: ci test run-test ganache-start ganache-stop
