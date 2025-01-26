package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/domain/model"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *model.RefreshToken) error
	Delete(ctx context.Context, tokenID uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.RefreshToken, error)
	GetByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	Update(ctx context.Context, token *model.RefreshToken) error
}
