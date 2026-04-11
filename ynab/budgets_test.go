package ynab

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListBudgets_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets" {
			t.Errorf("path = %s, want /budgets", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"budgets": []map[string]string{
					{"id": "b1", "name": "My Budget"},
					{"id": "b2", "name": "Other Budget"},
				},
			},
		})
	})
	budgets, err := c.ListBudgets()
	if err != nil {
		t.Fatal(err)
	}
	if len(budgets) != 2 {
		t.Fatalf("got %d budgets, want 2", len(budgets))
	}
	if budgets[0].Name != "My Budget" {
		t.Errorf("budgets[0].Name = %q, want %q", budgets[0].Name, "My Budget")
	}
}

func TestGetBudget_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1" {
			t.Errorf("path = %s, want /budgets/b1", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"budget": map[string]any{
					"id":   "b1",
					"name": "My Budget",
					"accounts": []map[string]string{
						{"id": "a1", "name": "Checking"},
					},
				},
			},
		})
	})
	budget, err := c.GetBudget("b1")
	if err != nil {
		t.Fatal(err)
	}
	if budget.ID != "b1" {
		t.Errorf("budget.ID = %q, want %q", budget.ID, "b1")
	}
	if len(budget.Accounts) != 1 {
		t.Errorf("got %d accounts, want 1", len(budget.Accounts))
	}
}

func TestGetBudget_APIError(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"id":     "404",
				"name":   "not_found",
				"detail": "Budget not found",
			},
		})
	})
	_, err := c.GetBudget("bad-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
