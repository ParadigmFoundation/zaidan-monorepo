# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
from google.protobuf import wrappers_pb2 as google_dot_protobuf_dot_wrappers__pb2
import types_pb2 as types__pb2


class OrderBookManagerStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.OrderBook = channel.unary_unary(
        '/OrderBookManager/OrderBook',
        request_serializer=types__pb2.OrderBookRequest.SerializeToString,
        response_deserializer=types__pb2.OrderBookResponse.FromString,
        )
    self.Updates = channel.unary_stream(
        '/OrderBookManager/Updates',
        request_serializer=types__pb2.OrderBookUpdatesRequest.SerializeToString,
        response_deserializer=types__pb2.OrderBookResponse.FromString,
        )


class OrderBookManagerServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def OrderBook(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Updates(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_OrderBookManagerServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'OrderBook': grpc.unary_unary_rpc_method_handler(
          servicer.OrderBook,
          request_deserializer=types__pb2.OrderBookRequest.FromString,
          response_serializer=types__pb2.OrderBookResponse.SerializeToString,
      ),
      'Updates': grpc.unary_stream_rpc_method_handler(
          servicer.Updates,
          request_deserializer=types__pb2.OrderBookUpdatesRequest.FromString,
          response_serializer=types__pb2.OrderBookResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'OrderBookManager', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))


class MakerStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.GetQuote = channel.unary_unary(
        '/Maker/GetQuote',
        request_serializer=types__pb2.GetQuoteRequest.SerializeToString,
        response_deserializer=types__pb2.GetQuoteResponse.FromString,
        )
    self.CheckQuote = channel.unary_unary(
        '/Maker/CheckQuote',
        request_serializer=types__pb2.CheckQuoteRequest.SerializeToString,
        response_deserializer=types__pb2.CheckQuoteResponse.FromString,
        )
    self.OrderStatusUpdate = channel.unary_unary(
        '/Maker/OrderStatusUpdate',
        request_serializer=types__pb2.OrderStatusUpdateRequest.SerializeToString,
        response_deserializer=types__pb2.OrderStatusUpdateResponse.FromString,
        )


class MakerServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def GetQuote(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CheckQuote(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def OrderStatusUpdate(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_MakerServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'GetQuote': grpc.unary_unary_rpc_method_handler(
          servicer.GetQuote,
          request_deserializer=types__pb2.GetQuoteRequest.FromString,
          response_serializer=types__pb2.GetQuoteResponse.SerializeToString,
      ),
      'CheckQuote': grpc.unary_unary_rpc_method_handler(
          servicer.CheckQuote,
          request_deserializer=types__pb2.CheckQuoteRequest.FromString,
          response_serializer=types__pb2.CheckQuoteResponse.SerializeToString,
      ),
      'OrderStatusUpdate': grpc.unary_unary_rpc_method_handler(
          servicer.OrderStatusUpdate,
          request_deserializer=types__pb2.OrderStatusUpdateRequest.FromString,
          response_serializer=types__pb2.OrderStatusUpdateResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'Maker', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))


class HotWalletStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.CreateOrder = channel.unary_unary(
        '/HotWallet/CreateOrder',
        request_serializer=types__pb2.CreateOrderRequest.SerializeToString,
        response_deserializer=types__pb2.CreateOrderResponse.FromString,
        )
    self.GetAllowance = channel.unary_unary(
        '/HotWallet/GetAllowance',
        request_serializer=types__pb2.GetAllowanceRequest.SerializeToString,
        response_deserializer=types__pb2.GetAllowanceResponse.FromString,
        )
    self.SetAllowance = channel.unary_unary(
        '/HotWallet/SetAllowance',
        request_serializer=types__pb2.SetAllowanceRequest.SerializeToString,
        response_deserializer=types__pb2.SetAllowanceResponse.FromString,
        )
    self.GetTokenBalance = channel.unary_unary(
        '/HotWallet/GetTokenBalance',
        request_serializer=types__pb2.GetBalanceRequest.SerializeToString,
        response_deserializer=types__pb2.GetBalanceResponse.FromString,
        )
    self.GetEtherBalance = channel.unary_unary(
        '/HotWallet/GetEtherBalance',
        request_serializer=types__pb2.GetBalanceRequest.SerializeToString,
        response_deserializer=types__pb2.GetBalanceResponse.FromString,
        )
    self.TransferEther = channel.unary_unary(
        '/HotWallet/TransferEther',
        request_serializer=types__pb2.TransferRequest.SerializeToString,
        response_deserializer=types__pb2.TransferResponse.FromString,
        )
    self.TransferToken = channel.unary_unary(
        '/HotWallet/TransferToken',
        request_serializer=types__pb2.TransferRequest.SerializeToString,
        response_deserializer=types__pb2.TransferResponse.FromString,
        )
    self.SendTransaction = channel.unary_unary(
        '/HotWallet/SendTransaction',
        request_serializer=types__pb2.SendTransactionRequest.SerializeToString,
        response_deserializer=types__pb2.SendTransactionResponse.FromString,
        )
    self.ExecuteZeroExTransaction = channel.unary_unary(
        '/HotWallet/ExecuteZeroExTransaction',
        request_serializer=types__pb2.ExecuteZeroExTransactionRequest.SerializeToString,
        response_deserializer=types__pb2.ExecuteZeroExTransactionResponse.FromString,
        )


class HotWalletServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def CreateOrder(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetAllowance(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SetAllowance(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetTokenBalance(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetEtherBalance(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def TransferEther(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def TransferToken(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SendTransaction(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def ExecuteZeroExTransaction(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_HotWalletServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'CreateOrder': grpc.unary_unary_rpc_method_handler(
          servicer.CreateOrder,
          request_deserializer=types__pb2.CreateOrderRequest.FromString,
          response_serializer=types__pb2.CreateOrderResponse.SerializeToString,
      ),
      'GetAllowance': grpc.unary_unary_rpc_method_handler(
          servicer.GetAllowance,
          request_deserializer=types__pb2.GetAllowanceRequest.FromString,
          response_serializer=types__pb2.GetAllowanceResponse.SerializeToString,
      ),
      'SetAllowance': grpc.unary_unary_rpc_method_handler(
          servicer.SetAllowance,
          request_deserializer=types__pb2.SetAllowanceRequest.FromString,
          response_serializer=types__pb2.SetAllowanceResponse.SerializeToString,
      ),
      'GetTokenBalance': grpc.unary_unary_rpc_method_handler(
          servicer.GetTokenBalance,
          request_deserializer=types__pb2.GetBalanceRequest.FromString,
          response_serializer=types__pb2.GetBalanceResponse.SerializeToString,
      ),
      'GetEtherBalance': grpc.unary_unary_rpc_method_handler(
          servicer.GetEtherBalance,
          request_deserializer=types__pb2.GetBalanceRequest.FromString,
          response_serializer=types__pb2.GetBalanceResponse.SerializeToString,
      ),
      'TransferEther': grpc.unary_unary_rpc_method_handler(
          servicer.TransferEther,
          request_deserializer=types__pb2.TransferRequest.FromString,
          response_serializer=types__pb2.TransferResponse.SerializeToString,
      ),
      'TransferToken': grpc.unary_unary_rpc_method_handler(
          servicer.TransferToken,
          request_deserializer=types__pb2.TransferRequest.FromString,
          response_serializer=types__pb2.TransferResponse.SerializeToString,
      ),
      'SendTransaction': grpc.unary_unary_rpc_method_handler(
          servicer.SendTransaction,
          request_deserializer=types__pb2.SendTransactionRequest.FromString,
          response_serializer=types__pb2.SendTransactionResponse.SerializeToString,
      ),
      'ExecuteZeroExTransaction': grpc.unary_unary_rpc_method_handler(
          servicer.ExecuteZeroExTransaction,
          request_deserializer=types__pb2.ExecuteZeroExTransactionRequest.FromString,
          response_serializer=types__pb2.ExecuteZeroExTransactionResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'HotWallet', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))


class WatcherStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.WatchTransaction = channel.unary_unary(
        '/Watcher/WatchTransaction',
        request_serializer=types__pb2.WatchTransactionRequest.SerializeToString,
        response_deserializer=types__pb2.WatchTransactionResponse.FromString,
        )


class WatcherServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def WatchTransaction(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_WatcherServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'WatchTransaction': grpc.unary_unary_rpc_method_handler(
          servicer.WatchTransaction,
          request_deserializer=types__pb2.WatchTransactionRequest.FromString,
          response_serializer=types__pb2.WatchTransactionResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'Watcher', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))


class ExchangeManagerStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.CreateOrder = channel.unary_unary(
        '/ExchangeManager/CreateOrder',
        request_serializer=types__pb2.ExchangeCreateOrderRequest.SerializeToString,
        response_deserializer=types__pb2.ExchangeOrderResponse.FromString,
        )
    self.GetOrder = channel.unary_unary(
        '/ExchangeManager/GetOrder',
        request_serializer=types__pb2.ExchangeOrderRequest.SerializeToString,
        response_deserializer=types__pb2.ExchangeOrderResponse.FromString,
        )
    self.GetOpenOrders = channel.unary_unary(
        '/ExchangeManager/GetOpenOrders',
        request_serializer=google_dot_protobuf_dot_wrappers__pb2.StringValue.SerializeToString,
        response_deserializer=types__pb2.ExchangeOrderArrayResponse.FromString,
        )
    self.CancelOrder = channel.unary_unary(
        '/ExchangeManager/CancelOrder',
        request_serializer=types__pb2.ExchangeOrderRequest.SerializeToString,
        response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
        )


class ExchangeManagerServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def CreateOrder(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetOrder(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetOpenOrders(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CancelOrder(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_ExchangeManagerServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'CreateOrder': grpc.unary_unary_rpc_method_handler(
          servicer.CreateOrder,
          request_deserializer=types__pb2.ExchangeCreateOrderRequest.FromString,
          response_serializer=types__pb2.ExchangeOrderResponse.SerializeToString,
      ),
      'GetOrder': grpc.unary_unary_rpc_method_handler(
          servicer.GetOrder,
          request_deserializer=types__pb2.ExchangeOrderRequest.FromString,
          response_serializer=types__pb2.ExchangeOrderResponse.SerializeToString,
      ),
      'GetOpenOrders': grpc.unary_unary_rpc_method_handler(
          servicer.GetOpenOrders,
          request_deserializer=google_dot_protobuf_dot_wrappers__pb2.StringValue.FromString,
          response_serializer=types__pb2.ExchangeOrderArrayResponse.SerializeToString,
      ),
      'CancelOrder': grpc.unary_unary_rpc_method_handler(
          servicer.CancelOrder,
          request_deserializer=types__pb2.ExchangeOrderRequest.FromString,
          response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'ExchangeManager', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))