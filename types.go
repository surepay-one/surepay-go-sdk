package surepay

import "github.com/surepay-one/surepay-go-sdk/internal/shared"

// ─── Balance ──────────────────────────────────────────────────────────────────

type Balance = shared.Balance

// ─── Deposits ─────────────────────────────────────────────────────────────────

type DepositStatus = shared.DepositStatus

const (
	DepositStatusPending    = shared.DepositStatusPending
	DepositStatusProcessing = shared.DepositStatusProcessing
	DepositStatusSuccess    = shared.DepositStatusSuccess
	DepositStatusFailed     = shared.DepositStatusFailed
	DepositStatusExpired    = shared.DepositStatusExpired
	DepositStatusCancelled  = shared.DepositStatusCancelled
)

type Deposit = shared.Deposit
type CreateDepositRequest = shared.CreateDepositRequest
type DepositListParams = shared.DepositListParams

// ─── Payouts ──────────────────────────────────────────────────────────────────

type PayoutStatus = shared.PayoutStatus

const (
	PayoutStatusPending    = shared.PayoutStatusPending
	PayoutStatusProcessing = shared.PayoutStatusProcessing
	PayoutStatusSuccess    = shared.PayoutStatusSuccess
	PayoutStatusFailed     = shared.PayoutStatusFailed
)

type Payout = shared.Payout
type CreatePayoutRequest = shared.CreatePayoutRequest
type PayoutListParams = shared.PayoutListParams

// ─── Shared ───────────────────────────────────────────────────────────────────

type ListResult[T any] = shared.ListResult[T]

// ─── Bank Inquiry ─────────────────────────────────────────────────────────────

type BankInquiryRequest = shared.BankInquiryRequest
type BankInquiryResult = shared.BankInquiryResult

// ─── Webhook events ───────────────────────────────────────────────────────────

type DepositEvent = shared.DepositEvent
type PayoutEvent = shared.PayoutEvent
