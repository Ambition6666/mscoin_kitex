package code_gen

type CacheRedis struct {
	Host string `yaml:"Host"`
	Node int `yaml:"Node"`
	Pass string `yaml:"Pass"`
}

type Mongo struct {
	Url string `yaml:"Url"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Database string `yaml:"Database"`
}

type Rocketmq struct {
	Addr string `yaml:"Addr"`
	WriteCap int `yaml:"WriteCap"`
	ReadCap int `yaml:"ReadCap"`
}

type Mysql struct {
	DataSource string `yaml:"DataSource"`
}

type config struct {
	CacheRedis CacheRedis `yaml:"CacheRedis"`
	Mongo Mongo `yaml:"Mongo"`
	Rocketmq Rocketmq `yaml:"Rocketmq"`
	Mysql Mysql `yaml:"Mysql"`
}
