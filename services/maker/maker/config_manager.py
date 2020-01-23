import os
import json

CONFIG_FILE = os.environ.get('CONFIG_FILE', 'maker/config.json')

class ConfigManager():

    address_to_ticker = {}
    ticker_to_pricing_data = {}

    def __init__(self) -> None:
        with open(CONFIG_FILE) as f:
            config_f = json.load(f)
            for config_manager in config_f['assets']:
                self.address_to_ticker[config_manager['address']] = config_manager['symbol']
                self.ticker_to_pricing_data[config_manager['symbol']] = config_manager['pricing_data']
                self.ticker_to_pricing_data[config_manager['symbol']]['decimals'] = config_manager['decimals']
            self.premium = float(config_f['premium'])
            self.validity_length = int(config_f['validity_length'])
            self.exchange_fees = config_f['exchange_fees']
            for exchange in self.exchange_fees.keys():
                self.exchange_fees[exchange] = float(self.exchange_fees[exchange])

    def get_ticker_with_address(self, address: str) -> str:
        return self.address_to_ticker[address]
