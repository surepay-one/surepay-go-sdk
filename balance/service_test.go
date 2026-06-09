package balance_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/surepay-one/surepay-go-sdk/balance"
	"github.com/surepay-one/surepay-go-sdk/internal/core"
)

func testService(srv *httptest.Server) *balance.Service {
	return balance.New(&core.Doer{
		APIKey:     "test-key",
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
		MaxRetries: 0,
		Log:        core.NoopLogger{},
	})
}

func TestGet_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/balance" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"balance": float64(1_000_000), "hold": float64(100_000),
				"available": float64(900_000), "currency": "VND",
			},
		})
	}))
	defer srv.Close()

	bal, err := testService(srv).Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bal.Available != 900_000 {
		t.Errorf("available = %d, want 900000", bal.Available)
	}
	if bal.Currency != "VND" {
		t.Errorf("currency = %q, want VND", bal.Currency)
	}
}
