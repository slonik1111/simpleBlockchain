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
	PrevBlockHash []byte
	Hash          []byte
	Nonce		int
}

func (b *Block) SetHash() {
	timestamp := Int64ToBytes(b.Timestamp)
	nonce := Int64ToBytes(int64(b.Nonce))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp, nonce}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}
	fmt.Printf("Создан PoW для %s\n", data)
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
