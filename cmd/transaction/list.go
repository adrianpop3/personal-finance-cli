package transaction

import (
	"fmt"
	"personal-finance-cli/db"
	"time"

	"github.com/spf13/cobra"
)

var (
	listID       int
	flagCategory string
	flagDesc     string
	flagFrom     string
	flagTo       string
	flagMin      float64
	flagMax      float64
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all transactions or filter by criteria",
	RunE: func(cmd *cobra.Command, args []string) error {
		var txs []db.Transaction
		var err error

		if listID > 0 {
			tx, err := db.GetTransactionByID(listID)
			if err != nil {
				return err
			}
			if tx == nil {
				fmt.Printf("Transaction with ID %d not found.\n", listID)
				return nil
			}
			txs = []db.Transaction{*tx}
		} else {
			// Parse optional filters
			var fromPtr, toPtr *time.Time
			var minPtr, maxPtr *float64

			if flagFrom != "" {
				if t, err := time.Parse("2006-01-02", flagFrom); err == nil {
					fromPtr = &t
				}
			}

			if flagTo != "" {
				if t, err := time.Parse("2006-01-02", flagTo); err == nil {
					toPtr = &t
				}
			}

			if flagMin != 0 {
				minPtr = &flagMin
			}

			if flagMax != 0 {
				maxPtr = &flagMax
			}

			txs, err = db.GetFilteredTransactions(flagCategory, flagDesc, fromPtr, toPtr, minPtr, maxPtr)
			if err != nil {
				return err
			}
		}

		if len(txs) == 0 {
			fmt.Println("No transactions found.")
			return nil
		}

		fmt.Println("ID | Amount | Category | Date | Description")
		for _, t := range txs {
			fmt.Printf("%d | %.2f | %s | %s | %s\n",
				t.ID, t.Amount, t.Category, t.Date.Format("2006-01-02"), t.Description)
		}

		return nil
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listID, "id", "i", 0, "ID of transaction to list (optional)")

	// Filter flags
	ListCmd.Flags().StringVar(&flagCategory, "category", "", "Filter by category")
	ListCmd.Flags().StringVar(&flagDesc, "desc", "", "Filter by description keyword")
	ListCmd.Flags().StringVar(&flagFrom, "from", "", "Filter from date YYYY-MM-DD")
	ListCmd.Flags().StringVar(&flagTo, "to", "", "Filter to date YYYY-MM-DD")
	ListCmd.Flags().Float64Var(&flagMin, "min", 0, "Minimum amount")
	ListCmd.Flags().Float64Var(&flagMax, "max", 0, "Maximum amount")

	TransactionCmd.AddCommand(ListCmd)
}
