package tools

import (
	"github.com/JadoJodo/ynab-mcp/ynab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all YNAB tools with the MCP server.
func RegisterAll(server *mcp.Server, client *ynab.Client) {
	registerBudgetTools(server, client)
	registerAccountTools(server, client)
	registerCategoryTools(server, client)
	registerTransactionTools(server, client)
	registerMonthTools(server, client)
	registerPayeeTools(server, client)
}

// textResult creates a CallToolResult with text content.
func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

// errResult creates a CallToolResult indicating an error.
func errResult(err error) *mcp.CallToolResult {
	r := &mcp.CallToolResult{}
	r.SetError(err)
	return r
}
