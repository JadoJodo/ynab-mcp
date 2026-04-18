package ynab

import "fmt"

type budgetsData struct {
	Budgets []BudgetSummary `json:"budgets"`
}

type budgetData struct {
	Budget BudgetDetail `json:"budget"`
}

// ListBudgets returns all budgets for the authenticated user.
func (c *Client) ListBudgets() ([]BudgetSummary, error) {
	data, err := doGet[budgetsData](c, "/budgets")
	if err != nil {
		return nil, err
	}
	return data.Budgets, nil
}

// GetBudget returns a single budget by ID.
func (c *Client) GetBudget(budgetID string) (*BudgetDetail, error) {
	data, err := doGet[budgetData](c, fmt.Sprintf("/budgets/%s", budgetID))
	if err != nil {
		return nil, err
	}
	return &data.Budget, nil
}
