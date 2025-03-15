package common

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

const SEND_BET_ACTION = "apuesta_enviada"
const ERROR_SEND_BET = "error al enviar apuesta"
const ERROR_RCV_CONF = "error al recibir confirmación de apuesta"

// Send a bet to server and wait for ack
func ApplySendBetProtocol(bet *Bet, clientId string, conn net.Conn) {
	parsedLen, parsed, errorOnParse := ParseBet(clientId, bet)
	if errorOnParse != nil {
		log.Errorf("action: parse_bet | result: fail | client_id: %v | error: %v",
			clientId,
			errorOnParse,
		)
	} else {
		errorSendLen := binary.Write(conn, binary.BigEndian, uint16(parsedLen))
		_, errorSendContent := io.WriteString(conn, parsed)
		if _checkErrorSendBet(clientId, errorSendLen, errorSendContent) {
			_, err := bufio.NewReader(conn).ReadString('\n')
			if _checkErrorRcvConf(clientId, err) {
				log.Infof("action: %v | result: success | dni: %v | numero: %v",
					SEND_BET_ACTION,
					bet.DNI,
					bet.Number,
				)
			} else {
				return
			}
		} else {
			return
		}
	}
}

// Returns false if an error occurs (log errors), true otherwise
func _checkErrorSendBet(clientId string, errSendLen error, errSendContent error) bool {
	if errSendLen != nil || errSendContent != nil {
		log.Debug("action: %v | result: fail | client_id: %v | error_send_len: %v | error_send_content: %v",
			ERROR_SEND_BET,
			clientId,
			errSendLen,
			errSendContent,
		)
		log.Errorf("action: %v | result: fail | client_id: %v | error: %v",
			ERROR_SEND_BET,
			clientId,
			ERROR_SEND_BET,
		)
		return false
	}
	return true
}

// Returns false if an error occurs (log the error), true otherwise
func _checkErrorRcvConf(clientId string, err error) bool {
	if err != nil {
		log.Errorf("action: %v | result: fail | client_id: %v | error: %v",
			ERROR_RCV_CONF,
			clientId,
			err,
		)
		return false
	}
	return true
}
