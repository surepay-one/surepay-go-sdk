package bankinquiry_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/surepay-one/surepay-go-sdk/bankinquiry"
	"github.com/surepay-one/surepay-go-sdk/internal/core"
	"github.com/surepay-one/surepay-go-sdk/internal/shared"
)

func testService(srv *httptest.Server) *bankinquiry.Service {
	return bankinquiry.New(&core.Doer{
		APIKey:     "test-key",
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
		MaxRetries: 0,
		Log:        core.NoopLogger{},
	})
}

func TestVerify(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/bank-inquiry" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"account_name":   "NGUYEN VAN A",
				"account_number": "1234567890",
				"bank_code":      "VCB",
			},
		})
	}))
	defer srv.Close()

	result, err := testService(srv).Verify(context.Background(), &shared.BankInquiryRequest{
		BankCode:      "VCB",
		AccountNumber: "1234567890",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AccountName != "NGUYEN VAN A" {
		t.Errorf("AccountName = %q, want NGUYEN VAN A", result.AccountName)
	}
	if gotBody["bank_code"] != "VCB" {
		t.Errorf("bank_code = %v, want VCB", gotBody["bank_code"])
	}
	if gotBody["account_number"] != "1234567890" {
		t.Errorf("account_number = %v, want 1234567890", gotBody["account_number"])
	}
}

func TestVerify_notFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"code": "not_found", "message": "bank account not found",
		})
	}))
	defer srv.Close()

	_, err := testService(srv).Verify(context.Background(), &shared.BankInquiryRequest{
		BankCode:      "VCB",
		AccountNumber: "0000000000",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
