// Package shared defines domain types shared across all resource sub-packages.
// External consumers access these via type aliases in the root surepay package.
package shared

import "time"

// ─── Balance ──────────────────────────────────────────────────────────────────

type Balance struct {
	Balance   int64  `json:"balance"`
	Hold      int64  `json:"hold"`
	Available int64  `json:"available"`
	Currency  string `json:"currency"`
}

// ─── Deposits ─────────────────────────────────────────────────────────────────

type DepositStatus string

const (
	DepositStatusPending    DepositStatus = "pending"
	DepositStatusProcessing DepositStatus = "processing"
	DepositStatusSuccess    DepositStatus = "success"
	DepositStatusFailed     DepositStatus = "failed"
	DepositStatusExpired    DepositStatus = "expired"
	DepositStatusCancelled  DepositStatus = "cancelled"
)

type Deposit struct {
	ID              string        `json:"id"`
	RefCode         string        `json:"ref_code"`
	RequestID       string        `json:"request_id"`
	Amount          int64         `json:"amount"`
	Currency        string        `json:"currency"`
	Status          DepositStatus `json:"status"`
	Reason          *string       `json:"reason"`
	Fee             int64         `json:"fee"`
	CheckoutURL     string        `json:"checkout_url"`
	QRCode          *string       `json:"qr_code"`
	AccountNumber   string        `json:"account_number"`
	AccountName     string        `json:"account_name"`
	BankBIN         string        `json:"bank_bin"`
	CreatedAt       time.Time     `json:"created_at"`
	ExpiresAt       time.Time     `json:"expires_at"`
	CompletedAt     *time.Time    `json:"completed_at"`
	IsOwnerVerified *bool         `json:"is_owner_verified"`
	SenderName      *string       `json:"sender_name"`
	SenderBankID    *string       `json:"sender_bank_id"`
}

type CreateDepositRequest struct {
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency,omitempty"`
	RequestID      string `json:"request_id,omitempty"`
	SenderBankID   string `json:"sender_bank_id,omitempty"`
	SenderBankName string `json:"sender_bank_name,omitempty"`
	SenderAccount  string `json:"sender_account,omitempty"`
	SenderName     string `json:"sender_name,omitempty"`
}

type DepositListParams struct {
	Page     int
	PageSize int
	Status   string
	Search   string
	FromDate string
	ToDate   string
}

// ─── Payouts ──────────────────────────────────────────────────────────────────

type PayoutStatus string

const (
	PayoutStatusPending    PayoutStatus = "pending"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusSuccess    PayoutStatus = "success"
	PayoutStatusFailed     PayoutStatus = "failed"
)

type Payout struct {
	ID          string       `json:"id"`
	RefCode     string       `json:"ref_code"`
	Amount      int64        `json:"amount"`
	Fee         int64        `json:"fee"`
	Status      PayoutStatus `json:"status"`
	Reason      *string      `json:"reason"`
	BankCode    string       `json:"bank_code"`
	BankAccount string       `json:"bank_account"`
	BankName    string       `json:"bank_name"`
	FullName    string       `json:"full_name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
}

type CreatePayoutRequest struct {
	Amount      int64  `json:"amount"`
	BankCode    string `json:"bank_code"`
	BankAccount string `json:"bank_account"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	BankName    string `json:"bank_name,omitempty"`
}

type PayoutListParams struct {
	Page     int
	PageSize int
	Status   string
	Search   string
	FromDate string
	ToDate   string
}

// ─── Shared pagination ────────────────────────────────────────────────────────

type ListResult[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

// ─── Bank Inquiry ─────────────────────────────────────────────────────────────

type BankInquiryRequest struct {
	BankCode      string `json:"bank_code"`
	AccountNumber string `json:"account_number"`
}

type BankInquiryResult struct {
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
}

// ─── Webhook events ───────────────────────────────────────────────────────────

type DepositEvent struct {
	Event           string  `json:"event"`
	ID              string  `json:"id"`
	RefCode         string  `json:"ref_code"`
	OrderCode       string  `json:"order_code"`
	TxnID           string  `json:"txn_id"`
	Amount          int64   `json:"amount"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`
	Reason          *string `json:"reason"`
	Fee             int64   `json:"fee"`
	CreatedAt       string  `json:"created_at"`
	ExpiresAt       string  `json:"expires_at"`
	IsOwnerVerified bool    `json:"is_owner_verified"`
	SenderName      *string `json:"sender_name"`
	SenderBankID    *string `json:"sender_bank_id"`
	Timestamp       int64   `json:"timestamp"`
	Signature       string  `json:"signature"`
}

type PayoutEvent struct {
	Event       string  `json:"event"`
	PayoutID    string  `json:"payout_id"`
	RefCode     string  `json:"ref_code"`
	Amount      int64   `json:"amount"`
	Fee         int64   `json:"fee"`
	Status      string  `json:"status"`
	Reason      *string `json:"reason"`
	BankCode    string  `json:"bank_code"`
	BankAccount string  `json:"bank_account"`
	FullName    string  `json:"full_name"`
	Currency    string  `json:"currency"`
	Timestamp   int64   `json:"timestamp"`
	Signature   string  `json:"signature"`
}
