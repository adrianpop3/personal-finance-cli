package tui

import (
	"personal-finance-cli/cmd/tui/budget"
	"personal-finance-cli/cmd/tui/transaction"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// RunMainMenu launches the main menu with arrow navigation and green theme
func RunMainMenu() error {
	app := tview.NewApplication()

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[::b][green]💰 Personal Finance CLI[::-]").
		SetDynamicColors(true)

	labels := []string{"Transactions", "Budgets", "Exit"}
	actions := []func(){
		func() { app.Suspend(func() { transaction.RunTUI() }) },
		func() { app.Suspend(func() { budget.RunTUI() }) },
		func() { app.Stop() },
	}

	current := 0
	buttonFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	buttons := []*tview.Button{}

	for i, label := range labels {
		idx := i
		btn := tview.NewButton("[green]" + label).
			SetSelectedFunc(actions[idx])
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

	return app.SetRoot(layout, true).EnableMouse(true).Run()
}
