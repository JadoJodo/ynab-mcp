package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/JadoJodo/ynab-mcp/ynab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type listAccountsInput struct {
	BudgetID string `json:"budget_id" jsonschema:"Budget ID or last-used"`
}

func registerAccountTools(server *mcp.Server, client *ynab.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_accounts",
		Description: "List all accounts in a YNAB budget",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input listAccountsInput) (*mcp.CallToolResult, struct{}, error) {
		accounts, err := client.ListAccounts(input.BudgetID)
		if err != nil {
			return errResult(err), struct{}{}, nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Found %d account(s):\n\n", len(accounts))
		for _, a := range accounts {
			if a.Deleted || a.Closed {
				continue
			}
			status := "Off-budget"
			if a.OnBudget {
				status = "On-budget"
			}
			fmt.Fprintf(&sb, "• %s (%s) — Balance: %s [%s]\n  ID: %s\n",
				a.Name, a.Type, FormatMilliunits(a.Balance), status, a.ID)
		}
		return textResult(sb.String()), struct{}{}, nil
	})
}
