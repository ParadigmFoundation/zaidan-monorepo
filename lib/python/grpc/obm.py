import grpc
import json
from services_pb2_grpc import OrderBookManagerStub
from types_pb2 import OrderBookRequest

# Connect to the OBM server
channel = grpc.insecure_channel('localhost:8000')
stub = OrderBookManagerStub(channel)

# Build the request
req = OrderBookRequest(exchange = "binance", symbol = "BTC/USDT")

# Call the server
response = stub.OrderBook(req)

bid_book = [[x.price, x.quantity] for x in response.bids]
ask_book = [[x.price, x.quantity] for x in response.asks]

print(json.dumps({"asks":ask_book, "bids":bid_book}))
