package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/uptrace/bun"
)

type RefreshTokenRepo struct {
	DB *bun.DB
}

func NewTokenRepository(db *bun.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		DB: db,
	}
}

func (t *RefreshTokenRepo) Save(ctx context.Context, token *model.RefreshToken) error {
	q := t.DB.NewInsert().Model(token)
	logger.GetContextLogger(ctx).Debugf("query: %s", q.String())
	_, err := q.Exec(ctx)
	return err
}

func (t *RefreshTokenRepo) Delete(ctx context.Context, token uuid.UUID) error {
	_, err := t.DB.NewDelete().Model((*model.RefreshToken)(nil)).Where("id = ?", token).Exec(ctx)
	return err
}

func (t *RefreshTokenRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.RefreshToken, error) {
	var token *model.RefreshToken
	err := t.DB.NewSelect().Model(token).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (t *RefreshTokenRepo) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	refreshToken := new(model.RefreshToken)

	q := t.DB.NewSelect().Model(refreshToken).Where("token = ?", token)
	logger.GetContextLogger(ctx).Debugf("query: %s", q.String())
	err := q.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (t *RefreshTokenRepo) Update(ctx context.Context, token *model.RefreshToken) error {
	_, err := t.DB.NewUpdate().Model(token).Where("id = ?", token.ID).Exec(ctx)
	return err
}
