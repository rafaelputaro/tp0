import logging
from abc import ABC
import socket
from common.utils import Bet, store_bets

class Protocol(ABC):

    ACTION_RECEIVE = "receive_message"

    ACTION_RESPONSE = "apuesta_recibida"

    ACTION_STORE_BET = "apuesta_almacenada"

    CODIFICATION = 'utf-8'

    BYTE_ORDER = 'big'

    AMOUNT_BYTES_LENGTH_MESSAGE = 2

    MSG_ERROR_ON_PARSE_BET = "Error on parse bet"

    BETS_DELIMITER = ';'    

    FIELDS_DELIMTER = ','

    @classmethod
    def apply_rcv_protocol(cls, client_sock: socket, agency = ""):
        """ Receive a bet from the client
            Returns:
            1) Bets, None, None  -> Keep reading
            2) None, None, id_agency -> First message with agency id
            3) None, amount_bets, None -> End of message
            4) Exception on parse or amount fields on message
        """
        length = int.from_bytes(client_sock.recv(Protocol.AMOUNT_BYTES_LENGTH_MESSAGE), byteorder=Protocol.BYTE_ORDER)
        msg = client_sock.recv(length).decode(Protocol.CODIFICATION).strip()
        splittedInBets = list(map(str, msg.split(Protocol.BETS_DELIMITER)))
        lengthSpInBets = len(splittedInBets)
        if lengthSpInBets == 1:
            return Protocol.parse_possible_not_bet(splittedInBets[0], agency)        
        else:
            return Protocol.parse_possibles_bets(splittedInBets, agency), None, None

    @classmethod
    def parse_possible_not_bet(cls, toEvaluate: str, agency = ""):
        """ Checks if the value passed by parameter is a bet or data about the message
        """
        splittedInFields = toEvaluate.split(Protocol.FIELDS_DELIMTER)
        lengthSpInFields = len(splittedInFields)
        if lengthSpInFields == 0:
            return None, None, None
        # ¿Is agencyId?        
        elif lengthSpInFields == 1:       
            return None, None, toEvaluate
        # ¿Is EOF?
        elif lengthSpInFields == 2:
            return None, splittedInFields[1], None
        # ¿Is a bet?
        else:
            bets = [Protocol.parse_bet(splittedInFields, agency)]
            return bets, None, None

    @classmethod
    def parse_bet(cls, betData: list[str], agency = ""):
        try:               
            return Bet(agency, betData[0], betData[1], betData[2], betData[3], betData[4])
        except Exception as e:
            #logging.debug(f'data: {len(betData)} | agency: {agency}')
            raise ValueError(f'{Protocol.MSG_ERROR_ON_PARSE_BET}')

    @classmethod
    def parse_possibles_bets(cls, possiblesBets: list[str], agency = ""):
        bets = []
        for possibleBet in possiblesBets:
            try:
                if len(possibleBet) > 0:
                    bets.append(Protocol.parse_bet(possibleBet.split(Protocol.FIELDS_DELIMTER), agency))
            except ValueError as e:
                logging.info(f'action: {Protocol.ACTION_PARSE_BET} | result: fail | msg: {e}')
        return bets 

    @classmethod
    def apply_store_bet(cls, bet: Bet):
        """ Support the bet locally
        """        
        store_bets([bet])
        logging.debug(f'action: {Protocol.ACTION_STORE_BET} | result: success | dni: {bet.document} | numero: {bet.number}')

    @classmethod
    def apply_store_bets(cls, bets: list[Bet]):
        """ Support the bets locally
        """        
        for bet in bets:
            Protocol.apply_store_bet(bet)

    @classmethod
    def apply_res_protocol(cls, client_sock: socket, amount_bets: str, amount_bets_expected: str):
        """ Send a confirmation message to the client if receive the amount bets expected
        """      
        client_sock.send("{}\n".format(amount_bets).encode('utf-8'))
        logging.debug(f'Leídos: {amount_bets} | Esperado: {amount_bets_expected}')
        result = "success" if (amount_bets == amount_bets_expected) else "failure"
        logging.info(f'action: {Protocol.ACTION_RESPONSE} | result: {result} | cantidad: {amount_bets}')