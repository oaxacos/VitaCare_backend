package user

import (
	"context"
	"github.com/google/uuid"

	"github.com/oaxacos/vitacare/internal/application/dto"
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
	err = u.UserRepo.Create(ctx, user)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := u.UserRepo.GetByEmail(ctx, email)
	log := logger.GetContextLogger(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return user, nil
}

func (u *UserService) Login(ctx context.Context, data dto.UserLoginDto) (*model.User, error) {
	user, err := u.UserRepo.GetByEmail(ctx, data.Email)
	log := logger.GetContextLogger(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	userPassword, err := u.UserRepo.GetPassword(ctx, user)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = user.Password.VerifyPassword(data.Password, userPassword.Hash)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// validate user password
	err = user.Password.VerifyPassword(data.Password, userPassword.Hash)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return user, nil

}

func (u *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := u.UserRepo.GetByID(ctx, id)
	log := logger.GetContextLogger(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return user, nil
}
