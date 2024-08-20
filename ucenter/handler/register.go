package handler

import (
	"context"
	"grpc_common/kitex_gen/ucenter"
)

type RegisterImpl struct{}

func (s *RegisterImpl) RegisterByPhone(ctx context.Context, req *ucenter.RegReq) (res *ucenter.RegRes, err error) {
	return
}
