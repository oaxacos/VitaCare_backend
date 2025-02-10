package main

import (
	"github.com/oaxacos/vitacare/internal/config"
	tokenRepository "github.com/oaxacos/vitacare/internal/domain/repository/token"
	userRepository "github.com/oaxacos/vitacare/internal/domain/repository/user"
	"github.com/oaxacos/vitacare/internal/domain/service/token"
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
    defer func ()  {
        logs.Info("closing connection")
        connection.Close()
    }()

	if err != nil {
		logs.Fatalf("failed to connect to database: %v", err)
	}

	userRepo := userRepository.NewUserRepository(connection.DB)
	tokenRepo := tokenRepository.NewTokenRepository(connection.DB)

	userSvc := userService.NewUserService(userRepo)
	tokenSvc := token.NewTokenService(conf, tokenRepo)

	s := server.NewServer(conf)

	http.NewUserController(s, userSvc, tokenSvc)

	err = s.Start()
	if err != nil {
		logs.Error(err)
	}
}
