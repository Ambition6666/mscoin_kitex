package config

const (
	EtcdAddr   = "192.168.40.134:2379"
	ServerName = "market"
	ServerAddr = "127.0.0.1:8082"
	MID        = 1
)

type Mysql struct {
	DataSource string `yaml:"DataSource"`
}

type CacheRedis struct {
	Host string `yaml:"Host"`
	Type int    `yaml:"Type"`
	Pass string `yaml:"Pass"`
}

type Mongo struct {
	Database string `yaml:"Database"`
	Url      string `yaml:"Url"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

type config struct {
	Mysql      Mysql      `yaml:"Mysql"`
	CacheRedis CacheRedis `yaml:"CacheRedis"`
	Mongo      Mongo      `yaml:"Mongo"`
}

var conf config

func GetConf() *config {
	return &conf
}
