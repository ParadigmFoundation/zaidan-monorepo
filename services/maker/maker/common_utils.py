from base64 import b64encode, b64decode
from gzip import compress, decompress
from json import dumps, loads
from uuid import uuid4, UUID

def is_valid_uuid(uuid_to_test: str, version=4) -> bool:
    try:
        uuid_obj = UUID(uuid_to_test, version=version)
    except ValueError:
        return False

    return str(uuid_obj) == uuid_to_test

def encode_to_bytes(data: object, str_encoding='utf-8') -> bytes:
    try:
        data_str = dumps(data)
        data_bytes = bytes(data_str, str_encoding)
        compressed = compress(data_bytes)
        encoded = b64encode(compressed)
        return encoded
    except Exception as error:
        raise Exception('failed to compress data: {}'.format(error.args))

def decode_from_bytes(data: bytes, str_encoding='utf-8') -> object:
    try:
        decoded = b64decode(data)
        decompressed = decompress(decoded)
        data_str = decompressed.decode(str_encoding)
        return loads(data_str)
    except Exception as error:
        raise Exception('failed to decompress data: {}'.format(error.args))
