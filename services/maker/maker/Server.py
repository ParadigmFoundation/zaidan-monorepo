import uuid
import math
import time
import maker_pb2
import maker_pb2_grpc
import grpc
from concurrent import futures
from PricingUtils import calculate_quote
from RiskUtils import risk_checks, order_status_update
from RedisInterface import RedisInterface
asset_pricing_data = {'ZRX': {'exchange_books': [('COINBASE', 'ZRX/USD'), ('BINANCE', 'ZRX/ETH')]},
                      'LINK': {'exchange_books': [('COINBASE', 'LINK/USD'), ('BINANCE', 'LINK/ETH')]},
                      'ETH': {'exchange_books': [('COINBASE', 'ETH/USD'), ('BINANCE', 'ETH/USDT'), ('GEMINI', 'ETH/USD')]},
                      'DAI': {'exchange_books': [('COINBASE', 'DAI/USDC')]}
                      }

VALIDITY_LENGTH = 15000

redis_interface = RedisInterface()


class MakerServicer(maker_pb2_grpc.MakerServicer):

    def GetQuote(self, request, context):
        quote_info = calculate_quote(request.maker_asset, request.taker_asset, request.maker_size, request.taker_size)
        expiry_timestamp = int(math.floor(time.time())) + int(VALIDITY_LENGTH)

        if False not in risk_checks(request.taker_asset, request.maker_asset, quote_info):
            quote = maker_pb2.GetQuoteResponse(quote_id=uuid.uuid4(), expiration=expiry_timestamp,
                                               taker_asset=request.taker_asset, maker_asset = request.maker_asset,
                                               taker_size=quote_info['taker_size'], maker_size=quote_info['maker_size'],
                                               status=200)
            redis_interface.add_quote(uuid.uuid4(), {'expiration':expiry_timestamp,
                                       'taker_asset':request.taker_asset, 'maker_asset':request.maker_asset,
                                       'taker_size':quote_info['taker_size'], 'maker_size':quote_info['maker_size']})
            return quote
        else:
            quote = maker_pb2.GetQuoteResponse(status=400)
            return quote

    def CheckQuote(self, request, context):
        quote = redis_interface.get_quote(request.quote_id)
        if quote['expration'] > time.time():
            return maker_pb2.CheckQuoteResponse(quote_id=request.quote_id, is_valid=False, status=1)

        redis_interface.fill_quote(request.quote_id)

        return maker_pb2.CheckQuoteResponse(quote_id=request.quote_id, is_valid=True, status=200)

    def OrderStatusUpdate(self, request, context):
        order_status_update(request.quote_id, request.status)
        return maker_pb2.OrderStatusUpdateResponse(status=200)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    maker_pb2_grpc.add_RouteGuideServicer_to_server(
        MakerServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()