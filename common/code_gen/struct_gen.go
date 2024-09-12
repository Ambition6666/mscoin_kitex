package code_gen

type Mysql struct {
	DataSource string `yaml:"DataSource"`
}

type CacheRedis struct {
	Host string `yaml:"Host"`
	Type string `yaml:"Type"`
	Pass string `yaml:"Pass"`
}

type Rocketmq struct {
	Addr string `yaml:"Addr"`
	WriteCap int `yaml:"WriteCap"`
	ReadCap int `yaml:"ReadCap"`
}

type config struct {
	Mysql Mysql `yaml:"Mysql"`
	CacheRedis CacheRedis `yaml:"CacheRedis"`
	Rocketmq Rocketmq `yaml:"Rocketmq"`
}
