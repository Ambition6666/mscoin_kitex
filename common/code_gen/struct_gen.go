package code_gen

type Mongo struct {
	Url string `yaml:"Url"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Database string `yaml:"Database"`
}

type Kafka struct {
	Addr string `yaml:"Addr"`
	WriteCap int `yaml:"WriteCap"`
	ReadCap int `yaml:"ReadCap"`
}

type CacheRedis struct {
	Host string `yaml:"Host"`
	Type string `yaml:"Type"`
	Pass string `yaml:"Pass"`
}

type Okx struct {
	Apikey string `yaml:"Apikey"`
	SecretKey string `yaml:"SecretKey"`
	Pass string `yaml:"Pass"`
	Host string `yaml:"Host"`
	Proxy string `yaml:"Proxy"`
}

type config struct {
	Mongo Mongo `yaml:"Mongo"`
	Kafka Kafka `yaml:"Kafka"`
	CacheRedis CacheRedis `yaml:"CacheRedis"`
	Okx Okx `yaml:"Okx"`
}
