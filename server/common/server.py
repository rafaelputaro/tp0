import socket
import logging
import sys
import signal

SIGNAL_HANDLER_ACTION="received_a_signal"
CLOSE_SERVER_SOCKET_ACTION="closing_server_socket"
CLOSE_SOCKET_ACTION="closing_a_client_socket"

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        # Initialize signal handling
        self.__init_sign_handling()

    """
    Initialization of signal handling
    """
    def __init_sign_handling(self):
        self._clients_sockets=[]
        signal.signal(signal.SIGTERM, self.__handle_a_signal)

    """
    Signal handling
    """
    def __handle_a_signal(self, signal_number, _stack):
        logging.info(f'action: {SIGNAL_HANDLER_ACTION} | result: success')
        self._server_socket.close()
        logging.debug(f'action: {CLOSE_SERVER_SOCKET_ACTION} | result: success')
        for socket in self._clients_sockets:
            socket.close()
            logging.debug(f'action: {CLOSE_SOCKET_ACTION} | result: success')
        sys.exit(0)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while True:
            client_sock = self.__accept_new_connection()
            self._clients_sockets.append(client_sock)    
            self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            msg = client_sock.recv(1024).rstrip().decode('utf-8')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            client_sock.send("{}\n".format(msg).encode('utf-8'))
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
            self._clients_sockets.remove(client_sock)

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
