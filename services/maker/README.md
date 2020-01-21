# Service: `maker`

The maker service is responsible for pricing and risk management.

## Services/Protobuf Code Generation

Generate the necessary protobuf/gRPC code with the following make target.

```sh
make gen
```
## Config file path

From the top level of the repository run the following:

```sh
export ASSET_CONFIG_FILE=$(pwd)/services/maker/maker/asset_config.json
```