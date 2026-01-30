package transaction

import (
	"fmt"

	"personal-finance-cli/internal/parser"

	"github.com/spf13/cobra"
)

var importFile string

var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import transactions from CSV/OFX/QFX file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if importFile == "" {
			return fmt.Errorf("--file is required")
		}

		parsed, err := parser.ParseFileByPath(importFile)
		if err != nil {
			return err
		}
		if len(parsed) == 0 {
			fmt.Println("No transactions parsed.")
			return nil
		}

		if err := parser.InsertParsedTransactions(parsed); err != nil {
			return err
		}

		fmt.Printf("Imported %d transactions.\n", len(parsed))

		var affected []string
		for _, p := range parsed {
			affected = append(affected, p.Category)
		}

		printBudgetAlertsForCategories(affected)

		return nil
	},
}

func init() {
	ImportCmd.Flags().StringVarP(&importFile, "file", "f", "", "Path to CSV/OFX/QFX file")
	_ = ImportCmd.MarkFlagRequired("file")
}
