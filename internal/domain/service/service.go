package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/application/dto"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/internal/domain/service/token"
)

type TokenService interface {
	GenerateToken(ctx context.Context, user *model.User) (string, string, error)
	ValidateRefreshToken(ctx context.Context, refreshToken string) (*model.RefreshToken, error)
	VerifyAccessToken(ctx context.Context, token string) (*token.AccessTokenClaims, error)
	DeleteRefreshTokenByUser(userID uuid.UUID) error
	DeleteRefreshToken(token string) error
}

type UserService interface {
	CreateUser(ctx context.Context, user dto.UserDto) (*model.User, error)
	ExistUser(ctx context.Context, email string) error
	LoginUser(ctx context.Context, data dto.UserLoginDto) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateUserRole(ctx context.Context, id uuid.UUID, role string) error
	UpdateUserInfo(ctx context.Context, id uuid.UUID, data dto.UpdateUserDto) error
}
