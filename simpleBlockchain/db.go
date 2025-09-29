package simpleBlockchain

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Connect() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, name, host, port,
	)
	fmt.Println("Строка подключения: " + connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("open")
	}

	fmt.Println("Подключение к БД")

	if err = db.Ping(); err != nil {
		log.Fatal("ping error:", err)
	}
	query, err := loadSQL("queries/Create.sql")
	if err != nil {
		log.Fatal("create")
	}

	db.Exec(query)

	fmt.Println("Таблицы созданы")

	return db
}

func loadSQL(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (bc *Blockchain) GetFreeTransactions() []*Transaction {
	freeTransactions := make([]*Transaction, 0)
	query, err := loadSQL("queries/GetFreeTransactions.sql")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := bc.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	tx := []byte{}
	for rows.Next() {
		rows.Scan(&tx)
		query, err = loadSQL("queries/GetInputsTx")
		check(err)
		ins, err := bc.db.Query(query, tx) 
		txIns := make([]TXInput, 0)
		check(err)
		for ins.Next() {
			input := TXInput{}
			ins.Scan(&input.Txid, &input.Vout, &input.ScriptSig)
			txIns = append(txIns, input)
		}
		query, err = loadSQL("queries/GetOutputsTx")
		check(err)
		outs, err := bc.db.Query(query, tx) 
		txOuts := make([]TXOutput, 0)
		check(err)
		for outs.Next() {
			output := TXOutput{}
			outs.Scan(&output.Value, &output.ScriptPubKey)
			txOuts = append(txOuts, output)
		}
		freeTransactions = append(freeTransactions, &Transaction{tx, txIns, txOuts})
	}
	return freeTransactions
}

func (bc *Blockchain) AddBlock(data string) {
	txs := bc.GetFreeTransactions()
	block := NewBlock(data, txs, bc.GetLastBlock().Hash)
	_, err := bc.db.Exec("INSERT INTO blocks (hash, prev_hash, data, timestamp, nonce) VALUES($1, $2, $3, $4, $5)",
		block.Hash, block.PrevBlockHash, block.Data, block.Timestamp, block.Nonce)
	if err != nil {
		log.Fatal(err)
	}

	_, err = bc.db.Exec("UPDATE transactions SET blockhash = $1 WHERE blockhash IS NULL", block.Hash) 
	check(err)

	_, err = bc.db.Exec("UPDATE tail SET hash = $1", block.Hash)
	if err != nil {
		log.Fatal(err)
	}
}

func (bc *Blockchain) GetBlocks() []*Block {
	blocks := make([]*Block, 0)
	rows, err := bc.db.Query("SELECT * FROM blocks")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var hash, prevHash, data []byte
		var timestamp int64
		var nonce int
		rows.Scan(&hash, &prevHash, &data, &timestamp, &nonce)
		blocks = append(blocks, &Block{
			Hash:          hash,
			PrevBlockHash: prevHash,
			Data:          data,
			Timestamp:     timestamp,
			Nonce:         nonce,
		})
	}
	return blocks
}

func (bc *Blockchain) GetBlockByHash(hash *[]byte) *Block {
	row := bc.db.QueryRow("SELECT * FROM blocks WHERE hash = $1", hash)
	var hash1, prevHash, data []byte
	var timestamp int64
	var nonce int
	row.Scan(&hash1, &prevHash, &data, &timestamp, &nonce)
	return &Block{
		Hash:          hash1,
		PrevBlockHash: prevHash,
		Data:          data,
		Timestamp:     timestamp,
		Nonce:         nonce,
	}
}

func (bc *Blockchain) GetLastBlock() *Block {
	var tail []byte
	row := bc.db.QueryRow("SELECT * FROM tail")
	row.Scan(&tail)
	return bc.GetBlockByHash(&tail)
}

func (bc *Blockchain) GetChain() []*Block {
	blocks := make([]*Block, 0)
	currentBlock := bc.GetLastBlock()
	for {
		blocks = append(blocks, currentBlock)
		if currentBlock.PrevBlockHash == nil {
			break
		}
		currentBlock = bc.GetBlockByHash(&currentBlock.PrevBlockHash)
	}
	return blocks
}

func (bc *Blockchain) Clear() {
	_, err := bc.db.Exec("DROP TABLE blocks;")
	check(err)
	_, err = bc.db.Exec("DROP TABLE tail;")
	check(err)
	_, err = bc.db.Exec("DROP TABLE transactions;")
	check(err)
	_, err = bc.db.Exec("DROP TABLE inputs;")
	check(err)
	_, err = bc.db.Exec("DROP TABLE outputs;")
	check(err)
}


func (bc *Blockchain) GetFreeOutputs(adr string) (int, []TXOutput, [][]byte, []int) {
	outputs := make([]TXOutput, 0)
	txns := make([][]byte, 0)
	index := make([]int, 0)
	sum := 0
	query, err := loadSQL("queries/GetFreeOutputs.sql")
	fmt.Println(query)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := bc.db.Query(query, adr)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var txId []byte
		var idx, amount int
		rows.Scan(&txId, &idx, &amount)
		sum += amount
		outputs = append(outputs, TXOutput{amount, adr})
		txns = append(txns, txId)
		index = append(index, idx)
	}
	fmt.Print(outputs)
	return sum, outputs, txns, index
}

func (bc *Blockchain) AddTransaction(from, to string, amount int) {
	tnx := NewTransaction(from, to, amount, bc)
	query := "INSERT INTO transactions VALUES ($1, $2)"
	_, err := bc.db.Exec(query, tnx.ID, nil)
	check(err)
	for i, input := range tnx.Vin {
		query = "INSERT INTO inputs VALUES($1, $2, $3, $4, $5)"
		_, err := bc.db.Exec(query, tnx.ID, i, input.Txid, input.Vout, input.ScriptSig)
		check(err)
	}
	for i, output := range tnx.Vout {
		query = "INSERT INTO outputs VALUES($1, $2, $3, $4)"
		_, err := bc.db.Exec(query, tnx.ID, i, output.ScriptPubKey, output.Value)
		check(err)
	}
	fmt.Println("Транзакция создана")
}
