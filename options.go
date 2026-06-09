package surepay

import (
	"net/http"
	"strings"
	"time"

	"github.com/surepay-one/surepay-go-sdk/internal/core"
)

// config holds construction-time settings before the Client is built.
type config struct {
	apiKey     string
	apiSecret  string
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	maxRetries int
	log        core.Logger
}

// Option configures a [Client].
type Option func(*config)

// WithBaseURL overrides the default API base URL.
// Sandbox: "http://localhost:8000/merchant/v1"
func WithBaseURL(u string) Option {
	return func(c *config) { c.baseURL = strings.TrimRight(u, "/") }
}

// WithHTTPClient sets a custom *http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *config) { c.httpClient = hc }
}

// WithTimeout sets the per-request timeout (default 30s).
func WithTimeout(d time.Duration) Option {
	return func(c *config) { c.timeout = d }
}

// WithMaxRetries sets the maximum retry count on 5xx and network errors (default 3).
// Pass 0 to disable retries.
func WithMaxRetries(n int) Option {
	return func(c *config) { c.maxRetries = n }
}

// WithLogger attaches a structured logger (*log/slog.Logger satisfies this directly).
func WithLogger(l Logger) Option {
	return func(c *config) { c.log = l }
}
