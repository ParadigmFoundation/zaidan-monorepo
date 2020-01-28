gen-go:
	go generate ../lib/go/grpc/generate.go

gen-maker:
	cd ../services/maker && make gen

gen: gen-go gen-maker
