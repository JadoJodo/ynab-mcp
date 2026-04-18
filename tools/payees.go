package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/JadoJodo/ynab-mcp/ynab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type listPayeesInput struct {
	BudgetID string `json:"budget_id" jsonschema:"Budget ID or last-used"`
}

func registerPayeeTools(server *mcp.Server, client *ynab.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_payees",
		Description: "List all payees in a YNAB budget",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input listPayeesInput) (*mcp.CallToolResult, struct{}, error) {
		payees, err := client.ListPayees(input.BudgetID)
		if err != nil {
			return errResult(err), struct{}{}, nil
		}

		var sb, body strings.Builder
		count := 0
		for _, p := range payees {
			if p.Deleted {
				continue
			}
			count++
			fmt.Fprintf(&body, "• %s (ID: %s)\n", p.Name, p.ID)
		}
		fmt.Fprintf(&sb, "Found %d payee(s):\n\n", count)
		sb.WriteString(body.String())
		return textResult(sb.String()), struct{}{}, nil
	})
}
