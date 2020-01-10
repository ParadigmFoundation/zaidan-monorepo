import grpc
from types_pb2_grpc import MakerStub
from types_pb2 import GetQuoteRequest, CheckQuoteRequest

# Connect to the OBM server
channel = grpc.insecure_channel('localhost:50051')
stub = MakerStub(channel)

# Build the request
req = GetQuoteRequest(maker_asset='ETH', taker_asset='DAI', maker_size='1')

# Call the server
response = stub.GetQuote(req)
print(response)

req = CheckQuoteRequest(quote_id='a')

response = stub.CheckQuote(req)
print(response)
print(response.is_valid)

