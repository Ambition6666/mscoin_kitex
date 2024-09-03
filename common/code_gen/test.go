package code_gen

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Ip       string `json:"ip,optional"`
}

type LoginRes struct {
	Username      string `json:"username"`
	Token         string `json:"token"`
	MemberLevel   string `json:"memberLevel"`
	RealName      string `json:"realName"`
	Country       string `json:"country"`
	Avatar        string `json:"avatar"`
	PromotionCode string `json:"promotionCode"`
	Id            int64  `json:"id"`
	LoginCount    int    `json:"loginCount"`
	SuperPartner  string `json:"superPartner"`
	MemberRate    int    `json:"memberRate"`
}
