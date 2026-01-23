package db

import (
	"database/sql"
	"fmt"
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
		date TEXT NOT NULL,
		external_id TEXT
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_transactions_external_id
	ON transactions(external_id)
	WHERE external_id IS NOT NULL;

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
	ExternalID  string
}

func InsertTransaction(tx Transaction) error {
	_, err := database.Exec(
		`INSERT INTO transactions 
		 (amount, description, category, date, external_id)
		 VALUES (?, ?, ?, ?, ?)`,
		tx.Amount,
		tx.Description,
		tx.Category,
		tx.Date.Format("2006-01-02"),
		nullIfEmpty(tx.ExternalID),
	)
	return err
}

func TransactionExistsByExternalID(externalID string) (bool, error) {
	if externalID == "" {
		return false, nil
	}

	var count int
	err := database.QueryRow(
		`SELECT COUNT(1) FROM transactions WHERE external_id = ?`,
		externalID,
	).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetTransactions() ([]Transaction, error) {
	rows, err := database.Query(`
		SELECT id, amount, description, category, date, external_id
		FROM transactions
		ORDER BY date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []Transaction
	for rows.Next() {
		var t Transaction
		var dateStr string
		var extID sql.NullString

		if err := rows.Scan(
			&t.ID,
			&t.Amount,
			&t.Description,
			&t.Category,
			&dateStr,
			&extID,
		); err != nil {
			return nil, err
		}

		t.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			t.Date = time.Time{}
		}

		if extID.Valid {
			t.ExternalID = extID.String
		}

		txs = append(txs, t)
	}
	return txs, nil
}

func GetFilteredTransactions(category, desc string, from, to *time.Time, min, max *float64) ([]Transaction, error) {
	if database == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `SELECT id, amount, description, category, date FROM transactions WHERE 1=1`
	args := []interface{}{}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	if desc != "" {
		query += " AND description LIKE ?"
		args = append(args, "%"+desc+"%")
	}

	if from != nil {
		query += " AND date >= ?"
		args = append(args, from.Format("2006-01-02"))
	}

	if to != nil {
		query += " AND date <= ?"
		args = append(args, to.Format("2006-01-02"))
	}

	if min != nil {
		query += " AND amount >= ?"
		args = append(args, *min)
	}

	if max != nil {
		query += " AND amount <= ?"
		args = append(args, *max)
	}

	query += " ORDER BY date DESC"

	rows, err := database.Query(query, args...)
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
		t.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			t.Date = time.Time{}
		}
		txs = append(txs, t)
	}
	return txs, nil
}

func UpdateTransaction(t Transaction) error {
	_, err := database.Exec(
		`UPDATE transactions
		 SET amount = ?, description = ?, category = ?, date = ?, external_id = ?
		 WHERE id = ?`,
		t.Amount,
		t.Description,
		t.Category,
		t.Date.Format("2006-01-02"),
		nullIfEmpty(t.ExternalID),
		t.ID,
	)
	return err
}

func DeleteTransaction(id int) error {
	_, err := database.Exec(`DELETE FROM transactions WHERE id = ?`, id)
	return err
}

func GetTransactionByID(id int) (*Transaction, error) {
	row := database.QueryRow(`
		SELECT id, amount, description, category, date, external_id
		FROM transactions
		WHERE id = ?
	`, id)

	var t Transaction
	var dateStr string
	var extID sql.NullString

	err := row.Scan(
		&t.ID,
		&t.Amount,
		&t.Description,
		&t.Category,
		&dateStr,
		&extID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	t.Date, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		t.Date = time.Time{}
	}

	if extID.Valid {
		t.ExternalID = extID.String
	}

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
		`INSERT INTO budgets (category, amount, period)
		 VALUES (?, ?, ?)`,
		b.Category,
		b.Amount,
		b.Period,
	)
	return err
}

func GetBudgets() ([]Budget, error) {
	rows, err := database.Query(`
		SELECT id, category, amount, period
		FROM budgets
		ORDER BY period DESC
	`)
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
	row := database.QueryRow(`
		SELECT id, category, amount, period
		FROM budgets
		WHERE id = ?
	`, id)

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
		`UPDATE budgets
		 SET category = ?, amount = ?, period = ?
		 WHERE id = ?`,
		b.Category,
		b.Amount,
		b.Period,
		b.ID,
	)
	return err
}

func DeleteBudget(id int) error {
	_, err := database.Exec(`DELETE FROM budgets WHERE id = ?`, id)
	return err
}

func GetBudgetRemaining(b Budget) (float64, error) {
	if database == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	period := b.Period
	if period == "monthly" || period == "" {
		period = time.Now().Format("2006-01")
	}

	query := `
	SELECT COALESCE(
		SUM(CASE WHEN amount < 0 THEN -amount ELSE 0 END),
		0
	)
	FROM transactions
	WHERE category = ?
	AND strftime('%Y-%m', date) = ?
	`

	var expenses float64
	err := database.QueryRow(query, b.Category, period).Scan(&expenses)
	if err != nil {
		return 0, err
	}

	return b.Amount - expenses, nil
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
