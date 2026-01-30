package budget

import (
	"fmt"
	"strconv"
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
		SetText("[::b][green]💰 Budgets Menu[::-]").
		SetDynamicColors(true)

	labels := []string{"List Budgets", "Add Budget", "Budget Status / Alerts", "Back"}
	actions := []func(){
		func() { app.Suspend(showBudgets) },
		func() { app.Suspend(AddInteractive) },
		func() { app.Suspend(showBudgetStatus) },
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

// ------------------ Budget Table -------------------

func showBudgets() {
	budgets, err := db.GetBudgets()
	if err != nil {
		fmt.Println("Error fetching budgets:", err)
		return
	}

	app := tview.NewApplication()
	table := tview.NewTable().SetSelectable(true, false)
	table.SetBorder(true).SetTitle("[green]Budgets (Enter=Edit/Delete, ESC=Back)").SetTitleAlign(tview.AlignCenter)

	headers := []string{"ID", "Category", "Amount", "Period"}
	for i, h := range headers {
		table.SetCell(0, i, tview.NewTableCell(fmt.Sprintf("[::b][green]%s[::-]", h)).SetSelectable(false))
	}

	refresh := func() {
		updated, err := db.GetBudgets()
		if err != nil {
			return
		}
		budgets = updated

		for r := 1; r < table.GetRowCount(); r++ {
			for c := 0; c < len(headers); c++ {
				table.SetCell(r, c, tview.NewTableCell(""))
			}
		}

		for r, b := range budgets {
			table.SetCell(r+1, 0, tview.NewTableCell(strconv.Itoa(b.ID)))
			table.SetCell(r+1, 1, tview.NewTableCell(b.Category))
			table.SetCell(r+1, 2, tview.NewTableCell(fmt.Sprintf("%.2f", b.Amount)))
			table.SetCell(r+1, 3, tview.NewTableCell(b.Period))
		}
	}

	refresh()

	table.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		b := budgets[row-1]
		showBudgetActions(b, table, app, refresh)
	})

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})

	_ = app.SetRoot(table, true).EnableMouse(true).Run()
}

func showBudgetActions(b db.Budget, parentTable *tview.Table, app *tview.Application, refresh func()) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("[green]Budget ID %d\nChoose an action[::-]", b.ID)).
		AddButtons([]string{"Edit", "Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Edit":
				app.Suspend(func() { UpdateInteractive(b) })
				refresh()
			case "Delete":
				_ = db.DeleteBudget(b.ID)
				refresh()
			case "Cancel":
			}
			app.SetRoot(parentTable, true)
		})

	app.SetRoot(modal, false)
}

// ------------------ Add / Update Forms -------------------

func AddInteractive() {
	app := tview.NewApplication()
	var form *tview.Form

	form = tview.NewForm().
		AddInputField("Category", "", 20, nil, nil).
		AddInputField("Amount", "", 20, nil, nil).
		AddInputField("Period", "monthly", 20, nil, nil).
		AddButton("Save", func() {
			category := form.GetFormItemByLabel("Category").(*tview.InputField).GetText()
			amountText := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
			period := form.GetFormItemByLabel("Period").(*tview.InputField).GetText()

			amount, err := strconv.ParseFloat(amountText, 64)
			if err != nil {
				shared.ShowError(app, "Invalid amount")
				return
			}

			b := db.Budget{Category: category, Amount: amount, Period: period}
			if err := db.InsertBudget(b); err != nil {
				shared.ShowError(app, "Error saving budget: "+err.Error())
				return
			}

			shared.ShowOK(app, "Budget added!", func() { app.Stop() })
		}).
		AddButton("Cancel", func() { app.Stop() })

	form.SetBorder(true).SetTitle("[green]Add Budget").SetTitleAlign(tview.AlignLeft)
	_ = app.SetRoot(form, true).EnableMouse(true).Run()
}

func UpdateInteractive(b db.Budget) {
	app := tview.NewApplication()
	var form *tview.Form

	form = tview.NewForm().
		AddInputField("Category", b.Category, 20, nil, nil).
		AddInputField("Amount", fmt.Sprintf("%.2f", b.Amount), 20, nil, nil).
		AddInputField("Period", b.Period, 20, nil, nil).
		AddButton("Save", func() {
			category := form.GetFormItemByLabel("Category").(*tview.InputField).GetText()
			amountText := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
			period := form.GetFormItemByLabel("Period").(*tview.InputField).GetText()

			amount, err := strconv.ParseFloat(amountText, 64)
			if err != nil {
				shared.ShowError(app, "Invalid amount")
				return
			}

			b.Category = category
			b.Amount = amount
			b.Period = period

			if err := db.UpdateBudget(b); err != nil {
				shared.ShowError(app, "Error updating budget: "+err.Error())
				return
			}

			shared.ShowOK(app, "Budget updated!", func() { app.Stop() })
		}).
		AddButton("Cancel", func() { app.Stop() })

	form.SetBorder(true).SetTitle(fmt.Sprintf("[green]Edit Budget ID %d", b.ID)).SetTitleAlign(tview.AlignLeft)
	_ = app.SetRoot(form, true).EnableMouse(true).Run()
}

// ------------------ Budget Status / Alerts -------------------

func showBudgetStatus() {
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
			if _, err := time.Parse("2006-01", m); err != nil {
				shared.ShowError(app, "Invalid month. Use YYYY-MM.")
				return
			}

			budgets, err := db.GetBudgets()
			if err != nil {
				shared.ShowError(app, "Error fetching budgets: "+err.Error())
				return
			}
			if len(budgets) == 0 {
				shared.ShowOK(app, "No budgets found.", func() { app.SetRoot(form, true) })
				return
			}

			table := tview.NewTable().SetSelectable(true, false)
			table.SetBorder(true).SetTitle("[green]Budget Status / Alerts (ESC=Back)").SetTitleAlign(tview.AlignCenter)

			headers := []string{"Category", "Limit", "Remaining", "Alert"}
			for i, h := range headers {
				table.SetCell(0, i, tview.NewTableCell(fmt.Sprintf("[::b][green]%s[::-]", h)).SetSelectable(false))
			}

			row := 1
			for _, b := range budgets {
				calcB := b
				if calcB.Period == "" || calcB.Period == "monthly" {
					calcB.Period = m
				}

				rem, err := db.GetBudgetRemaining(calcB)
				if err != nil {
					continue
				}

				alert := ""
				if rem < 0 {
					alert = "OVER"
				} else if b.Amount > 0 && (rem/b.Amount) <= 0.10 {
					alert = "LOW"
				}

				table.SetCell(row, 0, tview.NewTableCell(b.Category))
				table.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%.2f", b.Amount)))
				table.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%.2f", rem)))
				table.SetCell(row, 3, tview.NewTableCell(alert))
				row++
			}

			table.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEscape {
					app.SetRoot(form, true)
				}
			})

			app.SetRoot(table, true)
		}).
		AddButton("Back", func() { app.Stop() })

	form.SetBorder(true).SetTitle("[green]Budget Status / Alerts").SetTitleAlign(tview.AlignLeft)
	_ = app.SetRoot(form, true).EnableMouse(true).Run()
}
