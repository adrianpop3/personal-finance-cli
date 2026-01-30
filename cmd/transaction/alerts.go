package transaction

import (
	"fmt"
	"strings"
	"time"

	"personal-finance-cli/db"
)

func printBudgetAlertsForCategory(category string) {
	category = strings.TrimSpace(category)
	if category == "" {
		return
	}

	budgets, err := db.GetBudgets()
	if err != nil {
		return
	}

	month := time.Now().Format("2006-01")
	found := false

	for _, b := range budgets {
		if b.Category != category {
			continue
		}
		found = true

		calcB := b
		if calcB.Period == "" || calcB.Period == "monthly" {
			calcB.Period = month
		}

		rem, err := db.GetBudgetRemaining(calcB)
		if err != nil {
			continue
		}

		alert := ""
		if rem < 0 {
			alert = "OVER"
		} else if b.Amount > 0 && (rem/b.Amount) <= 0.10 {
			alert = "LOW (<=10%)"
		}

		if alert != "" {
			fmt.Printf("⚠ Budget alert: %s (%s) remaining %.2f / limit %.2f => %s\n",
				b.Category, calcB.Period, rem, b.Amount, alert)
		} else {
			fmt.Printf("Budget status: %s (%s) remaining %.2f / limit %.2f\n",
				b.Category, calcB.Period, rem, b.Amount)
		}
	}

	if !found {
		return
	}
}

func printBudgetAlertsForCategories(categories []string) {
	seen := map[string]struct{}{}
	for _, c := range categories {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		printBudgetAlertsForCategory(c)
	}
}
