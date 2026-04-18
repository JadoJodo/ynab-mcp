package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/JadoJodo/ynab-mcp/ynab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type listBudgetsInput struct{}

type getBudgetInput struct {
	BudgetID string `json:"budget_id" jsonschema:"Budget ID or last-used"`
}

func registerBudgetTools(server *mcp.Server, client *ynab.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_budgets",
		Description: "List all YNAB budgets for the authenticated user",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input listBudgetsInput) (*mcp.CallToolResult, struct{}, error) {
		budgets, err := client.ListBudgets()
		if err != nil {
			return errResult(err), struct{}{}, nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Found %d budget(s):\n\n", len(budgets))
		for _, b := range budgets {
			fmt.Fprintf(&sb, "• %s (ID: %s)\n", b.Name, b.ID)
		}
		return textResult(sb.String()), struct{}{}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_budget",
		Description: "Get details for a specific YNAB budget",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input getBudgetInput) (*mcp.CallToolResult, struct{}, error) {
		budget, err := client.GetBudget(input.BudgetID)
		if err != nil {
			return errResult(err), struct{}{}, nil
		}

		accountCount := 0
		for _, a := range budget.Accounts {
			if !a.Deleted {
				accountCount++
			}
		}
		groupCount := 0
		for _, g := range budget.CategoryGroups {
			if !g.Deleted {
				groupCount++
			}
		}
		payeeCount := 0
		for _, p := range budget.Payees {
			if !p.Deleted {
				payeeCount++
			}
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Budget: %s\n", budget.Name)
		fmt.Fprintf(&sb, "ID: %s\n", budget.ID)
		fmt.Fprintf(&sb, "Accounts: %d\n", accountCount)
		fmt.Fprintf(&sb, "Category groups: %d\n", groupCount)
		fmt.Fprintf(&sb, "Payees: %d\n", payeeCount)
		return textResult(sb.String()), struct{}{}, nil
	})
}
