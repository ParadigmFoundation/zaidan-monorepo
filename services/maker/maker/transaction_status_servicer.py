import types_pb2
import services_pb2_grpc
from risk_utils import RiskUtils

class TransactionStatusServicer(services_pb2_grpc.TransactionStatusServicer):

    def __init__(self) -> None:
        services_pb2_grpc.TransactionStatusServicer.__init__(self)
        self.risk_utils = RiskUtils()

    def TransactionStatusUpdate(self, request:object, context) -> object:
        self.risk_utils.order_status_update(request.quote_id, request.status)
        return types_pb2.TransactionStatusUpdateResponse(status=200)
