package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/JadoJodo/ynab-mcp/ynab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type listCategoriesInput struct {
	BudgetID string `json:"budget_id" jsonschema:"Budget ID or last-used"`
}

type getCategoryInput struct {
	BudgetID   string `json:"budget_id" jsonschema:"Budget ID or last-used"`
	CategoryID string `json:"category_id" jsonschema:"Category ID"`
}

func registerCategoryTools(server *mcp.Server, client *ynab.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_categories",
		Description: "List all category groups and categories in a YNAB budget",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input listCategoriesInput) (*mcp.CallToolResult, struct{}, error) {
		groups, err := client.ListCategories(input.BudgetID)
		if err != nil {
			return errResult(err), struct{}{}, nil
		}

		var sb strings.Builder
		for _, g := range groups {
			if g.Deleted || g.Hidden {
				continue
			}
			fmt.Fprintf(&sb, "## %s\n", g.Name)
			for _, c := range g.Categories {
				if c.Deleted || c.Hidden {
					continue
				}
				fmt.Fprintf(&sb, "  • %s — Budgeted: %s | Activity: %s | Balance: %s\n    ID: %s\n",
					c.Name, FormatMilliunits(c.Budgeted), FormatMilliunits(c.Activity), FormatMilliunits(c.Balance), c.ID)
			}
			sb.WriteString("\n")
		}
		return textResult(sb.String()), struct{}{}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_category",
		Description: "Get details for a specific category in a YNAB budget",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input getCategoryInput) (*mcp.CallToolResult, struct{}, error) {
		cat, err := client.GetCategory(input.BudgetID, input.CategoryID)
		if err != nil {
			return errResult(err), struct{}{}, nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Category: %s\n", cat.Name)
		if cat.CategoryGroupName != "" {
			fmt.Fprintf(&sb, "Group: %s\n", cat.CategoryGroupName)
		}
		fmt.Fprintf(&sb, "Budgeted: %s\n", FormatMilliunits(cat.Budgeted))
		fmt.Fprintf(&sb, "Activity: %s\n", FormatMilliunits(cat.Activity))
		fmt.Fprintf(&sb, "Balance: %s\n", FormatMilliunits(cat.Balance))
		if cat.GoalType != "" {
			fmt.Fprintf(&sb, "Goal type: %s\n", cat.GoalType)
			if cat.GoalTarget != 0 {
				fmt.Fprintf(&sb, "Goal target: %s\n", FormatMilliunits(cat.GoalTarget))
			}
			if cat.GoalPercentageComplete > 0 {
				fmt.Fprintf(&sb, "Goal progress: %d%%\n", cat.GoalPercentageComplete)
			}
		}
		return textResult(sb.String()), struct{}{}, nil
	})
}
