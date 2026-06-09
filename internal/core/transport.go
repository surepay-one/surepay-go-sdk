package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type idempotencyKey struct{}

// WithIdempotencyKey attaches an idempotency key to ctx.
// The value is forwarded as an Idempotency-Key header on the next request.
func WithIdempotencyKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, idempotencyKey{}, key)
}

// Doer executes authenticated SurePay API requests.
type Doer struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	HTTPClient *http.Client
	MaxRetries int
	Log        Logger
}

// Do executes one API call, retrying on 5xx and network errors.
// body is JSON-marshalled if non-nil; dst is decoded from response["data"].
func (d *Doer) Do(ctx context.Context, method, path string, body, dst any) error {
	return d.attempt(ctx, method, path, body, dst, 0)
}

func (d *Doer) attempt(ctx context.Context, method, path string, body, dst any, n int) error {
	bodyBytes, err := marshalBody(body)
	if err != nil {
		return err
	}

	req, err := d.buildRequest(ctx, method, path, bodyBytes)
	if err != nil {
		return err
	}

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		if n < d.MaxRetries {
			d.Log.Warn("surepay: network error, retrying", "attempt", n+1, "err", err)
			sleep(n)
			return d.attempt(ctx, method, path, body, dst, n+1)
		}
		return fmt.Errorf("surepay: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("surepay: read response: %w", err)
	}

	if resp.StatusCode >= 500 && n < d.MaxRetries {
		d.Log.Warn("surepay: server error, retrying", "attempt", n+1, "status", resp.StatusCode)
		sleep(n)
		return d.attempt(ctx, method, path, body, dst, n+1)
	}

	if resp.StatusCode >= 400 {
		return parseError(resp.StatusCode, respBytes)
	}

	return decodeSuccess(respBytes, dst)
}

func (d *Doer) buildRequest(ctx context.Context, method, path string, bodyBytes []byte) (*http.Request, error) {
	var bodyReader io.Reader
	if len(bodyBytes) > 0 {
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, d.BaseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("surepay: build request: %w", err)
	}

	req.Header.Set("X-API-Key", d.APIKey)
	req.Header.Set("Accept", "application/json")
	if len(bodyBytes) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	if d.APISecret != "" {
		u, _ := url.Parse(d.BaseURL + path)
		sig, ts := Sign(d.APISecret, method, u.Path, bodyBytes)
		req.Header.Set("X-Signature", sig)
		req.Header.Set("X-Timestamp", ts)
	}

	if key, ok := ctx.Value(idempotencyKey{}).(string); ok && key != "" {
		req.Header.Set("Idempotency-Key", key)
	}

	return req, nil
}

// ─── helpers ──────────────────────────────────────────────────────────────────

type envelope[T any] struct {
	Data T `json:"data"`
}

type errBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIError is defined here so transport can create it; the root package
// re-exports it as surepay.APIError via a type alias.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("surepay: HTTP %d %s — %s", e.StatusCode, e.Code, e.Message)
}

func parseError(status int, body []byte) error {
	var eb errBody
	_ = json.Unmarshal(body, &eb)
	return &APIError{StatusCode: status, Code: eb.Code, Message: eb.Message}
}

func decodeSuccess(respBytes []byte, dst any) error {
	if dst == nil {
		return nil
	}
	var raw json.RawMessage
	env := envelope[*json.RawMessage]{Data: &raw}
	if err := json.Unmarshal(respBytes, &env); err != nil {
		return json.Unmarshal(respBytes, dst)
	}
	return json.Unmarshal(raw, dst)
}

func marshalBody(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("surepay: marshal body: %w", err)
	}
	return b, nil
}

// AppendQuery appends url.Values to a path string.
func AppendQuery(path string, v url.Values) string {
	if len(v) == 0 {
		return path
	}
	if strings.Contains(path, "?") {
		return path + "&" + v.Encode()
	}
	return path + "?" + v.Encode()
}

// sleep waits 2^n × 500ms ± 25% jitter, capped at 30s.
func sleep(n int) {
	base := 500 * time.Millisecond
	d := time.Duration(math.Pow(2, float64(n))) * base
	if d > 30*time.Second {
		d = 30 * time.Second
	}
	jitter := time.Duration(rand.Int63n(int64(d)/4 + 1))
	if rand.Intn(2) == 0 {
		d += jitter
	} else {
		d -= jitter
	}
	time.Sleep(d)
}
