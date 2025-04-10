package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/uptrace/bun"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *model.RefreshToken) error
	Delete(ctx context.Context, tokenID uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.RefreshToken, error)
	GetByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	DeleteByToken(ctx context.Context, token string) error
}

type UserRepository interface {
	Save(ctx context.Context, tx *bun.Tx, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	WithTransaction(ctx context.Context, fn func(tx *bun.Tx) error) error
	Update(user *model.User) error
}

type PasswordRepository interface {
	VerifyPasswordText(ctx context.Context, userId uuid.UUID, plainText string) error
	Save(ctx context.Context, tx *bun.Tx, password *model.Password) error
}
