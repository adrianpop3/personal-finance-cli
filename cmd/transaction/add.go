package transaction

import (
	"fmt"
	"time"

	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var (
	addAmount      float64
	addDescription string
	addCategory    string
	addDate        string
)

// AddCmd represents the "transaction add" command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new transaction",
	RunE: func(cmd *cobra.Command, args []string) error {
		var txDate time.Time
		var err error
		if addDate == "" {
			txDate = time.Now()
		} else {
			txDate, err = time.Parse("2006-01-02", addDate)
			if err != nil {
				return fmt.Errorf("invalid date format: %w", err)
			}
		}

		tx := db.Transaction{
			Amount:      addAmount,
			Description: addDescription,
			Category:    addCategory,
			Date:        txDate,
		}

		if err := db.InsertTransaction(tx); err != nil {
			return err
		}

		fmt.Println("Transaction added.")
		return nil
	},
}

func init() {
	// Define flags
	AddCmd.Flags().Float64VarP(&addAmount, "amount", "a", 0, "Amount of transaction (required)")
	AddCmd.Flags().StringVarP(&addDescription, "description", "d", "", "Description")
	AddCmd.Flags().StringVarP(&addCategory, "category", "c", "Uncategorized", "Category")
	AddCmd.Flags().StringVarP(&addDate, "date", "", "", "Date YYYY-MM-DD (optional; defaults to today)")

	// Make amount required
	_ = AddCmd.MarkFlagRequired("amount")

	// Attach to parent TransactionCmd instead of root
	TransactionCmd.AddCommand(AddCmd)
}
