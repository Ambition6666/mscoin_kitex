package config

const (
	EtcdAddr   = "192.168.40.134:2379"
	ServerName = "exchange"
	ServerAddr = "192.168.40.1:8083"
	MID        = 1
	MARKET     = "market"
	UCENTER    = "ucenter"
)

type Mysql struct {
	DataSource string `yaml:"DataSource"`
}

type CacheRedis struct {
	Host string `yaml:"Host"`
	Node int    `yaml:"Node"`
	Pass string `yaml:"Pass"`
}

type Rocketmq struct {
	Addr     string `yaml:"Addr"`
	WriteCap int    `yaml:"WriteCap"`
	ReadCap  int    `yaml:"ReadCap"`
}

type config struct {
	Mysql      Mysql      `yaml:"Mysql"`
	CacheRedis CacheRedis `yaml:"CacheRedis"`
	Rocketmq   Rocketmq   `yaml:"Rocketmq"`
}

var conf config

func GetConf() *config {
	return &conf
}
