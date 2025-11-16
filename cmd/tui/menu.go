package tui

import (
	"personal-finance-cli/cmd/tui/transaction"

	"personal-finance-cli/cmd/tui/budget"

	"github.com/rivo/tview"
)

func RunMainMenu() error {
	app := tview.NewApplication()

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[::b][green]💰 Personal Finance CLI[::-]").
		SetDynamicColors(true)

	buttons := tview.NewList().
		AddItem("Transactions", "Manage your transactions", 't', func() {
			app.Suspend(func() { transaction.RunTUI() })
		}).
		AddItem("Budgets", "Manage your budgets", 'b', func() {
			app.Suspend(func() { budget.RunTUI() })
		}).
		AddItem("Exit", "Quit the application", 'q', func() {
			app.Stop()
		})

	buttons.SetBorder(true).SetTitle("Menu").SetTitleAlign(tview.AlignCenter)
	buttons.SetMainTextColor(tview.Styles.PrimaryTextColor)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 5, 1, false).
		AddItem(buttons, 0, 2, true)

	return app.SetRoot(flex, true).EnableMouse(true).Run()
}
