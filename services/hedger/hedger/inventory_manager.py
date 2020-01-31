import grpc
import os
from services_pb2_grpc import ExchangeManagerStub
from types_pb2 import ExchangeCreateOrderRequest, ExchangeOrderRequest, ExchangeOrderRequest, ExchangeOrder
from google.protobuf import wrappers_pb2 as wrappers


class InventoryManagerError(Exception):
    '''
    Defines an exception encountered during an inventory manager request.
    '''


EXCHANGEMANAGER_CHANNEL = os.environ.get('EXCHANGEMANAGER_CHANNEL', 'localhost:8000')

class InventoryManager():

    def __init__(self):
        self.initialize_exchangemanager_connection()

    def initialize_exchangemanager_connection(self) -> None:
        self.em_channel = grpc.insecure_channel(EXCHANGEMANAGER_CHANNEL)
        self.em_stub = ExchangeManagerStub(self.em_channel)

    def post_order(self, exchange: str, symbol: str, side: str, price: float, size: float) -> dict:
        '''
        Post an order to an exchange.

        :param exchange: The name of the exchange to post an order to.
        :param symbol: The market symbol to post an order to (of format BASE/QUOTE).
        :param side: The side of the order (either "buy" or "sell").
        :param price: The price of the order (in units of the quote asset).
        :param size: The size of the order (in units of the base asset).
        '''

        tickers = symbol.split('/')
        if len(tickers) != 2:
            raise ValueError(
                'invalid market symbol (must be BASE/QUOTE format)')

        if side not in ('buy', 'sell'):
            raise ValueError('invalid side (must be "buy" or "sell")')

        if side == 'buy':
            side_enum = ExchangeOrder.Side.Buy
        else:
            side_enum = ExchangeOrder.Side.Sell

        order = ExchangeOrder(price=str(price), symbol=symbol, amount=str(size), side=side_enum)

        req = ExchangeCreateOrderRequest(exchange=exchange, order=order)

        response = self.em_stub.CreateOrder(req)
        return response


    def get_order_status(self, exchange: str, order_id: str, symbol: str) -> dict:
        pass

    def cancel_order(self, exchange: str, order_id: str) -> bool:
        '''
        Cancel an open exchange order.

        :param exchange: The name of the exchange the order was place on.
        :param order_id: The UUID of the posted exchange order.
        '''

        req = ExchangeOrderRequest(exchange=exchange, id=order_id)
        response = self.em_stub.CancelOrder(req)
        return {'cancelled':True}

    def get_open_orders(self, exchange: str, symbol: str) -> list:
        '''
        Fetch all open exchange orders for a given exchange and market symbol.

        :param exchange: The name of the exchange.
        :param symbol: The market symbol to get open orders for (BASE/QUOTE).
        '''

        req = exchange
        print('received open orders request')
        response = self.em_stub.GetOpenOrders(wrappers.StringValue(req))
        print('successfully got open orders')
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


