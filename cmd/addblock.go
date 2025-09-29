package cmd

import (
	"fmt"

	"github.com/slonik1111/simpleBlockchain/simpleBlockchain"

	"github.com/spf13/cobra"
)

var addblockCmd = &cobra.Command{
	Use:   "addblock [data]",
	Short: "Добавить новый блок",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bc := simpleBlockchain.NewBlockchain()
		bc.AddBlock(args[0])
		fmt.Println("Блок добавлен:", args[0])
	},
}

func init() {
	rootCmd.AddCommand(addblockCmd)
}
