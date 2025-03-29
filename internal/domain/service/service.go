package service

import (
	"context"
	"github.com/oaxacos/vitacare/internal/application/dto"
	"github.com/oaxacos/vitacare/internal/domain/service/token"

	"github.com/oaxacos/vitacare/internal/domain/model"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, user *model.UserRole) (string, error)
	GenerateRefreshToken(ctx context.Context, user *model.UserRole) (string, error)
	VerifyAccessToken(ctx context.Context, token string) (*token.AccessTokenClaims, error)
	VerifyRefreshToken(ctx context.Context, token string) (any, error)
	RenewTokens(ctx context.Context, user *model.User) (string, string, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user dto.UserDto) error
	ExistUser(ctx context.Context, email string) error
	//GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	LoginUser(ctx context.Context, data dto.UserLoginDto) (*model.User, error)
	//GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}
