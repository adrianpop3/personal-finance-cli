package cmd

import (
	"os"

	"personal-finance-cli/cmd/budget"
	"personal-finance-cli/cmd/report"
	"personal-finance-cli/cmd/transaction"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "fincli",
	Short: "Personal Finance CLI Manager",
	Long:  "Track transactions, import statements, set budgets, and generate reports.",
}

func Execute() {
	_ = RootCmd.Execute()
}

func init() {
	RootCmd.SetOut(os.Stdout)
	RootCmd.SetErr(os.Stderr)

	RootCmd.AddCommand(transaction.TransactionCmd)
	RootCmd.AddCommand(budget.BudgetCmd)
	RootCmd.AddCommand(report.ReportCmd)
}
