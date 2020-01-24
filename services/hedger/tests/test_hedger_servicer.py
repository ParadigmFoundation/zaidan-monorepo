import unittest
import sys
import grpc
from concurrent import futures

try:
    sys.path.append('../hedger')
    import services_pb2_grpc
    from types_pb2 import HedgeOrderRequest
    from hedger_servicer import HedgerServicer
except:
    try:
        sys.path.append('../hedger/hedger')
        import services_pb2_grpc
        from types_pb2 import HedgeOrderRequest
        from hedger_servicer import HedgerServicer
    except Exception as error:
        raise Exception('failed to import server: {}'.format(error.args))


class TestHedgerServicer(unittest.TestCase):

    server_class = HedgerServicer
    port = 50052

    def set_up(self) -> None:
        self.server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        services_pb2_grpc.add_HedgerServicer_to_server(self.server_class(), self.server)
        self.server.add_insecure_port(f'[::]:{self.port}')
        self.server.start()
    
    def tear_down(self) -> None:
        self.server.stop(None)
    
    def test_server(self) -> None:
        self.set_up()
        with grpc.insecure_channel(f'localhost:{self.port}') as channel:
            stub = services_pb2_grpc.HedgerStub(channel)
            req = HedgeOrderRequest(id = 'test')
            response = stub.HedgeOrder(req)
        
        self.assertEqual(response.valid, True)
        self.tear_down()