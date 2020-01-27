# Python gRPC Types & Services (`libgrpc`)

## Generate
```bash
make gen
```

## Start OBM

```bash
# (from repository root)
make docker
docker run -d --rm -p 8000:8000 gcr.io/zaidan-io/obm --bind 0.0.0.0:8000 --exchange binance:BTC/USD
```

## Usage
```python
import grpc
import json
from libgrpc.services_pb2_grpc import OrderBookManagerStub
from libgrpc.types_pb2 import OrderBookRequest

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

```
