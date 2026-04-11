package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListPayeesTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"payees": []map[string]any{
					{"id": "p1", "name": "Grocery Store"},
					{"id": "p2", "name": "Electric Co"},
				},
			},
		})
	})
	text := callTool(t, cs, "list_payees", map[string]any{"budget_id": "b1"})
	if !strings.Contains(text, "Grocery Store") {
		t.Errorf("output missing 'Grocery Store': %s", text)
	}
	if !strings.Contains(text, "Electric Co") {
		t.Errorf("output missing 'Electric Co': %s", text)
	}
	if !strings.Contains(text, "2 payee(s)") {
		t.Errorf("output missing count: %s", text)
	}
}

func TestListPayeesTool_FiltersDeleted(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"payees": []map[string]any{
					{"id": "p1", "name": "Active Payee"},
					{"id": "p2", "name": "Deleted Payee", "deleted": true},
				},
			},
		})
	})
	text := callTool(t, cs, "list_payees", map[string]any{"budget_id": "b1"})
	if !strings.Contains(text, "Active Payee") {
		t.Errorf("output missing 'Active Payee': %s", text)
	}
	if strings.Contains(text, "Deleted Payee") {
		t.Errorf("output should not contain 'Deleted Payee': %s", text)
	}
}
