package handler

import (
	"context"
	login "grpc_common/kitex_gen/ucenter"
)

// LoginImpl implements the last service interface defined in the IDL.
type LoginImpl struct{}

// Login implements the LoginImpl interface.
func (s *LoginImpl) Login(ctx context.Context, req *login.LoginReq) (resp *login.LoginRes, err error) {
	// TODO: Your code here...
	return
}
