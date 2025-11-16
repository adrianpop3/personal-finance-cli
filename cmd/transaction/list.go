package transaction

import (
	"fmt"
	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var listID int

// ListCmd lists transactions; can filter by ID
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all transactions or a specific transaction by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		var txs []db.Transaction
		var err error

		if listID > 0 {
			// Fetch single transaction
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
			// Fetch all transactions
			txs, err = db.GetTransactions()
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
	TransactionCmd.AddCommand(ListCmd)
}
