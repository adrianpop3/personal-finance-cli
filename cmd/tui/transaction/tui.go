package transaction

import (
	"fmt"
	"os"
	"path/filepath"
	"personal-finance-cli/db"
	"personal-finance-cli/internal/parser"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RunTUI() {
	app := tview.NewApplication()

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[::b][green]💰 Transactions Menu[::-]").
		SetDynamicColors(true)

	labels := []string{"List Transactions", "Add Transaction", "Import From File", "Back"}
	actions := []func(){
		func() { app.Suspend(showTransactions) },
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

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		fmt.Println(err)
	}
}

// ------------------ Transaction Table -------------------

func showTransactions() {
	txs, err := db.GetTransactions()
	if err != nil {
		fmt.Println("Error fetching transactions:", err)
		return
	}

	app := tview.NewApplication()
	table := tview.NewTable().SetSelectable(true, false)
	table.SetBorder(true).SetTitle("[green]Transactions (Enter=Edit/Delete, ESC=Back)").SetTitleAlign(tview.AlignCenter)

	headers := []string{"ID", "Amount", "Category", "Date", "Description"}
	for i, h := range headers {
		table.SetCell(0, i, tview.NewTableCell(fmt.Sprintf("[::b][green]%s[::-]", h)).SetSelectable(false))
	}

	for r, t := range txs {
		table.SetCell(r+1, 0, tview.NewTableCell(strconv.Itoa(t.ID)))
		table.SetCell(r+1, 1, tview.NewTableCell(fmt.Sprintf("%.2f", t.Amount)))
		table.SetCell(r+1, 2, tview.NewTableCell(t.Category))
		table.SetCell(r+1, 3, tview.NewTableCell(t.Date.Format("2006-01-02")))
		table.SetCell(r+1, 4, tview.NewTableCell(t.Description))
	}

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

// ------------------ Transaction Modal -------------------

func showTransactionActions(tx db.Transaction, parentTable *tview.Table, app *tview.Application) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("[green]Transaction ID %d\nChoose an action[::-]", tx.ID)).
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

	form.SetBorder(true).SetTitle("[green]Add Transaction").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true).EnableMouse(true).Run()
}

func UpdateInteractive(tx db.Transaction) {
	app := tview.NewApplication()
	var form *tview.Form
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

	form.SetBorder(true).SetTitle(fmt.Sprintf("[green]Edit Transaction ID %d", tx.ID)).SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true).EnableMouse(true).Run()
}

// ------------------ Import From File flow -------------------

func ImportInteractive() {
	app := tview.NewApplication()
	var form *tview.Form
	form = tview.NewForm().
		AddInputField("File path", "", 60, nil, nil).
		AddButton("Import", func() {
			path := form.GetFormItemByLabel("File path").(*tview.InputField).GetText()
			if path == "" {
				fmt.Println("No path provided")
				return
			}
			f, err := os.Open(path)
			if err != nil {
				fmt.Println("Failed to open file:", err)
				return
			}
			defer f.Close()

			parsed, err := parser.DetectAndParse(f, filepath.Base(path))
			if err != nil {
				fmt.Println("Parse error:", err)
				return
			}
			if len(parsed) == 0 {
				fmt.Println("No transactions parsed")
				return
			}

			if err := parser.InsertParsedTransactions(parsed); err != nil {
				fmt.Println("Error inserting transactions:", err)
				return
			}

			catSet := map[string]struct{}{}
			for _, p := range parsed {
				if p.Amount < 0 {
					catSet[p.Category] = struct{}{}
				}
			}

			var summaryLines []string
			for cat := range catSet {
				budgets, err := db.GetBudgets()
				if err != nil {
					continue
				}
				for _, b := range budgets {
					if b.Category != cat {
						continue
					}
					rem, err := db.GetBudgetRemaining(b)
					if err != nil {
						continue
					}
					summaryLines = append(summaryLines, fmt.Sprintf("- %s (%s): remaining %.2f (limit %.2f)", b.Category, b.Period, rem, b.Amount))
				}
			}

			msg := fmt.Sprintf("Imported %d transactions\n", len(parsed))
			if len(summaryLines) > 0 {
				msg += "Updated budgets:\n"
				for _, l := range summaryLines {
					msg += l + "\n"
				}
			} else {
				msg += "No matching budgets were affected."
			}

			m := tview.NewModal().
				SetText("[green]" + msg + "[::-]").
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(i int, lbl string) {
					app.Stop()
				})
			app.SetRoot(m, false)
		}).
		AddButton("Cancel", func() { app.Stop() })

	form.SetBorder(true).SetTitle("[green]Import Transactions from File").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true).EnableMouse(true).Run()
}
