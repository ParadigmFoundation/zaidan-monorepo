import redis
from os import getenv
from base64 import b64encode, b64decode
from gzip import compress, decompress
from json import dumps, loads
from uuid import uuid4, UUID

REDIS_HOST = getenv("REDIS_HOST", default="localhost")
REDIS_PORT = getenv("REDIS_PORT", default="6379")
REDIS_PASSWORD = getenv("REDIS_PASSWORD", default=None)

class TestRedisInterface():

    quote_marks_key = "QUOTE_MARKS"
    pending_quote_marks_key = "PENDING_QUOTE_MARKS"
    pending_order_marks_key = "PENDING_ORDER_MARKS"
    unhedged_position_key = "UNHEDGED_POSITION"
    unhedged_positions = {}

    def __init__(self) -> None:
        self.redis_database = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, password=REDIS_PASSWORD)



    def get_quote(self, quote_id:str) -> object:
        if not is_valid_uuid(quote_id):
            raise ValueError("invalid quote ID", quote_id)

        raw_order_mark = self.redis_database.hget(self.quote_marks_key, quote_id)
        order_mark = decode_from_bytes(raw_order_mark)
        return order_mark['quote']

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


def is_valid_uuid(uuid_to_test: str, version: int = 4) -> bool:
    try:
        uuid_obj = UUID(uuid_to_test, version=version)
    except ValueError:
        return False

    return str(uuid_obj) == uuid_to_test

def encode_to_bytes(data: object, str_encoding: str = 'utf-8') -> bytes:
    try:
        data_str = dumps(data)
        data_bytes = bytes(data_str, str_encoding)
        compressed = compress(data_bytes)
        encoded = b64encode(compressed)
        return encoded
    except Exception as error:
        raise Exception('failed to compress data: {}'.format(error.args))

def decode_from_bytes(data: bytes, str_encoding: str = 'utf-8') -> object:
    try:
        decoded = b64decode(data)
        decompressed = decompress(decoded)
        data_str = decompressed.decode(str_encoding)
        return loads(data_str)
    except Exception as error:
        raise Exception('failed to decompress data: {}'.format(error.args))
