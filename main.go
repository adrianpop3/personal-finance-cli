package main

import (
	"log"
	"os"

	"personal-finance-cli/cmd"
	"personal-finance-cli/cmd/tui"
	"personal-finance-cli/db"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize DB:", err)
	}

	if len(os.Args) == 1 {
		if err := tui.RunMainMenu(); err != nil {
			log.Fatal(err)
		}
		return
	}

	cmd.Execute()
}
