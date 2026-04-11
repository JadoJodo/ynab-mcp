package main

import (
	"context"
	"fmt"
	"os"

	"github.com/JadoJodo/ynab-mcp/tools"
	"github.com/JadoJodo/ynab-mcp/ynab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	token := os.Getenv("YNAB_API_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "YNAB_API_TOKEN environment variable is required")
		os.Exit(1)
	}

	client := ynab.NewClient(token)

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "ynab-mcp",
			Version: "0.1.0",
		},
		nil,
	)

	tools.RegisterAll(server, client)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
