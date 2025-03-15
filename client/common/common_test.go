package common

import (
	"testing"
)

func TestBet(t *testing.T) {
	bet := NewBet("Rafael", "Putaro", "29752645", "1982-08-28", "2020")
	if (bet.Name != "Rafael") || (bet.LastName != "Putaro") || (bet.DNI != "29752645") || (bet.Birthday != "28-08-1982") || (bet.Number != "2020") {
		t.Fatal("The bet has not been created correctly")
	}
}

func TestParser(t *testing.T) {
	bet := NewBet("Rafael", "Putaro", "29752645", "1982-08-28", "2020")
	_, parsed, _ := ParseBet("125", bet)
	if parsed != "125,Rafael,Putaro,29752645,28-08-1982,2020" {
		t.Fatal("The bet has not been parsed correctly")
	}
}
