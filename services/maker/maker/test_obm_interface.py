import unittest
from obm_interface import OBMInterface

class TestOBMInterface(unittest.TestCase):

    def test_initialize_placeholder_env_get_book(self) -> None:
        obm_interface = OBMInterface()
        obm_interface.set_env("PLACEHOLDER")
        self.assertEqual(obm_interface.env, "PLACEHOLDER")
        self.assertEqual(obm_interface.get_order_book("COINBASE", "ETH/USD", "bids"), [[99, 100], [98, 200], [97, 300]])

if __name__ == '__main__':
    unittest.main()
