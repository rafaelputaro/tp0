package common

import (
	"errors"
	"fmt"
)

const MAX_LEN = 8192
const MSG_ERROR = "error: message too long"

// Get a string ("," delimiter) from a bet
func doParseBet(idAgency string, bet *Bet) string {
	return fmt.Sprintf(
		"%s,%s,%s,%s,%s,%v",
		idAgency,
		bet.Name,
		bet.LastName,
		bet.DNI,
		bet.Birthday,
		bet.Number,
	)
}

// // Get a string ("," delimiter) from a bet and check the size
func ParseBet(idAgency string, bet *Bet) (int, string, error) {
	parsed := doParseBet(idAgency, bet)
	parsedLen := len(parsed)
	if parsedLen > MAX_LEN {
		return parsedLen, parsed, errors.New(MSG_ERROR)
	}
	return parsedLen, parsed, nil
}
