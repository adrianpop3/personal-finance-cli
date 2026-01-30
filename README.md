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

<img width="1332" height="313" alt="Captură de ecran din 2026-01-31 la 00 34 48" src="https://github.com/user-attachments/assets/59ab7bbd-7b46-4efb-bb8f-fe13eb284363" />
<img width="1339" height="362" alt="Captură de ecran din 2026-01-31 la 00 37 45" src="https://github.com/user-attachments/assets/beae25fc-9b6d-4f1b-aa55-a62e8f04f225" />
<img width="1339" height="569" alt="Captură de ecran din 2026-01-31 la 00 38 00" src="https://github.com/user-attachments/assets/83fe06cd-6a28-4ea9-80b4-e163d7fb20a4" />
<img width="1332" height="569" alt="Captură de ecran din 2026-01-31 la 00 40 34" src="https://github.com/user-attachments/assets/89c1ba0e-dc99-42b5-b57f-5a2bbb45e5f4" />
<img width="1332" height="569" alt="Captură de ecran din 2026-01-31 la 00 41 06" src="https://github.com/user-attachments/assets/7b830e60-4dda-4242-a153-21902c2098a1" />
<img width="1332" height="254" alt="Captură de ecran din 2026-01-31 la 00 49 00" src="https://github.com/user-attachments/assets/96b95129-de58-476b-8d84-70805612ff18" />
<img width="1332" height="309" alt="Captură de ecran din 2026-01-31 la 00 41 24" src="https://github.com/user-attachments/assets/78ad2498-145a-4d19-b641-8081b20a5cae" />
<img width="1332" height="254" alt="Captură de ecran din 2026-01-31 la 01 22 00" src="https://github.com/user-attachments/assets/9d08662a-8a1c-4b7d-9fed-69d2e8c02324" />
<img width="1332" height="569" alt="Captură de ecran din 2026-01-31 la 00 42 10" src="https://github.com/user-attachments/assets/d051aaaf-aaff-4f94-b4f8-c7fc658ad0f1" />
<img width="1332" height="569" alt="Captură de ecran din 2026-01-31 la 00 42 20" src="https://github.com/user-attachments/assets/7ceb3ffe-2fc1-40ee-a55f-b0749f7642ab" />
<img width="1332" height="254" alt="Captură de ecran din 2026-01-31 la 00 42 41" src="https://github.com/user-attachments/assets/987efc62-8e5b-472c-bdaf-ea8f0ff5f282" />
<img width="1332" height="254" alt="Captură de ecran din 2026-01-31 la 00 42 51" src="https://github.com/user-attachments/assets/70336705-7ade-4f2e-bf6e-1399fa6588ab" />
<img width="1332" height="254" alt="Captură de ecran din 2026-01-31 la 00 43 03" src="https://github.com/user-attachments/assets/da08113d-fbb6-4788-a99a-38fb313bb7e3" />
<img width="1332" height="254" alt="Captură de ecran din 2026-01-31 la 00 43 12" src="https://github.com/user-attachments/assets/038aa0c8-b81e-41e8-99e7-6c7eb8b7a8bc" />
<img width="1332" height="254" alt="Captură de ecran din 2026-01-31 la 00 43 22" src="https://github.com/user-attachments/assets/c513f9a1-7a5d-414a-af29-7fde51fb1933" />

---

## License

This project is intended for educational purposes.
