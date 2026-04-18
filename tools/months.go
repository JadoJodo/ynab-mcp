package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/JadoJodo/ynab-mcp/ynab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type getMonthInput struct {
	BudgetID string `json:"budget_id" jsonschema:"Budget ID or last-used"`
	Month    string `json:"month" jsonschema:"Month in YYYY-MM-DD format (use first of month e.g. 2024-01-01)"`
}

func registerMonthTools(server *mcp.Server, client *ynab.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_month",
		Description: "Get budget summary for a specific month including category breakdowns",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input getMonthInput) (*mcp.CallToolResult, struct{}, error) {
		month, err := client.GetMonth(input.BudgetID, input.Month)
		if err != nil {
			return errResult(err), struct{}{}, nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Month: %s\n", month.Month)
		if month.Note != "" {
			fmt.Fprintf(&sb, "Note: %s\n", month.Note)
		}
		fmt.Fprintf(&sb, "Income: %s\n", FormatMilliunits(month.Income))
		fmt.Fprintf(&sb, "Budgeted: %s\n", FormatMilliunits(month.Budgeted))
		fmt.Fprintf(&sb, "Activity: %s\n", FormatMilliunits(month.Activity))
		fmt.Fprintf(&sb, "To Be Budgeted: %s\n", FormatMilliunits(month.ToBeBudgeted))

		if len(month.Categories) > 0 {
			sb.WriteString("\nCategories:\n")
			for _, c := range month.Categories {
				if c.Deleted || c.Hidden {
					continue
				}
				if c.Budgeted == 0 && c.Activity == 0 {
					continue
				}
				fmt.Fprintf(&sb, "  • %s — Budgeted: %s | Activity: %s | Balance: %s\n",
					c.Name, FormatMilliunits(c.Budgeted), FormatMilliunits(c.Activity), FormatMilliunits(c.Balance))
			}
		}
		return textResult(sb.String()), struct{}{}, nil
	})
}
