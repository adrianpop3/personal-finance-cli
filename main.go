package main

import (
	"log"
	"personal-finance-cli/cmd/tui"
	"personal-finance-cli/db"
)

func main() {
	// Initialize the database
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize DB:", err)
	}

	// Launch TUI
	if err := tui.RunMainMenu(); err != nil {
		log.Fatal(err)
	}
}
