package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListBudgetsTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"budgets": []map[string]string{
					{"id": "b1", "name": "My Budget"},
					{"id": "b2", "name": "Shared Budget"},
				},
			},
		})
	})
	text := callTool(t, cs, "list_budgets", nil)
	if !strings.Contains(text, "My Budget") {
		t.Errorf("output missing 'My Budget': %s", text)
	}
	if !strings.Contains(text, "Shared Budget") {
		t.Errorf("output missing 'Shared Budget': %s", text)
	}
	if !strings.Contains(text, "2 budget(s)") {
		t.Errorf("output missing count: %s", text)
	}
}

func TestGetBudgetTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"budget": map[string]any{
					"id":   "b1",
					"name": "My Budget",
					"accounts": []map[string]string{
						{"id": "a1", "name": "Checking"},
						{"id": "a2", "name": "Savings"},
					},
					"category_groups": []map[string]any{
						{"id": "g1", "name": "Bills", "categories": []any{}},
					},
					"payees": []map[string]string{
						{"id": "p1", "name": "Store"},
					},
				},
			},
		})
	})
	text := callTool(t, cs, "get_budget", map[string]any{"budget_id": "b1"})
	if !strings.Contains(text, "My Budget") {
		t.Errorf("output missing budget name: %s", text)
	}
	if !strings.Contains(text, "Accounts: 2") {
		t.Errorf("output missing account count: %s", text)
	}
	if !strings.Contains(text, "Category groups: 1") {
		t.Errorf("output missing category group count: %s", text)
	}
	if !strings.Contains(text, "Payees: 1") {
		t.Errorf("output missing payee count: %s", text)
	}
}

func TestGetBudgetTool_ExcludesDeleted(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"budget": map[string]any{
					"id":   "b1",
					"name": "My Budget",
					"accounts": []map[string]any{
						{"id": "a1", "name": "Checking"},
						{"id": "a2", "name": "Old", "deleted": true},
					},
					"category_groups": []map[string]any{
						{"id": "g1", "name": "Bills", "categories": []any{}},
						{"id": "g2", "name": "Stale", "deleted": true, "categories": []any{}},
					},
					"payees": []map[string]any{
						{"id": "p1", "name": "Store"},
						{"id": "p2", "name": "Gone", "deleted": true},
						{"id": "p3", "name": "Shop"},
					},
				},
			},
		})
	})
	text := callTool(t, cs, "get_budget", map[string]any{"budget_id": "b1"})
	if !strings.Contains(text, "Accounts: 1") {
		t.Errorf("output missing active account count: %s", text)
	}
	if !strings.Contains(text, "Category groups: 1") {
		t.Errorf("output missing active category group count: %s", text)
	}
	if !strings.Contains(text, "Payees: 2") {
		t.Errorf("output missing active payee count: %s", text)
	}
}
