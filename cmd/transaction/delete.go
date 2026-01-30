package transaction

import (
	"fmt"

	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var delID int

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a transaction by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if delID <= 0 {
			return fmt.Errorf("--id is required")
		}
		if err := db.DeleteTransaction(delID); err != nil {
			return err
		}
		fmt.Println("Transaction deleted.")
		return nil
	},
}

func init() {
	DeleteCmd.Flags().IntVarP(&delID, "id", "i", 0, "Transaction ID (required)")
	_ = DeleteCmd.MarkFlagRequired("id")
}
