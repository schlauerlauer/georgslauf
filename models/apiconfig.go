package models

type APIConfig struct {
	Server struct {
		Port	string	`yaml:"port" binding:"required"`
		Secret	string	`yaml:"secret" binding:"required"`
		Metrics	struct {
			Username	string	`yaml:"username"`
			Password	string	`yaml:"password" binding:"required"`
		} `yaml:"metrics"`
	} `yaml:"server"`
	Database struct {
		Mariadb SqlConfig `yaml:"mariadb" binding:"required"`
	} `yaml:"database"`
}

type SqlConfig struct {
	Username	string	`yaml:"username"`
	Password	string	`yaml:"password"`
	Hostname	string	`yaml:"hostname"`
	Port		string	`yaml:"port"`
	Database	string	`yaml:"database"`
}