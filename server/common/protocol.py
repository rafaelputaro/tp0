import logging
from abc import ABC
import socket
from common.utils import Bet, store_bets

class Protocol(ABC):

    ACTION_RECEIVE = "receive_message"

    ACTION_RESPONSE = "response_message"

    ACTION_STORE_BET = "apuesta_almacenada"

    CODIFICATION = 'utf-8'

    BYTE_ORDER = 'big'

    AMOUNT_BYTES_LENGTH_MESSAGE = 2

    @classmethod
    def apply_rcv_protocol(cls, client_sock: socket):
        """ Receive a bet from the client
        """        
        length = int.from_bytes(client_sock.recv(Protocol.AMOUNT_BYTES_LENGTH_MESSAGE), byteorder=Protocol.BYTE_ORDER)
        msg = client_sock.recv(length).decode(Protocol.CODIFICATION).strip()
        logging.debug(f'action: {Protocol.ACTION_RECEIVE} | result: success | msg: {msg}')
        bet: Bet = Bet(*msg.split(','))
        return bet

    @classmethod
    def apply_store_bet(cls, bet: Bet):
        """ Support the bet locally
        """        
        store_bets([bet])
        logging.info(f'action: {Protocol.ACTION_STORE_BET} | result: success | dni: {bet.document} | numero: {bet.document}')

    @classmethod
    def apply_res_protocol(cls, client_sock: socket, bet: Bet):
        """ Support the bet locally and send a confirmation message to the client
        """        
        Protocol.apply_store_bet(bet)
        client_sock.send("{}\n".format(bet.number).encode('utf-8'))
        logging.debug(f'action: {Protocol.ACTION_RESPONSE} | result: success | msg: {bet.number}')

    @classmethod
    def parse_data_to_message(cls, str_data: str):
        """ Generates a string as bytes in the format: <length in bytes str_data><str_data>. The data is encoded as UTF-8.
        Data length is 2 bytes integer as bigendian.            
        
            Examples
            
            >> str_data="Hola Mundo" <10 in bytes><Hola Mundo in bytes as Utf-8>

        """
        data = str_data.encode(Protocol.CODIFICATION)
        length = len(data)
        return length.to_bytes(Protocol.AMOUNT_BYTES_LENGTH_MESSAGE, Protocol.BYTE_ORDER)+data
