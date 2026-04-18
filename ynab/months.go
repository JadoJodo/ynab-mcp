package ynab

import "fmt"

type monthData struct {
	Month Month `json:"month"`
}

// GetMonth returns budget data for a specific month.
func (c *Client) GetMonth(budgetID, month string) (*Month, error) {
	data, err := doGet[monthData](c, fmt.Sprintf("/budgets/%s/months/%s", budgetID, month))
	if err != nil {
		return nil, err
	}
	return &data.Month, nil
}
