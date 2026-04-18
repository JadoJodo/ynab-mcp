package ynab

import "fmt"

type accountsData struct {
	Accounts []Account `json:"accounts"`
}

// ListAccounts returns all accounts for a budget.
func (c *Client) ListAccounts(budgetID string) ([]Account, error) {
	data, err := doGet[accountsData](c, fmt.Sprintf("/budgets/%s/accounts", budgetID))
	if err != nil {
		return nil, err
	}
	return data.Accounts, nil
}
