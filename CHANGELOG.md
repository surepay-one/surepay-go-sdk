# Changelog

All notable changes to this project will be documented in this file.
Format: [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
Versioning: [SemVer](https://semver.org/)

## [Unreleased]

## [0.1.0] - 2026-06-09

### Fixed
- Error code constants corrected to lowercase (`not_found`, `duplicate_request`, etc.) matching actual API responses — `IsNotFound`, `IsRateLimit`, `IsInsufficientBalance`, `IsDuplicate` now work correctly

### Added
- `Client` with functional-options constructor (`New`)
- `Balance.Get` — GET /balance
- `Deposits.List`, `Deposits.Create`, `Deposits.Get` — GET/POST /deposits, GET /deposits/{id}
- `Payouts.List`, `Payouts.Create`, `Payouts.Get` — GET/POST /payouts, GET /payouts/{id}
- `BankInquiry.Verify` — POST /bank-inquiry
- Automatic HMAC-SHA256 request signing when `apiSecret` is set
- Retry with exponential backoff on 5xx and network errors (default 3 attempts)
- `WithIdempotencyKey` context helper for safe POST retries
- `VerifyWebhookSignature` standalone function for inbound webhook verification
- Typed `APIError` with `Code`, `Message`, `StatusCode` fields
- Predicate helpers: `IsNotFound`, `IsRateLimit`, `IsInsufficientBalance`, `IsDuplicate`
- Generic `ListResult[T]` for paginated responses
- Zero non-stdlib dependencies
