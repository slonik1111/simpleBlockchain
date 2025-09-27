package simpleBlockchain

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
	createQuerry, err := loadSQL("queries/Create.sql")
	if err != nil {
		log.Fatal("create")
	}

	db.Exec(createQuerry)

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

func (bc *Blockchain) AddBlock(data string) {
	block := NewBlock(data, bc.GetLastBlock().Hash)
	_, err := bc.db.Exec("INSERT INTO blocks (hash, prev_hash, data, timestamp, nonce) VALUES($1, $2, $3, $4, $5)",
		block.Hash, block.PrevBlockHash, block.Data, block.Timestamp, block.Nonce)
	if err != nil {
		log.Fatal(err)
	}

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
		rows.Scan(&hash, &prevHash, &data, timestamp, nonce)
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
	row.Scan(&hash1, &prevHash, &data, timestamp, nonce)
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
	if err != nil {
		log.Fatal(err)
	}
	_, err = bc.db.Exec("DROP TABLE tail;")
	if err != nil {
		log.Fatal(err)
	}
}
