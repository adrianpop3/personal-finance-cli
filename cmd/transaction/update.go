package transaction

import (
	"fmt"
	"personal-finance-cli/db"
	"time"

	"github.com/spf13/cobra"
)

var (
	updateID          int
	updateAmount      float64
	updateDescription string
	updateCategory    string
	updateDate        string
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a transaction by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		txDate := time.Now()
		var err error
		if updateDate != "" {
			txDate, err = time.Parse("2006-01-02", updateDate)
			if err != nil {
				return fmt.Errorf("invalid date format: %w", err)
			}
		}

		tx := db.Transaction{
			ID:          updateID,
			Amount:      updateAmount,
			Description: updateDescription,
			Category:    updateCategory,
			Date:        txDate,
		}

		if err := db.UpdateTransaction(tx); err != nil {
			return err
		}
		fmt.Println("Transaction updated.")
		return nil
	},
}

func init() {
	UpdateCmd.Flags().IntVarP(&updateID, "id", "i", 0, "ID of transaction to update (required)")
	UpdateCmd.Flags().Float64VarP(&updateAmount, "amount", "a", 0, "New amount")
	UpdateCmd.Flags().StringVarP(&updateDescription, "description", "d", "", "New description")
	UpdateCmd.Flags().StringVarP(&updateCategory, "category", "c", "", "New category")
	UpdateCmd.Flags().StringVarP(&updateDate, "date", "", "", "New date YYYY-MM-DD")
	_ = UpdateCmd.MarkFlagRequired("id")

	TransactionCmd.AddCommand(UpdateCmd)
}
