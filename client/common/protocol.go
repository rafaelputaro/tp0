package common

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const ACTION_SEND_BET = "apuesta_enviada"
const ACTION_RCV_AMOUNT_BETS = "recibir_cantidad_apuestas"
const ACTION_ASK_WINNERS = "consulta_ganadores"
const ERROR_SEND = "error al mensaje al servidor"
const ERROR_RCV_AMOUNT_BETS = "error al recibir cantidad de apuestas"
const ERROR_RCV_WINNERS = "error al recibir ganadores"
const MAX_RETRY = 100
const TIME_RETRY = 100 //ms

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
	var checkError error = nil
	// if error on send retry
	for attemp := 0; attemp < MAX_RETRY; attemp++ {
		errorSendLen := binary.Write(conn, binary.BigEndian, uint16(parsedLen))
		lenWrited, errorSendContent := io.WriteString(conn, parsed)
		checkError = checkErrorSend(clientId, errorSendLen, errorSendContent, checkErrorShortWrite(parsedLen, lenWrited))
		if checkError == nil {
			break
		}
		time.Sleep(TIME_RETRY * time.Millisecond)
	}
	if checkError == nil {
		// clean buffer and update amount bets sended
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
		var checkError error = nil
		// if error on send retry
		for attemp := 0; attemp < MAX_RETRY; attemp++ {
			errorSendLen := binary.Write(conn, binary.BigEndian, uint16(parsedLen))
			lenWrited, errorSendContent := io.WriteString(conn, parsed)
			checkError = checkErrorSend(clientId, errorSendLen, errorSendContent, checkErrorShortWrite(parsedLen, lenWrited))
			if checkError == nil {
				break
			}
			time.Sleep(TIME_RETRY * time.Millisecond)
		}
		return checkError
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
		var checkError error = nil
		// if error on send retry
		for attemp := 0; attemp < MAX_RETRY; attemp++ {
			errorSendLen := binary.Write(conn, binary.BigEndian, uint16(parsedLen))
			lenWrited, errorSendContent := io.WriteString(conn, parsed)
			checkError = checkErrorSend(clientId, errorSendLen, errorSendContent, checkErrorShortWrite(parsedLen, lenWrited))
			if checkError == nil {
				break
			}
			time.Sleep(TIME_RETRY * time.Millisecond)
		}
		return checkError
	}
}

// Receive amount bets from server
func (protocol *Protocol) ApplyRecvAmountBetsProtocol(clientId string, conn net.Conn, amountBetseExpected int) error {
	readed, errorRcv := bufio.NewReader(conn).ReadString('\n')
	return checkRcvAmountBets(clientId, strings.TrimSuffix(readed, "\n"), strconv.Itoa(amountBetseExpected), errorRcv)
}

// Send amount bets to server and wait for ack
func (protocol *Protocol) ApplyRequestWinnersProtocol(clientId string, conn net.Conn) error {
	parsedLen, parsed := ParseRequestWinners()
	var checkError error = nil
	// if error on send retry
	for attemp := 0; attemp < MAX_RETRY; attemp++ {
		errorSendLen := binary.Write(conn, binary.BigEndian, uint16(parsedLen))
		lenWrited, errorSendContent := io.WriteString(conn, parsed)
		checkError = checkErrorSend(clientId, errorSendLen, errorSendContent, checkErrorShortWrite(parsedLen, lenWrited))
		if checkError == nil {
			break
		}
		time.Sleep(TIME_RETRY * time.Millisecond)
	}
	return checkError
}

// Receive winners from server. Returns if receive winners
func (protocol *Protocol) ApplyRecvWinnersProtocol(clientId string, conn net.Conn) (error, bool) {
	readed, errorRcv := bufio.NewReader(conn).ReadString('\n')
	return checkRcvWinners(clientId, strings.TrimSuffix(readed, "\n"), errorRcv)
}

// Returns error if short write
func checkErrorShortWrite(lenExpected int, lenWrited int) error {
	if lenExpected != lenWrited {
		return errors.New("error short write")
	} else {
		return nil
	}
}

// Returns false if an error occurs (log errors), true otherwise
func checkErrorSend(clientId string, errSendLen error, errSendContent error, errShortWrite error) error {
	checkErrSendLen := errSendLen != nil
	checkErrSendContent := errSendContent != nil
	checkErrShortWrite := errShortWrite != nil
	if checkErrSendLen || checkErrSendContent || checkErrShortWrite {
		log.Debugf("action: %v | result: fail | client_id: %v | error_send_len: %v | error_send_content: %v",
			ERROR_SEND,
			clientId,
			errSendLen,
			errSendContent,
		)
		log.Errorf("action: %v | result: fail | client_id: %v | error: %v",
			ERROR_SEND,
			clientId,
			ERROR_SEND,
		)
		if checkErrSendLen {
			return errSendLen
		} else {
			if checkErrShortWrite {
				return errShortWrite
			} else {
				return errSendContent
			}
		}
	}
	return nil
}

// Returns false if an error occurs (log errors), true otherwise
func checkRcvAmountBets(clientId string, amountBetsReaded string, amountBetseExpected string, errRcv error) error {
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

// Returns false if an error occurs (log errors) or not have a winner, true otherwise
func checkRcvWinners(clientId string, readed string, errRcv error) (error, bool) {
	if errRcv != nil {
		log.Debugf("action: %v | result: fail | client_id: %v | error_rcv: %v",
			ACTION_ASK_WINNERS,
			clientId,
			errRcv,
		)
		log.Errorf("action: %v | result: fail | client_id: %v | error: %v",
			ACTION_ASK_WINNERS,
			clientId,
			ERROR_RCV_WINNERS,
		)
	} else {
		haveWinners, winners := ParseWinnersMessage(readed)
		if haveWinners {
			log.Infof("action: %v | result: success | cant_ganadores: %v",
				ACTION_ASK_WINNERS,
				len(winners),
			)
			return nil, true
		}
		return nil, false
	}
	return errRcv, false
}
