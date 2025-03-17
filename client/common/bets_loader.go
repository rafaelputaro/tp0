package common

import (
	"encoding/csv"
	"fmt"
	"os"
)

const MSG_CANT_OPEN_FILE = "Can't open file"
const MSG_ERROR_ON_PARSE = "Can't parse bet readed"
const ACTION_OPEN_FILE = "open_file"
const ACTION_CLOSE_FILE = "close_file"

// const ACTION_READ_BEAT = "read_bet"
const ACTION_PARSE_BET_FROM_CSV = "parse_bet_from_csv"
const AMOUNT_FIELDS = 5

// Bet Entity that encapsulates how
type BetsLoader struct {
	file             *os.File
	reader           *csv.Reader
	filePath         string
	lastReaded       *Bet
	eof              bool
	amountBetsReaded int
}

// Initializes a new Bets Loader based on the parameters
func NewBetsLoader(filePath string) (*BetsLoader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("%s %s: %s", MSG_CANT_OPEN_FILE, filePath, err)
		betsLoader := &BetsLoader{
			file:             nil,
			reader:           nil,
			filePath:         filePath,
			lastReaded:       nil,
			eof:              true,
			amountBetsReaded: 0,
		}
		return betsLoader, err
	}
	log.Debugf("action: %s %s | result: success ", ACTION_OPEN_FILE, filePath)
	betsLoader := &BetsLoader{
		file:             file,
		reader:           csv.NewReader(file),
		filePath:         filePath,
		lastReaded:       nil,
		eof:              false,
		amountBetsReaded: 0,
	}
	return betsLoader, nil
}

// Close the associated file
func (betsLoader *BetsLoader) CloseFile() {
	betsLoader.file.Close()
	betsLoader.file = nil
	betsLoader.reader = nil
	log.Debugf("action: %s %s | result: success", ACTION_CLOSE_FILE, betsLoader.filePath)
}

// Read next bet from csv
func (betsLoader *BetsLoader) Next() (*Bet, error) {
	line, err := betsLoader.reader.Read()
	if err != nil {
		betsLoader.eof = true
		betsLoader.lastReaded = nil
		//	log.Debugf("action: %s %s | result: fail", ACTION_READ_BEAT, betsLoader.filePath)
		return nil, err
	}
	betsLoader.lastReaded, err = parseBetFromFields(line)
	if err == nil {
		betsLoader.amountBetsReaded += 1
	}
	return betsLoader.lastReaded, err
}

// Returns if the end of the file has been reached
func (betsLoader *BetsLoader) IsEof() bool {
	return betsLoader.eof
}

// Get the last bet readed from csv
func (betsLoader *BetsLoader) GetLast() *Bet {
	return betsLoader.lastReaded
}

// Parse bet from a csv line
func parseBetFromFields(fields []string) (*Bet, error) {
	if len(fields) == AMOUNT_FIELDS {
		bet := Bet{
			Name:     fields[0],
			LastName: fields[1],
			DNI:      fields[2],
			Birthday: fields[3],
			Number:   fields[4],
		}
		return &bet, nil
	}
	log.Errorf("action: %s | result: fail", ACTION_PARSE_BET_FROM_CSV)
	return nil, fmt.Errorf("%s", MSG_ERROR_ON_PARSE)
}
