import types_pb2
import services_pb2_grpc

class HedgerServicer(services_pb2_grpc.HedgerServicer):

    def __init__(self) -> None:
        services_pb2_grpc.HedgerServicer.__init__(self)
    
    def HedgeOrder(self, request: object, context) -> object:

        order_id = request.id
        
        if order_id:
            order_id_received = True
        else:
            order_id_received = False

        response = types_pb2.HedgeOrderResponse(valid=order_id_received)

        return response