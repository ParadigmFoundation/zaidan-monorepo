from requests import get, post


class InventoryManagerError(Exception):
    '''
    Defines an exception encountered during an inventory manager request.
    '''


class TestInventoryManager():
    '''
    Abstraction over the inventory-manager private API.
    '''

    open_orders = []

    def __init__(self):
        '''
        Create a new inventory manager wrapper.

        :param url: The host (and port) of the inventory manager service.
        '''


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

        request = {'price': price, 'size': size,
                   'side': side, 'symbol': symbol}

        return self._post('order/{}'.format(exchange), data=request)

    def get_order_status(self, exchange: str, order_id: str, symbol: str) -> dict:
        '''
        Get the status of a current order, by its exchange, ID, and market symbol.

        :param exchange: The name of the exchange the order was posted to.
        :param order_id: The UUID of the exchange order.
        :param symbol: The market symbol the order is for.
        '''

        tickers = symbol.split('/')
        if len(tickers) != 2:
            raise ValueError(
                'invalid market symbol (must be BASE/QUOTE format)')

        base_asset = tickers[0]
        quote_asset = tickers[1]

        return self._get('order/{}/{}/{}/{}'.format(exchange, order_id, base_asset, quote_asset))

    def cancel_order(self, exchange: str, order_id: str) -> bool:
        '''
        Cancel an open exchange order.

        :param exchange: The name of the exchange the order was place on.
        :param order_id: The UUID of the posted exchange order.
        '''

        return self._post('order/{}/{}/cancel'.format(exchange, order_id))

    def get_open_orders(self, exchange: str, symbol: str) -> list:
        '''
        Fetch all open exchange orders for a given exchange and market symbol.

        :param exchange: The name of the exchange.
        :param symbol: The market symbol to get open orders for (BASE/QUOTE).
        '''

        return self.open_orders

    def _get(self, endpoint: str, data=None) -> dict:
        '''
        Perform a GET call to an inventory manager endpoint.

        :param endpoint: The resource path.
        :param data: Additional data (query parameters).
        '''

        try:
            res = get(url='{}/{}'.format(self.url, endpoint),
                      params=data).json()
        except ValueError:
            raise InventoryManagerError('received invalid JSON response')
        except Exception as error:
            raise InventoryManagerError(
                'unknown exception encountered during request', error.args)

        if 'error' in res:
            raise InventoryManagerError(res['error'])

        return res

    def _post(self, endpoint: str, data=None) -> dict:
        '''
        Perform a POST call to an inventory manager endpoint.

        :param endpoint: The resource path.
        :param data: Additional data (request body).
        '''

        try:
            res = post(url='{}/{}'.format(self.url,
                                          endpoint), json=data).json()
        except ValueError:
            raise InventoryManagerError('received invalid JSON response')
        except Exception as error:
            raise InventoryManagerError(
                'unknown exception encountered during request', error.args)

        if 'error' in res:
            raise InventoryManagerError(res['error'])

        return res
