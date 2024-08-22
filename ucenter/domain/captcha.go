package domain

import (
	"common/tools"
	"encoding/json"
	"github.com/cloudwego/kitex/pkg/klog"
	"ucenter/config"
)

type CaptchaDomain struct {
}

type vaptchaReq struct {
	Id        string `json:"id"`
	Secretkey string `json:"secretkey"`
	Scene     int    `json:"scene"`
	Token     string `json:"token"`
	Ip        string `json:"ip"`
}

type vaptchaRsp struct {
	Success int    `json:"success"`
	Score   int    `json:"score"`
	Msg     string `json:"msg"`
}

func (d *CaptchaDomain) Verify(server string, token string, scene int, ip string) bool {
	req := vaptchaReq{
		Id:        config.GetConf().Captcha.Vid,
		Secretkey: config.GetConf().Captcha.Key,
		Scene:     scene,
		Token:     token,
		Ip:        ip,
	}
	resp, err := tools.Post(server, req)
	if err != nil {
		klog.Error("CaptchaDomain Verify post err : %s\n", err.Error())
		return false
	}

	var vaptchaRsp *vaptchaRsp
	err = json.Unmarshal(resp, &vaptchaRsp)
	if err != nil {
		klog.Errorf("CaptchaDomain Verify Unmarshal respBytes err : %s\n", err.Error())
		return false
	}

	if vaptchaRsp != nil && vaptchaRsp.Success == 1 {
		klog.Info("CaptchaDomain Verify success\n")
		return true
	}

	return false
}

func NewCaptchaDomain() *CaptchaDomain {
	return new(CaptchaDomain)
}
