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

    quote_marks_key = "QUOTE_MARKS"
    pending_quote_marks_key = "PENDING_QUOTE_MARKS"
    pending_order_marks_key = "PENDING_ORDER_MARKS"

    def __init__(self) -> None:
        self.redis_database = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, password=REDIS_PASSWORD)

    def _flush(self) -> None:
        self.redis_database.flushdb()

    def _set_pending_quote_size(self, asset: str, size: float) -> None:
        self.redis_database.hset(self.pending_quote_marks_key, asset, str(size))

    def _set_pending_order_size(self, asset: str, size: float) -> None:
        self.redis_database.hset(self.pending_order_marks_key, asset, str(size))

    def _set_quote(self, quote_id:str, quote:object, status:int = 0) -> None:
        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        order_mark = {'status': status, 'quote': quote}

        mark_compressed = encode_to_bytes(order_mark)
        self.redis_database.hset(self.quote_marks_key, quote_id, mark_compressed)

    def _update_quote_status(self, quote_id:str, new_status: int) -> None:
        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        order_mark_raw = self.redis_database.hget(self.quote_marks_key, quote_id)
        order_mark = decode_from_bytes(order_mark_raw)

        order_mark['status'] = new_status
        new_mark_raw = encode_to_bytes(order_mark)
        self.redis_database.hset(self.quote_marks_key, quote_id, new_mark_raw)

    def _get_quote_status(self, quote_id:str) -> int:
        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        raw_order_mark = self.redis_database.hget(self.quote_marks_key, quote_id)
        order_mark = decode_from_bytes(raw_order_mark)
        return int(order_mark['status'])

    def get_pending_quote_size(self, asset:str) -> float:
        result = self.redis_database.hget(self.pending_quote_marks_key, asset)
        if result:
            return float(result)
        else:
            self._set_pending_quote_size(asset, '0')
            return 0.0

    def get_pending_order_size(self, asset:str) -> float:
        result = self.redis_database.hget(self.pending_order_marks_key, asset)
        if result:
            return float(result)
        else:
            self._set_pending_order_size(asset, '0')
            return 0.0

    def get_quote(self, quote_id:str) -> object:
        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        raw_order_mark = self.redis_database.hget(self.quote_marks_key, quote_id)
        order_mark = decode_from_bytes(raw_order_mark)
        return order_mark['quote']

    def add_quote(self, quote_id:str, quote:object) -> None:
        curr = self.get_pending_quote_size(quote['maker_asset'])
        self._set_pending_quote_size(quote['maker_asset'], curr + quote['maker_size'])
        self._set_quote(quote_id, quote)

    def fill_quote(self, quote_id:str) -> None:
        quote = self.get_quote(quote_id)
        curr = self.get_pending_quote_size(quote['maker_asset'])
        self._set_pending_quote_size(quote['maker_asset'], curr - quote['maker_size'])
        curr = self.get_pending_order_size(quote['maker_asset'])
        self._set_pending_order_size(quote['maker_asset'], curr + quote['maker_size'])
        self._update_quote_status(quote_id, 1)

    def filled_order(self, quote_id:str) -> None:
        quote = self.get_quote(quote_id)
        curr = self.get_pending_order_size(quote['maker_asset'])
        self._set_pending_order_size(quote['maker_asset'], curr - quote['maker_size'])
        self._update_quote_status(quote_id, 2)

    def failed_order(self, quote_id:str) -> None:
        quote = self.get_quote(quote_id)
        curr = self.get_pending_order_size(quote['maker_asset'])
        self._set_pending_order_size(quote['maker_asset'], curr - quote['maker_size'])
        self._update_quote_status(quote_id, 3)
