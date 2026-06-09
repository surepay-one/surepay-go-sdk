package payouts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/surepay-one/surepay-go-sdk/internal/core"
	"github.com/surepay-one/surepay-go-sdk/internal/shared"
	"github.com/surepay-one/surepay-go-sdk/payouts"
)

func testService(srv *httptest.Server) *payouts.Service {
	return payouts.New(&core.Doer{
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
				"id": "pay-001", "amount": float64(200000), "status": "pending",
				"bank_code": "VCB", "bank_account": "1234567890", "full_name": "Nguyen Van A",
				"description": "salary", "created_at": "2026-01-01T00:00:00Z",
			},
		})
	}))
	defer srv.Close()

	pay, err := testService(srv).Create(context.Background(), &shared.CreatePayoutRequest{
		Amount: 200000, BankCode: "VCB", BankAccount: "1234567890",
		FullName: "Nguyen Van A", Description: "salary",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pay.ID != "pay-001" {
		t.Errorf("ID = %q, want pay-001", pay.ID)
	}
	if gotBody["bank_code"] != "VCB" {
		t.Errorf("bank_code = %v, want VCB", gotBody["bank_code"])
	}
}

func TestList_queryParams(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"items": []any{}, "total": 0, "page": 1, "page_size": 20, "total_pages": 0},
		})
	}))
	defer srv.Close()

	testService(srv).List(context.Background(), &shared.PayoutListParams{
		Page: 3, PageSize: 20, Status: "success",
	})

	for _, want := range []string{"page=3", "page_size=20", "status=success"} {
		if !strings.Contains(gotQuery, want) {
			t.Errorf("query %q missing %q", gotQuery, want)
		}
	}
}

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/payouts/pay-xyz" {
			t.Errorf("path = %q, want /payouts/pay-xyz", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "pay-xyz", "amount": float64(100000), "status": "success",
				"bank_code": "TCB", "bank_account": "9876543210", "full_name": "Tran Thi B",
				"description": "refund", "created_at": "2026-01-01T00:00:00Z",
			},
		})
	}))
	defer srv.Close()

	pay, err := testService(srv).Get(context.Background(), "pay-xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pay.ID != "pay-xyz" {
		t.Errorf("ID = %q, want pay-xyz", pay.ID)
	}
}
