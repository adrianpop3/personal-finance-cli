package report

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"personal-finance-cli/cmd/tui/shared"
	"personal-finance-cli/db"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RunTUI() {
	app := tview.NewApplication()

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[::b][green]📊 Reports Menu[::-]").
		SetDynamicColors(true)

	labels := []string{"Monthly Totals", "Category Breakdown (Chart)", "Back"}
	actions := []func(){
		func() { app.Suspend(monthlyReport) },
		func() { app.Suspend(categoryReport) },
		func() { app.Stop() },
	}

	current := 0
	buttonFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	buttons := []*tview.Button{}

	for i, label := range labels {
		idx := i
		btn := tview.NewButton("[green]" + label).SetSelectedFunc(actions[idx])
		btn.SetBorder(true)
		buttons = append(buttons, btn)
		buttonFlex.AddItem(btn, 3, 0, false)
	}

	highlight := func() {
		for i, btn := range buttons {
			if i == current {
				btn.SetLabel("[white][green]" + labels[i] + "[::-]")
			} else {
				btn.SetLabel("[green]" + labels[i] + "[::-]")
			}
		}
	}
	highlight()

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 5, 1, false).
		AddItem(buttonFlex, 0, 2, true)

	layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			current--
			if current < 0 {
				current = len(buttons) - 1
			}
			highlight()
			return nil
		case tcell.KeyDown:
			current++
			if current >= len(buttons) {
				current = 0
			}
			highlight()
			return nil
		case tcell.KeyEnter:
			actions[current]()
			return nil
		}
		return event
	})

	_ = app.SetRoot(layout, true).EnableMouse(true).Run()
}

func monthlyReport() {
	app := tview.NewApplication()
	monthDefault := time.Now().Format("2006-01")

	var form *tview.Form
	form = tview.NewForm().
		AddInputField("Month (YYYY-MM)", monthDefault, 10, nil, nil).
		AddButton("Show", func() {
			m := strings.TrimSpace(form.GetFormItemByLabel("Month (YYYY-MM)").(*tview.InputField).GetText())
			if m == "" {
				m = monthDefault
			}
			from, err := time.Parse("2006-01", m)
			if err != nil {
				shared.ShowError(app, "Invalid month (YYYY-MM)")
				return
			}
			to := from.AddDate(0, 1, -1)

			txs, err := db.GetFilteredTransactions("", "", &from, &to, nil, nil)
			if err != nil {
				shared.ShowError(app, err.Error())
				return
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

			text := fmt.Sprintf("[green]Monthly Totals: %s[::-]\n\nIncome:   %.2f\nExpenses: %.2f\nNet:      %.2f\n", m, income, expenses, net)

			tv := tview.NewTextView()
			tv.SetDynamicColors(true)
			tv.SetText(text)
			tv.SetBorder(true)
			tv.SetTitle("[green]Monthly Report (ESC=Back)")
			tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEscape {
					app.SetRoot(form, true)
					return nil
				}
				return event
			})

			app.SetRoot(tv, true)
		}).
		AddButton("Back", func() { app.Stop() })

	form.SetBorder(true).SetTitle("[green]Monthly Totals").SetTitleAlign(tview.AlignLeft)
	_ = app.SetRoot(form, true).EnableMouse(true).Run()
}

func categoryReport() {
	app := tview.NewApplication()
	monthDefault := time.Now().Format("2006-01")

	var form *tview.Form
	form = tview.NewForm().
		AddInputField("Month (YYYY-MM)", monthDefault, 10, nil, nil).
		AddButton("Show", func() {
			m := strings.TrimSpace(form.GetFormItemByLabel("Month (YYYY-MM)").(*tview.InputField).GetText())
			if m == "" {
				m = monthDefault
			}
			from, err := time.Parse("2006-01", m)
			if err != nil {
				shared.ShowError(app, "Invalid month (YYYY-MM)")
				return
			}
			to := from.AddDate(0, 1, -1)

			txs, err := db.GetFilteredTransactions("", "", &from, &to, nil, nil)
			if err != nil {
				shared.ShowError(app, err.Error())
				return
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
				shared.ShowOK(app, "No expenses found for "+m, func() { app.SetRoot(form, true) })
				return
			}

			type row struct {
				cat string
				val float64
			}
			var rows []row
			for k, v := range byCat {
				rows = append(rows, row{k, v})
			}
			sort.Slice(rows, func(i, j int) bool { return rows[i].val > rows[j].val })

			maxVal := rows[0].val
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("[green]Category Breakdown: %s (Expenses)[::-]\nTotal: %.2f\n\n", m, total))
			for _, r := range rows {
				pct := (r.val / total) * 100
				sb.WriteString(fmt.Sprintf("%-18s %10.2f  %6.1f%%  %s\n", r.cat, r.val, pct, asciiBar(r.val, maxVal, 28)))
			}

			tv := tview.NewTextView()
			tv.SetDynamicColors(true)
			tv.SetText(sb.String())
			tv.SetBorder(true)
			tv.SetTitle("[green]Category Report (ESC=Back)")
			tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEscape {
					app.SetRoot(form, true)
					return nil
				}
				return event
			})

			app.SetRoot(tv, true)
		}).
		AddButton("Back", func() { app.Stop() })

	form.SetBorder(true).SetTitle("[green]Category Breakdown (Chart)").SetTitleAlign(tview.AlignLeft)
	_ = app.SetRoot(form, true).EnableMouse(true).Run()
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
