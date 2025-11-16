package transaction

import (
	"fmt"
	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var deleteID int

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a transaction by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := db.DeleteTransaction(deleteID); err != nil {
			return err
		}
		fmt.Println("Transaction deleted.")
		return nil
	},
}

func init() {
	DeleteCmd.Flags().IntVarP(&deleteID, "id", "i", 0, "ID of transaction to delete (required)")
	_ = DeleteCmd.MarkFlagRequired("id")

	TransactionCmd.AddCommand(DeleteCmd)
}
