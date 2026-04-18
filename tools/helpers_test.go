package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JadoJodo/ynab-mcp/ynab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// setupToolServer creates an httptest server with the given handler, registers all
// MCP tools, and returns a ClientSession connected via in-memory transport.
func setupToolServer(t *testing.T, handler http.HandlerFunc) *mcp.ClientSession {
	t.Helper()

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	client := ynab.NewTestClient("test-token", ts.URL)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	RegisterAll(server, client)

	ct, st := mcp.NewInMemoryTransports()
	_, err := server.Connect(context.Background(), st, nil)
	if err != nil {
		t.Fatal(err)
	}

	mcpClient := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.1.0"}, nil)
	cs, err := mcpClient.Connect(context.Background(), ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

// callTool is a helper that calls a tool and returns the text content.
func callTool(t *testing.T, cs *mcp.ClientSession, name string, args map[string]any) string {
	t.Helper()
	result, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		t.Fatalf("CallTool(%s) error: %v", name, err)
	}
	if result.IsError {
		var texts []string
		for _, c := range result.Content {
			if tc, ok := c.(*mcp.TextContent); ok {
				texts = append(texts, tc.Text)
			}
		}
		t.Fatalf("CallTool(%s) returned error: %s", name, strings.Join(texts, "; "))
	}
	if len(result.Content) == 0 {
		t.Fatalf("CallTool(%s) returned no content", name)
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("CallTool(%s) content is %T, want *mcp.TextContent", name, result.Content[0])
	}
	return tc.Text
}

func TestTextResult(t *testing.T) {
	r := textResult("hello")
	if len(r.Content) != 1 {
		t.Fatalf("got %d content, want 1", len(r.Content))
	}
	tc, ok := r.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("content is %T, want *mcp.TextContent", r.Content[0])
	}
	if tc.Text != "hello" {
		t.Errorf("text = %q, want %q", tc.Text, "hello")
	}
}

func TestErrResult(t *testing.T) {
	r := errResult(json.Unmarshal([]byte("bad"), new(any)))
	if !r.IsError {
		t.Error("expected IsError = true")
	}
}
