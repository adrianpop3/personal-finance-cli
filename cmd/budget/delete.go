package budget

import (
	"fmt"
	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var deleteID int

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a budget by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := db.DeleteBudget(deleteID); err != nil {
			return err
		}
		fmt.Println("Budget deleted.")
		return nil
	},
}

func init() {
	DeleteCmd.Flags().IntVarP(&deleteID, "id", "i", 0, "ID of budget to delete (required)")
	_ = DeleteCmd.MarkFlagRequired("id")

	BudgetCmd.AddCommand(DeleteCmd)
}
