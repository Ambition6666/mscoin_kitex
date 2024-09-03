package code_gen

type Mysql struct {
	DataSource string `yaml:"DataSource"`
}

type CacheRedis struct {
	Pass string `yaml:"Pass"`
	Host string `yaml:"Host"`
	Type string `yaml:"Type"`
}

type JWT struct {
	AccessSecret string `yaml:"AccessSecret"`
	AccessExpire int `yaml:"AccessExpire"`
}

type config struct {
	Mysql Mysql `yaml:"Mysql"`
	CacheRedis CacheRedis `yaml:"CacheRedis"`
	JWT JWT `yaml:"JWT"`
}
