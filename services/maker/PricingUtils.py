import TestDealerCache from TestDealerCache
import math

asset_pricing_data = {'ZRX': {'exchange_books': [('COINBASE', 'ZRX/USD'), ('BINANCE', 'ZRX/ETH')], 'implied_pref': ('COINBASE', 'ZRX/USD')},
                      'LINK': {'exchange_books': [('COINBASE', 'LINK/USD'), ('BINANCE', 'LINK/ETH')], 'implied_pref': ('COINBASE', 'LINK/USD')},
                      'ETH': {'exchange_books': [('COINBASE', 'ETH/USD'), ('BINANCE', 'ETH/USDT'), ('GEMINI', 'ETH/USD')], 'implied_pref': ('COINBASE', 'ETH/USD')},
                      'DAI': {'exchange_books': [('COINBASE', 'DAI/USD')], 'implied_pref': ('COINBASE', 'DAI/USD'), 'constant_rate': 'PREF_INSIDE'}
                      }

cache = TestDealerCache()

def calculate_quote(maker_asset, taker_asset, maker_size=None, taker_size=None):
    maker_asset_pricing_data = asset_pricing_data[maker_asset]
    taker_asset_pricing_data = asset_pricing_data[taker_asset]
    order_books = []
    implied = True
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
                    order_books[book[0]] = (cache.get_order_book(book[0], book[1], 'bids'))
                if taker_size:
                    price = _get_price_from_book_base(order_books, taker_size, 'buy')
                    price = adjust_for_constant_rate(price, maker_asset_pricing_data['implied_pref'], 'maker_asset')
                    return {'maker_size': taker_size * price, 'taker_size': taker_size}
                else:
                    price = _get_price_from_book_quote(order_books, taker_size, 'buy')
                    price = adjust_for_constant_rate(price, maker_asset_pricing_data['implied_pref'], 'maker_asset')
                    return {'maker_size': maker_size, 'taker_size': maker_size*price}
            elif 'constant_rate' in taker_asset_pricing_data.keys():
                books_to_get = maker_asset_pricing_data['exchange_books']
                order_books = {}
                for book in books_to_get:
                    order_books[book[0]] = (cache.get_order_book(book[0], book[1], 'asks'))
                if maker_size:
                    price = _get_price_from_book_base(order_books, maker_size, 'sell')
                    price = adjust_for_constant_rate(price, taker_asset_pricing_data['implied_pref'], 'taker_asset')
                    return {'maker_size': maker_size, 'taker_size': maker_size*price}

                else:
                    price = _get_price_from_book_quote(order_books, maker_size, 'buy')
                    price = adjust_for_constant_rate(price, taker_asset_pricing_data['implied_pref'], 'taker_asset')
                    return {'maker_size': taker_size * price, 'taker_size': taker_size}
            else:
                if taker_size:
                    books_to_get = taker_asset_pricing_data['implied_pref']
                    order_books = {}
                    for book in books_to_get:
                        order_books[book[0]] = (cache.get_order_book(book[0], book[1], 'bids'))

                    taker_asset_price = _get_price_from_book_base(order_books, taker_size, 'buy')

                    books_to_get = maker_asset_pricing_data['implied_pref']
                    order_books = {}
                    for book in books_to_get:
                        order_books[book[0]] = (cache.get_order_book(book[0], book[1], 'asks'))

                    maker_asset_price = _get_price_from_book_quote(order_books, taker_size*taker_asset_price, 'sell')

                    return {'maker_size': maker_asset_price*taker_size*taker_asset_price, 'taker_size':taker_size}
                else:
                    books_to_get = maker_asset_pricing_data['implied_pref']
                    order_books = {}
                    for book in books_to_get:
                        order_books[book[0]] = (cache.get_order_book(book[0], book[1], 'bids'))

                    maker_asset_price = _get_price_from_book_base(order_books, taker_size, 'buy')

                    books_to_get = taker_asset_pricing_data['implied_pref']
                    order_books = {}
                    for book in books_to_get:
                        order_books[book[0]] = (cache.get_order_book(book[0], book[1], 'asks'))

                    taker_asset_price = _get_price_from_book_quote(order_books, maker_size*maker_asset_price, 'sell')

                    return {'maker_size': maker_size, 'taker_size':taker_asset_price*maker_size*maker_asset_price}


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
                exchange_sizes[EXCHANGES[0]] = exchange_sizes[EXCHANGES[0]] + size
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
                    inside_asks[exchange][0] = inside_asks[exchange][0] * c_rate
                    if inside_asks[exchange][0] * (1 + EXCHANGE_FEES[exchange]) < total_inside_ask[0] * (
                            1 + EXCHANGE_FEES[exchange]):
                        total_inside_ask = inside_asks[exchange]
                        total_inside_ask_exchange = exchange
                else:
                    exchange_exhausted[exchange] = True
            if False not in exchange_exhausted.values():
                exchange_sizes[EXCHANGES[0]] = exchange_sizes[EXCHANGES[0]] + size
                size = 0
        return max(exchange_prices.values())*(1 + PREMIUM)

    raise Exception('Error in calculating price')

def _get_price_from_book_quote(self, half_book, size, side):
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
                exchange_sizes[EXCHANGES[0]] = exchange_sizes[EXCHANGES[0]] + size
                size = 0

        total_price = 0

        for key in exchange_prices.keys():
            price = exchange_prices[key] * (1 + PREMIUM)
            quantity = exchange_sizes[key]
            total_price = total_price + (quantity / price)

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
                    if inside_asks[exchange][0] * (1 + EXCHANGE_FEES[exchange]) > total_inside_bid[0] * (
                            1 - EXCHANGE_FEES[exchange]):
                        total_inside_bid = inside_bids[exchange]
                        total_inside_bid_exchange = exchange
                else:
                    exchange_exhausted[exchange] = True

            if False not in exchange_exhausted.values():
                exchange_sizes[EXCHANGES[0]] = exchange_sizes[EXCHANGES[0]] + size
                size = 0

        total_price = 0

        for key in exchange_prices.keys():
            price = exchange_prices[key] * (1 - PREMIUM)
            quantity = exchange_sizes[key]
            total_price = total_price + (quantity / price)
    try:
        # make sure rounding right way
        final_price = math.ceil(total_price * 100000000) / 100000000
        return final_price
    except:
        raise Exception('Error in calculating price')