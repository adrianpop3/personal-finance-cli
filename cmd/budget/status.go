package budget

import (
	"fmt"
	"strings"
	"time"

	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var statusMonth string

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show budget remaining and alerts for a month",
	RunE: func(cmd *cobra.Command, args []string) error {
		m := strings.TrimSpace(statusMonth)
		if m == "" {
			m = time.Now().Format("2006-01")
		}
		if _, err := time.Parse("2006-01", m); err != nil {
			return fmt.Errorf("invalid --month, expected YYYY-MM")
		}

		budgets, err := db.GetBudgets()
		if err != nil {
			return err
		}
		if len(budgets) == 0 {
			fmt.Println("No budgets found.")
			return nil
		}

		fmt.Printf("Budget Status for %s\n", m)
		fmt.Println("ID | Category | Limit | Remaining | Alert")

		for _, b := range budgets {
			calcB := b
			if b.Period == "" || b.Period == "monthly" {
				calcB.Period = m
			}

			rem, err := db.GetBudgetRemaining(calcB)
			if err != nil {
				return err
			}

			alert := ""
			if rem < 0 {
				alert = "OVER"
			} else if b.Amount > 0 && rem/b.Amount <= 0.10 {
				alert = "LOW (<=10%)"
			}

			fmt.Printf("%d | %s | %.2f | %.2f | %s\n", b.ID, b.Category, b.Amount, rem, alert)
		}
		return nil
	},
}

func init() {
	StatusCmd.Flags().StringVarP(&statusMonth, "month", "m", "", "Month in format YYYY-MM (default: current month)")
}
