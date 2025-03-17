package common

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"strings"
)

const ACTION_SEND_BET = "apuesta_enviada"
const ACTION_RCV_AMOUNT_BETS = "recibir_cantidad_apuestas"

// const ACTION_MAX_BATCH_AMOUNT_REACHED = "Maximum batch length reached"
const ERROR_SEND_BET = "error al enviar apuesta"
const ERROR_RCV_AMOUNT_BETS = "error al recibir cantidad de apuestas"

// Bet Entity that encapsulates how
type Protocol struct {
	amountBetsSended int
	amountBetsBatch  int
	parsedLen        int
	betsParsed       []string
}

// Initializes a new Bets Loader based on the parameters
func NewProtocol() *Protocol {
	protocol := &Protocol{
		amountBetsSended: 0,
		amountBetsBatch:  0,
		parsedLen:        0,
	}
	return protocol
}

// It applies the protocol that consists of adding bets to a buffer and attempting to send
// them when this buffer is full.
// Bets that cannot be parsed are discarded.
// Returns error in case of comunication error
func (protocol *Protocol) ApplySendBetProtocol(clientId string, conn net.Conn, batchMaxAmount int, bet *Bet) error {
	parsedLen, parsed, errorOnParse := ParseBetWithEndDelimiter(bet)
	if errorOnParse != nil {
		log.Errorf("action: parse_bet | result: fail | client_id: %v | error: %v",
			clientId,
			errorOnParse,
		)
		return nil
	} else {
		parsedLenProtUpdated := protocol.parsedLen + parsedLen
		if (parsedLenProtUpdated < MAX_LEN) && (protocol.amountBetsBatch < batchMaxAmount) {
			protocol.appendBet(parsed, parsedLenProtUpdated)
			return nil
		} else {
			/*
				log.Debugf("action: %v | amount_bets_batch: %v | amount_bets_sended: %v",
					ACTION_MAX_BATCH_AMOUNT_REACHED,
					protocol.amountBetsBatch,
					protocol.amountBetsSended,
				)*/
			err := protocol.trySendBets(clientId, conn)
			if err == nil {
				protocol.appendBet(parsed, parsedLen)
			}
			return err
		}
	}
}

func (protocol *Protocol) trySendBets(clientId string, conn net.Conn) error {
	parsed := strings.Join(protocol.betsParsed, "")
	parsedLen := protocol.parsedLen
	errorSendLen := binary.Write(conn, binary.BigEndian, uint16(parsedLen))
	_, errorSendContent := io.WriteString(conn, parsed)
	checkError := _checkErrorSend(clientId, errorSendLen, errorSendContent)
	if checkError == nil {
		// clean buffer and uptad amount bets sended
		protocol.amountBetsSended += len(protocol.betsParsed)
		protocol.amountBetsBatch = 0
		protocol.betsParsed = protocol.betsParsed[:0]
		protocol.parsedLen = 0
		return nil
	}
	return checkError
}

// Send remaining bets
func (protocol *Protocol) SendRemainingBets(clientId string, conn net.Conn) error {
	if protocol.parsedLen > 0 {
		return protocol.trySendBets(clientId, conn)
	}
	return nil
}

// Append a bet to the protocol buffer
func (protocol *Protocol) appendBet(parsed string, parsedLenProtUpdated int) {
	protocol.betsParsed = append(protocol.betsParsed, parsed)
	protocol.amountBetsBatch += 1
	protocol.parsedLen = parsedLenProtUpdated
}

// Send amount agencyId to server
func (protocol *Protocol) ApplySendAgencyIdProtocol(clientId string, conn net.Conn) error {
	parsedLen, parsed, errorOnParse := ParseAgencyId(clientId)
	if errorOnParse != nil {
		log.Errorf("action: agency_id | result: fail | client_id: %v | error: %v",
			clientId,
			errorOnParse,
		)
		return errorOnParse
	} else {
		errorSendLen := binary.Write(conn, binary.BigEndian, uint16(parsedLen))
		_, errorSendContent := io.WriteString(conn, parsed)
		return _checkErrorSend(clientId, errorSendLen, errorSendContent)
	}
}

// Send amount bets to server and wait for ack
func (protocol *Protocol) ApplySendAmountBetsProtocol(clientId string, conn net.Conn, amountBets int) error {
	parsedLen, parsed, errorOnParse := ParseAmountBets(amountBets)
	if errorOnParse != nil {
		log.Errorf("action: amount_bets | result: fail | client_id: %v | error: %v",
			clientId,
			errorOnParse,
		)
		return errorOnParse
	} else {
		errorSendLen := binary.Write(conn, binary.BigEndian, uint16(parsedLen))
		_, errorSendContent := io.WriteString(conn, parsed)
		return _checkErrorSend(clientId, errorSendLen, errorSendContent)
	}
}

// Send amount bets to server and wait for ack
func (protocol *Protocol) ApplyRecvAmountBetsProtocol(clientId string, conn net.Conn, amountBetseExpected int) error {
	readed, errorRcv := bufio.NewReader(conn).ReadString('\n')
	return _checkRcvAmountBets(clientId, strings.TrimSuffix(readed, "\n"), strconv.Itoa(amountBetseExpected), errorRcv)
}

// Returns false if an error occurs (log errors), true otherwise
func _checkErrorSend(clientId string, errSendLen error, errSendContent error) error {
	checkErrSendLen := errSendLen != nil
	checkErrSendContent := errSendContent != nil
	if checkErrSendLen || checkErrSendContent {
		log.Debugf("action: %v | result: fail | client_id: %v | error_send_len: %v | error_send_content: %v",
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
		if checkErrSendLen {
			return errSendLen
		} else {
			return errSendContent
		}
	}
	return nil
}

// Returns false if an error occurs (log errors), true otherwise
func _checkRcvAmountBets(clientId string, amountBetsReaded string, amountBetseExpected string, errRcv error) error {
	if errRcv != nil {
		log.Debugf("action: %v | result: fail | client_id: %v | error_rcv: %v",
			ACTION_RCV_AMOUNT_BETS,
			clientId,
			errRcv,
		)
		log.Errorf("action: %v | result: fail | client_id: %v | error: %v",
			ACTION_RCV_AMOUNT_BETS,
			clientId,
			ERROR_RCV_AMOUNT_BETS,
		)
	} else {

		var check_success string
		if amountBetsReaded == amountBetseExpected {
			check_success = "success"
		} else {
			check_success = "fail"
		}
		log.Infof("action: %v | result: %v | client_id: %v",
			ACTION_RCV_AMOUNT_BETS,
			check_success,
			clientId,
		)
	}
	return errRcv

}
