# surepay-go-sdk

Official Go client library for the [SurePay](https://surepay.one) Merchant API.

[![Go Reference](https://pkg.go.dev/badge/github.com/surepay-one/surepay-go-sdk.svg)](https://pkg.go.dev/github.com/surepay-one/surepay-go-sdk)
[![CI](https://github.com/surepay-one/surepay-go-sdk/actions/workflows/ci.yml/badge.svg)](https://github.com/surepay-one/surepay-go-sdk/actions/workflows/ci.yml)

## Requirements

- Go 1.23+
- Zero non-stdlib dependencies

## Install

```bash
go get github.com/surepay-one/surepay-go-sdk
```

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "os"

    surepay "github.com/surepay-one/surepay-go-sdk"
)

func main() {
    client := surepay.New(
        os.Getenv("SUREPAY_API_KEY"),    // tpay_live_... or tpay_test_...
        os.Getenv("SUREPAY_API_SECRET"), // tpay_sec_... — enables auto HMAC signing
    )
    ctx := context.Background()

    // Check wallet balance
    bal, _ := client.Balance.Get(ctx)
    fmt.Printf("Available: %d VND\n", bal.Available)

    // Create a deposit order (thu hộ)
    dep, _ := client.Deposits.Create(ctx, &surepay.CreateDepositRequest{
        Amount:    100_000,
        RequestID: "ORD-20260609-001",
    })
    fmt.Println("Checkout URL:", dep.CheckoutURL)

    // Create a payout (chi hộ)
    pay, err := client.Payouts.Create(ctx, &surepay.CreatePayoutRequest{
        Amount:      500_000,
        BankCode:    "VCB",
        BankAccount: "1234567890",
        FullName:    "NGUYEN VAN A",
        Description: "Salary June 2026",
    })
    if err != nil {
        if surepay.IsInsufficientBalance(err) {
            fmt.Println("Not enough balance — top up first")
            return
        }
        panic(err)
    }
    fmt.Println("Payout ID:", pay.ID)
}
```

## Configuration

```go
client := surepay.New(apiKey, apiSecret,
    surepay.WithBaseURL("https://api.surepay.one/merchant/v1"),
    surepay.WithTimeout(15 * time.Second),
    surepay.WithMaxRetries(3),           // retries on 5xx + network errors
    surepay.WithLogger(slog.Default()),  // *log/slog.Logger satisfies the interface
)
```

| Option | Default | Description |
|--------|---------|-------------|
| `WithBaseURL(url)` | `https://api.surepay.one/merchant/v1` | Override base URL for local dev or staging |
| `WithTimeout(d)` | `30s` | HTTP request timeout |
| `WithMaxRetries(n)` | `3` | Retry attempts on 5xx and network errors |
| `WithHTTPClient(hc)` | — | Bring your own `*http.Client` |
| `WithLogger(l)` | no-op | Structured logger — any `*slog.Logger` works |

## Authentication

Every request requires an API key sent as an `X-API-Key` header. When `apiSecret` is provided to `New()`, every outgoing request is automatically signed with HMAC-SHA256 — the `X-Signature` and `X-Timestamp` headers are attached with no extra code needed.

## API reference

### Balance

#### `client.Balance.Get(ctx)`

Get current wallet balance. **Requires:** `balance:read` scope.

```go
bal, err := client.Balance.Get(ctx)
// bal.Balance   — total wallet balance in VND
// bal.Hold      — reserved for in-flight transactions
// bal.Available — bal.Balance - bal.Hold
// bal.Currency  — always "VND"
```

---

### Deposits

#### `client.Deposits.List(ctx, params)`

Paginated list of deposit (thu hộ) orders. **Requires:** `deposits:read` scope.

```go
result, err := client.Deposits.List(ctx, &surepay.DepositListParams{
    Page:     1,
    PageSize: 20,
    Status:   "success",        // pending|processing|success|failed|expired|cancelled
    Search:   "ORD-2026",
    FromDate: "2026-06-01",     // YYYY-MM-DD
    ToDate:   "2026-06-30",
})
// result.Items       — []surepay.Deposit
// result.Total       — total matching records
// result.TotalPages  — total pages
```

#### `client.Deposits.Create(ctx, req)`

Create a new deposit order. Returns a `checkout_url` (redirect) and `qr_code` (VietQR). **Requires:** `deposits:write` scope.

```go
dep, err := client.Deposits.Create(ctx, &surepay.CreateDepositRequest{
    Amount:         100_000,            // VND, required
    RequestID:      "ORD-20260609-001", // your order ID — optional, for idempotency
    Currency:       "VND",              // optional, default "VND"
    // Chính chủ verification — all optional:
    SenderBankID:   "970436",
    SenderBankName: "Vietcombank",
    SenderAccount:  "1234567890",
    SenderName:     "NGUYEN VAN A",
})
fmt.Println(dep.CheckoutURL)
fmt.Println(*dep.QRCode)
```

**Response fields:**

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | SurePay transaction UUID |
| `RequestID` | string | Your order ID |
| `Amount` | int64 | Amount in VND |
| `Status` | DepositStatus | pending, processing, success, failed, expired, cancelled |
| `Fee` | int64 | Transaction fee in VND |
| `CheckoutURL` | string | Redirect URL for payer |
| `QRCode` | \*string | VietQR data string |
| `AccountNumber` | string | Receiving bank account number |
| `AccountName` | string | Receiving account holder name |
| `BankBIN` | string | Receiving bank BIN |
| `ExpiresAt` | time.Time | Order expiry time |
| `IsOwnerVerified` | \*bool | Chính chủ result — nil until transfer received |

#### `client.Deposits.Get(ctx, id)`

Fetch a single deposit order by UUID. **Requires:** `deposits:read` scope.

```go
dep, err := client.Deposits.Get(ctx, "uuid-here")
// dep.Status: DepositStatusPending | DepositStatusSuccess | ...
```

---

### Payouts

#### `client.Payouts.List(ctx, params)`

Paginated list of payout (chi hộ) orders. **Requires:** `payouts:read` scope.

```go
result, err := client.Payouts.List(ctx, &surepay.PayoutListParams{
    Page:     1,
    PageSize: 20,
    Status:   "success",    // pending|processing|success|failed
    Search:   "REF-2026",
    FromDate: "2026-06-01",
    ToDate:   "2026-06-30",
})
```

#### `client.Payouts.Create(ctx, req)`

Initiate a payout bank transfer. Funds are deducted from your wallet immediately on success. **Requires:** `payouts:write` scope.

> Payouts are irreversible once status moves past `pending`. Verify bank details with `BankInquiry.Verify` first.

```go
pay, err := client.Payouts.Create(ctx, &surepay.CreatePayoutRequest{
    Amount:      500_000,       // VND, required
    BankCode:    "VCB",         // required (VCB, MB, TCB, ACB, ...)
    BankAccount: "1234567890",  // required
    FullName:    "NGUYEN VAN A",// required — use UPPERCASE
    Description: "Salary June", // required — transfer memo
    BankName:    "Vietcombank", // optional
})
```

#### `client.Payouts.Get(ctx, id)`

Fetch a single payout by UUID. **Requires:** `payouts:read` scope.

```go
pay, err := client.Payouts.Get(ctx, "uuid-here")
// pay.Status: PayoutStatusPending | PayoutStatusSuccess | ...
```

---

### Bank Inquiry

#### `client.BankInquiry.Verify(ctx, req)`

Look up the account holder name for a bank account. Call this before creating a payout to confirm the recipient. **Requires:** `payouts:read` scope.

```go
result, err := client.BankInquiry.Verify(ctx, &surepay.BankInquiryRequest{
    BankCode:      "VCB",
    AccountNumber: "1234567890",
})
if err == nil {
    fmt.Println("Account name:", result.AccountName)
}
```

---

## Idempotency

Attach an idempotency key to any `context.Context` before POST requests. The key is forwarded as an `Idempotency-Key` header — safe to retry on network errors without risk of duplicate transactions.

```go
ctx := surepay.WithIdempotencyKey(context.Background(), "ORD-20260609-001")
dep, err := client.Deposits.Create(ctx, &surepay.CreateDepositRequest{
    Amount:    100_000,
    RequestID: "ORD-20260609-001",
})
```

## Webhook verification

Every inbound webhook event from SurePay is HMAC-signed. Pass the **raw** request body (before `json.Unmarshal`) to `VerifyWebhookSignature`:

```go
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)

    if !surepay.VerifyWebhookSignature(os.Getenv("SUREPAY_API_SECRET"), body, "") {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    var base struct {
        Event string `json:"event"`
    }
    json.Unmarshal(body, &base)

    switch base.Event {
    case "deposit.success", "deposit.failed":
        var evt surepay.DepositEvent
        json.Unmarshal(body, &evt)
        fmt.Printf("deposit %s: %s\n", evt.Status, evt.ID)

    case "payout.success", "payout.failed":
        var evt surepay.PayoutEvent
        json.Unmarshal(body, &evt)
        fmt.Printf("payout %s: %s\n", evt.Status, evt.PayoutID)
    }

    w.WriteHeader(http.StatusOK)
}
```

## Error handling

```go
dep, err := client.Deposits.Create(ctx, req)
if err != nil {
    // Predicate helpers
    if surepay.IsNotFound(err)             { /* 404 */ }
    if surepay.IsRateLimit(err)            { /* 429 */ }
    if surepay.IsInsufficientBalance(err)  { /* 422 insufficient_balance */ }
    if surepay.IsDuplicate(err)            { /* 409 duplicate_request */ }

    // Full error details
    var apiErr *surepay.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d  code=%s  %s\n", apiErr.StatusCode, apiErr.Code, apiErr.Message)
    }
}
```

**Error codes:**

| HTTP | `Code` | Meaning |
|------|--------|---------|
| 400 | `validation_error` | Invalid request body or parameters |
| 401 | `unauthorized` | Missing or invalid API key |
| 401 | `signature_invalid` | HMAC signature failed or timestamp > 5 min |
| 403 | `permission_denied` | API key lacks required scope |
| 404 | `not_found` | Resource not found |
| 409 | `duplicate_request` | Idempotency key conflict |
| 422 | `insufficient_balance` | Top up wallet first |
| 422 | `invalid_state_transition` | Operation not allowed for current status |
| 500 | `internal_error` | Server error |

## HMAC signing

When `apiSecret` is set, all requests are signed automatically. The signing algorithm for manual use:

```
signingString = UNIX_TIMESTAMP + "\n" + METHOD + "\n" + PATH + "\n" + hex(sha256(body))
signature     = "sha256=" + hex(hmac-sha256(apiSecret, signingString))
```

Attach as headers: `X-Signature: <signature>` and `X-Timestamp: <unix_timestamp>`.

Signatures expire after **300 seconds** — generate per-request, never cache or reuse.

## License

MIT
