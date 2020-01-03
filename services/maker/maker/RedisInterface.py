import sys
import time
import pickle
from os import getenv
import os
import base64
import ast
import gzip
import json
from zaidan import Logger

import redis


REDIS_HOST = getenv("REDIS_HOST", default="localhost")
REDIS_PORT = getenv("REDIS_PORT", default="6379")
REDIS_PASSWORD = getenv("REDIS_PASSWORD", default=None)


class RedisInterface():
    """ DatabaseInterface class. """

    logger = Logger("pricer", "debug")

    def __init__(self):
        """ Initialize redis database connection. """
        self.redis_database = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, password=REDIS_PASSWORD)

    def set_pending_quote_size(self, asset, size):
        self.redis_database.set(asset + '_pending_quote_size', str(size))

    def get_pending_quote_size(self, asset):
        return float(self.redis_database.get(asset + '_pending_quote_size'))

    def set_pending_order_size(self, asset, size):
        self.redis_database.set(asset + '_pending_order_size', str(size))

    def get_pending_order_size(self, asset):
        return float(self.redis_database.get(asset + '_pending_order_size'))

    def set_quote(self, quote_id, quote, status=0):

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        order_mark = {'status': status, 'quote': quote}

        mark_compressed = encode_to_bytes(order_mark)
        self.db.hset(self.order_marks_key, quote_id, mark_compressed)

    def update_quote_status(self, quote_id, new_status):

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

    def get_quote(self, quote_id):

        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        if not self.db.hexists(self.order_marks_key, quote_id):
            raise NotFoundError("quote not found", quote_id)

        raw_order_mark = self.db.hget(self.order_marks_key, quote_id)
        order_mark = decode_from_bytes(raw_order_mark)
        return order_mark['quote']

    def add_quote(self, quote_id, quote):
        curr = self.get_pending_quote_size(quote['maker_asset'])
        self.set_pending_quote_size(quote['maker_asset'], curr + quote['maker_asset_size'])
        self.set_quote(quote_id, quote)

    def fill_quote(self, quote_id):
        quote = self.get_quote(quote_id)
        curr = self.get_pending_quote_size(quote['maker_asset'])
        self.set_pending_quote_size(quote['maker_asset'], curr - quote['maker_asset_size'])
        curr = self.get_pending_order_size(quote['maker_asset'])
        self.set_pending_order_size(quote['maker_asset'], curr + quote['maker_asset_size'])
        self.update_quote_status(quote_id, 1)

    def filled_order(self, quote_id):
        quote = self.get_quote(quote_id)
        curr = self.get_pending_order_size(quote['maker_asset'])
        self.set_pending_order_size(quote['maker_asset'], curr - quote['maker_asset_size'])
        self.update_quote_status(quote_id, 2)

    def failed_order(self, quote_id):
        quote = self.get_quote(quote_id)
        curr = self.get_pending_order_size(quote['maker_asset'])
        self.set_pending_order_size(quote['maker_asset'], curr - quote['maker_asset_size'])
        self.update_quote_status(quote_id, 3)

