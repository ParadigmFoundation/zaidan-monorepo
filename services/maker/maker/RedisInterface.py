import sys
import time
import pickle
from os import getenv
import os
import base64
import ast
import gzip
import json
from common_utils import is_valid_uuid, encode_to_bytes, decode_from_bytes
import redis
from os import getenv

REDIS_HOST = getenv("REDIS_HOST", default="localhost")
REDIS_PORT = getenv("REDIS_PORT", default="6379")
REDIS_PASSWORD = getenv("REDIS_PASSWORD", default=None)

class RedisInterface():

    order_marks_key = "ORDER_MARKS"

    def __init__(self) -> None:
        self.redis_database = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, password=REDIS_PASSWORD)

    def _set_pending_quote_size(self, asset: str, size: float) -> None:
        self.redis_database.set(asset + '_pending_quote_size', str(size))

    def _set_pending_order_size(self, asset: str, size: float) -> None:
        self.redis_database.set(asset + '_pending_order_size', str(size))

    def _set_quote(self, quote_id:str, quote, status=0):

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        order_mark = {'status': status, 'quote': quote}

        mark_compressed = encode_to_bytes(order_mark)
        self.redis_database.hset(self.order_marks_key, quote_id, mark_compressed)

    def _update_quote_status(self, quote_id, new_status):

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        order_mark_raw = self.redis_database.hget(self.order_marks_key, quote_id)
        order_mark = decode_from_bytes(order_mark_raw)

        order_mark['status'] = new_status
        new_mark_raw = encode_to_bytes(order_mark)
        self.redis_database.hset(self.order_marks_key, quote_id, new_mark_raw)

    def get_pending_quote_size(self, asset: str) -> float:
        result = self.redis_database.get(asset + '_pending_quote_size')
        if result:
            return float(result)
        else:
            self._set_pending_quote_size(asset + '_pending_quote_size', '0')
            return 0.0

    def get_pending_order_size(self, asset: str) -> float:
        result = self.redis_database.get(asset + '_pending_order_size')
        if result:
            return float(result)
        else:
            self._set_pending_order_size(asset + '_pending_order_size', '0')
            return 0.0

    def get_quote(self, quote_id):

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)


        raw_order_mark = self.redis_database.hget(self.order_marks_key, quote_id)
        order_mark = decode_from_bytes(raw_order_mark)
        return order_mark['quote']

    def add_quote(self, quote_id, quote):
        curr = self.get_pending_quote_size(quote['maker_asset'])
        self._set_pending_quote_size(quote['maker_asset'], curr + quote['maker_size'])
        self._set_quote(quote_id, quote)

    def fill_quote(self, quote_id):
        quote = self.get_quote(quote_id)
        curr = self.get_pending_quote_size(quote['maker_asset'])
        self._set_pending_quote_size(quote['maker_asset'], curr - quote['maker_size'])
        curr = self.get_pending_order_size(quote['maker_asset'])
        self._set_pending_order_size(quote['maker_asset'], curr + quote['maker_size'])
        self._update_quote_status(quote_id, 1)

    def filled_order(self, quote_id):
        quote = self.get_quote(quote_id)
        curr = self.get_pending_order_size(quote['maker_asset'])
        self._set_pending_order_size(quote['maker_asset'], curr - quote['maker_size'])
        self._update_quote_status(quote_id, 2)

    def failed_order(self, quote_id):
        quote = self.get_quote(quote_id)
        curr = self.get_pending_order_size(quote['maker_asset'])
        self._set_pending_order_size(quote['maker_asset'], curr - quote['maker_size'])
        self._update_quote_status(quote_id, 3)
