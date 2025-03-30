package config

import (
	"errors"
	"fmt"
	"github.com/flashlabs/rootpath/location"
	"github.com/knadh/koanf/parsers/yaml"
	env2 "github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"strings"
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

var (
	k            = koanf.New(".")
	configFolder = "config"
)

type Config struct {
	Server   Server   `koanf:"server"`
	Database Database `koanf:"database"`
	Cors     Cors     `koanf:"cors"`
	Token    Token    `koanf:"token"`
}

func getConfigFile(env []string) (*koanf.Koanf, error) {

	k := koanf.New(".")
	confile := "config/config.yaml"
	if len(env) > 0 && env[0] != "" {
		confile = fmt.Sprintf("config/config.%s.yaml", env[0])
		if err := location.Chdir(); err != nil {
			panic(err)
		}
	}

	if err := k.Load(file.Provider(confile), yaml.Parser()); err != nil {
		return nil, err
	}

	return k, nil
}

type Token struct {
	PrivateKeyAccessToken  string `koanf:"access-token-key"`
	PrivateKeyRefreshToken string `koan:"refresh-token-key"`
	AccessTimeExpiration   int    `koanf:"access-time-expiration"`
	RefreshTimeExpiration  int    `koanf:"refresh-time-expiration"`
}

var errConfigEmpty = errors.New("config file is empty")

func NewConfig(env ...string) (*Config, error) {
	k, err := getConfigFile(env)
	if err != nil {
		return nil, err
	}

	// load from env, we can override the config file
	err = k.Load(env2.Provider("", ".", func(s string) string {
		val := strings.ToLower(strings.Replace(s, "_", ".", -1))
		return val
	}), nil)
	if err != nil {
		return nil, err
	}
	var conf *Config
	err = k.Unmarshal("", &conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
