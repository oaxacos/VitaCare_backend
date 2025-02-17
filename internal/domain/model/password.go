package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	bun.BaseModel `bun:"user_passwords,alias:password"`
	PlainText     string    `bun:"-"`
	Hash          []byte    `bun:"password"`
	UserID        uuid.UUID `bun:"user_id" `
	ID            uuid.UUID `bun:"id,pk"`
	CreatedAt     time.Time `bun:"created_at"`
}

func NewPassword(plainText string, userID uuid.UUID) *Password {
	return &Password{
		PlainText: plainText,
		UserID:    userID,
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}
}

func (p *Password) SetHash() error {
	password := []byte(p.PlainText)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	p.Hash = hashedPassword
	return nil
}

func (p *Password) VerifyPassword(password string, hashedPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return err
	}
	return nil
}
