package cmd

import (
	"fmt"

	"github.com/slonik1111/simpleBlockchain/simpleBlockchain"

	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Показать все блоки",
	Run: func(cmd *cobra.Command, args []string) {
		bc := simpleBlockchain.NewBlockchain()
		for _, block := range bc.GetBlocks() {
			fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
			fmt.Printf("Data: %s\n", block.Data)
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
