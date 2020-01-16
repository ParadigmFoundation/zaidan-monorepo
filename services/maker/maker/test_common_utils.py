from uuid import uuid4
from common_utils import is_valid_uuid, encode_to_bytes, decode_from_bytes

import unittest

class CommonUtilsTest(unittest.TestCase):

    def test_is_valid_uuid(self) -> None:
        uuid_str = str(uuid4())
        self.assertEqual(is_valid_uuid(uuid_str), True)
    
    def test_encode_to_bytes(self) -> None:
        test_obj = {'key_1': 0, 'key_2': 2}
        test_bytes_det_str = '6tWyk6tjDdUslIw0FEAs42AbKNaADJGwwcYAAAA'
        self.assertEqual(encode_to_bytes(test_obj).decode("utf-8").split('/')[1], test_bytes_det_str)
    
    def test_decode_from_bytes(self) -> None:
        test_obj = {'key_1': 0, 'key_2': 2}
        test_bytes = b'H4sIAArDH14C/6tWyk6tjDdUslIw0FEAs42AbKNaADJGwwcYAAAA'
        self.assertEqual(decode_from_bytes(test_bytes), test_obj)
    
if __name__ == '__main__':
    unittest.main()
