# Personal Finance CLI Manager

A **command-line personal finance manager** for tracking personal income and expenses.  
The application allows you to import transactions from bank statements, categorize them automatically, set budgets, receive alerts when budgets are exceeded, and generate insightful reports — all directly from the terminal.

This project was developed as a **faculty assignment**, with a focus on clarity, correctness, and usability in a constrained CLI/TUI environment.

---

## Technologies Used

- **Go (Golang)** – Core language for building the CLI application
- **SQLite** – Lightweight local database for storing transactions, budgets, and categories
- **Cobra** – CLI framework for commands and subcommands (`add`, `update`, `delete`, `list`, `import`, `report`, etc.)
- **tview & tcell** – Libraries used to build an interactive, arrow-navigable Terminal UI (TUI)

---

## Features

### 1. Transactions & Budgets
- Full CRUD operations:
  - **Add Transactions / Budgets**
  - **Update Transactions / Budgets**
  - **Delete Transactions / Budgets**
  - **List Transactions / Budgets** (all or by ID)
- Import transactions from files:
  - **CSV**
  - **OFX**
- Data stored locally using SQLite
- Operations available through:
  - CLI commands (Cobra)
  - Interactive TUI views

---

### 2. Automatic Categorization
- Imported transactions (CSV / OFX) are automatically categorized using **regex-based rules** applied to transaction descriptions.
- Manual transactions added via CLI or TUI are automatically categorized **when the category field is left empty**.
- If no rule matches, the category defaults to `Uncategorized`.

---

### 3. Budget Tracking & Alerts
- Budgets can be set **per category**, typically on a monthly basis.
- Budget tracking automatically accounts for expenses in the selected period.
- **Alerts are generated when:**
  - A budget is exceeded
  - A budget is close to its limit (≤10% remaining)
- Alerts are shown:
  - Automatically after adding, updating, or importing transactions (CLI)
  - Via a dedicated **Budget Status / Alerts** view in the TUI

---

### 4. Search & Filter Transactions
- Transactions can be filtered using CLI flags:
  - By category
  - By description keywords
  - By date range
  - By minimum / maximum amount
- Multiple filters can be combined to narrow down results.

---

### 5. Reports & Insights
- Reports are generated directly in the terminal:
  - **Monthly totals** (income, expenses, net balance)
  - **Category breakdown** with ASCII bar charts
- Reports are available via:
  - CLI commands
  - Dedicated TUI views
- ASCII bar charts are used to visualize spending distribution per category.

---

### 6. Terminal User Interface (TUI)
- Interactive, arrow&mouse-navigable interface
- Main menu with:
  - Transactions
  - Budgets
  - Reports
  - Exit
- Blue-Green-themed styling throughout:
  - Buttons
  - Headers
  - Tables
  - Modals
- Features include:
  - Editable tables
  - Add / update forms with validation
  - Confirmation modals
  - Budget status and alert tables
  - Report views with charts

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/personal-finance-cli.git
cd personal-finance-cli
```

2. Install dependencies:

```bash
go mod tidy
```

3. Run the application:

```bash
go run main.go
```

Or build an executable:

```bash
go build -o fincli main.go
./fincli
```

---

## Usage

### CLI Examples (Cobra)

#### Transactions
```bash
transaction add --amount -50 --description "Groceries Lidl"
transaction update --id 1 --amount -60
transaction delete --id 1
transaction list
transaction list --id 1
```

#### Filtering Transactions
```bash
transaction list --category Food
transaction list --desc coffee
transaction list --from 2026-01-01 --to 2026-01-31 --min -100 --max -5
transaction list --category Coffee --desc Starbucks
```

#### Import Transactions
```bash
transaction import --file path/to/file.csv
transaction import --file path/to/file.ofx
```

---

#### Budgets
```bash
budget add --category Food --amount 300 --period monthly
budget update --id 1 --amount 350
budget delete --id 1
budget list
budget list --id 1
budget status
```

---

#### Reports
```bash
report monthly --month 2026-01
report categories --month 2026-01
```

---

## Terminal UI (TUI)

The application can also be fully used through an interactive TUI by simply running:

```bash
./fincli
```

### Example Screens

- Main menu (Transactions / Budgets / Reports)
- Transactions table with edit/delete modal
- Add / update transaction forms
- Search & filter transactions
- Budgets table and budget status alerts
- Monthly and category-based reports with ASCII charts

Screenshots showcasing these views can be found below.

---

## License

This project is intended for educational purposes.
