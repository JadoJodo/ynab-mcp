package ynab

import "time"

// Response is the generic YNAB API response envelope.
type Response[T any] struct {
	Data  T      `json:"data"`
	Error *Error `json:"error,omitempty"`
}

// Error represents a YNAB API error.
type Error struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

func (e *Error) Error() string {
	return e.Detail
}

// BudgetSummary is the abbreviated budget returned by list endpoints.
type BudgetSummary struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	LastModifiedOn *time.Time `json:"last_modified_on,omitempty"`
	DateFormat     *struct {
		Format string `json:"format"`
	} `json:"date_format,omitempty"`
	CurrencyFormat *struct {
		ISOCode          string `json:"iso_code"`
		DecimalDigits    int    `json:"decimal_digits"`
		DecimalSeparator string `json:"decimal_separator"`
		GroupSeparator   string `json:"group_separator"`
		CurrencySymbol   string `json:"currency_symbol"`
		SymbolFirst      bool   `json:"symbol_first"`
	} `json:"currency_format,omitempty"`
}

// BudgetDetail is the full budget returned by the get endpoint.
type BudgetDetail struct {
	BudgetSummary
	Accounts       []Account       `json:"accounts,omitempty"`
	CategoryGroups []CategoryGroup `json:"category_groups,omitempty"`
	Categories     []Category      `json:"categories,omitempty"`
	Payees         []Payee         `json:"payees,omitempty"`
	Months     []Month     `json:"months,omitempty"`
}

// Account represents a YNAB account.
type Account struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	OnBudget          bool   `json:"on_budget"`
	Closed            bool   `json:"closed"`
	Balance           int64  `json:"balance"`
	ClearedBalance    int64  `json:"cleared_balance"`
	UnclearedBalance  int64  `json:"uncleared_balance"`
	Deleted           bool   `json:"deleted"`
}

// CategoryGroup represents a YNAB category group with its categories.
type CategoryGroup struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Hidden     bool       `json:"hidden"`
	Deleted    bool       `json:"deleted"`
	Categories []Category `json:"categories"`
}

// Category represents a single YNAB category.
type Category struct {
	ID                      string `json:"id"`
	CategoryGroupID         string `json:"category_group_id"`
	CategoryGroupName       string `json:"category_group_name,omitempty"`
	Name                    string `json:"name"`
	Hidden                  bool   `json:"hidden"`
	Budgeted                int64  `json:"budgeted"`
	Activity                int64  `json:"activity"`
	Balance                 int64  `json:"balance"`
	GoalType                string `json:"goal_type,omitempty"`
	GoalTarget              int64  `json:"goal_target,omitempty"`
	GoalTargetMonth         string `json:"goal_target_month,omitempty"`
	GoalPercentageComplete  int    `json:"goal_percentage_complete,omitempty"`
	Deleted                 bool   `json:"deleted"`
}

// Transaction represents a YNAB transaction.
type Transaction struct {
	ID                string        `json:"id"`
	Date              string        `json:"date"`
	Amount            int64         `json:"amount"`
	Memo              string        `json:"memo,omitempty"`
	Cleared           string        `json:"cleared"`
	Approved          bool          `json:"approved"`
	FlagColor         string        `json:"flag_color,omitempty"`
	FlagName          string        `json:"flag_name,omitempty"`
	AccountID         string        `json:"account_id"`
	AccountName       string        `json:"account_name,omitempty"`
	PayeeID           string        `json:"payee_id,omitempty"`
	PayeeName         string        `json:"payee_name,omitempty"`
	CategoryID        string        `json:"category_id,omitempty"`
	CategoryName      string        `json:"category_name,omitempty"`
	TransferAccountID string        `json:"transfer_account_id,omitempty"`
	Deleted           bool          `json:"deleted"`
	Subtransactions   []Subtransaction `json:"subtransactions,omitempty"`
}

// Subtransaction represents a split within a transaction.
type Subtransaction struct {
	ID                string `json:"id"`
	TransactionID     string `json:"transaction_id"`
	Amount            int64  `json:"amount"`
	Memo              string `json:"memo,omitempty"`
	PayeeID           string `json:"payee_id,omitempty"`
	PayeeName         string `json:"payee_name,omitempty"`
	CategoryID        string `json:"category_id,omitempty"`
	CategoryName      string `json:"category_name,omitempty"`
	TransferAccountID string `json:"transfer_account_id,omitempty"`
	Deleted           bool   `json:"deleted"`
}

// SaveTransaction is used for creating transactions.
type SaveTransaction struct {
	AccountID  string `json:"account_id"`
	Date       string `json:"date"`
	Amount     int64  `json:"amount"`
	PayeeName  string `json:"payee_name,omitempty"`
	CategoryID string `json:"category_id,omitempty"`
	Memo       string `json:"memo,omitempty"`
	Cleared    string `json:"cleared,omitempty"`
	Approved   bool   `json:"approved"`
}

// UpdateTransaction is used for updating transactions.
type UpdateTransaction struct {
	AccountID  *string `json:"account_id,omitempty"`
	Date       *string `json:"date,omitempty"`
	Amount     *int64  `json:"amount,omitempty"`
	PayeeID    *string `json:"payee_id,omitempty"`
	PayeeName  *string `json:"payee_name,omitempty"`
	CategoryID *string `json:"category_id,omitempty"`
	Memo       *string `json:"memo,omitempty"`
	Cleared    *string `json:"cleared,omitempty"`
	Approved   *bool   `json:"approved,omitempty"`
	FlagColor  *string `json:"flag_color,omitempty"`
}

// Payee represents a YNAB payee.
type Payee struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	TransferAccountID string `json:"transfer_account_id,omitempty"`
	Deleted           bool   `json:"deleted"`
}

// Month represents a YNAB budget month.
type Month struct {
	Month      string     `json:"month"`
	Note       string     `json:"note,omitempty"`
	Income     int64      `json:"income"`
	Budgeted   int64      `json:"budgeted"`
	Activity   int64      `json:"activity"`
	ToBeBudgeted int64    `json:"to_be_budgeted"`
	Categories []Category `json:"categories,omitempty"`
	Deleted    bool       `json:"deleted"`
}
