import os
import json
from typing import Optional

CONFIG_FILE = os.environ.get('CONFIG_FILE', 'maker/config.json')
CHAIN_ID = int(os.environ.get('CHAIN_ID', 1337))

class ConfigManager():

    address_to_ticker = {}
    ticker_to_pricing_data = {}
    markets = []

    def __init__(self) -> None:
        with open(CONFIG_FILE) as f:
            config_f = json.load(f)
            all_markets = []
            for asset in config_f['assets']:
                for chain in asset['deployments']:
                    if chain['chain_id'] == CHAIN_ID:
                        asset['address'] = chain['address']
            for asset in config_f['assets']:
                self.address_to_ticker[asset['address']] = asset['symbol']
                self.ticker_to_pricing_data[asset['symbol']] = asset['pricing_data']
                self.ticker_to_pricing_data[asset['symbol']]['decimals'] = asset['decimals']
                if not asset['symbol'] == 'ETH':
                    markets_entry = {}
                    markets_entry['maker_asset_address'] = asset['address']
                    pair_assets = []
                    for pair_asset in config_f['assets']:
                        if not asset['address'] == pair_asset['address'] and not pair_asset['symbol'] == 'ETH':
                            pair_assets.append(pair_asset['address'])
                    markets_entry['taker_asset_addresses'] = pair_assets
                    markets_entry['min_size'] = asset['min_size']
                    markets_entry['max_size'] = asset['max_size']
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
        if not maker_asset_address and not taker_asset_address:
            for market in self.markets:
                markets_to_return.append(market)
        elif maker_asset_address and not taker_asset_address:
            for market in self.markets:
                if market['maker_asset_address'] == maker_asset_address:
                    return [market]
        elif taker_asset_address and not maker_asset_address:
            for market in self.markets:
                if taker_asset_address in market['taker_asset_addresses']:
                    markets_to_return.append(market)
        else:
            for market in self.markets:
                if market['maker_asset_address'] == maker_asset_address:
                    if taker_asset_address in market['taker_asset_addresses']:
                        return [market]

        return markets_to_return