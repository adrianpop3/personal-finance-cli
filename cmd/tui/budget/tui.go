package budget

import (
	"fmt"
	"personal-finance-cli/db"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ------------------ MAIN BUDGET TUI -------------------

func RunTUI() {
	app := tview.NewApplication()

	// Title
	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[::b][green]💰 Budgets Menu[::-]").
		SetDynamicColors(true)

	// Buttons
	buttons := tview.NewFlex().SetDirection(tview.FlexRow)
	buttons.AddItem(makeButton("List Budgets", func() {
		app.Suspend(func() { showBudgets() })
	}), 3, 0, true)

	buttons.AddItem(makeButton("Add Budget", func() {
		app.Suspend(func() { AddInteractive() })
	}), 3, 0, false)

	buttons.AddItem(makeButton("Back", func() {
		app.Stop()
	}), 3, 0, false)

	// Layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 5, 1, false).
		AddItem(buttons, 0, 2, true)

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		fmt.Println(err)
	}
}

// ------------------ HELPER -------------------

func makeButton(label string, selected func()) *tview.Button {
	btn := tview.NewButton(fmt.Sprintf("[::b]%s", label)).
		SetSelectedFunc(selected)
	btn.SetBorder(true)
	return btn
}

// ------------------ BUDGET TABLE -------------------

func showBudgets() {
	budgets, err := db.GetBudgets()
	if err != nil {
		fmt.Println("Error fetching budgets:", err)
		return
	}

	app := tview.NewApplication()
	table := tview.NewTable().SetSelectable(true, false)
	table.SetBorder(true).SetTitle("Budgets (Enter=Edit/Delete, ESC=Back)")

	headers := []string{"ID", "Category", "Amount", "Period"}
	for i, h := range headers {
		table.SetCell(0, i, tview.NewTableCell(fmt.Sprintf("[::b]%s", h)).SetSelectable(false))
	}

	for r, b := range budgets {
		table.SetCell(r+1, 0, tview.NewTableCell(strconv.Itoa(b.ID)))
		table.SetCell(r+1, 1, tview.NewTableCell(b.Category))
		table.SetCell(r+1, 2, tview.NewTableCell(fmt.Sprintf("%.2f", b.Amount)))
		table.SetCell(r+1, 3, tview.NewTableCell(b.Period))
	}

	// Row selected → show modal dialog
	table.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		b := budgets[row-1]
		showBudgetActions(b, table, app)
	})

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})

	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		fmt.Println(err)
	}
}

// ------------------ BUDGET ACTIONS -------------------

func showBudgetActions(b db.Budget, parentTable *tview.Table, app *tview.Application) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Budget ID %d\nChoose an action", b.ID)).
		AddButtons([]string{"Edit", "Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Edit":
				app.Suspend(func() { UpdateInteractive(b) })
			case "Delete":
				app.Suspend(func() {
					if err := db.DeleteBudget(b.ID); err != nil {
						fmt.Println("Delete error:", err)
					} else {
						fmt.Println("Budget deleted!")
					}
				})
			case "Cancel":
				// do nothing
			}
			app.SetRoot(parentTable, true)
		})

	app.SetRoot(modal, false)
}

// ------------------ ADD BUDGET FORM -------------------

func AddInteractive() {
	app := tview.NewApplication()
	var form *tview.Form

	form = tview.NewForm().
		AddInputField("Category", "", 20, nil, nil).
		AddInputField("Amount", "", 20, nil, nil).
		AddInputField("Period (YYYY-MM)", time.Now().Format("2006-01"), 20, nil, nil).
		AddButton("Save", func() {
			category := form.GetFormItemByLabel("Category").(*tview.InputField).GetText()
			amountText := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
			period := form.GetFormItemByLabel("Period (YYYY-MM)").(*tview.InputField).GetText()

			amount, err := strconv.ParseFloat(amountText, 64)
			if err != nil {
				fmt.Println("Invalid amount")
				return
			}

			b := db.Budget{
				Category: category,
				Amount:   amount,
				Period:   period,
			}

			if err := db.InsertBudget(b); err != nil {
				fmt.Println("Error saving budget:", err)
			} else {
				fmt.Println("Budget added!")
			}
			app.Stop()
		}).
		AddButton("Cancel", func() { app.Stop() })

	form.SetBorder(true).SetTitle("Add Budget").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		fmt.Println(err)
	}
}

// ------------------ UPDATE BUDGET FORM -------------------

func UpdateInteractive(b db.Budget) {
	app := tview.NewApplication()
	var form *tview.Form

	form = tview.NewForm().
		AddInputField("Category", b.Category, 20, nil, nil).
		AddInputField("Amount", fmt.Sprintf("%.2f", b.Amount), 20, nil, nil).
		AddInputField("Period (YYYY-MM)", b.Period, 20, nil, nil).
		AddButton("Save", func() {
			category := form.GetFormItemByLabel("Category").(*tview.InputField).GetText()
			amountText := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
			period := form.GetFormItemByLabel("Period (YYYY-MM)").(*tview.InputField).GetText()

			amount, err := strconv.ParseFloat(amountText, 64)
			if err != nil {
				fmt.Println("Invalid amount")
				return
			}

			b.Category = category
			b.Amount = amount
			b.Period = period

			if err := db.UpdateBudget(b); err != nil {
				fmt.Println("Error updating budget:", err)
			} else {
				fmt.Println("Budget updated!")
			}
			app.Stop()
		}).
		AddButton("Cancel", func() { app.Stop() })

	form.SetBorder(true).SetTitle(fmt.Sprintf("Edit Budget ID %d", b.ID)).SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		fmt.Println(err)
	}
}
