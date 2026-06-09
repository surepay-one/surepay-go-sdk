// Package payouts provides access to the SurePay /payouts endpoints.
package payouts

import (
	"context"
	"net/url"
	"strconv"

	"github.com/surepay-one/surepay-go-sdk/internal/core"
	"github.com/surepay-one/surepay-go-sdk/internal/shared"
)

// Service provides access to the /payouts endpoints.
type Service struct{ d *core.Doer }

// New creates a PayoutsService backed by the given Doer.
func New(d *core.Doer) *Service { return &Service{d: d} }

// List returns a paginated list of payout orders.
// Required scope: payouts:read
func (s *Service) List(ctx context.Context, params *shared.PayoutListParams) (*shared.ListResult[shared.Payout], error) {
	path := "/payouts"
	if params != nil {
		path = core.AppendQuery(path, encodePayoutParams(params))
	}
	var out shared.ListResult[shared.Payout]
	if err := s.d.Do(ctx, "GET", path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create initiates a payout (chi hộ) bank transfer from your available balance.
// Use [surepay.WithIdempotencyKey] to safely retry on network errors.
// Required scope: payouts:write
func (s *Service) Create(ctx context.Context, req *shared.CreatePayoutRequest) (*shared.Payout, error) {
	var out shared.Payout
	if err := s.d.Do(ctx, "POST", "/payouts", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves a single payout by UUID.
// Required scope: payouts:read
func (s *Service) Get(ctx context.Context, id string) (*shared.Payout, error) {
	var out shared.Payout
	if err := s.d.Do(ctx, "GET", "/payouts/"+id, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func encodePayoutParams(p *shared.PayoutListParams) url.Values {
	v := url.Values{}
	if p.Page > 0 {
		v.Set("page", strconv.Itoa(p.Page))
	}
	if p.PageSize > 0 {
		v.Set("page_size", strconv.Itoa(p.PageSize))
	}
	if p.Status != "" {
		v.Set("status", p.Status)
	}
	if p.Search != "" {
		v.Set("search", p.Search)
	}
	if p.FromDate != "" {
		v.Set("from_date", p.FromDate)
	}
	if p.ToDate != "" {
		v.Set("to_date", p.ToDate)
	}
	return v
}
