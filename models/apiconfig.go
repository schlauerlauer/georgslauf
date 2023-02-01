package models

type APIConfig struct {
    Server struct {
        Port    string  `yaml:"port" binding:"required"`
        Secret  string  `yaml:"secret" binding:"required"`
        Metrics struct {
            Username    string  `yaml:"username"`
            Password    string  `yaml:"password" binding:"required"`
        } `yaml:"metrics"`
    } `yaml:"server"`
    Database struct {
        Postgresql SqlConfig `yaml:"postgresql"`
    } `yaml:"database"`
}

type SqlConfig struct {
    Username    string  `yaml:"username" binding:"required"`
    Password    string  `yaml:"password" binding:"required"`
    Hostname    string  `yaml:"hostname" binding:"required"`
    Port        string  `yaml:"port" binding:"required"`
    Database    string  `yaml:"database" binding:"required"`
    SSL         string  `yaml:"ssl" binding:"required"`
    TZ          string  `yaml:"tz" binding:"required"`
}
