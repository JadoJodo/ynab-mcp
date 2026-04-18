package ynab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testClient creates a Client pointed at an httptest server running the given handler.
func testClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return NewTestClient("test-token", ts.URL)
}

func TestDoRequest_AuthHeader(t *testing.T) {
	var gotAuth string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(200)
	})
	_, _ = c.doRequest(http.MethodGet, "/test", nil)
	if gotAuth != "Bearer test-token" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer test-token")
	}
}

func TestDoRequest_ContentTypeOnBody(t *testing.T) {
	var gotCT string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotCT = r.Header.Get("Content-Type")
		w.WriteHeader(200)
	})

	// No body → no Content-Type
	_, _ = c.doRequest(http.MethodGet, "/test", nil)
	if gotCT != "" {
		t.Errorf("Content-Type without body = %q, want empty", gotCT)
	}

	// With body → Content-Type set
	_, _ = c.doRequest(http.MethodPost, "/test", map[string]string{"key": "val"})
	if gotCT != "application/json" {
		t.Errorf("Content-Type with body = %q, want %q", gotCT, "application/json")
	}
}

func TestDoRequest_APIError(t *testing.T) {
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
	_, err := c.doRequest(http.MethodGet, "/budgets/bad", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T: %v", err, err)
	}
	if apiErr.Detail != "Budget not found" {
		t.Errorf("error detail = %q, want %q", apiErr.Detail, "Budget not found")
	}
}

func TestDoRequest_APIError_NonJSON(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("internal server error"))
	})
	_, err := c.doRequest(http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*Error); ok {
		t.Fatal("expected generic error, got *Error")
	}
}

func TestDoGet_Success(t *testing.T) {
	type result struct {
		Value string `json:"value"`
	}
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"value": "hello"},
		})
	})
	got, err := doGet[result](c, "/test")
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != "hello" {
		t.Errorf("got %q, want %q", got.Value, "hello")
	}
}

func TestDoGet_InvalidJSON(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	})
	type result struct{}
	_, err := doGet[result](c, "/test")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDoPost_Success(t *testing.T) {
	type result struct {
		ID string `json:"id"`
	}
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"id": "abc"},
		})
	})
	got, err := doPost[result](c, "/test", map[string]string{"key": "val"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "abc" {
		t.Errorf("got %q, want %q", got.ID, "abc")
	}
}

func TestDoPut_Success(t *testing.T) {
	type result struct {
		ID string `json:"id"`
	}
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"id": "xyz"},
		})
	})
	got, err := doPut[result](c, "/test", map[string]string{"key": "val"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "xyz" {
		t.Errorf("got %q, want %q", got.ID, "xyz")
	}
}
