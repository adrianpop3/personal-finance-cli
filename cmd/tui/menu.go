// package tui

// import (
// 	"personal-finance-cli/cmd/tui/transaction"

// 	"github.com/rivo/tview"
// )

// // RunMainMenu launches the main menu
// func RunMainMenu() error {
// 	app := tview.NewApplication()

// 	list := tview.NewList().
// 		AddItem("Transactions", "Manage transactions", 't', func() {
// 			app.Suspend(func() {
// 				transaction.RunTUI()
// 			})
// 		}).
// 		AddItem("Exit", "Quit", 'q', func() {
// 			app.Stop()
// 		})

// 	list.SetBorder(true).SetTitle("Personal Finance CLI").SetTitleAlign(tview.AlignLeft)

//		return app.SetRoot(list, true).EnableMouse(true).Run()
//	}
package tui

import (
	"personal-finance-cli/cmd/tui/transaction"

	"github.com/rivo/tview"
)

// RunMainMenu launches the main menu with styled title and buttons
func RunMainMenu() error {
	app := tview.NewApplication()

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[::b][green]💰 Personal Finance CLI[::-]").
		SetDynamicColors(true)

	// Buttons
	buttons := tview.NewList().
		AddItem("Transactions", "Manage your transactions", 't', func() {
			app.Suspend(func() { transaction.RunTUI() })
		}).
		AddItem("Budgets", "Manage your budgets", 'b', func() {
			// TODO: launch budget TUI
		}).
		AddItem("Exit", "Quit the application", 'q', func() {
			app.Stop()
		})

	buttons.SetBorder(true).SetTitle("Menu").SetTitleAlign(tview.AlignCenter)
	buttons.SetMainTextColor(tview.Styles.PrimaryTextColor)

	// Layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 5, 1, false).
		AddItem(buttons, 0, 2, true)

	return app.SetRoot(flex, true).EnableMouse(true).Run()
}
