package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestGetMonthTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"month": map[string]any{
					"month":          "2024-01-01",
					"income":         500000,
					"budgeted":       400000,
					"activity":       -350000,
					"to_be_budgeted": 100000,
					"categories": []map[string]any{
						{"id": "c1", "name": "Rent", "budgeted": 100000, "activity": -100000, "balance": 0},
						{"id": "c2", "name": "Groceries", "budgeted": 50000, "activity": -30000, "balance": 20000},
					},
				},
			},
		})
	})
	text := callTool(t, cs, "get_month", map[string]any{"budget_id": "b1", "month": "2024-01-01"})
	if !strings.Contains(text, "Month: 2024-01-01") {
		t.Errorf("output missing month: %s", text)
	}
	if !strings.Contains(text, "Income: $500.00") {
		t.Errorf("output missing income: %s", text)
	}
	if !strings.Contains(text, "To Be Budgeted: $100.00") {
		t.Errorf("output missing TBB: %s", text)
	}
	if !strings.Contains(text, "Rent") {
		t.Errorf("output missing category: %s", text)
	}
}

func TestGetMonthTool_SkipsZeroAndHidden(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"month": map[string]any{
					"month":          "2024-01-01",
					"income":         500000,
					"budgeted":       400000,
					"activity":       -350000,
					"to_be_budgeted": 100000,
					"categories": []map[string]any{
						{"id": "c1", "name": "Active", "budgeted": 50000, "activity": -30000, "balance": 20000},
						{"id": "c2", "name": "ZeroCat", "budgeted": 0, "activity": 0, "balance": 0},
						{"id": "c3", "name": "HiddenCat", "budgeted": 10000, "activity": -5000, "balance": 5000, "hidden": true},
						{"id": "c4", "name": "DeletedCat", "budgeted": 10000, "activity": -5000, "balance": 5000, "deleted": true},
					},
				},
			},
		})
	})
	text := callTool(t, cs, "get_month", map[string]any{"budget_id": "b1", "month": "2024-01-01"})
	if !strings.Contains(text, "Active") {
		t.Errorf("output missing 'Active': %s", text)
	}
	if strings.Contains(text, "ZeroCat") {
		t.Errorf("output should not contain 'ZeroCat': %s", text)
	}
	if strings.Contains(text, "HiddenCat") {
		t.Errorf("output should not contain 'HiddenCat': %s", text)
	}
	if strings.Contains(text, "DeletedCat") {
		t.Errorf("output should not contain 'DeletedCat': %s", text)
	}
}
