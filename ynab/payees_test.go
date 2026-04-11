package ynab

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListPayees_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1/payees" {
			t.Errorf("path = %s, want /budgets/b1/payees", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"payees": []map[string]any{
					{"id": "p1", "name": "Grocery Store"},
					{"id": "p2", "name": "Electric Co"},
				},
			},
		})
	})
	payees, err := c.ListPayees("b1")
	if err != nil {
		t.Fatal(err)
	}
	if len(payees) != 2 {
		t.Fatalf("got %d payees, want 2", len(payees))
	}
	if payees[0].Name != "Grocery Store" {
		t.Errorf("payees[0].Name = %q, want %q", payees[0].Name, "Grocery Store")
	}
}
