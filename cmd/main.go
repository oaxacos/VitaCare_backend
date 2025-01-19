package main

import (
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/pkg/logger"
)

func main() {
	logs := logger.GetGlobalLogger()

	_, err := config.NewConfig()
	if err != nil {
		logs.Error(err)
	}

}
