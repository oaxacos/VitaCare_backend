package config

import (
	"errors"
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
	"path/filepath"
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
	k              = koanf.New(".")
	kTest          = koanf.New(".")
	configFile     = "config.yaml"
	configTestFile = "config.test.yaml"
	configFolder   = "config"
)

type Config struct {
	Server   Server   `koanf:"Server"`
	Database Database `koanf:"Database"`
	Cors     Cors     `koanf:"Cors"`
	Token    Token    `koanf:"Token"`
}

// todo: find a better way to get the root directory
func getDirectoryFile(fileName string) (string, error) {
	rootDir := os.Getenv("PROJECT_ROOT")
	if rootDir == "" {
		return "", errors.New("PROJECT_ROOT is not set, please run `make shell`")
	}
	configPath := filepath.Join(rootDir, configFolder, fileName)
	return configPath, nil
}

type Token struct {
	PrivateKeyAccessToken  string `koanf:"access-token-key"`
	PrivateKeyRefreshToken string `koan:"refresh-token-key"`
	AccessTimeExpiration   int    `koanf:"access-time-expiration"`
	RefreshTimeExpiration  int    `koanf:"refresh-time-expiration"`
}

var errConfigEmpty = errors.New("config file is empty")

func NewConfig() (*Config, error) {
	filePath, err := getDirectoryFile(configFile)
	if err != nil {
		return nil, err
	}
	err = k.Load(file.Provider(filePath), yaml.Parser())
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

func NewConfigTest() (*Config, error) {
	filePath, err := getDirectoryFile(configTestFile)
	if err != nil {
		return nil, err
	}
	fmt.Printf("file path: %s\n", filePath)
	err = kTest.Load(file.Provider(filePath), yaml.Parser())
	if err != nil {
		return nil, err
	}
	var conf *Config
	err = kTest.Unmarshal("", &conf)

	if err != nil {
		return nil, err
	}
	if conf == nil {
		return nil, errConfigEmpty
	}
	return conf, nil
}
