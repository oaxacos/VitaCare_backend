package config

import (
	"errors"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Server struct {
	Port   int  `koanf:"port"`
	Debug  bool `koanf:"debug"`
	Pretty bool `koanf:"pretty"`
}

type Cors struct {
	TrustedOrigins []string `koanf:"trusted-origins"`
}

type Database struct {
	DbName   string `koanf:"dbname"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	Username string `koanf:"username"`
}

var k = koanf.New(".")
var configFilePath = "config/config.yaml"

type Config struct {
	Server   Server   `koanf:"Server"`
	Database Database `koanf:"Database"`
	Cors     Cors     `koanf:"Cors"`
}

var errConfigEmpty = errors.New("config file is empty")

func NewConfig() (*Config, error) {
	err := k.Load(file.Provider(configFilePath), yaml.Parser())
	if err != nil {
		return nil, err
	}
	var conf *Config
	err = k.Unmarshal("", &conf)

	if err != nil {
		return nil, err
	}
	if conf == nil {
		return nil, errConfigEmpty
	}
	return conf, nil
}
