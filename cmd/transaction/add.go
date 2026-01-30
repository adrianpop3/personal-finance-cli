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
	addAmount      float64
	addDescription string
	addCategory    string
	addDate        string
)

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

		cat := strings.TrimSpace(addCategory)
		if cat == "" {
			cat = parser.InferCategory(addDescription)
		}

		tx := db.Transaction{
			Amount:      addAmount,
			Description: addDescription,
			Category:    cat,
			Date:        txDate,
		}

		if err := db.InsertTransaction(tx); err != nil {
			return err
		}

		fmt.Println("Transaction added.")

		printBudgetAlertsForCategory(cat)

		return nil
	},
}

func init() {
	AddCmd.Flags().Float64VarP(&addAmount, "amount", "a", 0, "Amount of transaction (required). Use negative for expenses.")
	AddCmd.Flags().StringVarP(&addDescription, "description", "d", "", "Description")
	AddCmd.Flags().StringVarP(&addCategory, "category", "c", "", "Category (optional; inferred if empty)")
	AddCmd.Flags().StringVarP(&addDate, "date", "", "", "Date YYYY-MM-DD (optional; defaults to today)")

	_ = AddCmd.MarkFlagRequired("amount")
}
