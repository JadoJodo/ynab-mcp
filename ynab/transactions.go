package ynab

import (
	"fmt"
	"net/http"
	"net/url"
)

type transactionsData struct {
	Transactions []Transaction `json:"transactions"`
}

type transactionData struct {
	Transaction Transaction `json:"transaction"`
}

// TransactionFilter specifies optional filters for listing transactions.
type TransactionFilter struct {
	SinceDate  string // YYYY-MM-DD
	Type       string // "uncategorized", "unapproved"
	AccountID  string
	CategoryID string
	PayeeID    string
}

// ListTransactions returns transactions for a budget, optionally filtered.
func (c *Client) ListTransactions(budgetID string, filter TransactionFilter) ([]Transaction, error) {
	var basePath string
	switch {
	case filter.AccountID != "":
		basePath = fmt.Sprintf("/budgets/%s/accounts/%s/transactions", budgetID, filter.AccountID)
	case filter.CategoryID != "":
		basePath = fmt.Sprintf("/budgets/%s/categories/%s/transactions", budgetID, filter.CategoryID)
	case filter.PayeeID != "":
		basePath = fmt.Sprintf("/budgets/%s/payees/%s/transactions", budgetID, filter.PayeeID)
	default:
		basePath = fmt.Sprintf("/budgets/%s/transactions", budgetID)
	}

	params := url.Values{}
	if filter.SinceDate != "" {
		params.Set("since_date", filter.SinceDate)
	}
	if filter.Type != "" {
		params.Set("type", filter.Type)
	}

	path := basePath
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	data, err := doGet[transactionsData](c, path)
	if err != nil {
		return nil, err
	}
	return data.Transactions, nil
}

// GetTransaction returns a single transaction by ID.
func (c *Client) GetTransaction(budgetID, transactionID string) (*Transaction, error) {
	data, err := doGet[transactionData](c, fmt.Sprintf("/budgets/%s/transactions/%s", budgetID, transactionID))
	if err != nil {
		return nil, err
	}
	return &data.Transaction, nil
}

// CreateTransaction creates a new transaction.
func (c *Client) CreateTransaction(budgetID string, txn SaveTransaction) (*Transaction, error) {
	body := map[string]any{"transaction": txn}
	data, err := doPost[transactionData](c, fmt.Sprintf("/budgets/%s/transactions", budgetID), body)
	if err != nil {
		return nil, err
	}
	return &data.Transaction, nil
}

// UpdateTransaction updates an existing transaction.
func (c *Client) UpdateTransaction(budgetID, transactionID string, txn UpdateTransaction) (*Transaction, error) {
	body := map[string]any{"transaction": txn}
	data, err := doPut[transactionData](c, fmt.Sprintf("/budgets/%s/transactions/%s", budgetID, transactionID), body)
	if err != nil {
		return nil, err
	}
	return &data.Transaction, nil
}

// DeleteTransaction deletes a transaction by ID.
func (c *Client) DeleteTransaction(budgetID, transactionID string) error {
	_, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/budgets/%s/transactions/%s", budgetID, transactionID), nil)
	return err
}
