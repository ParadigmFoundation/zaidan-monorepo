import unittest
from uuid import uuid4
import sys

try:
    sys.path.append('../maker')
    from risk_utils import RiskUtils
    from redis_interface import RedisInterface
except:
    try:
        sys.path.append('../maker/maker')
        from risk_utils import RiskUtils
        from redis_interface import RedisInterface
    except Exception as error:
        raise Exception('failed to import risk utils and redis interface: {}'.format(error.args))

class TestRiskUtils(unittest.TestCase):

    def test_risk_checks(self) -> None:
        redis_interface = RedisInterface()
        risk_utils = RiskUtils()
        redis_interface._set_pending_quote_size('DAI', 0)
        redis_interface._set_pending_quote_size('WETH', 0)
        self.assertNotIn(False, risk_utils.risk_checks('WETH', 'DAI', 1, {'maker_size': 100, 'taker_size': 1}, test=True))
        redis_interface._set_pending_quote_size('DAI', 1500)
        self.assertEqual(False, risk_utils.risk_checks('WETH', 'DAI', 1, {'maker_size': 100, 'taker_size': 1}, test=True)['pending_quote_size_check'])
        redis_interface._set_pending_order_size('WETH', 15)
        self.assertEqual(False, risk_utils.risk_checks('DAI', 'WETH', 1, {'maker_size': 1, 'taker_size': 150}, test=True)['pending_order_size_check'])
    
    def test_order_status_update(self) -> None:
        redis_interface = RedisInterface()
        redis_interface._flush()
        risk_utils = RiskUtils()
        uuid_str_1 = str(uuid4())
        uuid_str_2 = str(uuid4())
        test_quote = {'maker_asset': 'WETH', 'maker_size': 5.0}
        redis_interface.add_quote(uuid_str_1, test_quote)
        redis_interface.add_quote(uuid_str_2, test_quote)
        risk_utils.order_status_update(uuid_str_1, "filled")
        risk_utils.order_status_update(uuid_str_2, "failed")
        self.assertEqual(redis_interface._get_quote_status(uuid_str_1), 2)
        self.assertEqual(redis_interface._get_quote_status(uuid_str_2), 3)

if __name__ == '__main__':
    unittest.main()
