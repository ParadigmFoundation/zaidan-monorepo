package grpc

//go:generate protoc --go_out=. --proto_path=../../../proto types.proto
//go:generate protoc --go_out=plugins=grpc:. --proto_path=../../../proto services.proto
