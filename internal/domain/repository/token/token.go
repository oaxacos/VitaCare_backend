package tokenRepository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/internal/infrastructure/db"
	"github.com/oaxacos/vitacare/pkg/logger"
)

type RefreshTokenRepo struct {
	DB *db.DBRepository
}

func NewTokenRepository(db *db.DBRepository) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		DB: db,
	}
}

func (t *RefreshTokenRepo) Save(ctx context.Context, token *model.RefreshToken) error {
	_, err := t.DB.NewInsert().Model(token).Exec(ctx)
	return err
}

func (t *RefreshTokenRepo) Delete(ctx context.Context, token uuid.UUID) error {
	_, err := t.DB.NewDelete().Model((*model.RefreshToken)(nil)).Where("id = ?", token).Exec(ctx)
	return err
}

func (t *RefreshTokenRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.RefreshToken, error) {
	token := new(model.RefreshToken)
	q := t.DB.NewSelect().Model(token).Where("user_id = ?", userID)
	logger.GetContextLogger(ctx).Debugf("query: %s", q.String())
	err := q.Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return token, nil
}

func (t *RefreshTokenRepo) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	refreshToken := new(model.RefreshToken)
	q := t.DB.NewSelect().Model(refreshToken).Where("token = ?", token)
	err := q.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}
