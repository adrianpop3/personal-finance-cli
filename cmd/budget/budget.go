package budget

import "github.com/spf13/cobra"

var BudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Manage budgets",
}

func init() {
	BudgetCmd.AddCommand(AddCmd)
	BudgetCmd.AddCommand(ListCmd)
	BudgetCmd.AddCommand(UpdateCmd)
	BudgetCmd.AddCommand(DeleteCmd)
	BudgetCmd.AddCommand(StatusCmd)
}
