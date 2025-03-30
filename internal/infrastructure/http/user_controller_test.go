package http

import (
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/internal/domain/repository/password"
	tokenRepository "github.com/oaxacos/vitacare/internal/domain/repository/token"
	userRepository "github.com/oaxacos/vitacare/internal/domain/repository/user"
	"github.com/oaxacos/vitacare/internal/domain/service/token"
	"github.com/oaxacos/vitacare/internal/domain/service/user"
	"github.com/oaxacos/vitacare/internal/infrastructure/db"
	"github.com/oaxacos/vitacare/pkg/server"
	"testing"
)

func TestUserController(t *testing.T) {
	configTest, err := config.NewConfigTest()
	if err != nil {
		t.Fatalf("error loading config %v", err)
	}

	repoDb, err := db.NewConnection(configTest)
	if err != nil {
		t.Fatalf("error creating db connection %v", err)
	}

	passRepo := password.NewPasswordRepository(repoDb)
	tokenRepo := tokenRepository.NewTokenRepository(repoDb)
	userRepo := userRepository.NewUserRepository(repoDb)

	userService := user.NewUserService(userRepo, passRepo)
	tokenSvc := token.NewTokenService(configTest, tokenRepo)

	s := server.NewServer(configTest)

	NewUserController(s, userService, tokenSvc)

	err = s.Start()
	if err != nil {
		t.Fatal(err)
	}

}
