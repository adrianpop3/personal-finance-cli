package transaction

import (
	"fmt"
	"personal-finance-cli/db"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ------------------ MAIN TRANSACTION TUI -------------------

func RunTUI() {
	app := tview.NewApplication()

	// Title
	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[::b][green]💰 Transactions Menu[::-]").
		SetDynamicColors(true)

	// Buttons
	buttons := tview.NewFlex().SetDirection(tview.FlexRow)
	buttons.AddItem(makeButton("List Transactions", func() {
		app.Suspend(func() { showTransactions() })
	}), 3, 0, true)

	buttons.AddItem(makeButton("Add Transaction", func() {
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

// ------------------ TRANSACTION TABLE -------------------

func showTransactions() {
	txs, err := db.GetTransactions()
	if err != nil {
		fmt.Println("Error fetching transactions:", err)
		return
	}

	app := tview.NewApplication()
	table := tview.NewTable().SetSelectable(true, false)
	table.SetBorder(true).SetTitle("Transactions (Enter=Edit/Delete, ESC=Back)")

	headers := []string{"ID", "Amount", "Category", "Date", "Description"}
	for i, h := range headers {
		table.SetCell(0, i, tview.NewTableCell(fmt.Sprintf("[::b]%s", h)).SetSelectable(false))
	}

	for r, t := range txs {
		table.SetCell(r+1, 0, tview.NewTableCell(strconv.Itoa(t.ID)))
		table.SetCell(r+1, 1, tview.NewTableCell(fmt.Sprintf("%.2f", t.Amount)))
		table.SetCell(r+1, 2, tview.NewTableCell(t.Category))
		table.SetCell(r+1, 3, tview.NewTableCell(t.Date.Format("2006-01-02")))
		table.SetCell(r+1, 4, tview.NewTableCell(t.Description))
	}

	// Row selected → show modal dialog
	table.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		tx := txs[row-1]
		showTransactionActions(tx, table, app)
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

// ------------------ TRANSACTION ACTIONS -------------------

func showTransactionActions(tx db.Transaction, parentTable *tview.Table, app *tview.Application) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Transaction ID %d\nChoose an action", tx.ID)).
		AddButtons([]string{"Edit", "Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Edit":
				app.Suspend(func() { UpdateInteractive(tx) })
			case "Delete":
				app.Suspend(func() {
					if err := db.DeleteTransaction(tx.ID); err != nil {
						fmt.Println("Delete error:", err)
					} else {
						fmt.Println("Transaction deleted!")
					}
				})
			case "Cancel":
				// do nothing
			}
			app.SetRoot(parentTable, true) // return to table
		})

	app.SetRoot(modal, false)
}

// ------------------ ADD TRANSACTION FORM -------------------

func AddInteractive() {
	app := tview.NewApplication()
	var form *tview.Form // <-- declared before closure

	form = tview.NewForm().
		AddInputField("Amount", "", 20, nil, nil).
		AddInputField("Category", "Uncategorized", 20, nil, nil).
		AddInputField("Description", "", 50, nil, nil).
		AddInputField("Date (YYYY-MM-DD)", time.Now().Format("2006-01-02"), 20, nil, nil).
		AddButton("Save", func() {
			amountText := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
			category := form.GetFormItemByLabel("Category").(*tview.InputField).GetText()
			desc := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
			dateText := form.GetFormItemByLabel("Date (YYYY-MM-DD)").(*tview.InputField).GetText()

			amount, err := strconv.ParseFloat(amountText, 64)
			if err != nil {
				fmt.Println("Invalid amount")
				return
			}

			txDate, err := time.Parse("2006-01-02", dateText)
			if err != nil {
				fmt.Println("Invalid date")
				return
			}

			tx := db.Transaction{
				Amount:      amount,
				Category:    category,
				Description: desc,
				Date:        txDate,
			}

			if err := db.InsertTransaction(tx); err != nil {
				fmt.Println("Error saving transaction:", err)
			} else {
				fmt.Println("Transaction added!")
			}
			app.Stop()
		}).
		AddButton("Cancel", func() { app.Stop() })

	form.SetBorder(true).SetTitle("Add Transaction").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		fmt.Println(err)
	}
}

// ------------------ UPDATE TRANSACTION FORM -------------------

func UpdateInteractive(tx db.Transaction) {
	app := tview.NewApplication()
	var form *tview.Form // <-- declared before closure

	form = tview.NewForm().
		AddInputField("Amount", fmt.Sprintf("%.2f", tx.Amount), 20, nil, nil).
		AddInputField("Category", tx.Category, 20, nil, nil).
		AddInputField("Description", tx.Description, 50, nil, nil).
		AddInputField("Date (YYYY-MM-DD)", tx.Date.Format("2006-01-02"), 20, nil, nil).
		AddButton("Save", func() {
			amountText := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
			category := form.GetFormItemByLabel("Category").(*tview.InputField).GetText()
			desc := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
			dateText := form.GetFormItemByLabel("Date (YYYY-MM-DD)").(*tview.InputField).GetText()

			amount, err := strconv.ParseFloat(amountText, 64)
			if err != nil {
				fmt.Println("Invalid amount")
				return
			}

			txDate, err := time.Parse("2006-01-02", dateText)
			if err != nil {
				fmt.Println("Invalid date")
				return
			}

			tx.Amount = amount
			tx.Category = category
			tx.Description = desc
			tx.Date = txDate

			if err := db.UpdateTransaction(tx); err != nil {
				fmt.Println("Error updating transaction:", err)
			} else {
				fmt.Println("Transaction updated!")
			}
			app.Stop()
		}).
		AddButton("Cancel", func() { app.Stop() })

	form.SetBorder(true).SetTitle(fmt.Sprintf("Edit Transaction ID %d", tx.ID)).SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		fmt.Println(err)
	}
}
