import os
import json
from typing import Optional

CONFIG_FILE = os.environ.get('CONFIG_FILE', 'maker/config.json')

class ConfigManager():

    address_to_ticker = {}
    ticker_to_pricing_data = {}
    markets = []

    def __init__(self) -> None:
        with open(CONFIG_FILE) as f:
            config_f = json.load(f)
            all_markets = []
            for asset in config_f['assets']:
                self.address_to_ticker[asset['address']] = asset['symbol']
                self.ticker_to_pricing_data[asset['symbol']] = asset['pricing_data']
                self.ticker_to_pricing_data[asset['symbol']]['decimals'] = asset['decimals']
                markets_entry = []
                markets_entry.append(asset['address'])
                pair_assets = []
                for pair_asset in config_f['assets']:
                    if not asset['address'] == pair_asset['address']:
                        pair_assets.append(pair_asset['address'])
                markets_entry.append(pair_assets)
                markets_entry.append(asset['min_size'])
                markets_entry.append(asset['max_size'])
                all_markets.append(markets_entry)
            self.markets = all_markets
            self.premium = float(config_f['premium'])
            self.validity_length = int(config_f['validity_length'])
            self.exchange_fees = config_f['exchange_fees']
            for exchange in self.exchange_fees.keys():
                self.exchange_fees[exchange] = float(self.exchange_fees[exchange])



    def get_ticker_with_address(self, address: str) -> str:
        return self.address_to_ticker[address]

    def get_markets(self, maker_asset_address:Optional[str]=None, taker_asset_address:Optional[str]=None) -> list:
        markets_to_return = []
        for market in self.markets:
            if maker_asset_address:
                if market[0] == maker_asset_address:
                    if taker_asset_address:
                        if market[1] == taker_asset_address:
                            return [market]
                    else:
                        return [market]
            elif taker_asset_address:
                if taker_asset_address in market[1]:
                    markets_to_return.append(market)
            else:
                markets_to_return.append(market)

        return markets_to_return
