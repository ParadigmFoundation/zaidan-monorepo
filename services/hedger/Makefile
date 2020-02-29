DOCKER_IMAGE=hedger	
DOCKER_CONTEXT=../..	

include ../../rules/docker.mk	
include ../../rules/rules.mk

build:
	echo "No build in hedger"

ci:	
	echo "skipping non-go tests"	

lint:	
	echo "not linting for now (FIXME)" 

gen:
	python -m grpc_tools.protoc -I../../proto --python_out=./hedger --grpc_python_out=./hedger ../../proto/services.proto
	python -m grpc_tools.protoc -I../../proto --python_out=./hedger --grpc_python_out=./hedger ../../proto/types.proto
