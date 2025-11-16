package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func InitDB() error {
	var err error
	if database != nil {
		return nil
	}
	database, err = sql.Open("sqlite3", "./finance.db")
	if err != nil {
		return err
	}

	_, _ = database.Exec("PRAGMA journal_mode = WAL;")
	_, _ = database.Exec("PRAGMA foreign_keys = ON;")

	schema := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL,
		description TEXT,
		category TEXT,
		date TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS budgets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL,
		amount REAL NOT NULL,
		period TEXT NOT NULL,
		UNIQUE(category, period)
	);
	`

	_, err = database.Exec(schema)
	return err
}

type Transaction struct {
	ID          int
	Amount      float64
	Description string
	Category    string
	Date        time.Time
}

func InsertTransaction(tx Transaction) error {
	_, err := database.Exec(
		`INSERT INTO transactions (amount, description, category, date) VALUES (?, ?, ?, ?)`,
		tx.Amount, tx.Description, tx.Category, tx.Date.Format("2006-01-02"),
	)
	return err
}

func GetTransactions() ([]Transaction, error) {
	rows, err := database.Query(`SELECT id, amount, description, category, date FROM transactions ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []Transaction
	for rows.Next() {
		var t Transaction
		var dateStr string
		if err := rows.Scan(&t.ID, &t.Amount, &t.Description, &t.Category, &dateStr); err != nil {
			return nil, err
		}
		t.Date, _ = time.Parse("2006-01-02", dateStr)
		txs = append(txs, t)
	}
	return txs, nil
}

func UpdateTransaction(t Transaction) error {
	_, err := database.Exec(
		`UPDATE transactions SET amount = ?, description = ?, category = ?, date = ? WHERE id = ?`,
		t.Amount, t.Description, t.Category, t.Date.Format("2006-01-02"), t.ID,
	)
	return err
}

func DeleteTransaction(id int) error {
	_, err := database.Exec(`DELETE FROM transactions WHERE id = ?`, id)
	return err
}

func GetTransactionByID(id int) (*Transaction, error) {
	row := database.QueryRow(`SELECT id, amount, description, category, date FROM transactions WHERE id = ?`, id)

	var t Transaction
	var dateStr string
	err := row.Scan(&t.ID, &t.Amount, &t.Description, &t.Category, &dateStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	t.Date, _ = time.Parse("2006-01-02", dateStr)
	return &t, nil
}

type Budget struct {
	ID       int
	Category string
	Amount   float64
	Period   string
}

func InsertBudget(b Budget) error {
	_, err := database.Exec(
		`INSERT INTO budgets (category, amount, period) VALUES (?, ?, ?)`,
		b.Category, b.Amount, b.Period,
	)
	return err
}

func GetBudgets() ([]Budget, error) {
	rows, err := database.Query(`SELECT id, category, amount, period FROM budgets ORDER BY period DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []Budget
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.ID, &b.Category, &b.Amount, &b.Period); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}

func GetBudgetByID(id int) (*Budget, error) {
	row := database.QueryRow(`SELECT id, category, amount, period FROM budgets WHERE id = ?`, id)

	var b Budget
	err := row.Scan(&b.ID, &b.Category, &b.Amount, &b.Period)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func UpdateBudget(b Budget) error {
	_, err := database.Exec(
		`UPDATE budgets SET category = ?, amount = ?, period = ? WHERE id = ?`,
		b.Category, b.Amount, b.Period, b.ID,
	)
	return err
}

func DeleteBudget(id int) error {
	_, err := database.Exec(`DELETE FROM budgets WHERE id = ?`, id)
	return err
}
