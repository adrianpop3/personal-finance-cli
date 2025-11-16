package budget

import (
	"github.com/spf13/cobra"
)

var BudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Manage budgets",
}

func init() {
	// child commands will attach here (add, update, delete, list)
}
