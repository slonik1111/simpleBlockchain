package cmd

import (
	"fmt"

	"github.com/slonik1111/simpleBlockchain/simpleBlockchain"
	"github.com/spf13/cobra"
)

var dropCmd = &cobra.Command {
	Use: "clear",
	Short: "Очищает блокчейн",
	Run: func(cmd *cobra.Command, args []string) {
		bc := simpleBlockchain.NewBlockchain()
		bc.Clear()
		fmt.Println("Блокчейн очищен")
	},
}

func init() {
	rootCmd.AddCommand(dropCmd)
}