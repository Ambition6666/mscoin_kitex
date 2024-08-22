package handler

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"grpc_common/kitex_gen/ucenter"
	"ucenter/domain"
)

type RegisterImpl struct {
	CaptchaDomain *domain.CaptchaDomain
}

func (s *RegisterImpl) RegisterByPhone(ctx context.Context, req *ucenter.RegReq) (res *ucenter.RegRes, err error) {
	isVerify := s.CaptchaDomain.Verify(req.Captcha.Server, req.Captcha.Token, 2, req.Ip)
	if !isVerify {
		return nil, kerrors.NewBizStatusError(-1, "人机验证不通过")
	}
	return
}

func (s *RegisterImpl) SendCode(ctx context.Context, req *ucenter.CodeReq) (res *ucenter.NoRes, err error) {
	return
}

func NewRegisterImpl() *RegisterImpl {
	return &RegisterImpl{
		CaptchaDomain: domain.NewCaptchaDomain(),
	}
}
