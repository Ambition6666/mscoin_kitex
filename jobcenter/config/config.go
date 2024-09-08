package config

const (
	EtcdAddr   = "192.168.40.134:2379"
	ServerName = "jobcenter"
	MID        = 1
)

type Rocketmq struct {
	Addr     string `yaml:"Addr"`
	WriteCap int    `yaml:"WriteCap"`
	ReadCap  int    `yaml:"ReadCap"`
}

type CacheRedis struct {
	Host string `yaml:"Host"`
	Node int    `yaml:"Node"`
	Pass string `yaml:"Pass"`
}

type Okx struct {
	Apikey    string `yaml:"Apikey"`
	SecretKey string `yaml:"SecretKey"`
	Pass      string `yaml:"Pass"`
	Host      string `yaml:"Host"`
	Proxy     string `yaml:"Proxy"`
}

type Mongo struct {
	Url      string `yaml:"Url"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Database string `yaml:"Database"`
}

type config struct {
	Rocketmq   Rocketmq   `yaml:"Rocketmq"`
	CacheRedis CacheRedis `yaml:"CacheRedis"`
	Okx        Okx        `yaml:"Okx"`
	Mongo      Mongo      `yaml:"Mongo"`
}

var conf config

func GetConf() *config {
	return &conf
}
