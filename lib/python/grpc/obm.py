import grpc
import json
from services_pb2_grpc import OrderBookManagerStub
from types_pb2 import OrderBookRequest

# Connect to the OBM server
channel = grpc.insecure_channel('localhost:8001')
stub = OrderBookManagerStub(channel)

# Build the request
req = OrderBookRequest(exchange = "coinbase", symbol = "ETH/USD")

# Call the server
response = stub.OrderBook(req)

bid_book = [[x.price, x.quantity] for x in response.bids]
ask_book = [[x.price, x.quantity] for x in response.asks]

print(json.dumps({"asks":ask_book, "bids":bid_book}))