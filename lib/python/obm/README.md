# OBM Python client

## Generate
```bash
python -m grpc_tools.protoc -I../../../proto --python_out=. --grpc_python_out=. ../../../proto/obm.proto
```

## Usage
```python
import grpc
from obm_pb2_grpc import OrderBookManagerStub
from obm_pb2 import OrderBookRequest

# Connect to the OBM server
channel = grpc.insecure_channel('localhost:8000')
stub = OrderBookManagerStub(channel)

# Build the request
req = OrderBookRequest(exchange = "binance", symbol = "BTC/USDT")

# Call the server
response = stub.OrderBook(req)
print(response)
```
