import types_pb2
import services_pb2_grpc
from hedger import Hedger

class HedgerServicer(services_pb2_grpc.HedgerServicer):

    def __init__(self) -> None:
        services_pb2_grpc.HedgerServicer.__init__(self)
        self.hedger = Hedger()
    
    def HedgeOrder(self, request: object, context) -> object:

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