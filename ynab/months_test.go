package ynab

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetMonth_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1/months/2024-01-01" {
			t.Errorf("path = %s, want /budgets/b1/months/2024-01-01", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"month": map[string]any{
					"month":          "2024-01-01",
					"income":         500000,
					"budgeted":       400000,
					"activity":       -350000,
					"to_be_budgeted": 100000,
				},
			},
		})
	})
	month, err := c.GetMonth("b1", "2024-01-01")
	if err != nil {
		t.Fatal(err)
	}
	if month.Month != "2024-01-01" {
		t.Errorf("month = %q, want %q", month.Month, "2024-01-01")
	}
	if month.Income != 500000 {
		t.Errorf("income = %d, want 500000", month.Income)
	}
}
