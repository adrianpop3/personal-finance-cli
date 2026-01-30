package report

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate spending and income reports",
}

var month string

func init() {
	monthlyCmd := &cobra.Command{
		Use:   "monthly",
		Short: "Monthly totals (income, expenses, net)",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := month
			if strings.TrimSpace(m) == "" {
				m = time.Now().Format("2006-01")
			}

			from, err := time.Parse("2006-01", m)
			if err != nil {
				return fmt.Errorf("invalid --month, expected YYYY-MM")
			}
			to := from.AddDate(0, 1, -1)

			txs, err := db.GetFilteredTransactions("", "", &from, &to, nil, nil)
			if err != nil {
				return err
			}

			var income, expenses float64
			for _, t := range txs {
				if t.Amount >= 0 {
					income += t.Amount
				} else {
					expenses += -t.Amount
				}
			}
			net := income - expenses

			fmt.Printf("Monthly Report: %s\n", m)
			fmt.Printf("Income:   %.2f\n", income)
			fmt.Printf("Expenses: %.2f\n", expenses)
			fmt.Printf("Net:      %.2f\n", net)
			return nil
		},
	}
	monthlyCmd.Flags().StringVarP(&month, "month", "m", "", "Month in format YYYY-MM (default: current month)")

	categoriesCmd := &cobra.Command{
		Use:   "categories",
		Short: "Category breakdown for a month (expenses only) with ASCII bars",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := month
			if strings.TrimSpace(m) == "" {
				m = time.Now().Format("2006-01")
			}
			from, err := time.Parse("2006-01", m)
			if err != nil {
				return fmt.Errorf("invalid --month, expected YYYY-MM")
			}
			to := from.AddDate(0, 1, -1)

			txs, err := db.GetFilteredTransactions("", "", &from, &to, nil, nil)
			if err != nil {
				return err
			}

			byCat := map[string]float64{}
			var total float64
			for _, t := range txs {
				if t.Amount < 0 {
					byCat[t.Category] += -t.Amount
					total += -t.Amount
				}
			}

			if total == 0 {
				fmt.Printf("Category Report: %s\nNo expenses found.\n", m)
				return nil
			}

			type row struct {
				Cat   string
				Value float64
			}
			var rows []row
			for k, v := range byCat {
				rows = append(rows, row{Cat: k, Value: v})
			}
			sort.Slice(rows, func(i, j int) bool { return rows[i].Value > rows[j].Value })

			fmt.Printf("Category Report: %s (Expenses)\n", m)
			fmt.Printf("Total expenses: %.2f\n\n", total)

			maxVal := rows[0].Value
			for _, r := range rows {
				pct := (r.Value / total) * 100
				bar := asciiBar(r.Value, maxVal, 30)
				fmt.Printf("%-18s %10.2f  %6.1f%%  %s\n", r.Cat, r.Value, pct, bar)
			}
			return nil
		},
	}
	categoriesCmd.Flags().StringVarP(&month, "month", "m", "", "Month in format YYYY-MM (default: current month)")

	ReportCmd.AddCommand(monthlyCmd)
	ReportCmd.AddCommand(categoriesCmd)
}

func asciiBar(value, max float64, width int) string {
	if max <= 0 {
		return ""
	}
	ratio := value / max
	n := int(math.Round(ratio * float64(width)))
	if n < 0 {
		n = 0
	}
	if n > width {
		n = width
	}
	return strings.Repeat("█", n)
}
