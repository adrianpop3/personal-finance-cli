package transaction

import (
	"github.com/spf13/cobra"
)

var TransactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Manage transactions",
	Long:  "Create, list, update, and delete transactions in your personal finance manager.",
}
