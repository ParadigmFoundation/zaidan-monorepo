import unittest
import sys

try:
    sys.path.append('../maker')
    from config_manager import ConfigManager
except:
    try:
        sys.path.append('../maker/maker')
        from config_manager import ConfigManager
    except Exception as error:
        raise Exception('failed to import asset data: {}'.format(error.args))

class TestConfigManager(unittest.TestCase):

    def test_config_manager_get_ticker_with_address(self) -> None:
        WETH_ADDRESS = '0x0b1ba0af832d7c05fd64161e0db78e85978e8082'
        config_manager = ConfigManager()
        self.assertEqual(config_manager.get_ticker_with_address(WETH_ADDRESS), 'WETH')

if __name__ == '__main__':
    unittest.main()
