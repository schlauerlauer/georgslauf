package models

type APIConfig struct {
	Server struct {
		Port	string	`yaml:"port" binding:"required"`
		Secret	string	`yaml:"secret"`
	} `yaml:"server"`
	Database struct {
		Mysql MysqlConfig `yaml:"mysql" binding:"required"`
	} `yaml:"database"`
}

type MysqlConfig struct {
	Username	string	`yaml:"username"`
	Password	string	`yaml:"password"`
	Hostname	string	`yaml:"hostname"`
	Port		string	`yaml:"port"`
	Database	string	`yaml:"database"`
}