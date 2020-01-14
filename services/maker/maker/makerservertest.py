import grpc
from types_pb2_grpc import MakerStub
from types_pb2 import GetQuoteRequest, CheckQuoteRequest

# Connect to the OBM server
channel = grpc.insecure_channel('localhost:50051')
stub = MakerStub(channel)

# Build the request
req = GetQuoteRequest(maker_asset="0x0b1ba0af832d7c05fd64161e0db78e85978e8082",
                      taker_asset="0x34d402f14d58e001d8efbe6585051bf9706aa064", maker_size=str(1*10**18))

# Call the server
response = stub.GetQuote(req)
id = response.quote_id
print(response)

req = CheckQuoteRequest(quote_id='a')

response = stub.CheckQuote(req)
print(response)

req = CheckQuoteRequest(quote_id=id)

response = stub.CheckQuote(req)
print(response)

