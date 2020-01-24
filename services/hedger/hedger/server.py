import os
import grpc
import services_pb2_grpc
from concurrent import futures
from hedger_servicer import HedgerServicer

BIND_ADDRESS = os.environ.get("BIND_ADDRESS", "0.0.0.0:50052")

def serve() -> None:
    print('starting server')
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    services_pb2_grpc.add_HedgerServicer_to_server(
        HedgerServicer(), server)
    server.add_insecure_port(BIND_ADDRESS)
    server.start()
    print('server started')
    server.wait_for_termination()

if __name__ == '__main__':
    serve()