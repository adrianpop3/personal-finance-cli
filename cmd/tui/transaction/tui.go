package transaction

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"personal-finance-cli/cmd/tui/shared"
	"personal-finance-cli/db"
	"personal-finance-cli/internal/parser"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RunTUI() {
	app := tview.NewApplication()

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[::b][green]💰 Transactions Menu[::-]").
		SetDynamicColors(true)

	labels := []string{"List Transactions", "Search / Filter", "Add Transaction", "Import From File", "Back"}
	actions := []func(){
		func() { app.Suspend(showTransactions) },
		func() { app.Suspend(searchTransactions) },
		func() { app.Suspend(AddInteractive) },
		func() { app.Suspend(ImportInteractive) },
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

// ------------------ Table view -------------------

func showTransactions() {
	txs, err := db.GetTransactions()
	if err != nil {
		fmt.Println("Error fetching transactions:", err)
		return
	}
	showTransactionsTable(txs, func() ([]db.Transaction, error) {
		return db.GetTransactions()
	})
}

func showTransactionsTable(txs []db.Transaction, reload func() ([]db.Transaction, error)) {
	app := tview.NewApplication()
	table := tview.NewTable().SetSelectable(true, false)
	table.SetBorder(true).SetTitle("[green]Transactions (Enter=Edit/Delete, ESC=Back)").SetTitleAlign(tview.AlignCenter)

	headers := []string{"ID", "Amount", "Category", "Date", "Description"}
	for i, h := range headers {
		table.SetCell(0, i, tview.NewTableCell(fmt.Sprintf("[::b][green]%s[::-]", h)).SetSelectable(false))
	}

	refresh := func() {
		newTxs, err := reload()
		if err != nil {
			return
		}
		txs = newTxs

		for r := 1; r < table.GetRowCount(); r++ {
			for c := 0; c < len(headers); c++ {
				table.SetCell(r, c, tview.NewTableCell(""))
			}
		}

		for r, t := range txs {
			table.SetCell(r+1, 0, tview.NewTableCell(strconv.Itoa(t.ID)))
			table.SetCell(r+1, 1, tview.NewTableCell(fmt.Sprintf("%.2f", t.Amount)))
			table.SetCell(r+1, 2, tview.NewTableCell(t.Category))
			table.SetCell(r+1, 3, tview.NewTableCell(t.Date.Format("2006-01-02")))
			table.SetCell(r+1, 4, tview.NewTableCell(t.Description))
		}
	}

	refresh()

	table.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		tx := txs[row-1]
		showTransactionActions(tx, table, app, refresh)
	})

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})

	_ = app.SetRoot(table, true).EnableMouse(true).Run()
}

func showTransactionActions(tx db.Transaction, parentTable *tview.Table, app *tview.Application, refresh func()) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("[green]Transaction ID %d\nChoose an action[::-]", tx.ID)).
		AddButtons([]string{"Edit", "Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Edit":
				app.Suspend(func() { UpdateInteractive(tx) })
				refresh()
			case "Delete":
				_ = db.DeleteTransaction(tx.ID)
				refresh()
			case "Cancel":
			}
			app.SetRoot(parentTable, true)
		})
	app.SetRoot(modal, false)
}

// ------------------ Add / Update -------------------

func AddInteractive() {
	app := tview.NewApplication()
	var form *tview.Form

	form = tview.NewForm().
		AddInputField("Amount", "", 20, nil, nil).
		AddInputField("Category (optional)", "", 20, nil, nil).
		AddInputField("Description", "", 50, nil, nil).
		AddInputField("Date (YYYY-MM-DD)", time.Now().Format("2006-01-02"), 20, nil, nil).
		AddButton("Save", func() {
			amountText := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
			category := form.GetFormItemByLabel("Category (optional)").(*tview.InputField).GetText()
			desc := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
			dateText := form.GetFormItemByLabel("Date (YYYY-MM-DD)").(*tview.InputField).GetText()

			if strings.TrimSpace(category) == "" {
				category = parser.InferCategory(desc)
			}

			amount, err := strconv.ParseFloat(amountText, 64)
			if err != nil {
				shared.ShowError(app, "Invalid amount")
				return
			}

			txDate, err := time.Parse("2006-01-02", dateText)
			if err != nil {
				shared.ShowError(app, "Invalid date (use YYYY-MM-DD)")
				return
			}

			tx := db.Transaction{Amount: amount, Category: category, Description: desc, Date: txDate}
			if err := db.InsertTransaction(tx); err != nil {
				shared.ShowError(app, "Error saving transaction: "+err.Error())
				return
			}

			lines := shared.BudgetAlertLines(category)
			msg := "Transaction added!"
			if len(lines) > 0 {
				msg += "\n\nBudget status:\n" + strings.Join(lines, "\n")
			}
			shared.ShowOK(app, msg, func() { app.Stop() })
		}).
		AddButton("Cancel", func() { app.Stop() })

	form.SetBorder(true).SetTitle("[green]Add Transaction").SetTitleAlign(tview.AlignLeft)
	_ = app.SetRoot(form, true).EnableMouse(true).Run()
}

func UpdateInteractive(tx db.Transaction) {
	app := tview.NewApplication()
	var form *tview.Form

	form = tview.NewForm().
		AddInputField("Amount", fmt.Sprintf("%.2f", tx.Amount), 20, nil, nil).
		AddInputField("Category (optional)", tx.Category, 20, nil, nil).
		AddInputField("Description", tx.Description, 50, nil, nil).
		AddInputField("Date (YYYY-MM-DD)", tx.Date.Format("2006-01-02"), 20, nil, nil).
		AddButton("Save", func() {
			amountText := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
			category := form.GetFormItemByLabel("Category (optional)").(*tview.InputField).GetText()
			desc := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
			dateText := form.GetFormItemByLabel("Date (YYYY-MM-DD)").(*tview.InputField).GetText()

			if strings.TrimSpace(category) == "" {
				category = parser.InferCategory(desc)
			}

			amount, err := strconv.ParseFloat(amountText, 64)
			if err != nil {
				shared.ShowError(app, "Invalid amount")
				return
			}

			txDate, err := time.Parse("2006-01-02", dateText)
			if err != nil {
				shared.ShowError(app, "Invalid date (use YYYY-MM-DD)")
				return
			}

			tx.Amount = amount
			tx.Category = category
			tx.Description = desc
			tx.Date = txDate

			if err := db.UpdateTransaction(tx); err != nil {
				shared.ShowError(app, "Error updating transaction: "+err.Error())
				return
			}

			lines := shared.BudgetAlertLines(category)
			msg := "Transaction updated!"
			if len(lines) > 0 {
				msg += "\n\nBudget status:\n" + strings.Join(lines, "\n")
			}
			shared.ShowOK(app, msg, func() { app.Stop() })
		}).
		AddButton("Cancel", func() { app.Stop() })

	form.SetBorder(true).SetTitle(fmt.Sprintf("[green]Edit Transaction ID %d", tx.ID)).SetTitleAlign(tview.AlignLeft)
	_ = app.SetRoot(form, true).EnableMouse(true).Run()
}

// ------------------ Import -------------------

func ImportInteractive() {
	app := tview.NewApplication()
	var form *tview.Form

	form = tview.NewForm().
		AddInputField("File path", "", 60, nil, nil).
		AddButton("Import", func() {
			path := form.GetFormItemByLabel("File path").(*tview.InputField).GetText()
			if path == "" {
				shared.ShowError(app, "No path provided")
				return
			}

			f, err := os.Open(path)
			if err != nil {
				shared.ShowError(app, "Failed to open file: "+err.Error())
				return
			}
			defer f.Close()

			parsed, err := parser.DetectAndParse(f, filepath.Base(path))
			if err != nil {
				shared.ShowError(app, "Parse error: "+err.Error())
				return
			}
			if len(parsed) == 0 {
				shared.ShowOK(app, "No transactions parsed.", func() { app.SetRoot(form, true) })
				return
			}

			if err := parser.InsertParsedTransactions(parsed); err != nil {
				shared.ShowError(app, "Error inserting transactions: "+err.Error())
				return
			}

			var affected []string
			for _, p := range parsed {
				affected = append(affected, p.Category)
			}

			lines := shared.BudgetAlertLinesForCategories(affected)
			msg := fmt.Sprintf("Imported %d transactions.", len(parsed))
			if len(lines) > 0 {
				msg += "\n\nBudget status:\n" + strings.Join(lines, "\n")
			}

			shared.ShowOK(app, msg, func() { app.Stop() })
		}).
		AddButton("Cancel", func() { app.Stop() })

	form.SetBorder(true).SetTitle("[green]Import Transactions from File").SetTitleAlign(tview.AlignLeft)
	_ = app.SetRoot(form, true).EnableMouse(true).Run()
}

// ------------------ Search / Filter -------------------

func searchTransactions() {
	app := tview.NewApplication()
	var form *tview.Form

	form = tview.NewForm().
		AddInputField("Category", "", 20, nil, nil).
		AddInputField("Description contains", "", 30, nil, nil).
		AddInputField("From (YYYY-MM-DD)", "", 12, nil, nil).
		AddInputField("To (YYYY-MM-DD)", "", 12, nil, nil).
		AddInputField("Min amount", "", 12, nil, nil).
		AddInputField("Max amount", "", 12, nil, nil).
		AddButton("Search", func() {
			category := strings.TrimSpace(form.GetFormItemByLabel("Category").(*tview.InputField).GetText())
			desc := strings.TrimSpace(form.GetFormItemByLabel("Description contains").(*tview.InputField).GetText())
			fromStr := strings.TrimSpace(form.GetFormItemByLabel("From (YYYY-MM-DD)").(*tview.InputField).GetText())
			toStr := strings.TrimSpace(form.GetFormItemByLabel("To (YYYY-MM-DD)").(*tview.InputField).GetText())
			minStr := strings.TrimSpace(form.GetFormItemByLabel("Min amount").(*tview.InputField).GetText())
			maxStr := strings.TrimSpace(form.GetFormItemByLabel("Max amount").(*tview.InputField).GetText())

			var fromPtr, toPtr *time.Time
			if fromStr != "" {
				t, err := time.Parse("2006-01-02", fromStr)
				if err != nil {
					shared.ShowError(app, "Invalid from date (YYYY-MM-DD)")
					return
				}
				fromPtr = &t
			}
			if toStr != "" {
				t, err := time.Parse("2006-01-02", toStr)
				if err != nil {
					shared.ShowError(app, "Invalid to date (YYYY-MM-DD)")
					return
				}
				toPtr = &t
			}

			var minPtr, maxPtr *float64
			if minStr != "" {
				v, err := strconv.ParseFloat(minStr, 64)
				if err != nil {
					shared.ShowError(app, "Invalid min amount")
					return
				}
				minPtr = &v
			}
			if maxStr != "" {
				v, err := strconv.ParseFloat(maxStr, 64)
				if err != nil {
					shared.ShowError(app, "Invalid max amount")
					return
				}
				maxPtr = &v
			}

			txs, err := db.GetFilteredTransactions(category, desc, fromPtr, toPtr, minPtr, maxPtr)
			if err != nil {
				shared.ShowError(app, "Search error: "+err.Error())
				return
			}
			if len(txs) == 0 {
				shared.ShowOK(app, "No transactions found.", func() { app.SetRoot(form, true) })
				return
			}

			reload := func() ([]db.Transaction, error) {
				return db.GetFilteredTransactions(category, desc, fromPtr, toPtr, minPtr, maxPtr)
			}

			app.Suspend(func() {
				showTransactionsTable(txs, reload)
			})
		}).
		AddButton("Back", func() { app.Stop() })

	form.SetBorder(true).SetTitle("[green]Search / Filter Transactions").SetTitleAlign(tview.AlignLeft)
	_ = app.SetRoot(form, true).EnableMouse(true).Run()
}
