package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/oaxacos/vitacare/internal/application/dto"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/internal/domain/repository"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/uptrace/bun"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
)

type UserService struct {
	UserRepo     repository.UserRepository
	PasswordRepo repository.PasswordRepository
}

func NewUserService(userRepo repository.UserRepository, passwordRepo repository.PasswordRepository) *UserService {
	return &UserService{
		UserRepo:     userRepo,
		PasswordRepo: passwordRepo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, user dto.UserDto) error {
	newUser := model.NewPatientUser(user)
	log := logger.GetContextLogger(ctx)
	//save user
	err := u.UserRepo.WithTransaction(ctx, func(tx *bun.Tx) error {
		// save user
		err := u.UserRepo.Save(ctx, tx, newUser)
		if err != nil {
			log.Error(err)
			return err
		}
		// save user password
		err = u.PasswordRepo.Save(ctx, tx, newUser.Password)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (u *UserService) ExistUser(ctx context.Context, email string) error {
	_, err := u.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return ErrUserAlreadyExist
}
