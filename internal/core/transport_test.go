package core

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testDoer(srv *httptest.Server) *Doer {
	return &Doer{
		APIKey:     "test-key",
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
		MaxRetries: 0,
		Log:        NoopLogger{},
	}
}

func jsonHandler(status int, body any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(body)
	}
}

func TestDoer_success(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(200, map[string]any{
		"data": map[string]any{"balance": float64(1_000_000), "currency": "VND"},
	}))
	defer srv.Close()

	var out map[string]any
	err := testDoer(srv).Do(context.Background(), "GET", "/balance", nil, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["currency"] != "VND" {
		t.Errorf("currency = %v, want VND", out["currency"])
	}
}

func TestDoer_apiError(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(401, map[string]string{
		"code": "INVALID_API_KEY", "message": "not found",
	}))
	defer srv.Close()

	err := testDoer(srv).Do(context.Background(), "GET", "/balance", nil, nil)

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
	if apiErr.Code != "INVALID_API_KEY" {
		t.Errorf("Code = %q, want INVALID_API_KEY", apiErr.Code)
	}
}

func TestDoer_signingHeaders(t *testing.T) {
	var gotSig, gotTS string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get("X-Signature")
		gotTS = r.Header.Get("X-Timestamp")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	}))
	defer srv.Close()

	d := testDoer(srv)
	d.APISecret = "my-secret"
	d.Do(context.Background(), "GET", "/balance", nil, nil)

	if gotTS == "" {
		t.Error("X-Timestamp header missing")
	}
	if !strings.HasPrefix(gotSig, "sha256=") {
		t.Errorf("X-Signature = %q, want sha256= prefix", gotSig)
	}
}

func TestDoer_idempotencyKey(t *testing.T) {
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("Idempotency-Key")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	}))
	defer srv.Close()

	ctx := WithIdempotencyKey(context.Background(), "idem-123")
	testDoer(srv).Do(ctx, "POST", "/deposits", map[string]any{"amount": 1}, nil)

	if gotKey != "idem-123" {
		t.Errorf("Idempotency-Key = %q, want idem-123", gotKey)
	}
}

func TestAppendQuery(t *testing.T) {
	cases := []struct{ path, want string }{
		{"/deposits", "/deposits"},
	}
	_ = cases
	got := AppendQuery("/deposits", nil)
	if got != "/deposits" {
		t.Errorf("empty values: got %q, want /deposits", got)
	}
}
