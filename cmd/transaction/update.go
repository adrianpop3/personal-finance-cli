package transaction

import (
	"fmt"
	"strings"
	"time"

	"personal-finance-cli/db"
	"personal-finance-cli/internal/parser"

	"github.com/spf13/cobra"
)

var (
	updID          int
	updAmount      float64
	updDescription string
	updCategory    string
	updDate        string
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a transaction by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if updID <= 0 {
			return fmt.Errorf("--id is required")
		}

		existing, err := db.GetTransactionByID(updID)
		if err != nil {
			return err
		}
		if existing == nil {
			fmt.Printf("Transaction with ID %d not found.\n", updID)
			return nil
		}

		t := *existing

		if cmd.Flags().Changed("amount") {
			t.Amount = updAmount
		}
		if cmd.Flags().Changed("description") {
			t.Description = updDescription
		}
		if cmd.Flags().Changed("category") {
			t.Category = updCategory
		}
		if cmd.Flags().Changed("date") {
			d, err := time.Parse("2006-01-02", updDate)
			if err != nil {
				return fmt.Errorf("invalid --date (expected YYYY-MM-DD)")
			}
			t.Date = d
		}

		if strings.TrimSpace(t.Category) == "" {
			t.Category = parser.InferCategory(t.Description)
		}

		if err := db.UpdateTransaction(t); err != nil {
			return err
		}

		fmt.Println("Transaction updated.")

		printBudgetAlertsForCategory(t.Category)

		return nil
	},
}

func init() {
	UpdateCmd.Flags().IntVarP(&updID, "id", "i", 0, "Transaction ID (required)")
	UpdateCmd.Flags().Float64VarP(&updAmount, "amount", "a", 0, "New amount")
	UpdateCmd.Flags().StringVarP(&updDescription, "description", "d", "", "New description")
	UpdateCmd.Flags().StringVarP(&updCategory, "category", "c", "", "New category (can be empty to infer)")
	UpdateCmd.Flags().StringVar(&updDate, "date", "", "New date YYYY-MM-DD")

	_ = UpdateCmd.MarkFlagRequired("id")
}
