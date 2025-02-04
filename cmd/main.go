package main

import (
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/internal/domain/repository/user"
	userService "github.com/oaxacos/vitacare/internal/domain/service/user"
	"github.com/oaxacos/vitacare/internal/infrastructure/db"
	"github.com/oaxacos/vitacare/internal/infrastructure/http"
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

	connection, err := db.NewConnection(conf)
	if err != nil {
		logs.Fatalf("failed to connect to database: %v", err)
	}

	userRepo := user.NewUserRepository(connection.DB)

	userSvc := userService.NewUserService(userRepo)

	s := server.NewServer(conf)

	http.NewUserController(s.Mux, userSvc)

	err = s.Start()
	if err != nil {
		logs.Error(err)
	}
}
