import sys
import time
import pickle
from os import getenv
import os
import base64
import ast
import gzip
import json
from zaidan import Logger, DealerCache
from services_pb2_grpc import OrderBookManagerStub
from types_pb2 import OrderBookRequest
import grpc

ENVIRONMENT = getenv("ENVIRONMENT", 'TEST')
PRICE_PAIRS = {'WETH/DAI': {'COINBASE': 'ETH/USD'},
               'ZRX/DAI': {'COINBASE': 'ZRX/USD'}, 'USDC/USDT': {'BINANCE': 'USDC/USDT'},
               'DAI/USDC': {'COINBASE': 'DAI/USDC'}}
MAX_BOOK_AGES = {'ETH/USD': 10, 'ZRX/USD': 15, 'USDC/USDT': 15, 'DAI/USDC': 25}

OBM_CHANNEL = os.environ.get('OBM_CHANNEL', 'localhost:8000')

try:
    SUPPORTED_EXCHANGES = json.loads(os.environ['SUPPORTED_EXCHANGES'])
except:
    SUPPORTED_EXCHANGES = ["BINANCE", "COINBASE"]


class OrderBookException(Exception):
    ''' Exception class for order book errors. '''

    def __init__(self, value):
        ''' Initialize order book exception object. '''
        super(OrderBookException, self).__init__(value)
        self.value = value


class OrderBookWrapper():
    """ DatabaseInterface class. """

    def __init__(self, unit_test: bool = False) -> None:
        if unit_test:
            self.env = 'TEST'
        else:
            self.env = "LIVE"
            self.initialize_obm_connection()

    def initialize_obm_connection(self) -> None:
        self.obm_channel = grpc.insecure_channel(OBM_CHANNEL)
        self.obm_stub = OrderBookManagerStub(self.obm_channel)

    def set_env(self, env: str) -> None:
        self.env = env
        if self.env == 'PLACEHOLDER':
            pass
        else:
            self.initialize_obm_connection()

    def get_current_total_book(self, pair: str, age_check=True):
        """ Get current total book. """
        if self.env == "TEST":
            return self.get_placeholder_book(pair)
        book = {}
        bids_dict = {}
        asks_dict = {}
        for exchange in SUPPORTED_EXCHANGES:
            if exchange in PRICE_PAIRS[pair].keys():
                bids_dict[exchange.upper()] = self.get_order_book(exchange, PRICE_PAIRS[pair][exchange], 'bid')
                asks_dict[exchange.upper()] = self.get_order_book(exchange, PRICE_PAIRS[pair][exchange], 'ask')

        book["bids"] = bids_dict
        book["asks"] = asks_dict

        return book

    def get_placeholder_book(self, symbol):
        """ Get placeholder book. """
        placeholder_book = {}
        if symbol == 'ZRX/WETH':
            placeholder_book['bids'] = {'BINANCE': [[.0098, 100], [.0097, 200], [.0096, 300]]}
            placeholder_book['asks'] = {'BINANCE': [[.0102, 100], [.0103, 200], [.0104, 300]]}
        elif symbol == 'ZRX/DAI':
            placeholder_book['bids'] = {'BINANCE': [[.98, 100], [.97, 200], [.96, 300]],
                                        'COINBASE': [[.99, 100], [.98, 200], [.97, 300]]}
            placeholder_book['asks'] = {'BINANCE': [[1.02, 100], [1.03, 200], [1.04, 300]],
                                        'COINBASE': [[1.01, 100], [1.02, 200], [1.03, 300]]}
        elif symbol == 'USDC/USDT':
            placeholder_book['bids'] = {'BINANCE': [[.99, 100], [.989, 200], [.988, 300]]}
            placeholder_book['asks'] = {'BINANCE': [[.995, 100], [.996, 200], [.997, 300]]}
        elif symbol == 'DAI/USDC' or symbol == 'DAI/USD':
            placeholder_book['bids'] = {'COINBASE': [[1.001, 100], [1.002, 200], [1.003, 300]]}
            placeholder_book['asks'] = {'COINBASE': [[1.003, 100], [1.006, 200], [1.007, 300]]}
        else:
            placeholder_book['bids'] = {'BINANCE': [[98, 100], [97, 200], [96, 300]],
                                        'COINBASE': [[99, 100], [98, 200], [97, 300]]}
            placeholder_book['asks'] = {'BINANCE': [[102, 100], [103, 200], [104, 300]],
                                        'COINBASE': [[101, 100], [102, 200], [103, 300]]}

        return placeholder_book

    def get_order_book(self, exchange: str, symbol: str, side: str, max_age=20) -> list:

        if symbol == 'DAI/USD':
            symbol = 'DAI/USDC'

        # Build the request
        req = OrderBookRequest(exchange=exchange.lower(), symbol=symbol)

        # Call the server
        response = self.obm_stub.OrderBook(req)
        if side == 'bids':
            return [[x.price, x.quantity] for x in response.bids]
        else:
            return [[x.price, x.quantity] for x in response.asks]


