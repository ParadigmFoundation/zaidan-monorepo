build:
	go build -o watcher cmd/main.go

test:
	go test ./... -race

lint:
	golangci-lint run

ci:
	go test -race ./cmd ./server
