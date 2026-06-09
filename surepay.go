// Package surepay provides a Go client for the SurePay Merchant API.
//
// # Quick start
//
//	client := surepay.New(
//	    os.Getenv("SUREPAY_API_KEY"),
//	    os.Getenv("SUREPAY_API_SECRET"), // enables auto HMAC signing
//	)
//
//	bal, err  := client.Balance.Get(ctx)
//	dep, err  := client.Deposits.Create(ctx, &surepay.CreateDepositRequest{Amount: 100_000})
//	pay, err  := client.Payouts.Create(ctx, &surepay.CreatePayoutRequest{...})
//
// # Webhook verification
//
//	ok := surepay.VerifyWebhookSignature(secret, rawBody, r.Header.Get("X-Surepay-Signature"))
package surepay

import (
	"context"
	"net/http"
	"time"

	"github.com/surepay-one/surepay-go-sdk/balance"
	"github.com/surepay-one/surepay-go-sdk/bankinquiry"
	"github.com/surepay-one/surepay-go-sdk/deposits"
	"github.com/surepay-one/surepay-go-sdk/internal/core"
	"github.com/surepay-one/surepay-go-sdk/payouts"
)

const defaultBaseURL = "https://api.surepay.one/merchant/v1"

// Logger is a structured-log interface compatible with *log/slog.Logger.
type Logger = core.Logger

// Client is the SurePay API client. Safe for concurrent use.
// Create once with [New] and reuse across your application.
type Client struct {
	doer *core.Doer

	Balance     *balance.Service
	Deposits    *deposits.Service
	Payouts     *payouts.Service
	BankInquiry *bankinquiry.Service
}

// New creates a SurePay API client.
// apiSecret is optional — when non-empty, every request is HMAC-signed automatically.
func New(apiKey, apiSecret string, opts ...Option) *Client {
	cfg := &config{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		baseURL:    defaultBaseURL,
		timeout:    30 * time.Second,
		maxRetries: 3,
		log:        core.NoopLogger{},
	}
	for _, opt := range opts {
		opt(cfg)
	}
	hc := cfg.httpClient
	if hc == nil {
		hc = &http.Client{Timeout: cfg.timeout}
	}

	doer := &core.Doer{
		APIKey:     cfg.apiKey,
		APISecret:  cfg.apiSecret,
		BaseURL:    cfg.baseURL,
		HTTPClient: hc,
		MaxRetries: cfg.maxRetries,
		Log:        cfg.log,
	}
	return &Client{
		doer:        doer,
		Balance:     balance.New(doer),
		Deposits:    deposits.New(doer),
		Payouts:     payouts.New(doer),
		BankInquiry: bankinquiry.New(doer),
	}
}

// WithIdempotencyKey attaches an idempotency key to ctx. The key is forwarded
// as an Idempotency-Key header on the next request.
func WithIdempotencyKey(ctx context.Context, key string) context.Context {
	return core.WithIdempotencyKey(ctx, key)
}
