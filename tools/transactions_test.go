package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestListTransactionsTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transactions": []map[string]any{
					{
						"id":            "t1",
						"date":          "2024-01-15",
						"amount":        -12500,
						"payee_name":    "Grocery Store",
						"category_name": "Groceries",
						"account_name":  "Checking",
						"cleared":       "cleared",
						"memo":          "Weekly shopping",
					},
				},
			},
		})
	})
	text := callTool(t, cs, "list_transactions", map[string]any{"budget_id": "b1"})
	if !strings.Contains(text, "2024-01-15") {
		t.Errorf("output missing date: %s", text)
	}
	if !strings.Contains(text, "Grocery Store") {
		t.Errorf("output missing payee: %s", text)
	}
	if !strings.Contains(text, "-$12.50") {
		t.Errorf("output missing amount: %s", text)
	}
	if !strings.Contains(text, "Groceries") {
		t.Errorf("output missing category: %s", text)
	}
	if !strings.Contains(text, "Weekly shopping") {
		t.Errorf("output missing memo: %s", text)
	}
}

func TestGetTransactionTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transaction": map[string]any{
					"id":            "t1",
					"date":          "2024-01-15",
					"amount":        -12500,
					"payee_name":    "Grocery Store",
					"category_name": "Groceries",
					"account_name":  "Checking",
					"cleared":       "cleared",
					"approved":      true,
					"memo":          "Weekly shopping",
				},
			},
		})
	})
	text := callTool(t, cs, "get_transaction", map[string]any{"budget_id": "b1", "transaction_id": "t1"})
	if !strings.Contains(text, "ID: t1") {
		t.Errorf("output missing transaction ID: %s", text)
	}
	if !strings.Contains(text, "Approved: true") {
		t.Errorf("output missing approved: %s", text)
	}
}

func TestGetTransactionTool_WithSubtransactions(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transaction": map[string]any{
					"id":            "t1",
					"date":          "2024-01-15",
					"amount":        -25000,
					"payee_name":    "Store",
					"account_name":  "Checking",
					"cleared":       "cleared",
					"approved":      true,
					"subtransactions": []map[string]any{
						{"id": "st1", "amount": -15000, "category_name": "Groceries", "memo": "Food"},
						{"id": "st2", "amount": -10000, "category_name": "Household", "memo": "Supplies"},
					},
				},
			},
		})
	})
	text := callTool(t, cs, "get_transaction", map[string]any{"budget_id": "b1", "transaction_id": "t1"})
	if !strings.Contains(text, "Split transactions") {
		t.Errorf("output missing split header: %s", text)
	}
	if !strings.Contains(text, "-$15.00") {
		t.Errorf("output missing split amount: %s", text)
	}
	if !strings.Contains(text, "Groceries") {
		t.Errorf("output missing split category: %s", text)
	}
	if !strings.Contains(text, "Food") {
		t.Errorf("output missing split memo: %s", text)
	}
}

func TestCreateTransactionTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transaction": map[string]any{
					"id":            "t-new",
					"date":          "2024-01-20",
					"amount":        -12500,
					"payee_name":    "Coffee Shop",
					"category_name": "Dining Out",
					"account_name":  "Checking",
				},
			},
		})
	})
	text := callTool(t, cs, "create_transaction", map[string]any{
		"budget_id":  "b1",
		"account_id": "a1",
		"date":       "2024-01-20",
		"amount":     -12.50,
		"payee_name": "Coffee Shop",
	})
	if !strings.Contains(text, "created successfully") {
		t.Errorf("output missing success message: %s", text)
	}
	if !strings.Contains(text, "t-new") {
		t.Errorf("output missing transaction ID: %s", text)
	}
	if !strings.Contains(text, "-$12.50") {
		t.Errorf("output missing amount: %s", text)
	}
}

func TestUpdateTransactionTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transaction": map[string]any{
					"id":            "t1",
					"date":          "2024-01-15",
					"amount":        -15000,
					"payee_name":    "Updated Payee",
					"category_name": "Groceries",
					"account_name":  "Checking",
				},
			},
		})
	})
	text := callTool(t, cs, "update_transaction", map[string]any{
		"budget_id":      "b1",
		"transaction_id": "t1",
		"amount":         -15.00,
	})
	if !strings.Contains(text, "updated successfully") {
		t.Errorf("output missing success message: %s", text)
	}
	if !strings.Contains(text, "-$15.00") {
		t.Errorf("output missing amount: %s", text)
	}
}

func TestListTransactionsTool_RejectsMultipleEntityFilters(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("server should not be called; got %s %s", r.Method, r.URL.Path)
	})
	result, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "list_transactions",
		Arguments: map[string]any{
			"budget_id":   "b1",
			"account_id":  "a1",
			"category_id": "c1",
		},
	})
	if err != nil {
		t.Fatalf("CallTool error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected IsError = true for conflicting filters")
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("content is %T, want *mcp.TextContent", result.Content[0])
	}
	if !strings.Contains(tc.Text, "only one of") {
		t.Errorf("error text = %q, want message mentioning 'only one of'", tc.Text)
	}
}

func TestCreateTransactionTool_RendersFullFields(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transaction": map[string]any{
					"id":            "t-new",
					"date":          "2024-01-20",
					"amount":        -12500,
					"payee_name":    "Coffee Shop",
					"category_name": "Dining Out",
					"account_name":  "Checking",
					"cleared":       "uncleared",
					"approved":      true,
					"memo":          "",
				},
			},
		})
	})
	text := callTool(t, cs, "create_transaction", map[string]any{
		"budget_id":  "b1",
		"account_id": "a1",
		"date":       "2024-01-20",
		"amount":     -12.50,
	})
	for _, want := range []string{"Memo:", "Cleared:", "Approved:", "Flag:"} {
		if !strings.Contains(text, want) {
			t.Errorf("output missing %q: %s", want, text)
		}
	}
}

func TestUpdateTransactionTool_ClearsFlag(t *testing.T) {
	var capturedBody []byte
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transaction": map[string]any{
					"id":           "t1",
					"date":         "2024-01-15",
					"amount":       -10000,
					"account_name": "Checking",
				},
			},
		})
	})
	_ = callTool(t, cs, "update_transaction", map[string]any{
		"budget_id":      "b1",
		"transaction_id": "t1",
		"flag_color":     nil,
	})
	var envelope struct {
		Transaction map[string]json.RawMessage `json:"transaction"`
	}
	if err := json.Unmarshal(capturedBody, &envelope); err != nil {
		t.Fatalf("unmarshal body: %v (body: %s)", err, capturedBody)
	}
	flag, ok := envelope.Transaction["flag_color"]
	if !ok {
		t.Fatalf("flag_color missing from body: %s", capturedBody)
	}
	if string(flag) != "null" {
		t.Errorf("flag_color = %s, want null", flag)
	}
}

func TestUpdateTransactionTool_RejectsInvalidFlagColor(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("server should not be called; got %s %s", r.Method, r.URL.Path)
	})
	result, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "update_transaction",
		Arguments: map[string]any{
			"budget_id":      "b1",
			"transaction_id": "t1",
			"flag_color":     "magenta",
		},
	})
	if err != nil {
		t.Fatalf("CallTool error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected IsError = true for invalid flag_color")
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("content is %T, want *mcp.TextContent", result.Content[0])
	}
	if !strings.Contains(tc.Text, "flag_color") {
		t.Errorf("error text = %q, want message mentioning 'flag_color'", tc.Text)
	}
}

func TestUpdateTransactionTool_RendersFullFields(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transaction": map[string]any{
					"id":            "t1",
					"date":          "2024-01-15",
					"amount":        -15000,
					"payee_name":    "Updated Payee",
					"category_name": "Groceries",
					"account_name":  "Checking",
					"cleared":       "cleared",
					"approved":      true,
					"memo":          "verify",
					"flag_color":    "red",
				},
			},
		})
	})
	text := callTool(t, cs, "update_transaction", map[string]any{
		"budget_id":      "b1",
		"transaction_id": "t1",
		"flag_color":     "red",
		"memo":           "verify",
	})
	for _, want := range []string{"Memo: verify", "Cleared: cleared", "Approved: true", "Flag: red"} {
		if !strings.Contains(text, want) {
			t.Errorf("output missing %q: %s", want, text)
		}
	}
}

func TestDeleteTransactionTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/budgets/b1/transactions/t1" {
			t.Errorf("path = %s, want /budgets/b1/transactions/t1", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"transaction": map[string]any{"id": "t1"}},
		})
	})
	text := callTool(t, cs, "delete_transaction", map[string]any{
		"budget_id":      "b1",
		"transaction_id": "t1",
	})
	if !strings.Contains(text, "Transaction t1 deleted") {
		t.Errorf("output missing success message: %s", text)
	}
}
