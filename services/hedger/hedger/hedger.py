from os import getenv
import json
import collections
import time
import threading

from zaidan import Logger
from redis_interface import RedisInterface
from orderbook_wrapper import OrderBookWrapper
from inventory_manager import InventoryManager
from services_pb2_grpc import ExchangeManagerStub
from types_pb2 import GetOpenOrdersRequest
import grpc

EXCHANGEMANAGER_CHANNEL = getenv('EXCHANGEMANAGER_CHANNEL', default='em:8000')


Order = collections.namedtuple('Order', 'id, symbol, side, price, quantity, timestamp, exchange')

INVENTORY_MANAGER_URL = getenv('INVENTORY_MANAGER_URL', default='http://inventory-manager:8000')

REDIS_HOST = getenv('REDIS_HOST', default='localhost')
REDIS_PORT = getenv('REDIS_PORT', default='6379')
REDIS_PASSWORD = getenv('REDIS_PASSWORD', default=None)

MYSQL_HOST = getenv('MYSQL_HOST', default='127.0.0.1')
MYSQL_PORT = int(getenv('MYSQL_PORT', default='3306'))
MYSQL_DB = getenv('MYSQL_DB', default='dealer_db')
MYSQL_USER = getenv('MYSQL_USER', default=None)
MYSQL_PASSWORD = getenv('MYSQL_PASSWORD', default=None)

PRICE_PAIRS = {'WETH/DAI': {'BINANCE': 'ETH/USDT', 'COINBASE': 'ETH/USD'}, 'ZRX/WETH': {'BINANCE': 'ZRX/ETH'}, 'ZRX/DAI': {'COINBASE': 'ZRX/USD', 'BINANCE': 'ZRX/USDT'}}
DEFAULT_EXCHANGES = {'WETH/DAI': 'COINBASE', 'ZRX/WETH': 'BINANCE', 'ZRX/DAI': 'COINBASE'}
MIN_TRADE_SIZES = {'ETH/USDT': {'BINANCE': 10}, 'ETH/USD': {'COINBASE': 10}, 'ZRX/ETH': {'BINANCE': .03},
                   'ZRX/USD': {'COINBASE': 10}, 'ZRX/USDT': {'BINANCE': 10}}
EXCHANGE_FEES = {'BINANCE':.00075, 'COINBASE':.001}

PAIRS = ['WETH/DAI', 'ZRX/DAI']
ENVIRONMENT = getenv('ENVIRONMENT', 'TEST')

MAX_BOOK_AGES = {'ETH/USD': 10, 'ZRX/USD': 15, 'USDC/USDT': 15, 'DAI/USDC': 25}

DISTANCE_THRESHOLD = float(getenv('DISTANCE_THRESHOLD', '0.1'))  # distance threshold for order submission
ORDER_TIMEOUT = int(getenv('ORDER_TIMEOUT', '18000'))
TOKEN_GRANULARITY = int(getenv('TOKEN_GRANULARITY', '2'))
PRECISION = int(getenv('PRECISION', "2"))

OBM_DEBUG = getenv("OBM_DEBUG", "False")

if ENVIRONMENT == 'TEST':
    HEDGE_MARGIN = 0
else:
    HEDGE_MARGIN = .02

if OBM_DEBUG == "True":
    OBM_DEBUG = True
else:
    OBM_DEBUG = False

class Hedger():
    '''
    Dealer hedger implementation class.

    Relies on a private inventory-manager service API for exchange interaction.
    '''

    EXCHANGE_FEES = {}
    Units = {}

    open_orders_lock = threading.Lock()
    test_id = 0

    def __init__(self, im=None, cache=None, logger=None, test=False):
        '''
        Create a new hedger.

        :param logger: When runing alongside hedger server, pass logger here.
        '''

        if not logger:
            self.logger = Logger('Hedger')
        else:
            self.logger = logger

        if cache:
            self.cache = cache
        else:
            self.cache = RedisInterface()

        '''
        # mysql db for filled 0x and exchange orders
        if not test:
            self.db = DealerDatabase(MYSQL_HOST, MYSQL_PORT, MYSQL_DB, MYSQL_USER, MYSQL_PASSWORD)
        '''

        # inventory manager
        if im:
            self.im = im
        else:
            self.im = InventoryManager()
        self.logger.info("hedger created", {'ok': True})

        self.order_book_wrapper = OrderBookWrapper(unit_test=test)
        self.initialize_exchangemanager_connection()

    def initialize_exchangemanager_connection(self) -> None:
        chan = grpc.insecure_channel(EXCHANGEMANAGER_CHANNEL)
        self.em_stub = ExchangeManagerStub(chan)

    def events_callback(self, quote_id, order=None):
        """ Get called by."""

        if not order:
            order = self.cache.get_quote(quote_id)

        self.logger.info('received order', order)
        self.open_orders_lock.acquire()

        if order['taker_asset'] == 'WETH' or order['maker_asset'] == 'WETH':
            order['pair'] = 'WETH/DAI'
        elif order['taker_asset'] == 'ZRX' or order['maker_asset'] == 'ZRX':
            if order['taker_asset'] == 'WETH' or order['maker_asset'] == 'WETH':
                order['pair'] == 'ZRX/WETH'
            else:
                order['pair'] == 'ZRX/DAI'


        if order['pair'] == 'ZRX/WETH' and ENVIRONMENT != 'TEST':
            order['pair'] = 'ZRX/DAI'

        trade = {}
        if order['maker_asset'] == order['pair'].split('/')[0]:
            trade['size'] = float(order['maker_size'])
            trade['side'] = 'S'
            trade['price'] = float(order['taker_size']) / float(order['maker_size'])
        else:
            trade['size'] = float(order['taker_size'])
            trade['side'] = 'B'
            trade['price'] = float(order['maker_size'])/float(order['taker_size'])


        order_book = self.update_order_book(order['pair'])
        trade['pair'] = order['pair']
        new_orders, cancels = self.find_orders_to_place(trade, order_book)

        for new_order in new_orders:
            self.execute_order(new_order)

        for cancel in cancels:
            self.execute_cancel(cancel)

        self.open_orders_lock.release()

    def find_orders_to_place(self, trade, order_book):
        """ Receive a trade from the 0x order watcher and try to find an opposing buy or sell to place."""
        orders = []
        cancels = []
        pair = trade['pair']
        side = trade['side']
        size = round(trade['size'], PRECISION)
        price = trade['price']
        self.logger.info("Received trade", trade)
        unhedged_positions_update = 0
        position = self.get_current_position(pair)

        # checking if the order is large enough to constitute an acceptable exchange order
        self.logger.info("conditional info")
        self.logger.info(str(order_book['bids'][DEFAULT_EXCHANGES[pair]][0][0]))
        self.logger.info(str(size))
        self.logger.info(str(min_size_from_dealer_pair(pair)))
        if size * order_book['bids'][DEFAULT_EXCHANGES[pair]][0][0] > min_size_from_dealer_pair(pair):
            if side == 'B':
                self.logger.info('logic path', {'action': 'received buy'})
                unhedged_positions_update += (trade['size'] - size)
                if position >= 0:
                    self.logger.info('logic path', {'action': 'normal buy hedge'})
                    orders = self.find_sells_to_place(side, size, price, pair, order_book)
                elif (position < 0) and (position < -1*size):
                    self.logger.info('logic path', {'action': 'fully merging buy hedge'})
                    orders, cancels, leftover = self.decrement_orders(size, 'buys', order_book)
                    unhedged_positions_update += -leftover
                else:
                    self.logger.info('logic path', {'action': 'partially merging buy hedge'})
                    trade['size'] = size + position
                    d_orders, d_cancels, leftover = self.decrement_orders(size, 'buys', order_book, pair, cancel_all=True)
                    for new_order in d_orders:
                        self.execute_order(new_order)

                    for cancel in d_cancels:
                        self.execute_cancel(cancel)
                    return self.find_orders_to_place(trade, order_book)
            else:
                self.logger.info('logic path', {'action': 'received sell'})
                unhedged_positions_update += -1*(trade['size'] - size)
                if position <= 0:
                    self.logger.info('logic path', {'action': 'normal sell hedge'})
                    orders = self.find_buys_to_place(side, size, price, pair, order_book)
                elif (position > 0) and (position > size):
                    self.logger.info('logic path', {'action': 'fully merging sell hedge'})
                    orders, cancels, leftover = self.decrement_orders(size, 'sells', order_book)
                    unhedged_positions_update += leftover
                else:
                    self.logger.info('logic path', {'action': 'partially merging sell hedge'})
                    trade['size'] = size - position
                    d_orders, d_cancels, leftover = self.decrement_orders(size, 'sells', order_book, pair, cancel_all=True)
                    for new_order in d_orders:
                        self.execute_order(new_order)

                    for cancel in d_cancels:
                        self.execute_cancel(cancel)
                    return self.find_orders_to_place(trade, order_book)
        else:
            self.logger.info('logic path', {'action': 'trade too small to hedge'})
            if trade['side'] == 'B':
                unhedged_positions_update += trade['size']
            else:
                unhedged_positions_update += -1*trade['size']

        self.cache.set_unhedged_position(pair, self.cache.get_unhedged_position)
        return orders, cancels

    def update_order_book(self, pair):
        """ Fetch most recent order book. """
        return self.order_book_wrapper.get_current_total_book(pair)

    def get_current_position(self, pair):
        ''' Get current position on exchanges. '''
        self.logger.info('beginning of get current position function')
        position = 0
        if ENVIRONMENT == 'PRODUCTION' or ENVIRONMENT == 'DEBUG':
            for exchange in ['COINBASE']:
                self.logger.info('calling out to exchange manager in hedger')
                if ENVIRONMENT == 'DEBUG':
                    open_orders = self.get_open_orders_debug(exchange, PRICE_PAIRS[pair][exchange])
                else:
                    open_orders = self.im.get_open_orders(exchange, PRICE_PAIRS[pair][exchange])
                self.logger.info('received response from hedger')
                for order in open_orders:
                    if order['side'] == 'buy':
                        position = position - order['remaining']
                    else:
                        position = position + order['remaining']

        else:
            for exchange in PRICE_PAIRS[pair]:
                if ENVIRONMENT == 'DEBUG':
                    open_orders = self.get_open_orders_debug(exchange, PRICE_PAIRS[pair][exchange])
                else:
                    open_orders = self.im.get_open_orders(exchange, PRICE_PAIRS[pair][exchange])
                for order in open_orders:
                    if order['side'] == 'buy':
                        position = position - order['remaining']
                    else:
                        position = position + order['remaining']
        return position

    def get_open_orders(self, pair):
        ''' Get open orders for a given trading pair. '''
        open_orders = []
        for exchange in ['COINBASE']:
            if ENVIRONMENT == 'DEBUG':
                e_open_orders = self.get_open_orders_debug(exchange, PRICE_PAIRS[pair][exchange])
            else:
                e_open_orders = self.im.get_open_orders(exchange, PRICE_PAIRS[pair][exchange])
            for order in open_orders:
                order['exchange'] = exchange
            open_orders.extend(e_open_orders)

        return open_orders

    def decrement_orders(self, size, order_side, order_book, pair, cancel_all=False):
        """ Decrement orders. """
        cancels = []
        new_orders = []
        leftover = 0
        if True:
            size_outstanding = size
            open_orders_l = self.get_open_orders(pair)
            open_orders_l.sort(key=lambda x: x['price'])
            if order_side == 'buys':
                i = 0
                while size_outstanding > 0 and i < len(open_orders_l):
                    if open_orders_l[i]['remaining'] > size_outstanding:
                        order = open_orders_l[i]
                        # Checking if order size in ETH is large enough to place
                        if (open_orders_l[i]['remaining'] - size_outstanding) * order_book['asks'][order['exchange']][0][0] > float(
                                MIN_TRADE_SIZES[order['symbol']][order['exchange']]):
                            cancels.append(open_orders_l[i])
                            new_order = open_orders_l[i]
                            new_order = Order(new_order['id'], new_order['symbol'], 'B', new_order['price'],
                                              round(new_order['remaining'] - size_outstanding, PRECISION),
                                              new_order['timestamp'], new_order['exchange'])
                            new_orders.append(new_order)
                            size_outstanding = 0
                        else:
                            cancels.append(order)
                            leftover = open_orders_l[i]['remaining'] - size_outstanding
                            size_outstanding = 0
                    else:
                        cancels.append(open_orders_l[i])
                        size_outstanding = size_outstanding - open_orders_l[i]['remaining']
                        i = i + 1
            else:
                i = -1
                while size_outstanding > 0 and i > (-1 * len(open_orders_l) - 1):
                    if open_orders_l[i]['remaining'] > size_outstanding:
                        order = open_orders_l[i]
                        # Checking if order size in ETH is large enough to place
                        if (open_orders_l[i]['remaining'] - size_outstanding) * order_book['asks'][order['exchange']][0][0] > float(
                                MIN_TRADE_SIZES[order['symbol']][order['exchange']]):
                            cancels.append(open_orders_l[i])
                            new_order = open_orders_l[i]
                            new_order = Order(new_order['id'], new_order['symbol'], 'S', new_order['price'],
                                              round(new_order['remaining'] - size_outstanding, PRECISION),
                                              new_order['timestamp'], new_order['exchange'])
                            new_orders.append(new_order)
                            size_outstanding = 0
                        else:
                            cancels.append(order)
                            leftover = -1 * (open_orders_l[i]['remaining'] - size_outstanding)
                            size_outstanding = 0
                    else:
                        cancels.append(open_orders_l[i])
                        size_outstanding = size_outstanding - open_orders_l[i]['remaining']
                        i = i - 1

        return new_orders, cancels, leftover

    def find_buys_to_place(self, side, size, price, pair, order_book):
        """ Recieve a sell event and now look to place a buy to hedge."""

        buys = self.immediate_hedge_sell(size, order_book['asks'], pair)

        return buys

    def find_sells_to_place(self, side, size, price, pair, order_book):
        """ Receive a buy event and look to place a sell to hedge."""

        sells = self.immediate_hedge_buy(size, order_book['bids'], pair)

        return sells

    def immediate_hedge_sell(self, size, half_book, pair):
        '''Immedidately hedges a sell by taking orders on multiple exchanges.'''
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
            if inside_asks[exchange][0] < total_inside_ask[0]:
                total_inside_ask = inside_asks[exchange]
                total_inside_ask_exchange = exchange

        while size > 0:
            if total_inside_ask[1] > size:
                exchange_sizes[total_inside_ask_exchange] = exchange_sizes.get(total_inside_ask_exchange, 0) + size
                exchange_prices[total_inside_ask_exchange] = total_inside_ask[0]
                size = 0
            else:
                size = size - total_inside_ask[1]
                exchange_sizes[total_inside_ask_exchange] = \
                    exchange_sizes.get(total_inside_ask_exchange, 0) + total_inside_ask[1]
                exchange_prices[total_inside_ask_exchange] = total_inside_ask[0]
                exchange_levels[total_inside_ask_exchange] = exchange_levels[total_inside_ask_exchange] + 1

            total_inside_ask = [float('Inf'), float('Inf')]
            for exchange in half_book.keys():
                if not exchange_levels[exchange] >= len(half_book[exchange]):
                    inside_asks[exchange] = half_book[exchange][exchange_levels[exchange]]
                    if inside_asks[exchange][0] < total_inside_ask[0]:
                        total_inside_ask = inside_asks[exchange]
                        total_inside_ask_exchange = exchange
                else:
                    exchange_exhausted[exchange] = True
            if False not in exchange_exhausted.values():
                exchange_sizes[DEFAULT_EXCHANGES[pair]] = exchange_sizes[DEFAULT_EXCHANGES[pair]] + size
                size = 0

        orders = []
        for exchange in exchange_sizes:
            if exchange_sizes[exchange] > 0:
                orders.append(
                    Order(-1, PRICE_PAIRS[pair][exchange], 'B', exchange_prices[exchange] * (1 + HEDGE_MARGIN),
                          exchange_sizes[exchange], 0, exchange))

        return orders

    def immediate_hedge_buy(self, size, half_book, pair):
        '''Immediately hedges a buy by taking orders on multiple exchanges.'''
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
            if inside_bids[exchange][0] > total_inside_bid[0]:
                total_inside_bid = inside_bids[exchange]
                total_inside_bid_exchange = exchange

        while size > 0:
            if total_inside_bid[1] > size:
                exchange_sizes[total_inside_bid_exchange] = exchange_sizes.get(total_inside_bid_exchange, 0) + size
                exchange_prices[total_inside_bid_exchange] = total_inside_bid[0]
                size = 0
            else:
                size = size - total_inside_bid[1]
                exchange_sizes[total_inside_bid_exchange] = \
                    exchange_sizes.get(total_inside_bid_exchange, 0) + total_inside_bid[1]
                exchange_prices[total_inside_bid_exchange] = total_inside_bid[0]
                exchange_levels[total_inside_bid_exchange] = exchange_levels[total_inside_bid_exchange] + 1

            total_inside_bid = [float('-Inf'), float('-Inf')]
            for exchange in half_book.keys():
                if not exchange_levels[exchange] >= len(half_book[exchange]):
                    inside_bids[exchange] = half_book[exchange][exchange_levels[exchange]]
                    if inside_bids[exchange][0] > total_inside_bid[0]:
                        total_inside_bid = inside_bids[exchange]
                        total_inside_bid_exchange = exchange
                else:
                    exchange_exhausted[exchange] = True
            if False not in exchange_exhausted.values():
                exchange_sizes[DEFAULT_EXCHANGES[pair]] = exchange_sizes[DEFAULT_EXCHANGES[pair]] + size
                size = 0

        orders = []
        for exchange in exchange_sizes:
            if exchange_sizes[exchange] > 0:
                orders.append(
                    Order(-1, PRICE_PAIRS[pair][exchange], 'S', exchange_prices[exchange] * (1 - HEDGE_MARGIN),
                          exchange_sizes[exchange], 0, exchange))

        return orders

    def check_unhedged_positions(self, pair):
        """ Check unhedged positions. """
        new_orders = []
        DEFAULT_EXCHANGE = DEFAULT_EXCHANGES[pair]
        unhedged_position = self.cache.get_unhedged_position(pair)
        pair = PRICE_PAIRS[pair][DEFAULT_EXCHANGE]
        order_book = self.update_order_book(pair)
        bids = order_book['bids'][DEFAULT_EXCHANGE]
        asks = order_book['asks'][DEFAULT_EXCHANGE]
        if unhedged_position * order_book['bids'][DEFAULT_EXCHANGE][0][0] > float(
                MIN_TRADE_SIZES[pair][DEFAULT_EXCHANGE]):
            new_orders.append(
                Order(-1, pair, 'S', bids[0][0], round(unhedged_position, PRECISION), 0,
                      DEFAULT_EXCHANGE))
            unhedged_position = unhedged_position - round(
                unhedged_position, PRECISION)
            self.cache.set_unhedged_position(pair, unhedged_position)
        elif unhedged_position * order_book['bids'][DEFAULT_EXCHANGE][0][0] < (
                -1 * float(MIN_TRADE_SIZES[pair][DEFAULT_EXCHANGE])):
            new_orders.append(
                Order(-1, pair, 'B', asks[0][0], round(-1 * unhedged_position, PRECISION), 0,
                      DEFAULT_EXCHANGE))
            unhedged_position = unhedged_position + round(
                -1 * unhedged_position, PRECISION)
            self.cache.set_unhedged_position(pair, unhedged_position)

        return new_orders

    def execute_order(self, order, unhedged=False):
        """ Perform the execution of an order as well as adding it to the open_orders list."""

        if ENVIRONMENT == 'PRODUCTION':
            self.logger.info('executing order', order_to_dict(order))
            try:
                if order.side == 'B':
                    response = self.im.post_order(order.exchange, order.symbol, 'buy', order.price, order.quantity)
                else:
                    response = self.im.post_order(order.exchange, order.symbol, 'sell', order.price, order.quantity)
            except Exception as e:
                self.logger.error("Unknown exception in execute order", {'error': str(e)})
        else:
            self.logger.info('Order to be sent to exchange:', {'order': str(Order(self.test_id, order.symbol, order.side, order.price, order.quantity, int(time.time()), order.exchange))})
            self.test_id = self.test_id + 1

    def execute_cancel(self, order):
        """Cancel an active order."""

        if ENVIRONMENT == 'PRODUCTION':
            self.logger.info('canceling order', order_to_dict(order))
            response = self.im.cancel_order(order['exchange'], order['id'])
            if not response['cancelled']:
                raise Exception('Order ' + order.id + ' failed to cancel')

        else:
            print('Cancel to be sent to exchange:')
            print(order)

    def get_open_orders_debug(self, exchange, symbol):
        response = self.em_stub.GetOpenOrders(GetOpenOrdersRequest(exchange=exchange.lower()))
        self.logger.info('successfully got open orders')
        orders_list = []
        for order in response:
            if order.order.symbol == symbol:
                order_dict = {}
                if order.order.side == 0:
                    order_dict['side'] = 'buy'
                else:
                    order_dict['side'] = 'sell'
                order_dict['symbol'] = symbol
                order_dict['price'] = float(order.order.price)
                order_dict['id'] = order.order.id
                order_dict['timestamp'] = order.status.timestamp
                order_dict['remaining'] = float(order.order.size) - float(order.status.filled)
                orders_list.append(order_dict)

        return orders_list


def min_size_from_dealer_pair(pair):
    ''' Get the minimum trade size from a dealer pair. '''
    exchange_pairs = PRICE_PAIRS[pair].values()
    min_sizes = []
    for exchange_pair in exchange_pairs:
        min_sizes.extend(list(MIN_TRADE_SIZES[exchange_pair].values()))
    return min(min_sizes)


def order_to_dict(order):
    ''' Convert order NamedTuple to dict. '''
    return dict(id=order.id, symbol=order.symbol, side=order.side, price=order.price, quantity=order.quantity,
                timestamp=order.timestamp, exchange=order.exchange)


