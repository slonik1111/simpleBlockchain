package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/slonik1111/simpleBlockchain/simpleBlockchain"

	"github.com/spf13/cobra"
)

var addtxCmd = &cobra.Command{
	Use:   "addtx",
	Short: "Добавить транзакцию",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		bc := simpleBlockchain.NewBlockchain()
		amount, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Добавляется транзакция:", args[0], args[1], amount)
		bc.AddTransaction(args[0], args[1], amount)
		fmt.Println("Транзакция добавлена:", args[0], args[1], args[2])
	},
}

func init() {
	rootCmd.AddCommand(addtxCmd)
}