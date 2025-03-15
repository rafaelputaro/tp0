package common

// Bet Entity that encapsulates how
type Bet struct {
	Name     string
	LastName string
	DNI      string
	Birthday string
	Number   string
}

// Initializes a new Bet based on the parameters
func NewBet(name string, lastName string, dni string, birthday string, number string) *Bet {
	bet := &Bet{
		Name:     name,
		LastName: lastName,
		DNI:      dni,
		Birthday: birthday,
		Number:   number,
	}
	return bet
}

// Logs action over bet.
func logActionOverBet(bet *Bet, action string) {
	log.Debugf("action: %s | result: success | name: %s | last_name: %s | dni: %s | birth_day: %s | number: %s",
		action,
		bet.Name,
		bet.LastName,
		bet.DNI,
		bet.Birthday,
		bet.Number,
	)
}

// Load bet from config
func LoadBet(config *ClientConfig) *Bet {
	bet := Bet{
		Name:     config.Nombre,
		LastName: config.Apellido,
		DNI:      config.Documento,
		Birthday: config.Nacimiento,
		Number:   config.Numero,
	}
	logActionOverBet(&bet, "bet_from_config")
	return &bet
}
