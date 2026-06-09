// Package bankinquiry provides access to the SurePay /bank-inquiry endpoint.
package bankinquiry

import (
	"context"

	"github.com/surepay-one/surepay-go-sdk/internal/core"
	"github.com/surepay-one/surepay-go-sdk/internal/shared"
)

// Service provides access to the /bank-inquiry endpoint.
type Service struct{ d *core.Doer }

// New creates a BankInquiryService backed by the given Doer.
func New(d *core.Doer) *Service { return &Service{d: d} }

// Verify looks up the account holder name for a bank account.
// Call this before creating a payout to confirm the recipient's identity.
// Required scope: payouts:read
func (s *Service) Verify(ctx context.Context, req *shared.BankInquiryRequest) (*shared.BankInquiryResult, error) {
	var out shared.BankInquiryResult
	if err := s.d.Do(ctx, "POST", "/bank-inquiry", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
