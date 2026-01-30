package transaction

import (
	"fmt"
	"time"

	"personal-finance-cli/db"

	"github.com/spf13/cobra"
)

var (
	listID       int
	listCategory string
	listDesc     string
	listFrom     string
	listTo       string
	listMin      float64
	listMax      float64
	listHasMin   bool
	listHasMax   bool
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List transactions (all, by id, or filtered)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if listID > 0 {
			t, err := db.GetTransactionByID(listID)
			if err != nil {
				return err
			}
			if t == nil {
				fmt.Printf("Transaction with ID %d not found.\n", listID)
				return nil
			}
			printTransactions([]db.Transaction{*t})
			return nil
		}

		var fromPtr, toPtr *time.Time
		if listFrom != "" {
			t, err := time.Parse("2006-01-02", listFrom)
			if err != nil {
				return fmt.Errorf("invalid --from (expected YYYY-MM-DD)")
			}
			fromPtr = &t
		}
		if listTo != "" {
			t, err := time.Parse("2006-01-02", listTo)
			if err != nil {
				return fmt.Errorf("invalid --to (expected YYYY-MM-DD)")
			}
			toPtr = &t
		}

		var minPtr, maxPtr *float64
		if listHasMin {
			minPtr = &listMin
		}
		if listHasMax {
			maxPtr = &listMax
		}

		txs, err := db.GetFilteredTransactions(listCategory, listDesc, fromPtr, toPtr, minPtr, maxPtr)
		if err != nil {
			return err
		}
		if len(txs) == 0 {
			fmt.Println("No transactions found.")
			return nil
		}
		printTransactions(txs)
		return nil
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listID, "id", "i", 0, "List a specific transaction by ID")

	ListCmd.Flags().StringVar(&listCategory, "category", "", "Filter by category")
	ListCmd.Flags().StringVar(&listDesc, "desc", "", "Filter by description keyword (substring match)")
	ListCmd.Flags().StringVar(&listFrom, "from", "", "Filter from date YYYY-MM-DD")
	ListCmd.Flags().StringVar(&listTo, "to", "", "Filter to date YYYY-MM-DD")

	ListCmd.Flags().Float64Var(&listMin, "min", 0, "Filter min amount (requires flag set)")
	ListCmd.Flags().Float64Var(&listMax, "max", 0, "Filter max amount (requires flag set)")
	ListCmd.Flags().Lookup("min").NoOptDefVal = "0"
	ListCmd.Flags().Lookup("max").NoOptDefVal = "0"

	ListCmd.PreRun = func(cmd *cobra.Command, args []string) {
		listHasMin = cmd.Flags().Changed("min")
		listHasMax = cmd.Flags().Changed("max")
	}
}

func printTransactions(txs []db.Transaction) {
	fmt.Println("ID | Amount | Category | Date | Description")
	for _, t := range txs {
		fmt.Printf("%d | %.2f | %s | %s | %s\n",
			t.ID, t.Amount, t.Category, t.Date.Format("2006-01-02"), t.Description)
	}
}
