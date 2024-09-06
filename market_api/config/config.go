package config

const (
	EtcdAddr   = "192.168.40.134:2379"
	ServerName = "market"
	ServerAddr = "127.0.0.1:8889"
)

type Rocketmq struct {
	Addr     string `yaml:"Addr"`
	WriteCap int    `yaml:"WriteCap"`
	ReadCap  int    `yaml:"ReadCap"`
}

type config struct {
	Rocketmq Rocketmq `yaml:"Rocketmq"`
}

var conf config

func GetConf() *config {
	return &conf
}
