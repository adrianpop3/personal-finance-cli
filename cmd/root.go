package cmd

import (
	"fmt"
	"os"
	"personal-finance-cli/cmd/transaction"
	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "personal-finance-cli",
	Short: "Personal finance manager CLI",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initDatabase)

	RootCmd.AddCommand(transaction.TransactionCmd)
}

func initDatabase() {
	if err := db.InitDB(); err != nil {
		panic(err)
	}
}
