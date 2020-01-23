import unittest
import sys

try:
    sys.path.append('../maker')
    from asset_data import AssetData
except:
    try:
        sys.path.append('../maker/maker')
        from asset_data import AssetData
    except Exception as error:
        raise Exception('failed to import asset data: {}'.format(error.args))

class TestAssetData(unittest.TestCase):

    def test_asset_data_get_ticker_with_address(self) -> None:
        WETH_ADDRESS = '0x0b1ba0af832d7c05fd64161e0db78e85978e8082'
        asset_data = AssetData()
        self.assertEqual(asset_data.get_ticker_with_address(WETH_ADDRESS), 'WETH')

if __name__ == '__main__':
    unittest.main()