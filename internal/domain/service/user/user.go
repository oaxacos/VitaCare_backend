package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/application/dto"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/internal/domain/repository"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/uptrace/bun"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrNoUserWithEmail  = errors.New("invalid credentials")
	ErrNoUserWithID     = errors.New("user not found")
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

func (u *UserService) CreateUser(ctx context.Context, user dto.UserDto) (*model.User, error) {
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
		return nil, err
	}
	return newUser, nil
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

func (u *UserService) LoginUser(ctx context.Context, data dto.UserLoginDto) (*model.User, error) {
	user, err := u.UserRepo.GetByEmail(ctx, data.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoUserWithEmail
		}
		return nil, err
	}
	err = u.PasswordRepo.VerifyPasswordText(ctx, user.ID, data.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := u.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoUserWithID
		}
		return nil, err
	}
	return user, nil
}

func (u *UserService) UpdateUserRole(ctx context.Context, id uuid.UUID, role string) error {

	user, err := u.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoUserWithID
		}
		return err
	}
	err = user.UpdateRole(role)
	if err != nil {
		return err
	}
	logger.GetContextLogger(ctx).Infof("user %s updated to role %s", user.Email, user.Rol)
	err = u.UserRepo.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateUserInfo(ctx context.Context, id uuid.UUID, data dto.UpdateUserDto) error {
	if data.FirstName == "" && data.LastName == "" && data.Dni == "" && data.Phone == "" && data.BirthDate == "" {
		return errors.New("no data to update")
	}
	logger.GetContextLogger(ctx).Infof("updating user %s", id)

	user, err := u.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoUserWithID
		}
		return err
	}

	if data.FirstName != "" {
		user.FirstName = data.FirstName
	}
	if data.LastName != "" {
		user.LastName = data.LastName
	}
	if data.Dni != "" {
		user.DNI = data.Dni
	}
	if data.Phone != "" {
		user.Phone = data.Phone
	}
	//TODO: find a way to handle birthdate, can be iso or yy-mm-dd, or dd-mm-yy

	return u.UserRepo.Update(user)
}
