package common

import (
	"net"
	"os"
	"time"

	"github.com/op/go-logging"
)

const SIGNAL_ACTION = "received_a_sigterm"

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             string
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	DataPath       string
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(signalChannel chan os.Signal) {
	// Create the connection the server
	c.createClientSocket()
	// Reads data file
	loader, _ := NewBetsLoader(c.config.DataPath)
	// Create communication handler
	protocol := NewProtocol()
	// Send first message to server
	protocol.ApplySendAgencyIdProtocol(c.config.ID, c.conn)
loop_send_bets:
	// Wait for signal to end client
	for !loader.IsEof() {
		bet, errorReadFile := loader.Next()
		if errorReadFile == nil {
			protocol.ApplySendBetProtocol(c.config.ID, c.conn, c.config.BatchMaxAmount, bet)
		}
		// Receive a signal from channel
		if c.receiveSignalOrSleep(signalChannel) {
			break loop_send_bets
		}
	}
	// Send Remaining
	protocol.SendRemainingBets(c.config.ID, c.conn)
	// Send amount bets
	protocol.ApplySendAmountBetsProtocol(c.config.ID, c.conn, protocol.amountBetsSended)
	// Receive how many bets have been accepted by the server
	protocol.ApplyRecvAmountBetsProtocol(c.config.ID, c.conn, protocol.amountBetsSended)
	// Close File
	loader.CloseFile()
	log.Debugf("action: loop_send_bets_finished | result: success | client_id: %v", c.config.ID)
	// Request winners
	haveWinners := false
loop_winners:
	for !haveWinners {
		protocol.ApplyRequestWinnersProtocol(c.config.ID, c.conn)
		_, haveWinners = protocol.ApplyRecvWinnersProtocol(c.config.ID, c.conn)
		if !haveWinners {
			log.Debugf("action: keep_waiting_winners | result: success | client_id: %v", c.config.ID)
		}
		// Receive a signal from channel
		if c.receiveSignalOrSleep(signalChannel) {
			break loop_winners
		}
	}
	// Close connection
	c.conn.Close()
	log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
}

// Returns true if it receives a signal, otherwise returns false after sleeping
func (c *Client) receiveSignalOrSleep(signalChannel chan os.Signal) bool {
	// Receive a signal from channel or sleep
	select {
	case <-signalChannel:
		log.Debugf("action: %v | result: success | client_id: %v",
			SIGNAL_ACTION,
			c.config.ID,
		)
		return true
	case <-time.After(c.config.LoopPeriod):
		return false
	}
}
