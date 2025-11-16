package budget

import (
	"fmt"
	"time"

	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var (
	addCategory string
	addAmount   float64
	addPeriod   string
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new budget",
	RunE: func(cmd *cobra.Command, args []string) error {
		if addPeriod == "" {
			addPeriod = time.Now().Format("2006-01")
		}
		b := db.Budget{
			Category: addCategory,
			Amount:   addAmount,
			Period:   addPeriod,
		}
		if err := db.InsertBudget(b); err != nil {
			return err
		}
		fmt.Println("Budget added.")
		return nil
	},
}

func init() {
	AddCmd.Flags().StringVarP(&addCategory, "category", "c", "", "Category (required)")
	AddCmd.Flags().Float64VarP(&addAmount, "amount", "a", 0, "Budget amount (required)")
	AddCmd.Flags().StringVarP(&addPeriod, "period", "p", "", "Period YYYY-MM (optional; defaults to current month)")
	_ = AddCmd.MarkFlagRequired("category")
	_ = AddCmd.MarkFlagRequired("amount")

	BudgetCmd.AddCommand(AddCmd)
}
