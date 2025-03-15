package common

import (
	//"bufio"
	//"fmt"
	"net"
	"os"
	"time"

	"github.com/op/go-logging"
)

const SIGNAL_ACTION = "received_a_sigterm"

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	Nombre        string
	Apellido      string
	Documento     string
	Nacimiento    string
	Numero        string
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
func (c *Client) StartClientLoop(singalChannel chan os.Signal) {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
loop:
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Create the connection the server in every loop iteration.
		c.createClientSocket()
		// Send bet
		bet := LoadBet(&c.config)
		ApplySendBetProtocol(bet, c.config.ID, c.conn)
		// Close connection
		c.conn.Close()
		// Wait for signal to end client
		for {
			select {
			// Receive a signal from channel
			case <-singalChannel:
				log.Infof("action: %v | result: success | client_id: %v",
					SIGNAL_ACTION,
					c.config.ID,
				)
				break loop
			// Wait a time between sending one message and the next one
			case <-time.After(c.config.LoopPeriod):
			}
		}
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
