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
	s := server.NewServer(conf)
	err = s.Start()
	if err != nil {
		logs.Error(err)
	}
}
