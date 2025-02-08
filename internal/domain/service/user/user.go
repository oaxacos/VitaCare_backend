package user

import (
	"context"

	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/internal/domain/repository"
	"github.com/oaxacos/vitacare/pkg/logger"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, user *model.User) error {
	err := u.UserRepo.AlreadyExist(ctx, user.Email)
	log := logger.GetContextLogger(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	err = u.UserRepo.Save(ctx, user)

	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
