# Personal Finance CLI Manager

A **command-line personal finance manager** a for tracking personal income and expenses. Import transactions from
bank statements, categorize them automatically, set budgets, and generate insightful reports in a simple manner.
At the moment the functionalities available in this project are CRUD operations for transactions and budgets.
These operations can be made straight from the terminal using the defined commands or using the TUI views that were added for better user exerience.

---

## Technologies Used

- **Go (Golang)** – Core language for building the CLI.
- **SQLite** – Lightweight local database for storing transactions, budgets, and categories.
- **Cobra** – CLI framework for building commands and subcommands (`add`, `update`, `delete`, `list`, etc.).
- **tview & tcell** – Libraries used to build an interactive, arrow-navigable Terminal UI (TUI).

---

## Features

### 1. Transactions & Budgets
- Full CRUD operations:
  - **Add Transaction**
  - **Update Transaction**
  - **Delete Transaction**
  - **List Transactions** (All or by ID)
- Interactive TUI:
  - Arrow navigation
  - Green-themed buttons
  - Edit/Delete modal for each transaction
  - Add/Update forms fully functional

### 3. Terminal UI
- Main menu with:
  - Transactions
  - Budgets
  - Exit
- Arrow navigation for all menus
- Green-colored styling throughout (buttons, headers, modals)
- Forms for adding/editing items with proper validation
- Modals for edit/delete confirmation

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

3. Run the CLI:

```bash
go run main.go
```

Or, make an executable by running these two commands after previous point 2. :

```bash
go build -o fincli main.go
./fincli
```

## Usage

# Example of CLI Commands (via Cobra)

- transaction add --amount 50 --category Food --description "Groceries"
- transaction update --id 1 --amount 60
- transaction delete --id 1
- transaction list
- transaction list --id 1

- budget add --category Food --amount 200 --period monthly
- budget update --id 1 --amount 250
- budget delete --id 1
- budget list
- budget list --id 1

# Example of TUI views

<img width="1071" height="210" alt="Captură de ecran din 2025-11-16 la 20 47 11" src="https://github.com/user-attachments/assets/52f7eab3-5c17-47c5-9647-9487e345c9cd" />

<img width="1071" height="210" alt="Captură de ecran din 2025-11-16 la 20 47 17" src="https://github.com/user-attachments/assets/3ed261e6-aa29-4a6d-8862-49d7b1552f85" />

<img width="1071" height="210" alt="Captură de ecran din 2025-11-16 la 20 54 38" src="https://github.com/user-attachments/assets/908c799e-f967-49b8-8a8d-36e1b4859ac9" />

<img width="1071" height="210" alt="Captură de ecran din 2025-11-16 la 20 47 37" src="https://github.com/user-attachments/assets/3e240f6d-175e-4464-b8dd-b3e679315828" />

<img width="1071" height="210" alt="Captură de ecran din 2025-11-16 la 20 48 02" src="https://github.com/user-attachments/assets/6d38eff9-96eb-4598-8243-39832a6932ea" />

<img width="1071" height="210" alt="Captură de ecran din 2025-11-16 la 20 57 01" src="https://github.com/user-attachments/assets/8d66473c-5673-4220-9842-78fe2be1eea3" />

<img width="1071" height="210" alt="Captură de ecran din 2025-11-16 la 20 48 26" src="https://github.com/user-attachments/assets/847e819c-b4b7-4e80-9140-f645b1a30324" />




