package ynab

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListAccounts_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1/accounts" {
			t.Errorf("path = %s, want /budgets/b1/accounts", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"accounts": []map[string]any{
					{"id": "a1", "name": "Checking", "type": "checking", "balance": 50000},
					{"id": "a2", "name": "Savings", "type": "savings", "balance": 100000},
				},
			},
		})
	})
	accounts, err := c.ListAccounts("b1")
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 2 {
		t.Fatalf("got %d accounts, want 2", len(accounts))
	}
	if accounts[0].Name != "Checking" {
		t.Errorf("accounts[0].Name = %q, want %q", accounts[0].Name, "Checking")
	}
}

func TestListAccounts_Empty(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"accounts": []any{},
			},
		})
	})
	accounts, err := c.ListAccounts("b1")
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 0 {
		t.Errorf("got %d accounts, want 0", len(accounts))
	}
}
