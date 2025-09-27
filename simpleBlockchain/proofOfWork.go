package simpleBlockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
)

const targetBits = 1

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{block, target}
}

func Int64ToBytes(val int64) []byte {
	result := make([]byte, 0)
	result = binary.BigEndian.AppendUint64(result, uint64(val))
	return result
}

func (pow *ProofOfWork) Validate() bool {
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	Hash := hash[:]
	hashValue := big.NewInt(0).SetBytes(Hash)
	return hashValue.Cmp(pow.target) == -1
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.Data,
		pow.block.PrevBlockHash,
		Int64ToBytes(pow.block.Timestamp),
		Int64ToBytes(int64(nonce)),
	}, []byte{})
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	nonce := 0
	var data []byte
	var Hash []byte
	for nonce < math.MaxInt32 {
		data = pow.prepareData(nonce)
		hash := sha256.Sum256(data)
		Hash = hash[:]
		hashValue := big.NewInt(0).SetBytes(Hash)
		if hashValue.Cmp(pow.target) == -1 {
			break
		}
		nonce++
	}
	return nonce, Hash
}
