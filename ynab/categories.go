package ynab

import "fmt"

type categoryGroupsData struct {
	CategoryGroups []CategoryGroup `json:"category_groups"`
}

type categoryData struct {
	Category Category `json:"category"`
}

// ListCategories returns all category groups and their categories for a budget.
func (c *Client) ListCategories(budgetID string) ([]CategoryGroup, error) {
	data, err := doGet[categoryGroupsData](c, fmt.Sprintf("/budgets/%s/categories", budgetID))
	if err != nil {
		return nil, err
	}
	return data.CategoryGroups, nil
}

// GetCategory returns a single category by ID.
func (c *Client) GetCategory(budgetID, categoryID string) (*Category, error) {
	data, err := doGet[categoryData](c, fmt.Sprintf("/budgets/%s/categories/%s", budgetID, categoryID))
	if err != nil {
		return nil, err
	}
	return &data.Category, nil
}
