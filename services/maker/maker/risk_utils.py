from redis_interface import RedisInterface
from pricing_utils import calculate_quote

MAX_PENDING_QUOTE_SIZE = {'ZRX': 1000, 'WETH': 10, 'DAI': 1000}
MAX_PENDING_ORDER_SIZE = {'ZRX': 1000, 'WETH': 10, 'DAI': 1000}

class RiskUtils():

    redis_interface = RedisInterface()

    def risk_checks(self, taker_asset:str, maker_asset:str, taker_size_request:float, sizes:object, test=None) -> object:
        checks = {}
        pending_maker_quote_size = self.redis_interface.get_pending_quote_size(maker_asset)
        pending_maker_order_size = self.redis_interface.get_pending_order_size(maker_asset)

        if taker_size_request:
            reverse_quote = calculate_quote(maker_asset, taker_asset, None, sizes['maker_size'], test)
            if reverse_quote['taker_size'] < taker_size_request:
                checks['crossed_quote_check'] = False
            else:
                checks['crossed_quote_check'] = True
        else:
            reverse_quote = calculate_quote(maker_asset, taker_asset, sizes['taker_size'], None, test)
            if reverse_quote['maker_size'] > sizes['maker_size']:
                checks['crossed_quote_check'] = False
            else:
                checks['crossed_quote_check'] = True

        if pending_maker_quote_size + sizes['maker_size'] < MAX_PENDING_QUOTE_SIZE[maker_asset]:
            checks['pending_quote_size_check'] = True
        else:
            checks['pending_quote_size_check'] = False

        if pending_maker_order_size + sizes['maker_size'] < MAX_PENDING_ORDER_SIZE[maker_asset]:
            checks['pending_order_size_check'] = True
        else:
            checks['pending_order_size_check'] = False

        # todo: Checks to add with other services: inventory manager call for balances, taker profiler call, volatility check

        return checks

    def order_status_update(self, quote_id:str, status:str) -> None:
        if status == 'filled':
            self.redis_interface.filled_order(quote_id)
        else:
            self.redis_interface.failed_order(quote_id)
