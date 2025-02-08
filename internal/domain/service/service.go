package service

import (
	"context"
	"github.com/oaxacos/vitacare/internal/domain/model"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, user *model.UserRole) (string, error)
	GenerateRefreshToken(ctx context.Context, user *model.UserRole) (string, error)
	VerifyAccessToken(ctx context.Context, token string) (any, error)
	VerifyRefreshToken(ctx context.Context, token string) (any, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
}
