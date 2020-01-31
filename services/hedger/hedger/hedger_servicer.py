import types_pb2
import services_pb2_grpc
from hedger import Hedger
import logging

class HedgerServicer(services_pb2_grpc.HedgerServicer):

    def __init__(self) -> None:
        services_pb2_grpc.HedgerServicer.__init__(self)
        self.hedger = Hedger()
        self.logger = logging.Logger('hedger-logger')

    
    def HedgeOrder(self, request: object, context) -> object:
        self.logger.info('received request to hedge order: ' + str(request.id))

        order_id = request.id

        if order_id:
            order_id_received = True
        else:
            order_id_received = False
        try:
            if order_id:
               self.hedger.events_callback(order_id)
        except Exception as e:
            print('failed to hedge trade')
            print(e)

        response = types_pb2.HedgeOrderResponse(valid=order_id_received)

        return response