package deposits_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/surepay-one/surepay-go-sdk/deposits"
	"github.com/surepay-one/surepay-go-sdk/internal/core"
	"github.com/surepay-one/surepay-go-sdk/internal/shared"
)

func testService(srv *httptest.Server) *deposits.Service {
	return deposits.New(&core.Doer{
		APIKey:     "test-key",
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
		MaxRetries: 0,
		Log:        core.NoopLogger{},
	})
}

func TestCreate(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "dep-001", "amount": float64(100000), "status": "pending",
				"checkout_url": "https://checkout.surepay.one/dep-001",
				"created_at": "2026-01-01T00:00:00Z", "expires_at": "2026-01-01T00:15:00Z",
			},
		})
	}))
	defer srv.Close()

	dep, err := testService(srv).Create(context.Background(), &shared.CreateDepositRequest{
		Amount: 100000, RequestID: "ORD-001",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dep.ID != "dep-001" {
		t.Errorf("ID = %q, want dep-001", dep.ID)
	}
	if gotBody["request_id"] != "ORD-001" {
		t.Errorf("request_id = %v, want ORD-001", gotBody["request_id"])
	}
}

func TestList_queryParams(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"items": []any{}, "total": 0, "page": 1, "page_size": 10, "total_pages": 0},
		})
	}))
	defer srv.Close()

	testService(srv).List(context.Background(), &shared.DepositListParams{
		Page: 2, PageSize: 10, Status: "success",
	})

	for _, want := range []string{"page=2", "page_size=10", "status=success"} {
		if !strings.Contains(gotQuery, want) {
			t.Errorf("query %q missing %q", gotQuery, want)
		}
	}
}

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/deposits/dep-abc" {
			t.Errorf("path = %q, want /deposits/dep-abc", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "dep-abc", "amount": float64(50000), "status": "success",
				"created_at": "2026-01-01T00:00:00Z", "expires_at": "2026-01-01T00:15:00Z",
			},
		})
	}))
	defer srv.Close()

	dep, err := testService(srv).Get(context.Background(), "dep-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dep.ID != "dep-abc" {
		t.Errorf("ID = %q, want dep-abc", dep.ID)
	}
}
