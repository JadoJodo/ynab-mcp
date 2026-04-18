package ynab

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListCategories_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1/categories" {
			t.Errorf("path = %s, want /budgets/b1/categories", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"category_groups": []map[string]any{
					{
						"id":   "g1",
						"name": "Bills",
						"categories": []map[string]any{
							{"id": "c1", "name": "Rent", "budgeted": 100000},
						},
					},
				},
			},
		})
	})
	groups, err := c.ListCategories("b1")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(groups))
	}
	if groups[0].Name != "Bills" {
		t.Errorf("group name = %q, want %q", groups[0].Name, "Bills")
	}
	if len(groups[0].Categories) != 1 {
		t.Fatalf("got %d categories, want 1", len(groups[0].Categories))
	}
}

func TestGetCategory_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1/categories/c1" {
			t.Errorf("path = %s, want /budgets/b1/categories/c1", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"category": map[string]any{
					"id":       "c1",
					"name":     "Rent",
					"budgeted": 100000,
					"activity": -50000,
					"balance":  50000,
				},
			},
		})
	})
	cat, err := c.GetCategory("b1", "c1")
	if err != nil {
		t.Fatal(err)
	}
	if cat.Name != "Rent" {
		t.Errorf("category name = %q, want %q", cat.Name, "Rent")
	}
	if cat.Budgeted != 100000 {
		t.Errorf("budgeted = %d, want 100000", cat.Budgeted)
	}
}
