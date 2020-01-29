# -*- coding: utf-8 -*-
"""
Hedger Tests - contain unittests for the hedger.py class

@author: ubuntu
"""
import sys
import time
import os
import collections

try:
    from test.TestInventoryManager import TestInventoryManager
    from test.TestRedisInterface import TestRedisInterface
except:
    from TestInventoryManager import TestInventoryManager
    from TestRedisInterface import TestRedisInterface

try:
    sys.path.append('hedger')
    from hedger import Hedger
except:
    try:
        sys.path.append('../hedger/hedger')
        from hedger import Hedger
    except:
        raise Exception('failed imports')

Order = collections.namedtuple('Order', 'id, symbol, side, price, quantity, timestamp, exchange')
tim = TestInventoryManager()
tri = TestRedisInterface()

hedger = Hedger(im=tim, cache=tri, test=True)

def test_hedge_sell():
    """Test hedger sell."""

    half_book = hedger.update_order_book('WETH/DAI')['asks']
    assert hedger.immediate_hedge_sell(50, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USD', side='B', price=101, quantity=50, timestamp=0, exchange='COINBASE')]

    assert hedger.immediate_hedge_sell(100, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USD', side='B', price=101, quantity=100, timestamp=0, exchange='COINBASE')]

    assert hedger.immediate_hedge_sell(350, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USDT', side='B', price=102, quantity=100, timestamp=0, exchange='BINANCE'), Order(id=-1, symbol='ETH/USD', side='B', price=102, quantity=250, timestamp=0, exchange='COINBASE')]

    assert hedger.immediate_hedge_sell(600, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USDT', side='B', price=103, quantity=300, timestamp=0, exchange='BINANCE'), Order(id=-1, symbol='ETH/USD', side='B', price=102, quantity=300, timestamp=0, exchange='COINBASE')]

    assert hedger.immediate_hedge_sell(900, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USDT', side='B', price=103, quantity=300, timestamp=0, exchange='BINANCE'), Order(id=-1, symbol='ETH/USD', side='B', price=103, quantity=600, timestamp=0, exchange='COINBASE')]

    assert hedger.immediate_hedge_sell(1200, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USDT', side='B', price=104, quantity=600, timestamp=0, exchange='BINANCE'), Order(id=-1, symbol='ETH/USD', side='B', price=103, quantity=600, timestamp=0, exchange='COINBASE')]

    half_book = hedger.update_order_book('ZRX/DAI')['asks']
    assert hedger.immediate_hedge_sell(50, half_book, 'ZRX/DAI') == [Order(id=-1, symbol='ZRX/USD', side='B', price=1.01, quantity=50, timestamp=0, exchange='COINBASE')]

    half_book = hedger.update_order_book('ZRX/WETH')['asks']
    assert hedger.immediate_hedge_sell(50, half_book, 'ZRX/WETH') == [Order(id=-1, symbol='ZRX/ETH', side='B', price=0.0102, quantity=50, timestamp=0, exchange='BINANCE')]

def test_hedge_buy():
    """Test hedge buy."""

    half_book = hedger.update_order_book('WETH/DAI')['bids']
    assert hedger.immediate_hedge_buy(50, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USD', side='S', price=99, quantity=50, timestamp=0, exchange='COINBASE')]

    assert hedger.immediate_hedge_buy(100, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USD', side='S', price=99, quantity=100, timestamp=0, exchange='COINBASE')]

    assert hedger.immediate_hedge_buy(350, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USDT', side='S', price=98, quantity=100, timestamp=0, exchange='BINANCE'), Order(id=-1, symbol='ETH/USD', side='S', price=98, quantity=250, timestamp=0, exchange='COINBASE')]
    assert hedger.immediate_hedge_buy(600, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USDT', side='S', price=97, quantity=300, timestamp=0, exchange='BINANCE'), Order(id=-1, symbol='ETH/USD', side='S', price=98, quantity=300, timestamp=0, exchange='COINBASE')]
    assert hedger.immediate_hedge_buy(900, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USDT', side='S', price=97, quantity=300, timestamp=0, exchange='BINANCE'), Order(id=-1, symbol='ETH/USD', side='S', price=97, quantity=600, timestamp=0, exchange='COINBASE')]
    assert hedger.immediate_hedge_buy(1200, half_book, 'WETH/DAI') == [Order(id=-1, symbol='ETH/USDT', side='S', price=96, quantity=600, timestamp=0, exchange='BINANCE'), Order(id=-1, symbol='ETH/USD', side='S', price=97, quantity=600, timestamp=0, exchange='COINBASE')]

    half_book = hedger.update_order_book('ZRX/DAI')['bids']
    assert hedger.immediate_hedge_buy(50, half_book, 'ZRX/DAI') == [Order(id=-1, symbol='ZRX/USD', side='S', price=0.99, quantity=50, timestamp=0, exchange='COINBASE')]

    half_book = hedger.update_order_book('ZRX/WETH')['bids']
    assert hedger.immediate_hedge_buy(50, half_book, 'ZRX/WETH') == [Order(id=-1, symbol='ZRX/ETH', side='S', price=0.0098, quantity=50, timestamp=0, exchange='BINANCE')]



def test_decrement_orders():
    order_book = hedger.update_order_book('WETH/DAI')
    test_time = 1
    hedger.im.open_orders.append({'id':'123abc', 'price': order_book['bids']['COINBASE'][0][0], 'symbol': 'ETH/USD', 'side': 'buy', 'remaining': 300, 'timestamp':test_time, 'exchange': 'COINBASE'})
    new_orders, cancels, leftover = hedger.decrement_orders(100, 'buys', order_book, 'WETH/DAI')
    assert new_orders == [Order(id='123abc', symbol='ETH/USD', side='B', price=99, quantity=200, timestamp=test_time, exchange='COINBASE')]
    assert cancels == [{'id': '123abc', 'price': 99, 'symbol': 'ETH/USD', 'side': 'buy', 'remaining': 300, 'timestamp': test_time, 'exchange': 'COINBASE'}]
    assert leftover == 0

    new_orders, cancels, leftover = hedger.decrement_orders(300, 'buys', order_book, 'WETH/DAI')
    assert new_orders == []
    assert cancels[0]['id'] == '123abc'
    assert leftover == 0


def test_check_unhedged_positions():
    hedger.cache.unhedged_positions['WETH/DAI'] = 10
    assert hedger.check_unhedged_positions('WETH/DAI') == [Order(id=-1, symbol='ETH/USD', side='S', price=99, quantity=10, timestamp=0, exchange='COINBASE')]

    hedger.cache.unhedged_positions['WETH/DAI'] = .000001
    assert hedger.check_unhedged_positions('WETH/DAI') == []

def test_find_orders_to_place():
    order_book = hedger.update_order_book('WETH/DAI')
    hedger.cache.unhedged_positions['WETH/DAI'] = 0
    test_trade = {}
    test_trade['size'] = 50
    test_trade['side'] = 'S'
    test_trade['pair'] = 'WETH/DAI'
    test_trade['price'] = 100
    assert hedger.find_orders_to_place(test_trade, order_book) == ([Order(id=-1, symbol='ETH/USD', side='B', price=101, quantity=50, timestamp=0, exchange='COINBASE')], [])

def test_events_callback():
    order = {}
    order['maker_size'] = 50
    order['side'] = 'S'
    order['pair'] = 'WETH/DAI'
    order['maker_price'] = 100
    hedger.events_callback(1, order)


test_hedge_sell()
test_hedge_buy()
test_decrement_orders()
test_check_unhedged_positions()
test_find_orders_to_place()
test_events_callback()
