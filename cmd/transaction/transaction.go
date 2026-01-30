package transaction

import "github.com/spf13/cobra"

var TransactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Manage transactions",
	Long:  "Create, list, update, delete, import, and search/filter transactions.",
}

func init() {
	TransactionCmd.AddCommand(AddCmd)
	TransactionCmd.AddCommand(ListCmd)
	TransactionCmd.AddCommand(UpdateCmd)
	TransactionCmd.AddCommand(DeleteCmd)
	TransactionCmd.AddCommand(ImportCmd)
}
