import maker_pb2
import maker_pb2_grpc
import grpc
from concurrent import futures
from PricingUtils import calculate_price

asset_pricing_data = {'ZRX': {'exchange_books': [('COINBASE', 'ZRX/USD'), ('BINANCE', 'ZRX/ETH')]},
                      'LINK': {'exchange_books': [('COINBASE', 'LINK/USD'), ('BINANCE', 'LINK/ETH')]},
                      'ETH': {'exchange_books': [('COINBASE', 'ETH/USD'), ('BINANCE', 'ETH/USDT'), ('GEMINI', 'ETH/USD')]},
                      'DAI': {'exchange_books': [('COINBASE', 'DAI/USDC')]}
                      }


class MakerServicer(maker_pb2_grpc.MakerServicer):

    def GetQuote(self, request, context):
        price = calculate_price(request.makerAsset, request.takerAsset, request.makerSize, request.takerSize)




def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    maker_pb2_grpc.add_RouteGuideServicer_to_server(
        MakerServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()