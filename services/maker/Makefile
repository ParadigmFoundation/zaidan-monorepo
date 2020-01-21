DOCKER_IMAGE=maker	
DOCKER_CONTEXT=../..	

include ../../rules/docker.mk	
include ../../rules/rules.mk	

ci:	
	echo "skipping non-go tests"	

lint:	
	echo "not linting for now (FIXME)" 

gen:
	python -m grpc_tools.protoc -I../../proto --python_out=. --grpc_python_out=. ../../proto/services.proto
	python -m grpc_tools.protoc -I../../proto --python_out=. --grpc_python_out=. ../../proto/types.proto