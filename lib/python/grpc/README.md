# Python gRPC Types

## Generate
```bash
python -m grpc_tools.protoc -I../../../proto --python_out=. --grpc_python_out=. ../../../proto/types.proto
```

## Usage
```python
import grpc
from types_pb2_grpc import OrderBookManagerStub
from types_pb2 import OrderBookRequest

# Connect to the OBM server
channel = grpc.insecure_channel('localhost:8000')
stub = OrderBookManagerStub(channel)

# Build the request
req = OrderBookRequest(exchange = "binance", symbol = "BTC/USDT")

# Call the server
response = stub.OrderBook(req)
print(response)
```
