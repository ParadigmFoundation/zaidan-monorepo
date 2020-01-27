import os
import grpc
import services_pb2_grpc
from concurrent import futures
from maker_servicer import MakerServicer
from transaction_status_servicer import TransactionStatusServicer

BIND_ADDRESS = os.environ.get("BIND_ADDRESS", "0.0.0.0:50051")

def serve() -> None:
    print('starting server')
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    services_pb2_grpc.add_MakerServicer_to_server(
        MakerServicer(), server)
    services_pb2_grpc.add_TransactionStatusServicer_to_server(
        TransactionStatusServicer(), server)
    server.add_insecure_port(BIND_ADDRESS)
    server.start()
    print('server started')
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
