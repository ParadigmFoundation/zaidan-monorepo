import unittest
import sys

try:
    sys.path.append('../maker')
    import grpc
    from concurrent import futures
    import services_pb2_grpc
    from types_pb2 import GetQuoteRequest, CheckQuoteRequest
    from maker_servicer import MakerServicer
except:
    try:
        sys.path.append('../maker/maker')
        import grpc
        from concurrent import futures
        import services_pb2_grpc
        from types_pb2 import GetQuoteRequest, CheckQuoteRequest
        from maker_servicer import MakerServicer
    except Exception as error:
        raise Exception('failed to import server: {}'.format(error.args))

class test_server(unittest.TestCase):

    server_class = MakerServicer
    port = 50051

    def set_up(self) -> None:
        self.server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        services_pb2_grpc.add_MakerServicer_to_server(self.server_class(test=True), self.server)
        self.server.add_insecure_port(f'[::]:{self.port}')
        self.server.start()
    
    def tear_down(self) -> None:
        self.server.stop(None)
    
    def test_server(self) -> None:
        self.set_up()
        with grpc.insecure_channel(f'localhost:{self.port}') as channel:
            stub = services_pb2_grpc.MakerStub(channel)
            req_1 = GetQuoteRequest(maker_asset="0x0b1ba0af832d7c05fd64161e0db78e85978e8082", taker_asset="0x34d402f14d58e001d8efbe6585051bf9706aa064", maker_size=str(1*10**18))
            response_1 = stub.GetQuote(req_1)
            
            req_2 = CheckQuoteRequest(quote_id='a')
            response_2 = stub.CheckQuote(req_2)
            
            req_3 = CheckQuoteRequest(quote_id=response_1.quote_id)
            response_3 = stub.CheckQuote(req_3)
        
        self.assertEqual(type(response_1.quote_id), str)
        self.assertEqual(response_2.status, 1)
        self.assertEqual(response_3.status, 200)
        self.tear_down()
