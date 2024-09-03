package config

const (
	EtcdAddr   = "192.168.40.134:2379"
	ServerName = "ucenter"
	ServerAddr = "127.0.0.1:8081"
	MID        = 1
)

type Mysql struct {
	DataSource string `yaml:"DataSource"`
}

type CacheRedis struct {
	Pass string `yaml:"Pass"`
	Host string `yaml:"Host"`
	Type string `yaml:"Type"`
}

type Captcha struct {
	Key string `yaml:"Key"`
	Vid string `yaml:"Vid"`
}

type JWT struct {
	AccessSecret string `yaml:"AccessSecret"`
	AccessExpire int64  `yaml:"AccessExpire"`
}

type config struct {
	Mysql      Mysql      `yaml:"Mysql"`
	CacheRedis CacheRedis `yaml:"CacheRedis"`
	Captcha    Captcha    `yaml:"Captcha"`
	JWT        JWT        `yaml:"JWT"`
}

var conf config

func GetConf() *config {
	return &conf
}
