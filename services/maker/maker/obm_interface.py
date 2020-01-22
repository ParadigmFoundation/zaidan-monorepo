from time import time
import grpc
from services_pb2_grpc import OrderBookManagerStub
from types_pb2 import OrderBookRequest
import os

OBM_CHANNEL = os.environ.get('OBM_CHANNEL', 'localhost:8000')

class DealerCacheError(Exception):
    ''' Signifies an error encountered while interacting with the cache. '''

class NotFoundError(DealerCacheError):
    ''' Signifies an error arising when a requested item does not exist. '''

class OutOfDateError(DealerCacheError):
    ''' Signifies a record is out-of-date according to specified parameters. '''

class OBMInterface():

    env = 'LIVE'

    def __int__(self, unit_test:bool=False) -> None:
        if unit_test:
            self.env = 'PLACEHOLDER'
        else:
            self.initialize_obm_connection()
    
    def initialize_obm_connection(self) -> None:
        self.obm_channel = grpc.insecure_channel(OBM_CHANNEL)
        self.obm_stub = OrderBookManagerStub(self.obm_channel)

    def set_env(self, env:str) -> None:
        self.env = env
        if self.env == 'PLACEHOLDER':
            pass
        else:
            self.initialize_obm_connection()

    def get_order_book(self, exchange: str, symbol: str, side: str, max_age=20) -> list:

        if symbol == 'DAI/USD':
            symbol = 'DAI/USDC'

        # record call time to use for expiration check
        if self.env == 'PLACEHOLDER':
            placeholder_book = {}
            if symbol == 'ZRX/ETH':
                placeholder_book['bids'] = [[.0098, 100], [.0097, 200], [.0096, 300]]
                placeholder_book['asks'] = [[.0102, 100], [.0103, 200], [.0104, 300]]
            elif symbol == 'ZRX/USD':
                placeholder_book['bids'] = [[.29, 100], [.28, 200], [.27, 300]]
                placeholder_book['asks'] = [[.31, 100], [.32, 200], [.33, 300]]
            elif symbol == 'DAI/USDC':
                placeholder_book['bids'] = [[1.001, 100], [1.002, 200], [1.003, 300]]
                placeholder_book['asks'] = [[1.003, 100], [1.006, 200], [1.007, 300]]
            elif symbol == 'ETH/USD':
                placeholder_book['bids'] = [[99, 100], [98, 200], [97, 300]]
                placeholder_book['asks'] = [[101, 100], [102, 200], [103, 300]]
            elif symbol == 'LINK/USD':
                placeholder_book['bids'] = [[1.99, 100], [1.98, 200], [1.97, 300]]
                placeholder_book['asks'] = [[2.01, 100], [2.02, 200], [2.03, 300]]

            return placeholder_book[side]

        else:

            # Build the request
            req = OrderBookRequest(exchange=exchange.lower(), symbol=symbol)

            # Call the server
            response = self.obm_stub.OrderBook(req)
            if side == 'bids':

                return [[x.price, x.quantity] for x in response.bids]
            else:
                return [[x.price, x.quantity] for x in response.asks]
