package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListAccountsTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"accounts": []map[string]any{
					{"id": "a1", "name": "Checking", "type": "checking", "balance": 50000, "on_budget": true},
					{"id": "a2", "name": "Credit Card", "type": "creditCard", "balance": -25000, "on_budget": true},
				},
			},
		})
	})
	text := callTool(t, cs, "list_accounts", map[string]any{"budget_id": "b1"})
	if !strings.Contains(text, "Checking") {
		t.Errorf("output missing 'Checking': %s", text)
	}
	if !strings.Contains(text, "$50.00") {
		t.Errorf("output missing balance: %s", text)
	}
	if !strings.Contains(text, "On-budget") {
		t.Errorf("output missing status: %s", text)
	}
}

func TestListAccountsTool_FiltersDeletedAndClosed(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"accounts": []map[string]any{
					{"id": "a1", "name": "Active", "type": "checking", "balance": 50000, "on_budget": true},
					{"id": "a2", "name": "Deleted", "type": "checking", "balance": 0, "deleted": true},
					{"id": "a3", "name": "Closed", "type": "checking", "balance": 0, "closed": true},
				},
			},
		})
	})
	text := callTool(t, cs, "list_accounts", map[string]any{"budget_id": "b1"})
	if !strings.Contains(text, "Active") {
		t.Errorf("output missing 'Active': %s", text)
	}
	if strings.Contains(text, "Deleted") {
		t.Errorf("output should not contain 'Deleted': %s", text)
	}
	if strings.Contains(text, "Closed") {
		t.Errorf("output should not contain 'Closed': %s", text)
	}
}
