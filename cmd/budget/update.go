package budget

import (
	"fmt"
	"time"

	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var (
	updateID       int
	updateCategory string
	updateAmount   float64
	updatePeriod   string
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a budget by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		b := db.Budget{
			ID:       updateID,
			Category: updateCategory,
			Amount:   updateAmount,
			Period:   updatePeriod,
		}

		if b.Period == "" {
			b.Period = time.Now().Format("2006-01")
		}

		if err := db.UpdateBudget(b); err != nil {
			return err
		}
		fmt.Println("Budget updated.")
		return nil
	},
}

func init() {
	UpdateCmd.Flags().IntVarP(&updateID, "id", "i", 0, "ID of budget to update (required)")
	UpdateCmd.Flags().StringVarP(&updateCategory, "category", "c", "", "New category")
	UpdateCmd.Flags().Float64VarP(&updateAmount, "amount", "a", 0, "New budget amount")
	UpdateCmd.Flags().StringVarP(&updatePeriod, "period", "p", "", "New period YYYY-MM")
	_ = UpdateCmd.MarkFlagRequired("id")

	BudgetCmd.AddCommand(UpdateCmd)
}
