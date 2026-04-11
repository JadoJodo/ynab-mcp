package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
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
	if !strings.Contains(text, "Transaction: t1") {
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
