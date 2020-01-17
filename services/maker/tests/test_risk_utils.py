import risk_utils
from redis_interface import RedisInterface

redis_interface = RedisInterface()

def test_risk_checks():
    redis_interface._set_pending_quote_size('DAI', 0)
    redis_interface._set_pending_quote_size('WETH', 0)
    assert False not in (risk_utils.risk_checks('WETH', 'DAI', 1, {'maker_size': 100, 'taker_size': 1}))
    redis_interface._set_pending_quote_size('DAI', 1500)
    assert not (risk_utils.risk_checks('WETH', 'DAI', 1, {'maker_size': 100, 'taker_size': 1}))['pending_quote_size_check']
    redis_interface._set_pending_order_size('WETH', 15)
    assert not (risk_utils.risk_checks('DAI', 'WETH', 1, {'maker_size': 1, 'taker_size': 150}))['pending_order_size_check']

test_risk_checks()