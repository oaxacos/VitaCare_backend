package token

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/config"
	repository "github.com/oaxacos/vitacare/internal/domain/repository/token"
	"github.com/oaxacos/vitacare/internal/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTokenService(t *testing.T) {
	conf, err := config.NewConfigTest()
	assert.NoError(t, err)
	fmt.Printf("config: %+v\n", conf.Database)
	DB, err := db.NewConnection(conf)
	assert.NoError(t, err)
	repo := repository.NewTokenRepository(DB.DB)

	tokenService := NewTokenService(conf, repo)

	userTest := DommyUser{
		ID:    uuid.MustParse("f3b27c73-c844-476e-917f-56278406579f"),
		Email: "test@test.com",
		Role:  "admin",
	}
	ctx := context.Background()

	t.Run("generate access token", func(t *testing.T) {
		token, err := tokenService.GenerateAccessToken(ctx, userTest)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generate refresh token", func(t *testing.T) {

		refreshToken, err := tokenService.GenerateRefreshToken(ctx, userTest)
		assert.NoError(t, err)
		assert.NotEmpty(t, refreshToken)

		tokenInDB, err := tokenService.repo.GetByToken(ctx, refreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, tokenInDB)

		err = tokenService.VerifyRefreshToken(ctx, refreshToken)
		assert.NoError(t, err)

	})

}

func TestTokenExpiration(t *testing.T) {
	var token = "KjIZ52lspWj0xBHp7zYQxizcCHGT9w5tRsRRUZ8ViGQ="
	conf, err := config.NewConfigTest()
	assert.NoError(t, err)
	fmt.Printf("config: %+v\n", conf.Database)
	DB, err := db.NewConnection(conf)
	assert.NoError(t, err)
	repo := repository.NewTokenRepository(DB.DB)

	tokenService := NewTokenService(conf, repo)
	err = tokenService.VerifyRefreshToken(context.Background(), token)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidToken)
}
