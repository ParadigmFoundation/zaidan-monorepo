import os
import json

ASSET_CONFIG_FILE = os.environ.get('ASSET_CONFIG_FILE', 'asset_config.json')
path_to_current_file = os.path.realpath(__file__)
current_directory = os.path.split(path_to_current_file)[0]
path_to_file = os.path.join(current_directory, ASSET_CONFIG_FILE)

class AssetData():

    address_to_ticker = {}
    ticker_to_pricing_data = {}

    def __init__(self) -> None:
        try:
            with open(path_to_file) as f:
                asset_config_f = json.load(f)
                for asset_data in asset_config_f['assets']:
                    self.address_to_ticker[asset_data['address']] = asset_data['symbol']
                    self.ticker_to_pricing_data[asset_data['symbol']] = asset_data['pricing_data']
                    self.ticker_to_pricing_data[asset_data['symbol']]['decimals'] = asset_data['decimals']
        except:
            with open(ASSET_CONFIG_FILE) as f:
                asset_config_f = json.load(f)
                for asset_data in asset_config_f['assets']:
                    self.address_to_ticker[asset_data['address']] = asset_data['symbol']
                    self.ticker_to_pricing_data[asset_data['symbol']] = asset_data['pricing_data']
                    self.ticker_to_pricing_data[asset_data['symbol']]['decimals'] = asset_data['decimals']

    def get_ticker_with_address(self, address: str) -> str:
        return self.address_to_ticker[address]
