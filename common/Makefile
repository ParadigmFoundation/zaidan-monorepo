SNAPSHOT_IMAGE_URL=gcr.io/zaidan-eth-net/0x-snapshot:latest

test: ganache-start run-test ganache-stop

run-test:
	go test -v ./...

ganache-start:
	docker run -d --rm -p 8545:8545 -p 8546:8546 --name zrx $(SNAPSHOT_IMAGE_URL)
	sleep 0.5

ganache-stop:
	docker kill zrx
