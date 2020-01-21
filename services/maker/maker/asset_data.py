import os
import json

ASSET_CONFIG_FILE = os.environ.get('ASSET_CONFIG_FILE', 'asset_config.json')

class AssetData():

    address_to_ticker = {}

    def __init__(self) -> None:
        with open(ASSET_CONFIG_FILE) as f:
            asset_config_f = json.load(f)
            for asset_data in asset_config_f['assets']:
                self.address_to_ticker[asset_data['address']] = asset_data['symbol']

    def get_ticker_with_address(self, address: str) -> str:
        return self.address_to_ticker[address]
