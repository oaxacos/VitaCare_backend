package config

import (
	"errors"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type server struct {
	Port   int  `koanf:"port"`
	Debug  bool `koanf:"debug"`
	Pretty bool `koanf:"pretty"`
}

type database struct {
	DbName   string `koanf:"dbname"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	Username string `koanf:"username"`
}

var k = koanf.New(".")
var configFilePath = "config/config.yaml"

type Config struct {
	Server   server   `koanf:"server"`
	Database database `koanf:"database"`
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
