import unittest
from uuid import uuid4
from typing import Tuple
import sys

try:
    sys.path.append('../maker')
    from redis_interface import RedisInterface
except:
    try:
        sys.path.append('../maker/maker')
        from redis_interface import RedisInterface
    except Exception as error:
        raise Exception('failed to import redis interface: {}'.format(error.args))

class RedisInterfaceTest(unittest.TestCase):

    def test_redis_connection(self) -> None:
        try:
            redis_interface = RedisInterface()
        except Exception as error:
            raise Exception('failed to connect to redis {}'.format(error.args))

    def test_flush(self) -> None:
        redis_interface = RedisInterface()
        redis_interface.redis_database.set('test_key', 'test_value')
        self.assertEqual(redis_interface.redis_database.get('test_key'), b'test_value')
        redis_interface._flush()
        flushed_value = redis_interface.redis_database.get('test_key')
        self.assertEqual(flushed_value, None)

    def test_set_get_pending_quote_size(self) -> None:
        redis_interface = RedisInterface()
        redis_interface._set_pending_quote_size('WETH', 5)
        self.assertEqual(redis_interface.get_pending_quote_size('WETH'), 5.0)
    
    def test_set_get_pending_order_size(self) -> None:
        redis_interface = RedisInterface()
        redis_interface._set_pending_order_size('WETH', 5)
        self.assertEqual(redis_interface.get_pending_order_size('WETH'), 5.0)
    
    def test_set_get_quote(self) -> None:
        redis_interface = RedisInterface()
        uuid_str = str(uuid4())
        test_obj = {'key_1': 0, 'key_2': 2}
        redis_interface._set_quote(uuid_str, test_obj)
        self.assertEqual(redis_interface.get_quote(uuid_str), test_obj)
    
    def test_update_get_quote_status(self) -> None:
        redis_interface = RedisInterface()
        uuid_str = str(uuid4())
        test_obj = {'key_1': 0, 'key_2': 2}
        redis_interface._set_quote(uuid_str, test_obj)
        redis_interface._update_quote_status(uuid_str, 2)
        self.assertEqual(redis_interface._get_quote_status(uuid_str), 2)

    def get_test_quote_and_id(self) -> Tuple[str, object]:
        uuid_str = str(uuid4())
        test_quote = {'maker_asset': 'WETH', 'maker_size': 5.0}
        return (uuid_str, test_quote)

    def test_add_fill_quote_filled_failed_order(self) -> None:
        redis_interface = RedisInterface()
        redis_interface._flush()
        test_quote_and_id_1 = self.get_test_quote_and_id()
        redis_interface.add_quote(test_quote_and_id_1[0], test_quote_and_id_1[1])
        self.assertEqual(redis_interface.get_quote(test_quote_and_id_1[0]), test_quote_and_id_1[1])
        self.assertEqual(redis_interface.get_pending_quote_size(test_quote_and_id_1[1]["maker_asset"]), 5.0)
        test_quote_and_id_2 = self.get_test_quote_and_id()
        redis_interface.add_quote(test_quote_and_id_2[0], test_quote_and_id_2[1])
        self.assertEqual(redis_interface.get_quote(test_quote_and_id_2[0]), test_quote_and_id_2[1])
        self.assertEqual(redis_interface.get_pending_quote_size(test_quote_and_id_2[1]["maker_asset"]), 10.0)
        redis_interface.fill_quote(test_quote_and_id_1[0])
        self.assertEqual(redis_interface.get_pending_quote_size(test_quote_and_id_1[1]["maker_asset"]), 5.0)
        self.assertEqual(redis_interface.get_pending_order_size(test_quote_and_id_1[1]["maker_asset"]), 5.0)
        self.assertEqual(redis_interface._get_quote_status(test_quote_and_id_1[0]), 1)
        redis_interface.filled_order(test_quote_and_id_1[0])
        self.assertEqual(redis_interface.get_pending_order_size(test_quote_and_id_1[1]["maker_asset"]), 0.0)
        self.assertEqual(redis_interface._get_quote_status(test_quote_and_id_1[0]), 2)
        test_quote_and_id_3 = self.get_test_quote_and_id()
        redis_interface.add_quote(test_quote_and_id_3[0], test_quote_and_id_3[1])
        redis_interface.fill_quote(test_quote_and_id_3[0])
        redis_interface.failed_order(test_quote_and_id_3[0])
        self.assertEqual(redis_interface.get_pending_order_size(test_quote_and_id_3[1]["maker_asset"]), 0.0)
        self.assertEqual(redis_interface._get_quote_status(test_quote_and_id_3[0]), 3)
    
if __name__ == '__main__':
    unittest.main()
