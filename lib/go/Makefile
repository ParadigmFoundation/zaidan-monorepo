SNAPSHOT_IMAGE_URL=gcr.io/zaidan-eth-net/0x-snapshot:v3

ci: export ETHEREUM_JSONRPC_URL = http://0x-snapshot:8545
ci: run-test

lint:
	golangci-lint run

test: ganache-start run-test ganache-stop

run-test:
	go test -v ./...

ganache-start:
	docker run -d --rm -p 8545:8545 -p 8546:8546 --name zrx $(SNAPSHOT_IMAGE_URL)
	sleep 1

ganache-stop:
	docker kill zrx

gen:
	go generate ./...

.PHONY: ci test run-test ganache-start ganache-stop lint gen
