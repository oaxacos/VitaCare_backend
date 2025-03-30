package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/application/dto"
	"github.com/oaxacos/vitacare/internal/domain/model"
)

type TokenService interface {
	GenerateToken(ctx context.Context, user *model.User) (string, string, error)
	ValidateRefreshToken(ctx context.Context, refreshToken string) (*model.RefreshToken, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user dto.UserDto) (*model.User, error)
	ExistUser(ctx context.Context, email string) error
	LoginUser(ctx context.Context, data dto.UserLoginDto) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}
