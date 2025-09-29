package simpleBlockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	Transactions   []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce		  int
}

func (b *Block) HashTransactions() []byte {
	data := make([][]byte, 0)
	for _, tx := range b.Transactions {
		data = append(data, tx.ID)
	}
	sum := bytes.Join(data, []byte{})
	hash := sha256.Sum256(sum)
	return hash[:]
}

func (b *Block) SetHash() {
	timestamp := Int64ToBytes(b.Timestamp)
	nonce := Int64ToBytes(int64(b.Nonce))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.HashTransactions() ,b.Data, timestamp, nonce}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func NewBlock(data string, txs []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Transactions: txs,
	}
	fmt.Printf("Создан PoW для %s\n", data)
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

func NewGenesisBlock(tx *Transaction) *Block {
	return NewBlock("Genesis Block", []*Transaction{tx} ,[]byte{})
}
