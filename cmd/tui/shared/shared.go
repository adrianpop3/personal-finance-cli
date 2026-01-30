package shared

import (
	"fmt"
	"strings"
	"time"

	"personal-finance-cli/db"

	"github.com/rivo/tview"
)

func ShowOK(app *tview.Application, text string, done func()) {
	showModal(app, text, done)
}

func ShowError(app *tview.Application, text string) {
	showModal(app, "Error: "+text, nil)
}

func showModal(app *tview.Application, text string, done func()) {
	m := tview.NewModal().
		SetText("[green]" + text + "[::-]").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(i int, lbl string) {
			if done != nil {
				done()
			}
		})
	app.SetRoot(m, false)
}

func BudgetAlertLines(category string) []string {
	category = strings.TrimSpace(category)
	if category == "" {
		return nil
	}

	budgets, err := db.GetBudgets()
	if err != nil {
		return nil
	}

	month := time.Now().Format("2006-01")
	var lines []string

	for _, b := range budgets {
		if b.Category != category {
			continue
		}

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
			lines = append(lines, fmt.Sprintf("⚠ %s (%s): remaining %.2f / limit %.2f => %s",
				b.Category, calcB.Period, rem, b.Amount, alert))
		} else {
			lines = append(lines, fmt.Sprintf("%s (%s): remaining %.2f / limit %.2f",
				b.Category, calcB.Period, rem, b.Amount))
		}
	}

	return lines
}

func BudgetAlertLinesForCategories(categories []string) []string {
	seen := map[string]struct{}{}
	var all []string
	for _, c := range categories {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		all = append(all, BudgetAlertLines(c)...)
	}
	return all
}
