package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JadoJodo/ynab-mcp/ynab"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var validFlagColors = map[string]bool{
	"red": true, "orange": true, "yellow": true,
	"green": true, "blue": true, "purple": true,
}

// normalizeFlagColor validates and normalizes a flag_color JSON input.
// Returns (nil, nil) if the field should be omitted (not provided).
// Returns (json.RawMessage("null"), nil) to clear the flag.
// Returns (`"red"` etc, nil) for a valid color.
// Returns an error for invalid colors.
func normalizeFlagColor(raw json.RawMessage) (json.RawMessage, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, nil
	}
	if bytes.Equal(trimmed, []byte("null")) {
		return json.RawMessage("null"), nil
	}
	var s string
	if err := json.Unmarshal(trimmed, &s); err != nil {
		return nil, fmt.Errorf("flag_color must be a string or null, got %s", string(trimmed))
	}
	if !validFlagColors[s] {
		return nil, fmt.Errorf("flag_color %q is invalid; must be one of red, orange, yellow, green, blue, purple, or null", s)
	}
	return json.Marshal(s)
}

type listTransactionsInput struct {
	BudgetID   string `json:"budget_id" jsonschema:"Budget ID or last-used"`
	SinceDate  string `json:"since_date,omitempty" jsonschema:"Only return transactions on or after this date (YYYY-MM-DD)"`
	Type       string `json:"type,omitempty" jsonschema:"Filter by type: uncategorized or unapproved"`
	AccountID  string `json:"account_id,omitempty" jsonschema:"Filter by account ID (cannot combine with category_id or payee_id)"`
	CategoryID string `json:"category_id,omitempty" jsonschema:"Filter by category ID (cannot combine with account_id or payee_id)"`
	PayeeID    string `json:"payee_id,omitempty" jsonschema:"Filter by payee ID (cannot combine with account_id or category_id)"`
}

type getTransactionInput struct {
	BudgetID      string `json:"budget_id" jsonschema:"Budget ID or last-used"`
	TransactionID string `json:"transaction_id" jsonschema:"Transaction ID"`
}

type createTransactionInput struct {
	BudgetID   string  `json:"budget_id" jsonschema:"Budget ID or last-used"`
	AccountID  string  `json:"account_id" jsonschema:"Account ID for the transaction"`
	Date       string  `json:"date" jsonschema:"Transaction date (YYYY-MM-DD)"`
	Amount     float64 `json:"amount" jsonschema:"Amount in dollars (negative for outflow e.g. -12.50)"`
	PayeeName  string  `json:"payee_name,omitempty" jsonschema:"Payee name"`
	CategoryID string  `json:"category_id,omitempty" jsonschema:"Category ID"`
	Memo       string  `json:"memo,omitempty" jsonschema:"Transaction memo"`
	Cleared    string  `json:"cleared,omitempty" jsonschema:"Cleared status: cleared or uncleared"`
}

type updateTransactionInput struct {
	BudgetID      string   `json:"budget_id" jsonschema:"Budget ID or last-used"`
	TransactionID string   `json:"transaction_id" jsonschema:"Transaction ID to update"`
	Amount        *float64 `json:"amount,omitempty" jsonschema:"New amount in dollars (negative for outflow)"`
	Date          *string  `json:"date,omitempty" jsonschema:"New date (YYYY-MM-DD)"`
	PayeeName     *string  `json:"payee_name,omitempty" jsonschema:"New payee name"`
	CategoryID    *string  `json:"category_id,omitempty" jsonschema:"New category ID"`
	Memo          *string  `json:"memo,omitempty" jsonschema:"New memo"`
	Cleared       *string  `json:"cleared,omitempty" jsonschema:"New cleared status: cleared or uncleared"`
	Approved      *bool    `json:"approved,omitempty" jsonschema:"Whether to approve the transaction"`
	FlagColor     json.RawMessage `json:"flag_color,omitempty" jsonschema:"Flag color: red, orange, yellow, green, blue, purple — or null to clear"`
}

type deleteTransactionInput struct {
	BudgetID      string `json:"budget_id" jsonschema:"Budget ID or last-used"`
	TransactionID string `json:"transaction_id" jsonschema:"Transaction ID to delete"`
}

func renderTransaction(t *ynab.Transaction) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "ID: %s\n", t.ID)
	fmt.Fprintf(&sb, "Date: %s\n", t.Date)
	fmt.Fprintf(&sb, "Amount: %s\n", FormatMilliunits(t.Amount))
	fmt.Fprintf(&sb, "Payee: %s\n", t.PayeeName)
	fmt.Fprintf(&sb, "Category: %s\n", t.CategoryName)
	fmt.Fprintf(&sb, "Account: %s\n", t.AccountName)
	fmt.Fprintf(&sb, "Memo: %s\n", t.Memo)
	fmt.Fprintf(&sb, "Cleared: %s\n", t.Cleared)
	fmt.Fprintf(&sb, "Approved: %v\n", t.Approved)
	fmt.Fprintf(&sb, "Flag: %s\n", t.FlagColor)
	if len(t.Subtransactions) > 0 {
		sb.WriteString("\nSplit transactions:\n")
		for _, s := range t.Subtransactions {
			fmt.Fprintf(&sb, "  • %s | %s", FormatMilliunits(s.Amount), s.CategoryName)
			if s.Memo != "" {
				fmt.Fprintf(&sb, " | %s", s.Memo)
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func registerTransactionTools(server *mcp.Server, client *ynab.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_transactions",
		Description: "List transactions in a YNAB budget. Optionally filter by date, type, account, category, or payee (only one entity filter at a time).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input listTransactionsInput) (*mcp.CallToolResult, struct{}, error) {
		entityFilters := 0
		if input.AccountID != "" {
			entityFilters++
		}
		if input.CategoryID != "" {
			entityFilters++
		}
		if input.PayeeID != "" {
			entityFilters++
		}
		if entityFilters > 1 {
			return errResult(fmt.Errorf("only one of account_id, category_id, or payee_id may be set")), struct{}{}, nil
		}

		txns, err := client.ListTransactions(input.BudgetID, ynab.TransactionFilter{
			SinceDate:  input.SinceDate,
			Type:       input.Type,
			AccountID:  input.AccountID,
			CategoryID: input.CategoryID,
			PayeeID:    input.PayeeID,
		})
		if err != nil {
			return errResult(err), struct{}{}, nil
		}

		var sb, body strings.Builder
		count := 0
		for _, t := range txns {
			if t.Deleted {
				continue
			}
			count++
			payee := t.PayeeName
			if payee == "" {
				payee = "(no payee)"
			}
			fmt.Fprintf(&body, "• %s | %s | %s", t.Date, payee, FormatMilliunits(t.Amount))
			if t.CategoryName != "" {
				fmt.Fprintf(&body, " | %s", t.CategoryName)
			}
			if t.Memo != "" {
				fmt.Fprintf(&body, " | %s", t.Memo)
			}
			fmt.Fprintf(&body, "\n  Account: %s | Status: %s | ID: %s\n", t.AccountName, t.Cleared, t.ID)
		}
		fmt.Fprintf(&sb, "Found %d transaction(s):\n\n", count)
		sb.WriteString(body.String())
		return textResult(sb.String()), struct{}{}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_transaction",
		Description: "Get details for a specific transaction",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input getTransactionInput) (*mcp.CallToolResult, struct{}, error) {
		t, err := client.GetTransaction(input.BudgetID, input.TransactionID)
		if err != nil {
			return errResult(err), struct{}{}, nil
		}
		return textResult(renderTransaction(t)), struct{}{}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_transaction",
		Description: "Create a new transaction in a YNAB budget",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input createTransactionInput) (*mcp.CallToolResult, struct{}, error) {
		txn := ynab.SaveTransaction{
			AccountID:  input.AccountID,
			Date:       input.Date,
			Amount:     MilliunitsFromDollars(input.Amount),
			PayeeName:  input.PayeeName,
			CategoryID: input.CategoryID,
			Memo:       input.Memo,
			Cleared:    input.Cleared,
			Approved:   true,
		}

		created, err := client.CreateTransaction(input.BudgetID, txn)
		if err != nil {
			return errResult(err), struct{}{}, nil
		}
		return textResult("Transaction created successfully!\n" + renderTransaction(created)), struct{}{}, nil
	})

	updateSchema, err := jsonschema.For[updateTransactionInput](&jsonschema.ForOptions{})
	if err != nil {
		panic(fmt.Errorf("update_transaction schema: %w", err))
	}
	if updateSchema.Properties != nil {
		updateSchema.Properties["flag_color"] = &jsonschema.Schema{
			Types:       []string{"string", "null"},
			Description: "Flag color: red, orange, yellow, green, blue, purple — or null to clear",
		}
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_transaction",
		Description: "Update an existing transaction in a YNAB budget",
		InputSchema: updateSchema,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input updateTransactionInput) (*mcp.CallToolResult, struct{}, error) {
		flagColor, err := normalizeFlagColor(input.FlagColor)
		if err != nil {
			return errResult(err), struct{}{}, nil
		}
		update := ynab.UpdateTransaction{
			Date:       input.Date,
			PayeeName:  input.PayeeName,
			CategoryID: input.CategoryID,
			Memo:       input.Memo,
			Cleared:    input.Cleared,
			Approved:   input.Approved,
			FlagColor:  flagColor,
		}
		if input.Amount != nil {
			m := MilliunitsFromDollars(*input.Amount)
			update.Amount = &m
		}

		updated, err := client.UpdateTransaction(input.BudgetID, input.TransactionID, update)
		if err != nil {
			return errResult(err), struct{}{}, nil
		}
		return textResult("Transaction updated successfully!\n" + renderTransaction(updated)), struct{}{}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_transaction",
		Description: "Delete a transaction from a YNAB budget",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input deleteTransactionInput) (*mcp.CallToolResult, struct{}, error) {
		if err := client.DeleteTransaction(input.BudgetID, input.TransactionID); err != nil {
			return errResult(err), struct{}{}, nil
		}
		return textResult(fmt.Sprintf("Transaction %s deleted.", input.TransactionID)), struct{}{}, nil
	})
}
