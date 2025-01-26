package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type RefreshToken struct {
	bun.BaseModel `bun:"tokens,alias:tokens"`
	ID            uuid.UUID `bun:"id,pk"`
	Token         string    `bun:"token"`
	UserID        uuid.UUID `bun:"user_id"`
	CreatedAt     time.Time `bun:"created_at"`
	ExpiredAt     time.Time `bun:"expired_at"`
}

func NewRefreshToken(token string, userID uuid.UUID, expiredAt time.Duration) *RefreshToken {
	return &RefreshToken{
		ID:        uuid.New(),
		Token:     token,
		UserID:    userID,
		ExpiredAt: time.Now().Add(expiredAt),
		CreatedAt: time.Now(),
	}
}
