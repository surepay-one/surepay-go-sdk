// Package deposits provides access to the SurePay /deposits endpoints.
package deposits

import (
	"context"
	"net/url"
	"strconv"

	"github.com/surepay-one/surepay-go-sdk/internal/core"
	"github.com/surepay-one/surepay-go-sdk/internal/shared"
)

// Service provides access to the /deposits endpoints.
type Service struct{ d *core.Doer }

// New creates a DepositsService backed by the given Doer.
func New(d *core.Doer) *Service { return &Service{d: d} }

// List returns a paginated list of deposit orders.
// Required scope: deposits:read
func (s *Service) List(ctx context.Context, params *shared.DepositListParams) (*shared.ListResult[shared.Deposit], error) {
	path := "/deposits"
	if params != nil {
		path = core.AppendQuery(path, encodeDepositParams(params))
	}
	var out shared.ListResult[shared.Deposit]
	if err := s.d.Do(ctx, "GET", path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a new deposit order (thu hộ).
// Required scope: deposits:write
func (s *Service) Create(ctx context.Context, req *shared.CreateDepositRequest) (*shared.Deposit, error) {
	var out shared.Deposit
	if err := s.d.Do(ctx, "POST", "/deposits", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves a single deposit order by UUID.
// Required scope: deposits:read
func (s *Service) Get(ctx context.Context, id string) (*shared.Deposit, error) {
	var out shared.Deposit
	if err := s.d.Do(ctx, "GET", "/deposits/"+id, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func encodeDepositParams(p *shared.DepositListParams) url.Values {
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
