# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: services.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
from google.protobuf import wrappers_pb2 as google_dot_protobuf_dot_wrappers__pb2
import types_pb2 as types__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='services.proto',
  package='',
  syntax='proto3',
  serialized_options=_b('Z\004grpc'),
  serialized_pb=_b('\n\x0eservices.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1egoogle/protobuf/wrappers.proto\x1a\x0btypes.proto2\x81\x01\n\x10OrderBookManager\x12\x32\n\tOrderBook\x12\x11.OrderBookRequest\x1a\x12.OrderBookResponse\x12\x39\n\x07Updates\x12\x18.OrderBookUpdatesRequest\x1a\x12.OrderBookResponse0\x01\x32\xbb\x01\n\x05Maker\x12/\n\x08GetQuote\x12\x10.GetQuoteRequest\x1a\x11.GetQuoteResponse\x12\x35\n\nCheckQuote\x12\x12.CheckQuoteRequest\x1a\x13.CheckQuoteResponse\x12J\n\x11OrderStatusUpdate\x12\x19.OrderStatusUpdateRequest\x1a\x1a.OrderStatusUpdateResponse2\x8a\x05\n\tHotWallet\x12\x38\n\x0b\x43reateOrder\x12\x13.CreateOrderRequest\x1a\x14.CreateOrderResponse\x12>\n\rValidateOrder\x12\x15.ValidateOrderRequest\x1a\x16.ValidateOrderResponse\x12;\n\x0cGetAllowance\x12\x14.GetAllowanceRequest\x1a\x15.GetAllowanceResponse\x12;\n\x0cSetAllowance\x12\x14.SetAllowanceRequest\x1a\x15.SetAllowanceResponse\x12:\n\x0fGetTokenBalance\x12\x12.GetBalanceRequest\x1a\x13.GetBalanceResponse\x12:\n\x0fGetEtherBalance\x12\x12.GetBalanceRequest\x1a\x13.GetBalanceResponse\x12\x34\n\rTransferEther\x12\x10.TransferRequest\x1a\x11.TransferResponse\x12\x34\n\rTransferToken\x12\x10.TransferRequest\x1a\x11.TransferResponse\x12\x44\n\x0fSendTransaction\x12\x17.SendTransactionRequest\x1a\x18.SendTransactionResponse\x12_\n\x18\x45xecuteZeroExTransaction\x12 .ExecuteZeroExTransactionRequest\x1a!.ExecuteZeroExTransactionResponse2R\n\x07Watcher\x12G\n\x10WatchTransaction\x12\x18.WatchTransactionRequest\x1a\x19.WatchTransactionResponse2\x9a\x02\n\x0f\x45xchangeManager\x12\x42\n\x0b\x43reateOrder\x12\x1b.ExchangeCreateOrderRequest\x1a\x16.ExchangeOrderResponse\x12\x39\n\x08GetOrder\x12\x15.ExchangeOrderRequest\x1a\x16.ExchangeOrderResponse\x12J\n\rGetOpenOrders\x12\x1c.google.protobuf.StringValue\x1a\x1b.ExchangeOrderArrayResponse\x12<\n\x0b\x43\x61ncelOrder\x12\x15.ExchangeOrderRequest\x1a\x16.google.protobuf.Empty2?\n\x06Hedger\x12\x35\n\nHedgeOrder\x12\x12.HedgeOrderRequest\x1a\x13.HedgeOrderResponseB\x06Z\x04grpcb\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_empty__pb2.DESCRIPTOR,google_dot_protobuf_dot_wrappers__pb2.DESCRIPTOR,types__pb2.DESCRIPTOR,])



_sym_db.RegisterFileDescriptor(DESCRIPTOR)


DESCRIPTOR._options = None

_ORDERBOOKMANAGER = _descriptor.ServiceDescriptor(
  name='OrderBookManager',
  full_name='OrderBookManager',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=93,
  serialized_end=222,
  methods=[
  _descriptor.MethodDescriptor(
    name='OrderBook',
    full_name='OrderBookManager.OrderBook',
    index=0,
    containing_service=None,
    input_type=types__pb2._ORDERBOOKREQUEST,
    output_type=types__pb2._ORDERBOOKRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='Updates',
    full_name='OrderBookManager.Updates',
    index=1,
    containing_service=None,
    input_type=types__pb2._ORDERBOOKUPDATESREQUEST,
    output_type=types__pb2._ORDERBOOKRESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_ORDERBOOKMANAGER)

DESCRIPTOR.services_by_name['OrderBookManager'] = _ORDERBOOKMANAGER


_MAKER = _descriptor.ServiceDescriptor(
  name='Maker',
  full_name='Maker',
  file=DESCRIPTOR,
  index=1,
  serialized_options=None,
  serialized_start=225,
  serialized_end=412,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetQuote',
    full_name='Maker.GetQuote',
    index=0,
    containing_service=None,
    input_type=types__pb2._GETQUOTEREQUEST,
    output_type=types__pb2._GETQUOTERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='CheckQuote',
    full_name='Maker.CheckQuote',
    index=1,
    containing_service=None,
    input_type=types__pb2._CHECKQUOTEREQUEST,
    output_type=types__pb2._CHECKQUOTERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='OrderStatusUpdate',
    full_name='Maker.OrderStatusUpdate',
    index=2,
    containing_service=None,
    input_type=types__pb2._ORDERSTATUSUPDATEREQUEST,
    output_type=types__pb2._ORDERSTATUSUPDATERESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_MAKER)

DESCRIPTOR.services_by_name['Maker'] = _MAKER


_HOTWALLET = _descriptor.ServiceDescriptor(
  name='HotWallet',
  full_name='HotWallet',
  file=DESCRIPTOR,
  index=2,
  serialized_options=None,
  serialized_start=415,
  serialized_end=1065,
  methods=[
  _descriptor.MethodDescriptor(
    name='CreateOrder',
    full_name='HotWallet.CreateOrder',
    index=0,
    containing_service=None,
    input_type=types__pb2._CREATEORDERREQUEST,
    output_type=types__pb2._CREATEORDERRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ValidateOrder',
    full_name='HotWallet.ValidateOrder',
    index=1,
    containing_service=None,
    input_type=types__pb2._VALIDATEORDERREQUEST,
    output_type=types__pb2._VALIDATEORDERRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='GetAllowance',
    full_name='HotWallet.GetAllowance',
    index=2,
    containing_service=None,
    input_type=types__pb2._GETALLOWANCEREQUEST,
    output_type=types__pb2._GETALLOWANCERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='SetAllowance',
    full_name='HotWallet.SetAllowance',
    index=3,
    containing_service=None,
    input_type=types__pb2._SETALLOWANCEREQUEST,
    output_type=types__pb2._SETALLOWANCERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='GetTokenBalance',
    full_name='HotWallet.GetTokenBalance',
    index=4,
    containing_service=None,
    input_type=types__pb2._GETBALANCEREQUEST,
    output_type=types__pb2._GETBALANCERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='GetEtherBalance',
    full_name='HotWallet.GetEtherBalance',
    index=5,
    containing_service=None,
    input_type=types__pb2._GETBALANCEREQUEST,
    output_type=types__pb2._GETBALANCERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='TransferEther',
    full_name='HotWallet.TransferEther',
    index=6,
    containing_service=None,
    input_type=types__pb2._TRANSFERREQUEST,
    output_type=types__pb2._TRANSFERRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='TransferToken',
    full_name='HotWallet.TransferToken',
    index=7,
    containing_service=None,
    input_type=types__pb2._TRANSFERREQUEST,
    output_type=types__pb2._TRANSFERRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='SendTransaction',
    full_name='HotWallet.SendTransaction',
    index=8,
    containing_service=None,
    input_type=types__pb2._SENDTRANSACTIONREQUEST,
    output_type=types__pb2._SENDTRANSACTIONRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ExecuteZeroExTransaction',
    full_name='HotWallet.ExecuteZeroExTransaction',
    index=9,
    containing_service=None,
    input_type=types__pb2._EXECUTEZEROEXTRANSACTIONREQUEST,
    output_type=types__pb2._EXECUTEZEROEXTRANSACTIONRESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_HOTWALLET)

DESCRIPTOR.services_by_name['HotWallet'] = _HOTWALLET


_WATCHER = _descriptor.ServiceDescriptor(
  name='Watcher',
  full_name='Watcher',
  file=DESCRIPTOR,
  index=3,
  serialized_options=None,
  serialized_start=1067,
  serialized_end=1149,
  methods=[
  _descriptor.MethodDescriptor(
    name='WatchTransaction',
    full_name='Watcher.WatchTransaction',
    index=0,
    containing_service=None,
    input_type=types__pb2._WATCHTRANSACTIONREQUEST,
    output_type=types__pb2._WATCHTRANSACTIONRESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_WATCHER)

DESCRIPTOR.services_by_name['Watcher'] = _WATCHER


_EXCHANGEMANAGER = _descriptor.ServiceDescriptor(
  name='ExchangeManager',
  full_name='ExchangeManager',
  file=DESCRIPTOR,
  index=4,
  serialized_options=None,
  serialized_start=1152,
  serialized_end=1434,
  methods=[
  _descriptor.MethodDescriptor(
    name='CreateOrder',
    full_name='ExchangeManager.CreateOrder',
    index=0,
    containing_service=None,
    input_type=types__pb2._EXCHANGECREATEORDERREQUEST,
    output_type=types__pb2._EXCHANGEORDERRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='GetOrder',
    full_name='ExchangeManager.GetOrder',
    index=1,
    containing_service=None,
    input_type=types__pb2._EXCHANGEORDERREQUEST,
    output_type=types__pb2._EXCHANGEORDERRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='GetOpenOrders',
    full_name='ExchangeManager.GetOpenOrders',
    index=2,
    containing_service=None,
    input_type=google_dot_protobuf_dot_wrappers__pb2._STRINGVALUE,
    output_type=types__pb2._EXCHANGEORDERARRAYRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='CancelOrder',
    full_name='ExchangeManager.CancelOrder',
    index=3,
    containing_service=None,
    input_type=types__pb2._EXCHANGEORDERREQUEST,
    output_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_EXCHANGEMANAGER)

DESCRIPTOR.services_by_name['ExchangeManager'] = _EXCHANGEMANAGER


_HEDGER = _descriptor.ServiceDescriptor(
  name='Hedger',
  full_name='Hedger',
  file=DESCRIPTOR,
  index=5,
  serialized_options=None,
  serialized_start=1436,
  serialized_end=1499,
  methods=[
  _descriptor.MethodDescriptor(
    name='HedgeOrder',
    full_name='Hedger.HedgeOrder',
    index=0,
    containing_service=None,
    input_type=types__pb2._HEDGEORDERREQUEST,
    output_type=types__pb2._HEDGEORDERRESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_HEDGER)

DESCRIPTOR.services_by_name['Hedger'] = _HEDGER

# @@protoc_insertion_point(module_scope)
