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


def is_valid_uuid(uuid_to_test: str, version=4) -> bool:
    """
    Check if uuid_to_test is a valid UUID.

    :param uuid_to_test: The UUID being validated.
    :param version: Optionally specify UUID version (default: 4).
    """
    try:
        uuid_obj = UUID(uuid_to_test, version=version)
    except ValueError:
        return False

    return str(uuid_obj) == uuid_to_test


def encode_to_bytes(data: object, str_encoding='utf-8') -> bytes:
    '''
    Encode and compress structured data to a base-64 encoded bytestring.

    :param data: The structured data to encode and compress.
    '''

    try:
        data_str = dumps(data)
        data_bytes = bytes(data_str, str_encoding)
        compressed = compress(data_bytes)
        encoded = b64encode(compressed)
        return encoded
    except Exception as error:
        raise Exception('failed to compress data: {}'.format(error.args))


def decode_from_bytes(data: bytes, str_encoding='utf-8') -> object:
    '''
    Decode and decompress structured data from a base-64 encoded bytestring.

    :param data: The encoded and compressed structured data as a bytestring.
    '''

    try:
        decoded = b64decode(data)
        decompressed = decompress(decoded)
        data_str = decompressed.decode(str_encoding)
        return loads(data_str)
    except Exception as error:
        raise Exception('failed to decompress data: {}'.format(error.args))



class DealerCacheError(Exception):
    ''' Signifies an error encountered while interacting with the cache. '''


class NotFoundError(DealerCacheError):
    ''' Signifies an error arising when a requested item does not exist. '''


class OutOfDateError(DealerCacheError):
    ''' Signifies a record is out-of-date according to specified parameters. '''


class TestDealerCache():
    '''
    Abstraction over a Redis database for dealer quotes and orders.
    Provides compression/encoding for storing structured data in redis.
    '''

    # redis key for quotes (order mark) hash table
    order_marks_key = "ORDER_MARKS"

    # redis key for per-symbol un-hedged positions hash table
    unhedged_position_key = "UNHEDGED_POSITION"
    unhedged_positions = {}

    def __init__(self, host='::', port=6379, password=None):
        '''
        Create a new DealerCache instance.
        :param host: The hostname of the Redis server.
        :param port: The port of the Redis server (default: 6379)
        :param password: The password for the Redis server (default: None)
        '''

        self.db = Redis(host=host, port=port, password=password)

    def set_unhedged_position(self, symbol: str, size: float) -> None:
        '''
        Set a new value for the per-symbol un-hedged position.
        :param symbol: The market symbol (BASE/QUOTE) format.
        :param size: The floating-point value of the un-hedged position.
        '''

        self.unhedged_positions[symbol] = size

    def get_unhedged_position(self, symbol: str) -> float:
        '''
        Get the value for the per-symbol un-hedged position. If no value is
        present, float(0.0) will be returned.
        :param symbol: The market symbol (BASE/QUOTE) format.
        '''

        return self.unhedged_positions[symbol]

    def set_order_book(self, exchange: str, symbol: str, side: str, levels: list) -> None:
        '''
        Encode, compress, and set an order book.
        Will also update the coresponding timestamp key for the book.
        :param exchange: The name of the exchange to set the book for.
        :param symbol: The currency pair of the book.
        :param side: The side of the market book represents (bid or ask).
        :param levels: The nested list of price levels ([[price_level, qty],]).
        '''

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

    def get_order_book(self, exchange: str, symbol: str, side: str, max_age=20, env='PROD') -> list:
        '''
        Fetch and decode an order book from the cache by exchage/size/side.
        If the book is out-of-date according to the max_age parameter, an
        OutOfDateError is raised. Set a max_age of 0 to skip the age check.
        :param exchange: The name of the exchange hosting the market.
        :param symbol: The currency pair (BASE/QUOTE) of the market to get.
        :param side: The side (bid or ask) of the book to get.
        :param max_age: The maximum age (in seconds) of the book data.
        '''

        if symbol == 'DAI/USD':
            symbol = 'DAI/USDC'

        # record call time to use for expiration check
        if env == 'TEST':
            placeholder_book = {}
            if symbol == 'ZRX/ETH':
                placeholder_book['bids'] = [[.0098, 100], [.0097, 200], [.0096, 300]]
                placeholder_book['asks'] = [[.0102, 100], [.0103, 200], [.0104, 300]]
            elif symbol == 'ZRX/USD':
                placeholder_book['bids'] = [[.29, 100], [.28, 200], [.27, 300]]
                placeholder_book['asks'] = [[.31, 100], [.32, 200], [.33, 300]]
            elif symbol == 'DAI/USD':
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
            # Connect to the OBM server
            channel = grpc.insecure_channel('localhost:8000')
            stub = OrderBookManagerStub(channel)

            # Build the request
            req = OrderBookRequest(exchange=exchange.lower(), symbol=symbol)

            # Call the server
            response = stub.OrderBook(req)
            if side == 'bids':

                return [[x.price, x.quantity] for x in response.bids]
            else:
                return [[x.price, x.quantity] for x in response.asks]

    def set_quote(self, quote_id: str, quote: object, status=0) -> None:
        '''
        Store an order mark object by its quote UUID.
        Used to initially store quotes, and to update their statuses.
        Status codes:
        - 0: Quote generated
        - 1: Validated and submitted for settlement
        - 2: Filled, sent to hedger
        :param quote_id: The quote UUID string.
        :param order_mark: The quote object.
        :param status: The new status of the quote.
        '''

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        order_mark = {'status': status, 'quote': quote}

        mark_compressed = encode_to_bytes(order_mark)
        self.db.hset(self.order_marks_key, quote_id, mark_compressed)

    def update_quote_status(self, quote_id: str, new_status: int) -> None:
        '''
        Update an existing quote's status.
        :param quote_id: The ID of the quote to update.
        :param new_status: The new status code of the quote.
        '''

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        if not self.db.hexists(self.order_marks_key, quote_id):
            raise NotFoundError('quote with specified ID not found')

        order_mark_raw = self.db.hget(self.order_marks_key, quote_id)
        order_mark = decode_from_bytes(order_mark_raw)

        if 'status' not in order_mark:
            raise DealerCacheError('malformed order mark; no known status')

        order_mark['status'] = new_status
        new_mark_raw = encode_to_bytes(order_mark)
        self.db.hset(self.order_marks_key, quote_id, new_mark_raw)

    def get_quote(self, quote_id: str) -> object:
        '''
        Fetch an order mark by it's quote ID (if it exists).
        Will raise a NotFoundError if the quote does not exist, and a CacheError
        if the quote_id is invalid.
        :param quote_id: The quote ID generated on initial request.
        '''

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        if not self.db.hexists(self.order_marks_key, quote_id):
            raise NotFoundError("quote not found", quote_id)

        raw_order_mark = self.db.hget(self.order_marks_key, quote_id)
        order_mark = decode_from_bytes(raw_order_mark)
        return order_mark['quote']

    def get_all_order_marks(self, only_valid=True) -> dict:
        '''
        Fetch a dictionary with all order marks (price quotes and statuses).
        :param only_valid: If set, (default) only non-expired result included.
        '''

        call_time = time()

        order_marks = {}
        raw_marks = self.db.hgetall(self.order_marks_key)

        # iterate over encoded/compressed order_marks and decode
        for mark_id in raw_marks.keys():
            decoded_mark = decode_from_bytes(raw_marks[mark_id])

            # when only_valid is set, only return un-expired quotes
            quote_expiration = int(decoded_mark['quote']['expiration'])
            if only_valid and call_time >= quote_expiration:
                continue

            order_marks[mark_id.decode('utf-8')] = decoded_mark

        return order_marks

    def remove_order_mark(self, quote_id: str) -> None:
        '''
        Remove an order mark by its ID.
        :param quote_id: The UUID included when the quote was generated.
        '''

        self.db.hdel(self.order_marks_key, quote_id)

    def get_quote_ids(self) -> list:
        '''
        Fetch an array of all quote ID's in the cache.
        '''

        return [mark_id.decode() for mark_id in self.db.hkeys(self.order_marks_key)]

