import maker_pb2
import maker_pb2_grpc
import grpc
from concurrent import futures


class MakerServicer(maker_pb2_grpc.MakerServicer):

    def __init__(self):
        pass



def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    maker_pb2_grpc.add_RouteGuideServicer_to_server(
        MakerServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()