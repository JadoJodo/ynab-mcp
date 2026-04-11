package ynab

import "fmt"

type payeesData struct {
	Payees []Payee `json:"payees"`
}

// ListPayees returns all payees for a budget.
func (c *Client) ListPayees(budgetID string) ([]Payee, error) {
	data, err := doGet[payeesData](c, fmt.Sprintf("/budgets/%s/payees", budgetID))
	if err != nil {
		return nil, err
	}
	return data.Payees, nil
}
