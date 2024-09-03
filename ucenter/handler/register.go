package handler

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"grpc_common/kitex_gen/ucenter"
	"ucenter/domain"
)

type RegisterImpl struct {
	CaptchaDomain *domain.CaptchaDomain
	MemberDomain  *domain.MemberDomain
}

func (s *RegisterImpl) RegisterByPhone(ctx context.Context, req *ucenter.RegReq) (res *ucenter.RegRes, err error) {
	isVerify := s.CaptchaDomain.Verify(req.Captcha.Server, req.Captcha.Token, 2, req.Ip)
	if !isVerify {
		return nil, kerrors.NewBizStatusError(-1, "人机验证不通过")
	}
	mem := s.MemberDomain.FindMemberByPhone(ctx, req.Phone)
	if mem != nil {
		return nil, kerrors.NewBizStatusError(-1, "手机号已经被注册")
	}
	err = s.MemberDomain.Register(
		ctx,
		req.Username,
		req.Phone,
		req.Password,
		req.Country,
		req.Promotion,
		req.SuperPartner,
	)
	if err != nil {
		klog.Error("注册失败: ", err)
		return nil, kerrors.NewBizStatusError(-1, "注册失败")
	}
	return
}

func (s *RegisterImpl) SendCode(ctx context.Context, req *ucenter.CodeReq) (res *ucenter.NoRes, err error) {
	return
}

func NewRegisterImpl() *RegisterImpl {
	return &RegisterImpl{
		CaptchaDomain: domain.NewCaptchaDomain(),
		MemberDomain:  domain.NewMemberDomain(),
	}
}
