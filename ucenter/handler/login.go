package handler

import (
	"common/tools"
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/golang-jwt/jwt/v4"
	login "grpc_common/kitex_gen/ucenter"
	"time"
	"ucenter/config"
	"ucenter/domain"
)

// LoginImpl implements the last service interface defined in the IDL.
type LoginImpl struct {
	CaptchaDomain *domain.CaptchaDomain
	MemberDomain  *domain.MemberDomain
}

// Login implements the LoginImpl interface.
func (s *LoginImpl) Login(ctx context.Context, req *login.LoginReq) (resp *login.LoginRes, err error) {
	// TODO: Your code here...
	//校验人机
	isVerify := s.CaptchaDomain.Verify(req.Captcha.Server, req.Captcha.Token, 2, req.Ip)
	if !isVerify {
		return nil, kerrors.NewBizStatusError(-1, "人机验证不通过")
	}
	//查询salt
	mem := s.MemberDomain.FindMemberByPhone(ctx, req.Username)
	if mem == nil {
		return nil, kerrors.NewBizStatusError(-1, "用户不存在")
	}
	salt := mem.Salt
	verify := tools.Verify(req.Password, salt, mem.Password, nil)
	if !verify {
		return nil, kerrors.NewBizStatusError(-1, "账号密码不正确")
	}
	accessExpire := config.GetConf().JWT.AccessExpire
	accessSecret := config.GetConf().JWT.AccessSecret
	token, err := s.getJwtToken(accessSecret, time.Now().Unix(), accessExpire, mem.Id)
	if err != nil {
		return nil, kerrors.NewBizStatusError(-1, "未知错误，请联系管理员")
	}
	loginCount := mem.LoginCount + 1
	go func() {
		s.MemberDomain.UpdateLoginCount(mem.Id, 1)
	}()
	return &login.LoginRes{
		Token:         token,
		Id:            mem.Id,
		Username:      mem.Username,
		MemberLevel:   mem.MemberLevelStr(),
		MemberRate:    mem.MemberRate(),
		RealName:      mem.RealName,
		Country:       mem.Country,
		Avatar:        mem.Avatar,
		PromotionCode: mem.PromotionCode,
		SuperPartner:  mem.SuperPartner,
		LoginCount:    int32(loginCount),
	}, nil
	return
}

func (l *LoginImpl) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func NewLoginImpl() *LoginImpl {
	return &LoginImpl{
		CaptchaDomain: domain.NewCaptchaDomain(),
		MemberDomain:  domain.NewMemberDomain(),
	}
}
