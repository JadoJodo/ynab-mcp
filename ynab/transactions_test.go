package ynab

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestListTransactions_NoFilter(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1/transactions" {
			t.Errorf("path = %s, want /budgets/b1/transactions", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transactions": []map[string]any{
					{"id": "t1", "date": "2024-01-15", "amount": -12500},
				},
			},
		})
	})
	txns, err := c.ListTransactions("b1", TransactionFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(txns) != 1 {
		t.Fatalf("got %d transactions, want 1", len(txns))
	}
}

func TestListTransactions_WithAccountID(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1/accounts/a1/transactions" {
			t.Errorf("path = %s, want /budgets/b1/accounts/a1/transactions", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"transactions": []any{}},
		})
	})
	_, err := c.ListTransactions("b1", TransactionFilter{AccountID: "a1"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListTransactions_WithCategoryID(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1/categories/c1/transactions" {
			t.Errorf("path = %s, want /budgets/b1/categories/c1/transactions", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"transactions": []any{}},
		})
	})
	_, err := c.ListTransactions("b1", TransactionFilter{CategoryID: "c1"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListTransactions_WithPayeeID(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1/payees/p1/transactions" {
			t.Errorf("path = %s, want /budgets/b1/payees/p1/transactions", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"transactions": []any{}},
		})
	})
	_, err := c.ListTransactions("b1", TransactionFilter{PayeeID: "p1"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListTransactions_WithQueryParams(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("since_date"); got != "2024-01-01" {
			t.Errorf("since_date = %q, want %q", got, "2024-01-01")
		}
		if got := r.URL.Query().Get("type"); got != "unapproved" {
			t.Errorf("type = %q, want %q", got, "unapproved")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"transactions": []any{}},
		})
	})
	_, err := c.ListTransactions("b1", TransactionFilter{
		SinceDate: "2024-01-01",
		Type:      "unapproved",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTransaction_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/budgets/b1/transactions/t1" {
			t.Errorf("path = %s, want /budgets/b1/transactions/t1", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transaction": map[string]any{
					"id":     "t1",
					"date":   "2024-01-15",
					"amount": -12500,
				},
			},
		})
	})
	txn, err := c.GetTransaction("b1", "t1")
	if err != nil {
		t.Fatal(err)
	}
	if txn.ID != "t1" {
		t.Errorf("ID = %q, want %q", txn.ID, "t1")
	}
}

func TestCreateTransaction_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/budgets/b1/transactions" {
			t.Errorf("path = %s, want /budgets/b1/transactions", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var envelope map[string]json.RawMessage
		json.Unmarshal(body, &envelope)
		if _, ok := envelope["transaction"]; !ok {
			t.Error("body missing 'transaction' wrapper")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transaction": map[string]any{
					"id":     "t-new",
					"date":   "2024-01-20",
					"amount": -5000,
				},
			},
		})
	})
	txn, err := c.CreateTransaction("b1", SaveTransaction{
		AccountID: "a1",
		Date:      "2024-01-20",
		Amount:    -5000,
		Approved:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if txn.ID != "t-new" {
		t.Errorf("ID = %q, want %q", txn.ID, "t-new")
	}
}

func TestDeleteTransaction_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/budgets/b1/transactions/t1" {
			t.Errorf("path = %s, want /budgets/b1/transactions/t1", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	if err := c.DeleteTransaction("b1", "t1"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateTransaction_Success(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/budgets/b1/transactions/t1" {
			t.Errorf("path = %s, want /budgets/b1/transactions/t1", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var envelope map[string]json.RawMessage
		json.Unmarshal(body, &envelope)
		if _, ok := envelope["transaction"]; !ok {
			t.Error("body missing 'transaction' wrapper")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"transaction": map[string]any{
					"id":     "t1",
					"date":   "2024-01-15",
					"amount": -10000,
				},
			},
		})
	})
	memo := "updated memo"
	txn, err := c.UpdateTransaction("b1", "t1", UpdateTransaction{
		Memo: &memo,
	})
	if err != nil {
		t.Fatal(err)
	}
	if txn.ID != "t1" {
		t.Errorf("ID = %q, want %q", txn.ID, "t1")
	}
}
