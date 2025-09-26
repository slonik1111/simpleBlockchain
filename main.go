package main

import (
	"fmt"
	"strconv"

	"github.com/slonik1111/simpleBlockchain/simpleBlockchain"
)

func main() {
	fmt.Println("Started program")
	bc := simpleBlockchain.NewBlockchain()
	fmt.Println("Created blockchain")
	bc.AddBlock("Send 1 BTC to Ivan")
	fmt.Println("Created first block")
	bc.AddBlock("Send 2 more BTC to Ivan")
	fmt.Println("Created second block")

	for _, block := range bc.GetBlocks() {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
		pow := simpleBlockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
