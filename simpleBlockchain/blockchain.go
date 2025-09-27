package simpleBlockchain

import (
	"database/sql"
	"fmt"
	"log"
)

type Blockchain struct {
	db *sql.DB
}

func NewBlockchain() *Blockchain {
	db := Connect() 

	var isEmpty bool

	db.QueryRow("SELECT COUNT(*) = 0 FROM blocks;").Scan(&isEmpty)

	fmt.Println(isEmpty)

	if isEmpty {
		genesisBlock := NewGenesisBlock()

		_, err := db.Exec("INSERT INTO blocks (hash, prev_hash, data, timestamp, nonce) VALUES($1, $2, $3, $4, $5)", 
			genesisBlock.Hash, nil, genesisBlock.Data, genesisBlock.Timestamp, genesisBlock.Nonce)

		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("INSERT INTO tail VALUES($1)", genesisBlock.Hash)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Генезис блок в БД")
	}

	return &Blockchain{db}
}
