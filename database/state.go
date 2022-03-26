package database

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Tx
	dbFile    *os.File
}

func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genesisFilePath := filepath.Join(cwd, "database", "genesis.json")
	genesisFile, err := loadGenesis(genesisFilePath)
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range genesisFile.Balances {
		balances[account] = balance
	}

	txDbFilePath := filepath.Join(cwd, "database", "tx.db")
	txDbFile, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(txDbFile)

	state := &State{
		Balances:  balances,
		txMempool: make([]Tx, 0),
		dbFile:    txDbFile,
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}
