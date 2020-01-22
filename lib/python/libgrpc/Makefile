gen: build-types build-services

build-types:
	python -m grpc_tools.protoc -I../../../proto --python_out=. --grpc_python_out=. ../../../proto/types.proto

build-services:
	python -m grpc_tools.protoc -I../../../proto --python_out=. --grpc_python_out=. ../../../proto/services.proto