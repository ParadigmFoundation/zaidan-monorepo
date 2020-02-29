DOCKER_IMAGE=maker	
DOCKER_CONTEXT=../..	

include ../../rules/docker.mk	
include ../../rules/rules.mk	

build:
	echo "No build in maker"

ci:
	echo "skipping non-go tests"	

lint:	
	echo "not linting for now (FIXME)" 

gen:
	python -m grpc_tools.protoc -I../../proto --python_out=./maker --grpc_python_out=./maker ../../proto/services.proto
	python -m grpc_tools.protoc -I../../proto --python_out=./maker --grpc_python_out=./maker ../../proto/types.proto
