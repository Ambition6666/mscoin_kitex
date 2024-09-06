package handler

import (
	"context"

	"grpc_common/kitex_gen/market"
	"market/domain"
)

// ExchangeRateImpl implements the last service interface defined in the IDL.
type ExchangeRateImpl struct {
	exchangeRateDomain *domain.ExchangeRateDomain
}

// UsdRate implements the ExchangeRateImpl interface.
func (s *ExchangeRateImpl) UsdRate(ctx context.Context, req *market.RateReq) (resp *market.RateRes, err error) {
	// TODO: Your code here...
	usdtRate := s.exchangeRateDomain.GetUsdRate(req.GetUnit())
	return &market.RateRes{
		Rate: usdtRate,
	}, nil
}

func NewExchangeRateImpl() *ExchangeRateImpl {
	return &ExchangeRateImpl{
		exchangeRateDomain: domain.NewExchangeRateDomain(),
	}
}
