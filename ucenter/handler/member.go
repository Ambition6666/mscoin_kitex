package handler

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/jinzhu/copier"
	"grpc_common/kitex_gen/ucenter"
	"ucenter/domain"
)

type MemberImpl struct {
	memberDomain *domain.MemberDomain
}

func (m MemberImpl) FindMemberById(ctx context.Context, req *ucenter.MemberReq) (res *ucenter.MemberRes, err error) {
	//TODO implement me
	mem := m.memberDomain.FindMemberById(ctx, req.MemberId)
	err = copier.Copy(res, mem)
	if err != nil {
		klog.Error("FindMemberById: ", err)
		return nil, err
	}

	return
}

func NewMemberImpl() *MemberImpl {
	return &MemberImpl{
		memberDomain: domain.NewMemberDomain(),
	}
}
