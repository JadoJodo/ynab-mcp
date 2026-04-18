package tools

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListCategoriesTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"category_groups": []map[string]any{
					{
						"id":   "g1",
						"name": "Bills",
						"categories": []map[string]any{
							{"id": "c1", "name": "Rent", "budgeted": 100000, "activity": -100000, "balance": 0},
							{"id": "c2", "name": "Electric", "budgeted": 15000, "activity": -12000, "balance": 3000},
						},
					},
				},
			},
		})
	})
	text := callTool(t, cs, "list_categories", map[string]any{"budget_id": "b1"})
	if !strings.Contains(text, "## Bills") {
		t.Errorf("output missing group header: %s", text)
	}
	if !strings.Contains(text, "Rent") {
		t.Errorf("output missing 'Rent': %s", text)
	}
	if !strings.Contains(text, "$100.00") {
		t.Errorf("output missing budgeted amount: %s", text)
	}
}

func TestListCategoriesTool_FiltersHiddenDeleted(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"category_groups": []map[string]any{
					{
						"id":   "g1",
						"name": "Active Group",
						"categories": []map[string]any{
							{"id": "c1", "name": "Visible", "budgeted": 10000},
							{"id": "c2", "name": "HiddenCat", "budgeted": 5000, "hidden": true},
							{"id": "c3", "name": "DeletedCat", "budgeted": 0, "deleted": true},
						},
					},
					{
						"id":         "g2",
						"name":       "Hidden Group",
						"hidden":     true,
						"categories": []map[string]any{},
					},
					{
						"id":         "g3",
						"name":       "Deleted Group",
						"deleted":    true,
						"categories": []map[string]any{},
					},
				},
			},
		})
	})
	text := callTool(t, cs, "list_categories", map[string]any{"budget_id": "b1"})
	if !strings.Contains(text, "Visible") {
		t.Errorf("output missing 'Visible': %s", text)
	}
	if strings.Contains(text, "HiddenCat") {
		t.Errorf("output should not contain 'HiddenCat': %s", text)
	}
	if strings.Contains(text, "DeletedCat") {
		t.Errorf("output should not contain 'DeletedCat': %s", text)
	}
	if strings.Contains(text, "Hidden Group") {
		t.Errorf("output should not contain 'Hidden Group': %s", text)
	}
	if strings.Contains(text, "Deleted Group") {
		t.Errorf("output should not contain 'Deleted Group': %s", text)
	}
}

func TestGetCategoryTool_Success(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"category": map[string]any{
					"id":                        "c1",
					"name":                      "Rent",
					"category_group_name":       "Bills",
					"budgeted":                  100000,
					"activity":                  -100000,
					"balance":                   0,
					"goal_type":                 "TB",
					"goal_target":               120000,
					"goal_percentage_complete":   83,
				},
			},
		})
	})
	text := callTool(t, cs, "get_category", map[string]any{"budget_id": "b1", "category_id": "c1"})
	if !strings.Contains(text, "Rent") {
		t.Errorf("output missing category name: %s", text)
	}
	if !strings.Contains(text, "Bills") {
		t.Errorf("output missing group name: %s", text)
	}
	if !strings.Contains(text, "Goal type: TB") {
		t.Errorf("output missing goal type: %s", text)
	}
	if !strings.Contains(text, "$120.00") {
		t.Errorf("output missing goal target: %s", text)
	}
	if !strings.Contains(text, "83%") {
		t.Errorf("output missing goal progress: %s", text)
	}
}

func TestGetCategoryTool_NoGoal(t *testing.T) {
	cs := setupToolServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"category": map[string]any{
					"id":       "c1",
					"name":     "Groceries",
					"budgeted": 50000,
					"activity": -30000,
					"balance":  20000,
				},
			},
		})
	})
	text := callTool(t, cs, "get_category", map[string]any{"budget_id": "b1", "category_id": "c1"})
	if strings.Contains(text, "Goal") {
		t.Errorf("output should not contain goal info: %s", text)
	}
}
