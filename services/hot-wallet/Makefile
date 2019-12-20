ci: export ETHEREUM_JSONRPC_URL = http://0x-snapshot:8545
ci: run-test

lint:
	golangci-lint run

test: run-test

run-test:
	go test -v ./...

.PHONY: ci lint test run-test
