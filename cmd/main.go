package main

import (
	_ "github.com/oaxacos/vitacare/docs"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/internal/domain/repository/password"
	tokenRepository "github.com/oaxacos/vitacare/internal/domain/repository/token"
	userRepository "github.com/oaxacos/vitacare/internal/domain/repository/user"
	"github.com/oaxacos/vitacare/internal/domain/service/token"
	"github.com/oaxacos/vitacare/internal/domain/service/user"
	"github.com/oaxacos/vitacare/internal/infrastructure/db"
	"github.com/oaxacos/vitacare/internal/infrastructure/http"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/server"
)

// @title VitaCare API
// @version 0.0
// @description This the service of Vitacare.
//@securityDefinitions.apikey Token
//@in header
//@name Authorization

func main() {
	logs := logger.GetGlobalLogger()
	defer logger.CloseLogger()

	conf, err := config.NewConfig("")
	if err != nil {
		logs.Fatal(err)
	}

	dbRepo, err := db.NewConnection(conf)
	if err != nil {
		logs.Fatal(err)
	}

	logs.Info("connected to database")

	defer func() {
		logs.Info("closing connection")
		dbRepo.Close()
	}()
	passRepo := password.NewPasswordRepository(dbRepo)
	userRepo := userRepository.NewUserRepository(dbRepo)
	tokenRepo := tokenRepository.NewTokenRepository(dbRepo)

	userSvc := user.NewUserService(userRepo, passRepo)
	tokenSvc := token.NewTokenService(conf, tokenRepo)

	s := server.NewServer(conf)

	http.NewUserController(s, userSvc, tokenSvc)

	err = s.Start()
	if err != nil {
		logs.Fatal(err)
	}
}
