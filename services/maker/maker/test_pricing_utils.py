import unittest
from pricing_utils import PricingUtils

class TestPricingUtils(unittest.TestCase):
    
    def test_fixed_rate_price(self) -> None:
        pricing_utils = PricingUtils()
        self.assertEqual(pricing_utils.calculate_quote('ETH', 'DAI', 1, None, True), {'maker_size': 1, 'taker_size': 101.202})
        self.assertEqual(pricing_utils.calculate_quote('DAI', 'ETH', None, 1, True), {'maker_size': 98.40797407776671, 'taker_size': 1})
        self.assertEqual(pricing_utils.calculate_quote('ETH', 'DAI', None, 100, True), {'maker_size': 0.9881227643722457, 'taker_size': 100})
        self.assertEqual(pricing_utils.calculate_quote('DAI', 'ETH', 100, None, True), {'maker_size': 100, 'taker_size': 1.0161778142183395})
    
    def test_direct_price(self) -> None:
        pricing_utils = PricingUtils()
        self.assertEqual(pricing_utils.calculate_quote('ZRX', 'DAI', 100, None, True), {'maker_size': 100, 'taker_size': 31.062})
        self.assertEqual(pricing_utils.calculate_quote('DAI', 'ZRX', None, 100, True), {'maker_size': 28.82657826520439, 'taker_size': 100})
        self.assertEqual(pricing_utils.calculate_quote('ZRX', 'DAI', None, 10, True), {'maker_size': 32.19367716180542, 'taker_size': 10})
        self.assertEqual(pricing_utils.calculate_quote('DAI', 'ZRX', 10, None, True), {'maker_size': 10, 'taker_size': 34.69020814055711})
    
    def test_fully_implied_price(self) -> None:
        pricing_utils = PricingUtils()
        self.assertEqual(pricing_utils.calculate_quote('LINK', 'ZRX', 10, None, True), {'maker_size': 10, 'taker_size': 69.31034482758619})
        self.assertEqual(pricing_utils.calculate_quote('ZRX', 'LINK', None, 10, True), {'maker_size': 63.809539890281634, 'taker_size': 10})
        self.assertEqual(pricing_utils.calculate_quote('LINK', 'ZRX', None, 50, True), {'maker_size': 7.170776314578263, 'taker_size': 50})
        self.assertEqual(pricing_utils.calculate_quote('ZRX', 'LINK', 50, None, True), {'maker_size': 50, 'taker_size': 7.7889447236180915})

if __name__ == '__main__':
    unittest.main()
