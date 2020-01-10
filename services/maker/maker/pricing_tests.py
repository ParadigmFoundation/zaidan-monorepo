import PricingUtils as pricing_utils



def test_fixed_rate_price():
    assert pricing_utils.calculate_quote('ETH', 'DAI', 1, None) == {'maker_size': 1, 'taker_size': 101.202}
    assert pricing_utils.calculate_quote('DAI', 'ETH', None, 1) == {'maker_size': 98.40797407776671, 'taker_size': 1}
    assert pricing_utils.calculate_quote('ETH', 'DAI', None, 100) == {'maker_size': 0.9881227643722457, 'taker_size': 100}
    assert pricing_utils.calculate_quote('DAI', 'ETH', 100, None) == {'maker_size': 100, 'taker_size': 1.0161778142183395}


def test_direct_price():
    assert pricing_utils.calculate_quote('ZRX', 'DAI', 100, None) == {'maker_size': 100, 'taker_size': 31.062}
    assert pricing_utils.calculate_quote('DAI', 'ZRX', None, 100) == {'maker_size': 28.82657826520439, 'taker_size': 100}
    assert pricing_utils.calculate_quote('ZRX', 'DAI', None, 10) == {'maker_size': 32.19367716180542, 'taker_size': 10}
    assert pricing_utils.calculate_quote('DAI', 'ZRX', 10, None) == {'maker_size': 10, 'taker_size': 34.69020814055711}


def test_fully_implied_price():
    assert pricing_utils.calculate_quote('LINK', 'ZRX', 10, None) == {'maker_size': 10, 'taker_size': 69.31034482758619}
    assert pricing_utils.calculate_quote('ZRX', 'LINK', None, 10) == {'maker_size': 63.809539890281634, 'taker_size': 10}
    assert pricing_utils.calculate_quote('LINK', 'ZRX', None, 50) == {'maker_size': 7.170776314578263, 'taker_size': 50}
    assert pricing_utils.calculate_quote('ZRX', 'LINK', 50, None) == {'maker_size': 50, 'taker_size': 7.7889447236180915}


#test_direct_price()
#test_fixed_rate_price()
#test_fully_implied_price()

print(pricing_utils.calculate_quote('ETH', 'DAI', 1, None))