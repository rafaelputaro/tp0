package common

import (
	"errors"
	"fmt"
)

const MAX_LEN = 7000
const MSG_ERROR = "error: message too long"

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
