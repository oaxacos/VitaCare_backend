package service

import (
	"context"
	"github.com/oaxacos/vitacare/internal/application/dto"
	"github.com/oaxacos/vitacare/internal/domain/model"
)

type TokenService interface {
	GenerateToken(ctx context.Context, user *model.User) (string, string, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user dto.UserDto) (*model.User, error)
	ExistUser(ctx context.Context, email string) error
	LoginUser(ctx context.Context, data dto.UserLoginDto) (*model.User, error)
}
