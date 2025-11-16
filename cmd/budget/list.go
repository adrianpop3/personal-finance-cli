package budget

import (
	"fmt"
	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var listID int

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all budgets or a specific budget by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		var budgets []db.Budget
		var err error

		if listID > 0 {
			b, err := db.GetBudgetByID(listID)
			if err != nil {
				return err
			}
			if b == nil {
				fmt.Printf("Budget with ID %d not found.\n", listID)
				return nil
			}
			budgets = []db.Budget{*b}
		} else {
			budgets, err = db.GetBudgets()
			if err != nil {
				return err
			}
		}

		if len(budgets) == 0 {
			fmt.Println("No budgets found.")
			return nil
		}

		fmt.Println("ID | Category | Amount | Period")
		for _, b := range budgets {
			fmt.Printf("%d | %s | %.2f | %s\n", b.ID, b.Category, b.Amount, b.Period)
		}
		return nil
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listID, "id", "i", 0, "ID of budget to list (optional)")
	BudgetCmd.AddCommand(ListCmd)
}
