package handler

import (
	"context"
	exchange "grpc_common/kitex_gen/exchange"
)

// OrderImpl implements the last service interface defined in the IDL.
type OrderImpl struct{}

// FindOrderHistory implements the OrderImpl interface.
func (s *OrderImpl) FindOrderHistory(ctx context.Context, req *exchange.OrderReq) (resp *exchange.OrderRes, err error) {
	// TODO: Your code here...
	return
}

// FindOrderCurrent implements the OrderImpl interface.
func (s *OrderImpl) FindOrderCurrent(ctx context.Context, req *exchange.OrderReq) (resp *exchange.OrderRes, err error) {
	// TODO: Your code here...
	return
}

// Add implements the OrderImpl interface.
func (s *OrderImpl) Add(ctx context.Context, req *exchange.OrderReq) (resp *exchange.AddOrderRes, err error) {
	// TODO: Your code here...
	return
}

// FindByOrderId implements the OrderImpl interface.
func (s *OrderImpl) FindByOrderId(ctx context.Context, req *exchange.OrderReq) (resp *exchange.ExchangeOrderOrigin, err error) {
	// TODO: Your code here...
	return
}

// CancelOrder implements the OrderImpl interface.
func (s *OrderImpl) CancelOrder(ctx context.Context, req *exchange.OrderReq) (resp *exchange.CancelOrderRes, err error) {
	// TODO: Your code here...
	return
}

func NewOrderImpl() *OrderImpl {
	return &OrderImpl{}
}
