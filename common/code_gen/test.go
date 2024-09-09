package code_gen

type AssetReq struct {
	CoinName string `json:"coinName,optional" path:"coinName,optional"`
	Ip       string `json:"ip,optional"`
}
