package service

import (
	"context"
	"github.com/oaxacos/vitacare/internal/domain/service/token"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, user token.DommyUser) (string, error)
	GenerateRefreshToken(ctx context.Context, user token.DommyUser) (string, error)
	VerifyAccessToken(ctx context.Context, token string) (any, error)
	VerifyRefreshToken(ctx context.Context, token string) (any, error)
}
