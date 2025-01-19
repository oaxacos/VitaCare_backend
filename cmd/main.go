package main

import (
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/server"
)

func main() {
	logs := logger.GetGlobalLogger()
	defer logger.CloseLogger()

	conf, err := config.NewConfig()
	if err != nil {
		logs.Error(err)
	}
	s := server.NewServer()
	err = s.Start(conf)
	if err != nil {
		logs.Error(err)
	}
}
