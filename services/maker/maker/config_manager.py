import os
import json

CONFIG_FILE = os.environ.get('CONFIG_FILE', 'maker/config.json')

class ConfigManager():

    address_to_ticker = {}
    ticker_to_pricing_data = {}

    def __init__(self) -> None:
        with open(CONFIG_FILE) as f:
            asset_config_f = json.load(f)
            for asset_data in asset_config_f['assets']:
                self.address_to_ticker[asset_data['address']] = asset_data['symbol']
                self.ticker_to_pricing_data[asset_data['symbol']] = asset_data['pricing_data']
                self.ticker_to_pricing_data[asset_data['symbol']]['decimals'] = asset_data['decimals']
            self.premium = float(asset_config_f['premium'])
            self.validity_length = int(asset_config_f['validity_length'])
            self.exchange_fees = asset_config_f['exchange_fees']
            for exchange in self.exchange_fees.keys():
                self.exchange_fees[exchange] = float(self.exchange_fees[exchange])

    def get_ticker_with_address(self, address: str) -> str:
        return self.address_to_ticker[address]
