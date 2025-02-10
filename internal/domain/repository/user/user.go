package userRepository

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

func (u *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
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

func (u *UserRepo) savePassword(ctx context.Context, password *model.Password) error {
	q := u.DB.NewInsert().Model(password)
	_, err := q.Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) Create(ctx context.Context, user *model.User) error {
	tx, err := u.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = u.Save(ctx, user)
	if err != nil {
		return err
	}
	err = u.savePassword(ctx, user.Password)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) GetPassword(ctx context.Context, user *model.User) (*model.Password, error) {
	password := new(model.Password)
	q := u.DB.NewSelect().Model(password).Where("user_id = ?", user.ID)
	err := q.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return password, nil
}
