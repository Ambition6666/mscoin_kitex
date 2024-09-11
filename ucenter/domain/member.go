package domain

import (
	"common/tools"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"ucenter/dao"
	"ucenter/model"
)

type MemberDomain struct {
	MemberDao *dao.MemberDao
}

func (d *MemberDomain) FindMemberByPhone(ctx context.Context, phone string) *model.Member {
	mem, err := d.MemberDao.FindByPhone(ctx, phone)
	if err != nil {
		return nil
	}
	return mem
}

func (d *MemberDomain) Register(ctx context.Context, username string, phone string, password string, country string, promotion string, partner string) error {
	mem := model.NewMember()
	tools.Default(mem)
	mem.Id = 0

	//首先处理密码 密码要md5加密，但是md5不安全，所以我们给加上一个salt值

	salt, pwd := tools.Encode(password, nil)
	mem.Salt = salt
	mem.Password = pwd
	mem.MobilePhone = phone
	mem.Username = username
	mem.Country = country
	mem.PromotionCode = promotion
	mem.FillSuperPartner(partner)
	mem.MemberLevel = model.GENERAL
	mem.Avatar = "https://mszlu.oss-cn-beijing.aliyuncs.com/mscoin/defaultavatar.png"
	err := d.MemberDao.Save(ctx, mem)
	return err
}

func (d *MemberDomain) UpdateLoginCount(id int64, incr int) {
	err := d.MemberDao.UpdateLoginCount(context.Background(), id, incr)
	if err != nil {
		klog.Error(err)
	}
}

func (d *MemberDomain) FindMemberById(ctx context.Context, id int64) *model.Member {
	mem, err := d.MemberDao.FindById(ctx, id)
	if err != nil {
		return nil
	}
	return mem
}

func NewMemberDomain() *MemberDomain {
	return &MemberDomain{
		MemberDao: dao.NewMemberDao(),
	}
}
