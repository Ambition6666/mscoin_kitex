package config

const (
	EtcdAddr   = "192.168.40.134:2379"
	ServerName = "market_api"
	ServerAddr = "192.168.40.1:8889"
	MID        = 1
)

type Rocketmq struct {
	WriteCap int    `yaml:"WriteCap"`
	ReadCap  int    `yaml:"ReadCap"`
	Addr     string `yaml:"Addr"`
}

type config struct {
	Rocketmq Rocketmq `yaml:"Rocketmq"`
}

var conf config

func GetConf() *config {
	return &conf
}
