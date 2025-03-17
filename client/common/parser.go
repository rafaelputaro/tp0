package common

import (
	"errors"
	"fmt"
	"strings"
)

const MAX_LEN = 7000
const MSG_ERROR = "error: message too long"
const WINNERS_REQUEST = "WINNERS"
const KEEP_WAITING_WINNERS_TAG = "KEEP WAITING WINNERS"
const WINNERS_DELIMITER = ";"

// Get a string ("," delimiter) from a bet
func doParseBet(bet *Bet) string {
	return fmt.Sprintf(
		"%s,%s,%s,%s,%v",
		bet.Name,
		bet.LastName,
		bet.DNI,
		bet.Birthday,
		bet.Number,
	)
}

// Get a string ("," delimiter) and (";" delimitr between bets)  from a bet and check the size
func ParseBetWithEndDelimiter(bet *Bet) (int, string, error) {
	parsed := fmt.Sprintf("%s;", doParseBet(bet))
	return checkLen(parsed)
}

// Get a string ("," delimiter) from a bet and check the size
func ParseBet(bet *Bet) (int, string, error) {
	parsed := doParseBet(bet)
	return checkLen(parsed)
}

// Get a string from amountBets in format "EOF,<amount_bets_string>"
func ParseAmountBets(amountBets int) (int, string, error) {
	parsed := fmt.Sprintf("EOF,%v", amountBets)
	return checkLen(parsed)
}

// Get a string to request winners"
func ParseRequestWinners() (int, string) {
	return len(WINNERS_REQUEST), WINNERS_REQUEST
}

// Parse and check length
func ParseAgencyId(idAgency string) (int, string, error) {
	parsed := idAgency
	return checkLen(parsed)
}

// Check length message
func checkLen(parsed string) (int, string, error) {
	parsedLen := len(parsed)
	if parsedLen > MAX_LEN {
		return parsedLen, parsed, errors.New(MSG_ERROR)
	}
	return parsedLen, parsed, nil
}

// Returns true if the messages contains winners
func ParseWinnersMessage(msg string) (bool, []string) {
	if strings.Contains(msg, KEEP_WAITING_WINNERS_TAG) {
		return false, nil
	} else {
		toReturn := []string{}
		if len(msg) > 0 {
			toReturn = strings.Split(msg, WINNERS_DELIMITER)
		}
		return true, toReturn
	}
}
