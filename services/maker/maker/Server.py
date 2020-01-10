import uuid
import sys
import math
import time
import types_pb2
import services_pb2_grpc
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


class MakerServicer(services_pb2_grpc.MakerServicer):

    def GetQuote(self, request, context):
        quote_info = calculate_quote(request.maker_asset, request.taker_asset, request.maker_size, request.taker_size)
        expiry_timestamp = int(math.floor(time.time())) + int(VALIDITY_LENGTH)


        if False not in risk_checks(request.taker_asset, request.maker_asset, request.taker_size, quote_info):
            id = str(uuid.uuid4())
            quote = types_pb2.GetQuoteResponse(quote_id=id, expiration=str(expiry_timestamp),
                                               taker_asset=request.taker_asset, maker_asset=request.maker_asset,
                                               taker_size=str(quote_info['taker_size']),
                                               maker_size=str(quote_info['maker_size']))
            redis_interface.add_quote(id, {'expiration':expiry_timestamp,
                                       'taker_asset':request.taker_asset, 'maker_asset':request.maker_asset,
                                       'taker_size':quote_info['taker_size'], 'maker_size':quote_info['maker_size']})
            return quote
        else:
            quote = types_pb2.GetQuoteResponse(status=400)
            return quote

        return quote

    def CheckQuote(self, request, context):
        try:
            quote = redis_interface.get_quote(request.quote_id)
        except ValueError:
            return types_pb2.CheckQuoteResponse(quote_id=request.quote_id, is_valid=False, status=1)

        if quote['expiration'] < time.time():
            return types_pb2.CheckQuoteResponse(quote_id=request.quote_id, is_valid=False, status=2)

        redis_interface.fill_quote(request.quote_id)

        return types_pb2.CheckQuoteResponse(quote_id=request.quote_id, is_valid=True, status=200)

    def OrderStatusUpdate(self, request, context):
        order_status_update(request.quote_id, request.status)
        return types_pb2.OrderStatusUpdateResponse(status=200)


def serve():
    print('starting server')
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    services_pb2_grpc.add_MakerServicer_to_server(
        MakerServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print('server started')
    server.wait_for_termination()


if __name__ == '__main__':
    serve()