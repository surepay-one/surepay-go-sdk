// Package balance provides access to the SurePay /balance endpoint.
package balance

import (
	"context"

	"github.com/surepay-one/surepay-go-sdk/internal/core"
	"github.com/surepay-one/surepay-go-sdk/internal/shared"
)

// Service provides access to the /balance endpoint.
type Service struct{ d *core.Doer }

// New creates a BalanceService backed by the given Doer.
func New(d *core.Doer) *Service { return &Service{d: d} }

// Get returns the current wallet balance, hold, and available amount.
// Required scope: balance:read
func (s *Service) Get(ctx context.Context) (*shared.Balance, error) {
	var out shared.Balance
	if err := s.d.Do(ctx, "GET", "/balance", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
