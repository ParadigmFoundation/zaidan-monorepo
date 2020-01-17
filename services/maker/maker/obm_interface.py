from sys import path
from time import time

from redis import Redis
import grpc
import json
from services_pb2_grpc import OrderBookManagerStub
from types_pb2 import OrderBookRequest
from base64 import b64encode, b64decode
from gzip import compress, decompress
from json import dumps, loads
from uuid import uuid4, UUID
import os

OBM_CHANNEL = os.environ.get('OBM_CHANNEL', 'localhost:8000')


class DealerCacheError(Exception):
    ''' Signifies an error encountered while interacting with the cache. '''


class NotFoundError(DealerCacheError):
    ''' Signifies an error arising when a requested item does not exist. '''


class OutOfDateError(DealerCacheError):
    ''' Signifies a record is out-of-date according to specified parameters. '''


class OBMInterface():
    '''
    Abstraction over a Redis database for dealer quotes and orders.
    Provides compression/encoding for storing structured data in redis.
    '''

    # redis key for quotes (order mark) hash table
    order_marks_key = "ORDER_MARKS"

    # redis key for per-symbol un-hedged positions hash table
    unhedged_position_key = "UNHEDGED_POSITION"
    unhedged_positions = {}

    obm_channel = grpc.insecure_channel(OBM_CHANNEL)
    obm_stub = OrderBookManagerStub(obm_channel)

    env = 'LIVE'

    def __int__(self, unit_test=False):
        if unit_test:
            self.env = 'PLACEHOLDER'

    def set_env(self, env):
        self.env = env


    def set_order_book(self, exchange: str, symbol: str, side: str, levels: list) -> None:

        updated_timestamp = time()

        if side not in ('bid', 'ask'):
            raise ValueError('side must be "bid" or "ask"')

        symbols = symbol.split('/')
        if len(symbols) != 2:
            raise ValueError('symbol must be BASE_TICKER/QUOTE_TICKER format')

        base_key = f'{symbol.upper()}_{exchange.lower()}_{side}'
        timestamp_key = f'{base_key}_timestamp'

        compressed_book = encode_to_bytes(levels)
        self.db.set(base_key, compressed_book)
        self.db.set(timestamp_key, str(updated_timestamp))

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
