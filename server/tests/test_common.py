from common.protocol import Protocol
from common.lottery import Lottery, KEY_AGENCIES_WAITING
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
        agencies_waiting = {KEY_AGENCIES_WAITING: []}
        lottery = Lottery(agencies_waiting, 1)
        # First message - agency_id
        (_,_,agency_id) = Protocol.apply_rcv_bets_protocol(socket_client)
        self.assertEqual(agency_id, "125")                
        # First bet
        (bets,_,_) = Protocol.apply_rcv_bets_protocol(socket_client, agency_id)
        self.assertEqual(bets[0].first_name, "Santiago Lionel")
        self.assertEqual(bets[0].document, "30904465")
        self.assertEqual(bets[0].birthdate.year, 1999)
        self.assertEqual(bets[0].number, 7574)
        Protocol.apply_store_bet(lottery, bets[0])
        # Second bet
        (bets,_,_) = Protocol.apply_rcv_bets_protocol(socket_client, agency_id)
        self.assertEqual(bets[0].agency, 125)
        self.assertEqual(bets[0].first_name, "Pedro Alberto")
        self.assertEqual(bets[0].document, "4090446")
        self.assertEqual(bets[0].birthdate.year, 1970)
        self.assertEqual(bets[0].number, 9987)
        Protocol.apply_store_bet(lottery, bets[0])
        # Third bet
        self.assertEqual(bets[1].agency, 125)
        self.assertEqual(bets[1].first_name, "Lionel Andrés")
        self.assertEqual(bets[1].document, "3019985")
        self.assertEqual(bets[1].birthdate.year, 1985)
        self.assertEqual(bets[1].number, 7574)
        Protocol.apply_store_bet(lottery, bets[1])
        # Exception on amount fields on message
        try:
            (_,_,_) = Protocol.apply_rcv_bets_protocol(socket_client, agency_id)
        except ValueError as e:
            self.assertIsNotNone(e)
        # No more bets
        (_,amount_bets,_) = Protocol.apply_rcv_bets_protocol(socket_client, agency_id)
        self.assertEqual(amount_bets, "3")
        # Response
        Protocol.apply_res_amount_bets_protocol(socket_client, amount_bets,0)
        # Get Winners
        # Is waiting?
        self.assertFalse(lottery.agency_is_waiting(agency_id))
        # Add agency
        lottery.add_agency(agency_id)
        # Is waiting?
        self.assertTrue(lottery.agency_is_waiting(agency_id))
        # Amount waiting
        self.assertEqual(lottery.get_amount_agency_wating(), 1)
        Protocol.apply_winners_protocol(socket_client, agency_id, lottery)
        self.assertEqual(len(lottery.get_winners_from_agency(agency_id)), 2)
        # Close socket
        socket_client.close()
        # Assert that the recv method was called
        mock_client_socket.recv.assert_called()
        # Only test the server:
        #self.assertTrue(False)

class MockClientFunc:

    VALUES = [
        '125',
        'Santiago Lionel,Lorca,30904465,1999-03-17,7574',
        'Pedro Alberto,Pascual,4090446,1970-05-12,9987;Lionel Andrés,Messi,3019985,1985-05-24,7574',
        'Clubber,Lang,9519985,1952-05-21,7891,pugilista',
        'EOF,3',
        'WINNERS'
    ]

    def __init__(self):
        self.index = 0
        self.values = bytearray()
        for value in self.VALUES :
            self.values += self.parse_data_to_message(value)

    def return_value(self, bytes: int):
        to_return = self.values[self.index: self.index+bytes]        
        self.index += bytes
        return to_return

    def send_value(self, data: bytes):
        return True
    
    def parse_data_to_message(self, str_data: str):
        """ Generates a string as bytes in the format: <length in bytes str_data><str_data>. The data is encoded as UTF-8.
        Data length is 2 bytes integer as bigendian.            
        
            Examples
            
            str_data="Hola Mundo" <10 in bytes><Hola Mundo in bytes as Utf-8>

        """
        data = str_data.encode(Protocol.CODIFICATION)
        length = len(data)
        to_send = length.to_bytes(Protocol.AMOUNT_BYTES_LENGTH_MESSAGE, Protocol.BYTE_ORDER)+data[:]
        return to_send

if __name__ == '__main__':
    unittest.main()
