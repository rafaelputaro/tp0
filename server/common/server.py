import socket
import time
import logging
import sys
import signal
import multiprocessing
from common.protocol import Protocol
from common.lottery import Lottery, KEY_AGENCIES_WAITING

SIGNAL_HANDLER_ACTION="received_a_signal"
CLOSE_SERVER_SOCKET_ACTION="closing_server_socket"
JOIN_PROCESS_ACTION="join_process"
READ_BET_ACTION="read_beat"
SLEEP_POLLING=1

class Server:
    def __init__(self, port, listen_backlog, amount_clients):        
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._amount_clients = amount_clients
        # Process management
        self._childs_processses = []
        self._processes_manager = multiprocessing.Manager()
        # Locks
        self.lock_lottery = self._processes_manager.Lock()
        # Creates lottery
        self._shared_data = self._processes_manager.dict({KEY_AGENCIES_WAITING: []})
        self.lottery = Lottery(self._shared_data, amount_clients)     
        # Initialize signal handling
        self.__init_sign_handling()

    """
    Initialization of signal handling
    """
    def __init_sign_handling(self):
        signal.signal(signal.SIGTERM, self.__handle_a_signal)

    """
    Signal handling
    """
    def __handle_a_signal(self, signal_number, _stack):
        logging.info(f'action: {SIGNAL_HANDLER_ACTION} | result: success')
        self._server_socket.close()
        logging.debug(f'action: {CLOSE_SERVER_SOCKET_ACTION} | result: success')
        # Join childs
        for process in self._childs_processses:
            process.join()
            logging.debug(f'action: {JOIN_PROCESS_ACTION} | result: success')
        sys.exit(0)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # Handle signal to graceful shutdown the server
        while True:
            client_sock = self.__accept_new_connection()
            process = multiprocessing.Process(target=self.__handle_client_connection, args=(client_sock,))
            process.start()
            self._childs_processses.append(process)                  

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # Read agency_id
            (_, _, agency_id) = Protocol.apply_rcv_bets_protocol(client_sock)

            bets_counter = 0
            amount_bets_expected = None
            logging.debug(f'action: read_agency_id | result: success | agency_id: {agency_id}')
            # Read bet loop
            while (amount_bets_expected == None):
                try:
                    (bets, amount_bets_expected, _) = Protocol.apply_rcv_bets_protocol(client_sock, agency_id)
                    with self.lock_lottery:
                        Protocol.apply_store_bets(self.lottery, bets)
                    bets_counter += len(bets)
                except ValueError as e:
                    logging.info(f'action: {READ_BET_ACTION} | result: fail | cantidad: {e}')
                except TypeError as e:
                    logging.debug(f'action: stop_rcv_and_store_bets | result: success | msg: no more bets')
                    break
            # Response with amount bets
            Protocol.apply_res_amount_bets_protocol(client_sock, str(bets_counter), amount_bets_expected)            
            # Add agency waiting for lottery
            with self.lock_lottery:                    
                self.lottery.add_agency(agency_id)
            # Winners polling 
            keep_polling = True
            while keep_polling:
                with self.lock_lottery:
                    keep_polling = not Protocol.apply_winners_protocol(client_sock, agency_id, self.lottery)
                time.sleep(SLEEP_POLLING*self._amount_clients)
        except OSError as e:
            logging.error(f'action: receive_message | result: fail | error: {e}')
        except TypeError as e:
            logging.error(f'action: parse_message | result: fail | error: {e}')
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
