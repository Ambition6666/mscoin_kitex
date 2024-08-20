package code_gen

type Server struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
}

type Settings struct {
	Timeout int `yaml:"timeout"`
	MaxConnections int `yaml:"max_connections"`
}

type Database struct {
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Settings Settings `yaml:"settings"`
}

type config struct {
	Server Server `yaml:"server"`
	Database Database `yaml:"database"`
}
