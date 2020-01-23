import uuid
import math
import time
import types_pb2
import services_pb2_grpc
from pricing_utils import PricingUtils
from risk_utils import RiskUtils
from redis_interface import RedisInterface
from asset_data import AssetData

VALIDITY_LENGTH = 15000

class MakerServicer(services_pb2_grpc.MakerServicer):

    def __init__(self, test:bool=False) -> None:
        services_pb2_grpc.MakerServicer.__init__(self)
        self.asset_data = AssetData()
        self.redis_interface = RedisInterface()
        self.risk_utils = RiskUtils()
        self.pricing_utils = PricingUtils()
        self.test = test

    def GetQuote(self, request: object, context) -> object:

        maker_asset = self.asset_data.get_ticker_with_address(request.maker_asset)
        taker_asset = self.asset_data.get_ticker_with_address(request.taker_asset)

        trading_maker_size = PricingUtils.convert_to_trading_units(maker_asset, request.maker_size)
        trading_taker_size = PricingUtils.convert_to_trading_units(taker_asset, request.taker_size)

        quote_info = self.pricing_utils.calculate_quote(maker_asset,
                                     taker_asset, trading_maker_size, trading_taker_size, self.test)
        expiry_timestamp = int(math.floor(time.time())) + int(VALIDITY_LENGTH)

        if False not in self.risk_utils.risk_checks(taker_asset, maker_asset, trading_taker_size, quote_info, self.test):
            quote_id = str(uuid.uuid4())

            self.redis_interface.add_quote(quote_id, {'expiration':expiry_timestamp,
                                       'taker_asset':request.taker_asset, 'maker_asset':request.maker_asset,
                                       'taker_size':quote_info['taker_size'], 'maker_size':quote_info['maker_size']})

            quote_info = PricingUtils.format_quote(taker_asset, maker_asset, quote_info)

            quote = types_pb2.GetQuoteResponse(quote_id=quote_id, expiration=expiry_timestamp,
                                               taker_asset=request.taker_asset, maker_asset=request.maker_asset,
                                               taker_size=str(quote_info['taker_size']),
                                               maker_size=str(quote_info['maker_size']))
            return quote
        else:
            response = types_pb2.GetQuoteResponse(status=400)
            return response


    def CheckQuote(self, request:object, context) -> object:
        try:
            quote = self.redis_interface.get_quote(request.quote_id)
        except ValueError:
            return types_pb2.CheckQuoteResponse(quote_id=request.quote_id, is_valid=False, status=1)

        if quote['expiration'] < time.time():
            return types_pb2.CheckQuoteResponse(quote_id=request.quote_id, is_valid=False, status=2)

        self.redis_interface.fill_quote(request.quote_id)

        return types_pb2.CheckQuoteResponse(quote_id=request.quote_id, is_valid=True, status=200)

    def OrderStatusUpdate(self, request:object, context) -> object:
        self.risk_utils.order_status_update(request.quote_id, request.status)
        return types_pb2.OrderStatusUpdateResponse(status=200)