import sys
import time
import pickle
from os import getenv
import os
import base64
import ast
import gzip
import json
from base64 import b64encode, b64decode
from gzip import compress, decompress
from json import dumps, loads
from uuid import uuid4, UUID
import redis
from os import getenv


REDIS_HOST = getenv("REDIS_HOST", default="localhost")
REDIS_PORT = getenv("REDIS_PORT", default="6379")
REDIS_PASSWORD = getenv("REDIS_PASSWORD", default=None)

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


class RedisInterface():
    """ DatabaseInterface class. """

    order_marks_key = "ORDER_MARKS"

    def __init__(self):
        """ Initialize redis database connection. """
        self.redis_database = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, password=REDIS_PASSWORD)

    def set_pending_quote_size(self, asset, size):
        self.redis_database.set(asset + '_pending_quote_size', str(size))

    def get_pending_quote_size(self, asset):
        result = self.redis_database.get(asset + '_pending_quote_size')
        if result:
            return float(result)
        else:
            self.set_pending_quote_size(asset + '_pending_quote_size', '0')
            return 0

    def set_pending_order_size(self, asset, size):
        self.redis_database.set(asset + '_pending_order_size', str(size))

    def get_pending_order_size(self, asset):
        result = self.redis_database.get(asset + '_pending_order_size')
        if result:
            return float(result)
        else:
            self.set_pending_order_size(asset + '_pending_order_size', '0')
            return 0

    def set_quote(self, quote_id, quote, status=0):

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        order_mark = {'status': status, 'quote': quote}

        mark_compressed = encode_to_bytes(order_mark)
        self.redis_database.hset(self.order_marks_key, quote_id, mark_compressed)

    def update_quote_status(self, quote_id, new_status):

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        order_mark_raw = self.redis_database.hget(self.order_marks_key, quote_id)
        order_mark = decode_from_bytes(order_mark_raw)

        order_mark['status'] = new_status
        new_mark_raw = encode_to_bytes(order_mark)
        self.redis_database.hset(self.order_marks_key, quote_id, new_mark_raw)

    def get_quote(self, quote_id):

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)


        raw_order_mark = self.redis_database.hget(self.order_marks_key, quote_id)
        order_mark = decode_from_bytes(raw_order_mark)
        return order_mark['quote']

    def add_quote(self, quote_id, quote):
        curr = self.get_pending_quote_size(quote['maker_asset'])
        self.set_pending_quote_size(quote['maker_asset'], curr + quote['maker_size'])
        self.set_quote(quote_id, quote)

    def fill_quote(self, quote_id):
        quote = self.get_quote(quote_id)
        curr = self.get_pending_quote_size(quote['maker_asset'])
        self.set_pending_quote_size(quote['maker_asset'], curr - quote['maker_size'])
        curr = self.get_pending_order_size(quote['maker_asset'])
        self.set_pending_order_size(quote['maker_asset'], curr + quote['maker_size'])
        self.update_quote_status(quote_id, 1)

    def filled_order(self, quote_id):
        quote = self.get_quote(quote_id)
        curr = self.get_pending_order_size(quote['maker_asset'])
        self.set_pending_order_size(quote['maker_asset'], curr - quote['maker_size'])
        self.update_quote_status(quote_id, 2)

    def failed_order(self, quote_id):
        quote = self.get_quote(quote_id)
        curr = self.get_pending_order_size(quote['maker_asset'])
        self.set_pending_order_size(quote['maker_asset'], curr - quote['maker_size'])
        self.update_quote_status(quote_id, 3)




