package common

import (
	"testing"
)

func TestBet(t *testing.T) {
	bet := NewBet("Rafael", "Putaro", "29752645", "1982-08-28", "2020")
	if (bet.Name != "Rafael") || (bet.LastName != "Putaro") || (bet.DNI != "29752645") || (bet.Birthday != "1982-08-28") || (bet.Number != "2020") {
		t.Fatal("The bet has not been created correctly")
	}
}

func TestParser(t *testing.T) {
	bet := NewBet("Rafael", "Putaro", "29752645", "1982-08-28", "2020")
	_, parsed, _ := ParseBet(bet)
	print(parsed)
	if parsed != "Rafael,Putaro,29752645,1982-08-28,2020" {
		t.Fatal("The bet has not been parsed correctly")
	}
	_, parsed, _ = ParseAgencyId("1")
	if parsed != "1" {
		t.Fatal("The agencyId has not been parsed correctly")
	}
	_, parsed, _ = ParseAmountBets(11)
	if parsed != "EOF,11" {
		t.Fatal("The amount_bets not been parsed correctly")
	}
}
