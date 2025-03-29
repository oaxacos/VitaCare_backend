package password

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/internal/infrastructure/db"
	"github.com/uptrace/bun"
)

type PasswordRepo struct {
	DB db.DBRepository
}

func NewPasswordRepository(db *db.DBRepository) *PasswordRepo {
	return &PasswordRepo{
		DB: *db,
	}
}

func (p *PasswordRepo) Save(ctx context.Context, tx *bun.Tx, password *model.Password) error {
	q := tx.NewInsert().Model(password)
	_, err := q.Exec(ctx)
	return err
}

func (p *PasswordRepo) getByUserID(ctx context.Context, userID uuid.UUID) (*model.Password, error) {
	password := new(model.Password)

	q := p.DB.NewSelect().Model(password).Where("user_id = ?", userID)
	err := q.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return password, nil
}

func (p *PasswordRepo) VerifyPasswordText(ctx context.Context, userId uuid.UUID, plainText string) error {
	password, err := p.getByUserID(ctx, userId)
	if err != nil {
		return err
	}
	if password == nil {
		return sql.ErrNoRows
	}
	return password.VerifyPassword(plainText, password.Hash)
}
