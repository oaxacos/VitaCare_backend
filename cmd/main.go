package main

import (
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/pkg/logger"
)

func main() {
	logs := logger.GetGlobalLogger()
	conf, err := config.NewConfig()
	if err != nil {
		logs.Fatalf(err.Error())
	}
	logs.Infof("Starting server on %d", conf.Server.Port)
}
