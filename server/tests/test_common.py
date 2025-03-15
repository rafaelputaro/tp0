from common.protocol import Protocol
from common.utils import *
import os
import unittest
from unittest.mock import patch
import socket

class TestUtils(unittest.TestCase):

    def tearDown(self):
        if os.path.exists(STORAGE_FILEPATH):
            os.remove(STORAGE_FILEPATH)

    def test_bet_init_must_keep_fields(self):
        b = Bet('1', 'first', 'last', '10000000','2000-12-20', 7500)
        self.assertEqual(1, b.agency)
        self.assertEqual('first', b.first_name)
        self.assertEqual('last', b.last_name)
        self.assertEqual('10000000', b.document)
        self.assertEqual(datetime.date(2000, 12, 20), b.birthdate)
        self.assertEqual(7500, b.number)

    def test_has_won_with_winner_number_must_be_true(self):
        b = Bet('1', 'first', 'last', 10000000,'2000-12-20', LOTTERY_WINNER_NUMBER)
        self.assertTrue(has_won(b))

    def test_has_won_with_winner_number_must_be_true(self):
        b = Bet('1', 'first', 'last', 10000000,'2000-12-20', LOTTERY_WINNER_NUMBER + 1)
        self.assertFalse(has_won(b))

    def test_store_bets_and_load_bets_keeps_fields_data(self):
        to_store = [Bet('1', 'first', 'last', '10000000','2000-12-20', 7500)]
        store_bets(to_store)
        from_load = list(load_bets())

        self.assertEqual(1, len(from_load))
        self._assert_equal_bets(to_store[0], from_load[0])


    def test_store_bets_and_load_bets_keeps_registry_order(self):
        to_store = [
            Bet('0', 'first_0', 'last_0', '10000000','2000-12-20', 7500),
            Bet('1', 'first_1', 'last_1', '10000001','2000-12-21', 7501),
        ]
        store_bets(to_store)
        from_load = list(load_bets())

        self.assertEqual(2, len(from_load))
        self._assert_equal_bets(to_store[0], from_load[0])
        self._assert_equal_bets(to_store[1], from_load[1])

    def _assert_equal_bets(self, b1, b2):
        self.assertEqual(b1.agency, b2.agency)
        self.assertEqual(b1.first_name, b2.first_name)
        self.assertEqual(b1.last_name, b2.last_name)
        self.assertEqual(b1.document, b2.document)
        self.assertEqual(b1.birthdate, b2.birthdate)
        self.assertEqual(b1.number, b2.number)

    @patch('socket.socket')
    def test_protocol(self, mock_socket):
        mock_client_func = MockClientFunc()     
        # Create a mock socket object
        mock_client_socket = mock_socket.return_value
        mock_client_socket.recv.side_effect = mock_client_func.return_value
        mock_client_socket.send.side_effect = mock_client_func.send_value
        # Connect
        socket_client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        socket_client.connect(('localhost', 12345))
        # Use protocol
        bet: Bet = Protocol.apply_rcv_protocol(socket_client)
        self.assertEqual(bet.agency, 125)
        self.assertEqual(bet.first_name, "Santiago Lionel")
        self.assertEqual(bet.document, "30904465")
        self.assertEqual(bet.birthdate.year, 1999)
        self.assertEqual(bet.number, 7574)
        bet = Protocol.apply_rcv_protocol(socket_client)
        Protocol.apply_res_protocol(socket_client, bet)
        self.assertEqual(bet.agency, 125)
        self.assertEqual(bet.first_name, "Pedro Alberto")
        self.assertEqual(bet.document, "4090446")
        self.assertEqual(bet.birthdate.year, 1970)
        self.assertEqual(bet.number, 9987)
        socket_client.close()
        # Assert that the recv method was called
        mock_client_socket.recv.assert_called()
        # Assert the expected result - Trick to only run test's
        #self.assertTrue(False)

class MockClientFunc:

    VALUES = ['125,Santiago Lionel,Lorca,30904465,1999-03-17,7574',
              '125,Pedro Alberto,Pascual,4090446,1970-05-12,9987',
              '125,Lionel Andrés,Messi,20119985,1985-05-24,2022']

    def __init__(self):
        self.index = 0
        self.values = bytearray()
        for value in self.VALUES :
            self.values += Protocol.parse_data_to_message(value)

    def return_value(self, bytes: int):

        to_return = self.values[self.index: self.index+bytes]        
        self.index += bytes
        return to_return

    def send_value(self, data: bytes):
        return True

if __name__ == '__main__':
    unittest.main()

