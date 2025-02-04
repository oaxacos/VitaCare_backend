package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/uptrace/bun"
)

type UserRepo struct {
	DB *bun.DB
}

var (
	ErrorAlreadyExist = errors.New("user already exist")
)

func NewUserRepository(db *bun.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (u *UserRepo) Save(ctx context.Context, user *model.User) error {
	q := u.DB.NewInsert().Model(user)
	_, err := q.Exec(ctx)
	return err
}

func (u *UserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	user := new(model.User)

	q := u.DB.NewSelect().Model(user).Where("id = ?", id)
	err := q.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	q := u.DB.NewSelect().Model(user).Where("email = ?", email)
	err := q.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) AlreadyExist(ctx context.Context, email string) error {
	user := new(model.User)
	q := u.DB.NewSelect().Model(user).Where("email = ?", email)
	err := q.Scan(ctx)
	log := logger.GetContextLogger(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("user not found")
			return nil
		}
		return err
	}
	if user.ID == uuid.Nil {
		return ErrorAlreadyExist
	}

	return nil
}
