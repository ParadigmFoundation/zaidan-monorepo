from TestDealerCache import TestDealerCache
from RedisInterface import RedisInterface
import math

asset_pricing_data = {'ZRX': {'exchange_books': [('COINBASE', 'ZRX/USD'), ('BINANCE', 'ZRX/ETH')], 'implied_pref': ('COINBASE', 'ZRX/USD')},
                      'LINK': {'exchange_books': [('COINBASE', 'LINK/USD'), ('BINANCE', 'LINK/ETH')], 'implied_pref': ('COINBASE', 'LINK/USD')},
                      'ETH': {'exchange_books': [('COINBASE', 'ETH/USD'), ('BINANCE', 'ETH/USDT')], 'implied_pref': ('COINBASE', 'ETH/USD')},
                      'DAI': {'exchange_books': [('COINBASE', 'DAI/USD')], 'implied_pref': ('COINBASE', 'DAI/USD'), 'constant_rate': 'PREF_INSIDE'}
                      }

EXCHANGE_FEES = {'BINANCE':.00075, 'COINBASE':.002, 'GEMINI':.001}
PREMIUM = .001
cache = TestDealerCache()
redis_interface = RedisInterface()

def calculate_quote(maker_asset, taker_asset, maker_size=None, taker_size=None):
    #pending_maker_size = redis_interface.get_pending_quote_size(maker_asset)
    #pending_taker_size = redis_interface.get_pending_quote_size(taker_asset)
    maker_asset_pricing_data = asset_pricing_data[maker_asset]
    taker_asset_pricing_data = asset_pricing_data[taker_asset]
    order_books = []
    implied = True
    if maker_size:
        maker_size = float(maker_size)
    if taker_size:
        taker_size = float(taker_size)
    for market in maker_asset_pricing_data['exchange_books']:
        if maker_asset + '/' + taker_asset == market[1]:
            order_books.append(maker_asset + '/' + taker_asset)
            implied = False
        if taker_asset + '/' + maker_asset == market[1]:
            order_books.append(taker_asset + '/' + maker_asset)
            implied = False

    if implied:
        if maker_asset_pricing_data['implied_pref'][1].split('/')[1] == taker_asset_pricing_data['implied_pref'][1].split('/')[1]:
            if 'constant_rate' in maker_asset_pricing_data.keys():
                books_to_get = taker_asset_pricing_data['exchange_books']
                order_books = {}
                for book in books_to_get:
                    if book[1].split('/')[1] == maker_asset_pricing_data['implied_pref'][1].split('/')[1]:
                        order_books[book[0]] = (cache.get_order_book(book[0], book[1], 'bids'))
                if taker_size:
                    price = _get_price_from_book_base(order_books, taker_size, 'buy')
                    price = adjust_for_constant_rate(price, maker_asset_pricing_data['implied_pref'], 'maker_asset', side_spef='taker')
                    return {'maker_size': taker_size * price, 'taker_size': taker_size}
                else:
                    price = _get_price_from_book_quote(order_books, maker_size, 'buy')
                    price = adjust_for_constant_rate(price, maker_asset_pricing_data['implied_pref'], 'maker_asset', side_spef='maker')
                    return {'maker_size': maker_size, 'taker_size': maker_size/price}
            elif 'constant_rate' in taker_asset_pricing_data.keys():
                books_to_get = maker_asset_pricing_data['exchange_books']
                order_books = {}
                for book in books_to_get:
                    if book[1].split('/')[1] == maker_asset_pricing_data['implied_pref'][1].split('/')[1]:
                        order_books[book[0]] = (cache.get_order_book(book[0], book[1], 'asks'))
                if maker_size:
                    price = _get_price_from_book_base(order_books, maker_size, 'sell')
                    price = adjust_for_constant_rate(price, taker_asset_pricing_data['implied_pref'], 'taker_asset', side_spef='maker')
                    return {'maker_size': maker_size, 'taker_size': maker_size*price}

                else:
                    price = _get_price_from_book_quote(order_books, taker_size, 'sell')
                    price = adjust_for_constant_rate(price, taker_asset_pricing_data['implied_pref'], 'taker_asset', side_spef='taker')
                    return {'maker_size': taker_size/price, 'taker_size': taker_size}
            else:
                if taker_size:
                    imp_book = taker_asset_pricing_data['implied_pref']
                    order_books = {}
                    order_books[imp_book[0]] = (cache.get_order_book(imp_book[0], imp_book[1], 'bids'))

                    taker_asset_price = _get_price_from_book_base(order_books, taker_size, 'buy')

                    imp_book = maker_asset_pricing_data['implied_pref']
                    order_books = {}
                    order_books[imp_book[0]] = (cache.get_order_book(imp_book[0], imp_book[1], 'asks'))

                    maker_asset_price = _get_price_from_book_quote(order_books, taker_size*taker_asset_price, 'sell')

                    return {'maker_size': taker_size*taker_asset_price/maker_asset_price, 'taker_size':taker_size}
                else:
                    imp_book = maker_asset_pricing_data['implied_pref']
                    order_books = {}
                    order_books[imp_book[0]] = (cache.get_order_book(imp_book[0], imp_book[1], 'asks'))

                    maker_asset_price = _get_price_from_book_base(order_books, maker_size, 'sell')

                    imp_book = taker_asset_pricing_data['implied_pref']
                    order_books = {}
                    order_books[imp_book[0]] = (cache.get_order_book(imp_book[0], imp_book[1], 'bids'))

                    taker_asset_price = _get_price_from_book_quote(order_books, maker_size*maker_asset_price, 'sell')

                    return {'maker_size': maker_size, 'taker_size':maker_size*maker_asset_price/taker_asset_price}


        else:
            raise Exception


def _get_price_from_book_base(half_book, size, side):
    ''' Return price. '''


    if side == 'buy':
        inside_bids = {}
        exchange_levels = {}
        total_inside_bid = [float('-Inf'), float('-Inf')]
        exchange_sizes = {}
        exchange_prices = {}
        exchange_exhausted = {}

        for exchange in half_book.keys():
            exchange_exhausted[exchange] = False
            exchange_levels[exchange] = 0
            exchange_sizes[exchange] = 0
            inside_bids[exchange] = half_book[exchange][0]
            inside_bids[exchange][0] = inside_bids[exchange][0]
            if inside_bids[exchange][0] * (1 - EXCHANGE_FEES[exchange]) > total_inside_bid[0] * (
                    1 - EXCHANGE_FEES[exchange]):
                total_inside_bid = inside_bids[exchange]
                total_inside_bid_exchange = exchange

        while size > 0:
            if total_inside_bid[1] > size:
                exchange_sizes[total_inside_bid_exchange] = exchange_sizes.get(total_inside_bid_exchange, 0) + size
                exchange_prices[total_inside_bid_exchange] = total_inside_bid[0] * (
                            1 - EXCHANGE_FEES[total_inside_bid_exchange])
                size = 0
            else:
                size = size - total_inside_bid[1]
                exchange_sizes[total_inside_bid_exchange] = \
                    exchange_sizes.get(total_inside_bid_exchange, 0) + total_inside_bid[1]
                exchange_prices[total_inside_bid_exchange] = total_inside_bid[0] * (
                            1 - EXCHANGE_FEES[total_inside_bid_exchange])
                exchange_levels[total_inside_bid_exchange] = exchange_levels[total_inside_bid_exchange] + 1

            total_inside_bid = [float('-Inf'), float('-Inf')]
            for exchange in half_book.keys():
                if not exchange_levels[exchange] >= len(half_book[exchange]):
                    inside_bids[exchange] = half_book[exchange][exchange_levels[exchange]]
                    inside_bids[exchange][0] = inside_bids[exchange][0]
                    if inside_bids[exchange][0] * (1 - EXCHANGE_FEES[exchange]) > total_inside_bid[0] * (
                            1 - EXCHANGE_FEES[exchange]):
                        total_inside_bid = inside_bids[exchange]
                        total_inside_bid_exchange = exchange
                else:
                    exchange_exhausted[exchange] = True
            if False not in exchange_exhausted.values():
                exchange_sizes[half_book.keys()[0]] = exchange_sizes[half_book.keys()[0]] + size
                size = 0

        return min(exchange_prices.values())*(1 - PREMIUM)

    if side == 'sell':
        inside_asks = {}
        exchange_levels = {}
        total_inside_ask = [float('Inf'), float('Inf')]
        exchange_sizes = {}
        exchange_prices = {}
        exchange_exhausted = {}


        for exchange in half_book.keys():
            exchange_exhausted[exchange] = False
            exchange_levels[exchange] = 0
            exchange_sizes[exchange] = 0
            inside_asks[exchange] = half_book[exchange][0]
            inside_asks[exchange][0] = inside_asks[exchange][0]
            if inside_asks[exchange][0] * (1 + EXCHANGE_FEES[exchange]) < total_inside_ask[0] * (
                    1 + EXCHANGE_FEES[exchange]):
                total_inside_ask = inside_asks[exchange]
                total_inside_ask_exchange = exchange

        while size > 0:
            if total_inside_ask[1] > size:
                exchange_sizes[total_inside_ask_exchange] = exchange_sizes.get(total_inside_ask_exchange, 0) + size
                exchange_prices[total_inside_ask_exchange] = total_inside_ask[0] * (
                            1 + EXCHANGE_FEES[total_inside_ask_exchange])
                size = 0
            else:
                size = size - total_inside_ask[1]
                exchange_sizes[total_inside_ask_exchange] = \
                    exchange_sizes.get(total_inside_ask_exchange, 0) + total_inside_ask[1]
                exchange_prices[total_inside_ask_exchange] = total_inside_ask[0] * (
                            1 + EXCHANGE_FEES[total_inside_ask_exchange])
                exchange_levels[total_inside_ask_exchange] = exchange_levels[total_inside_ask_exchange] + 1

            total_inside_ask = [float('Inf'), float('Inf')]

            for exchange in half_book.keys():
                if not exchange_levels[exchange] >= len(half_book[exchange]):
                    inside_asks[exchange] = half_book[exchange][exchange_levels[exchange]]
                    inside_asks[exchange][0] = inside_asks[exchange][0]
                    if inside_asks[exchange][0] * (1 + EXCHANGE_FEES[exchange]) < total_inside_ask[0] * (
                            1 + EXCHANGE_FEES[exchange]):
                        total_inside_ask = inside_asks[exchange]
                        total_inside_ask_exchange = exchange
                else:
                    exchange_exhausted[exchange] = True
            if False not in exchange_exhausted.values():
                exchange_sizes[half_book.keys()[0]] = exchange_sizes[half_book.keys()[0]] + size
                size = 0
        return max(exchange_prices.values())*(1 + PREMIUM)

    raise Exception('Error in calculating price')

def _get_price_from_book_quote(half_book, size, side):
    ''' Return price. '''


    exchange_levels = {}
    exchange_sizes = {}
    exchange_prices = {}
    exchange_exhausted = {}

    if side == 'sell':
        inside_asks = {}
        total_inside_ask = [float('Inf'), float('Inf')]
        for exchange in half_book.keys():
            exchange_exhausted[exchange] = False
            exchange_levels[exchange] = 0
            exchange_sizes[exchange] = 0
            inside_asks[exchange] = half_book[exchange][0]
            inside_asks[exchange][0] = inside_asks[exchange][0]
            if inside_asks[exchange][0] * (1 + EXCHANGE_FEES[exchange]) < total_inside_ask[0] * (
                    1 + EXCHANGE_FEES[exchange]):
                total_inside_ask = inside_asks[exchange]
                total_inside_ask_exchange = exchange

        while size > 0:
            total_available = total_inside_ask[0] * total_inside_ask[1]
            if total_available > size:
                exchange_sizes[total_inside_ask_exchange] = exchange_sizes.get(total_inside_ask_exchange, 0) + size
                exchange_prices[total_inside_ask_exchange] = total_inside_ask[0] * (
                            1 + EXCHANGE_FEES[total_inside_ask_exchange])
                size = 0
            else:
                size = size - total_available
                exchange_sizes[total_inside_ask_exchange] = \
                    exchange_sizes.get(total_inside_ask_exchange, 0) + total_available
                exchange_prices[total_inside_ask_exchange] = total_inside_ask[0] * (
                            1 + EXCHANGE_FEES[total_inside_ask_exchange])
                exchange_levels[total_inside_ask_exchange] = exchange_levels[total_inside_ask_exchange] + 1

            total_inside_ask = [float('Inf'), float('Inf')]
            for exchange in half_book.keys():
                if not exchange_levels[exchange] >= len(half_book[exchange]):
                    inside_asks[exchange] = half_book[exchange][exchange_levels[exchange]]
                    inside_asks[exchange][0] = inside_asks[exchange][0]
                    if inside_asks[exchange][0] * (1 + EXCHANGE_FEES[exchange]) < total_inside_ask[0] * (
                            1 + EXCHANGE_FEES[exchange]):
                        total_inside_ask = inside_asks[exchange]
                        total_inside_ask_exchange = exchange
                else:
                    exchange_exhausted[exchange] = True

            if False not in exchange_exhausted.values():
                exchange_sizes[half_book.keys()[0]] = exchange_sizes[half_book.keys()[0]] + size
                size = 0

        return max(exchange_prices.values()) * (1 + PREMIUM)

    if side == 'buy':
        inside_bids = {}
        total_inside_bid = [float('-Inf'), float('-Inf')]
        for exchange in half_book.keys():
            exchange_exhausted[exchange] = False
            exchange_levels[exchange] = 0
            exchange_sizes[exchange] = 0
            inside_bids[exchange] = half_book[exchange][0]
            inside_bids[exchange][0] = inside_bids[exchange][0]
            if inside_bids[exchange][0] * (1 - EXCHANGE_FEES[exchange]) > total_inside_bid[0] * (
                    1 + EXCHANGE_FEES[exchange]):
                total_inside_bid = inside_bids[exchange]
                total_inside_bid_exchange = exchange

        while size > 0:
            total_available = total_inside_bid[0] * total_inside_bid[1]
            if total_available > size:
                exchange_sizes[total_inside_bid_exchange] = exchange_sizes.get(total_inside_bid_exchange, 0) + size
                exchange_prices[total_inside_bid_exchange] = total_inside_bid[0] * (
                            1 - EXCHANGE_FEES[total_inside_bid_exchange])
                size = 0
            else:
                size = size - total_available
                exchange_sizes[total_inside_bid_exchange] = \
                    exchange_sizes.get(total_inside_bid_exchange, 0) + total_available
                exchange_prices[total_inside_bid_exchange] = total_inside_bid[0] * (
                            1 - EXCHANGE_FEES[total_inside_bid_exchange])
                exchange_levels[total_inside_bid_exchange] = exchange_levels[total_inside_bid_exchange] + 1

            total_inside_bid = [float('-Inf'), float('-Inf')]
            for exchange in half_book.keys():
                if not exchange_levels[exchange] >= len(half_book[exchange]):
                    inside_bids[exchange] = half_book[exchange][exchange_levels[exchange]]
                    inside_bids[exchange][0] = inside_bids[exchange][0]
                    if inside_bids[exchange][0] * (1 + EXCHANGE_FEES[exchange]) > total_inside_bid[0] * (
                            1 - EXCHANGE_FEES[exchange]):
                        total_inside_bid = inside_bids[exchange]
                        total_inside_bid_exchange = exchange
                else:
                    exchange_exhausted[exchange] = True

            if False not in exchange_exhausted.values():
                exchange_sizes[half_book.keys()[0]] = exchange_sizes[half_book.keys()[0]] + size
                size = 0

        total_price = 0

        return min(exchange_prices.values())*(1 - PREMIUM)


def adjust_for_constant_rate(price, book, asset_side, side_spef):
    if asset_side == 'maker_asset':
        half_book = cache.get_order_book(book[0], book[1], 'asks')
        return price/(half_book[0][0])
    if asset_side == 'taker_asset':
        half_book = cache.get_order_book(book[0], book[1], 'bids')
        return price/(half_book[0][0])

def calculate_fee(fee_asset):
    gas_price = redis_interface.get_gas_price()
    gas_limit = redis_interface.get_gas_limit()
    gas_fee = gas_limit * (gas_price * (10 ** -9))

    # fast path: taker token is weth, so no rate conversion necessary
    if fee_asset == 'ETH':
        return gas_fee
    else:
        pseudo_quote = calculate_quote('ETH', fee_asset, gas_fee, None)
        return round(pseudo_quote['taker_size'], 15)

