from RedisInterface import RedisInterface
from PricingUtils import calculate_quote

MAX_PENDING_QUOTE_SIZE = {'ZRX': 1000, 'ETH': 10, 'DAI': 1000}
MAX_PENDING_ORDER_SIZE = {'ZRX': 1000, 'ETH': 10, 'DAI': 1000}

redis_interface = RedisInterface()

def risk_checks(taker_asset, maker_asset, taker_size_request, sizes):
    checks = {}
    pending_maker_quote_size = redis_interface.get_pending_quote_size(maker_asset)
    pending_maker_order_size = redis_interface.get_pending_order_size(maker_asset)


    if taker_size_request:
        taker_size_request = float(taker_size_request)
        reverse_quote = calculate_quote(maker_asset, taker_asset, None, sizes['maker_size'])
        if reverse_quote['taker_size'] < taker_size_request:
            checks['crossed_quote_check'] = False
        else:
            checks['crossed_quote_check'] = True
    else:
        reverse_quote = calculate_quote(maker_asset, taker_asset, sizes['taker_size'], None)
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

    ##Checks to add with other services: inventory manager call for balances, taker profiler call, volatility check

    return checks

def order_status_update(quote_id, status):
    if status == 'filled':
        redis_interface.filled_order(quote_id)
    else:
        redis_interface.failed_order(quote_id)
