package userRepository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"

	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/internal/infrastructure/db"
)

type UserRepo struct {
	DB *db.DBRepository
}

var (
	ErrorAlreadyExist = errors.New("user already exist")
)

func NewUserRepository(db *db.DBRepository) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (u *UserRepo) Save(ctx context.Context, tx *bun.Tx, user *model.User) error {
	q := tx.NewInsert().Model(user)
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
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) AlreadyExist(ctx context.Context, email string) error {
	user := new(model.User)
	q := u.DB.NewSelect().Model(user).Where("email = ?", email)
	err := q.Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	if user.ID == uuid.Nil {
		return ErrorAlreadyExist
	}

	return nil
}

func (u *UserRepo) WithTransaction(ctx context.Context, fn func(tx *bun.Tx) error) error {
	return u.DB.WithTransaction(ctx, fn)
}

func (u *UserRepo) Update(user *model.User) error {
	user.UpdateAt = time.Now()
	q := u.DB.NewUpdate().Model(user).WherePK()
	_, err := q.Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}
