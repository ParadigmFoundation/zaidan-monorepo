import unittest
from risk_utils import RiskUtils
from redis_interface import RedisInterface

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

if __name__ == '__main__':
    unittest.main()
