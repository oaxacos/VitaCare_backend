package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/application/dto"
	"github.com/uptrace/bun"
)

type UserRole string

var (
	ErrInvalidRole = errors.New("invalid role")
)

var (
	AdminRole     UserRole = "admin"
	DoctorRole    UserRole = "doctor"
	PatientRole   UserRole = "patient"
	SecretaryRole UserRole = "secretary"
)

type User struct {
	bun.BaseModel `bun:"users,alias:users"`
	ID            uuid.UUID    `bun:"id,pk"`
	Email         string       `bun:"email"`
	FirstName     string       `bun:"first_name"`
	LastName      string       `bun:"last_name"`
	Rol           UserRole     `bun:"rol"`
	DNI           string       `bun:"dni"`
	Birthdate     time.Time    `bun:"birthdate"`
	Phone         string       `bun:"phone"`
	IsActive      bool         `bun:"is_active"`
	DeceasedAt    sql.NullTime `bun:"deceased_at"`
	CreatedAt     time.Time    `bun:"created_at"`
	UpdateAt      time.Time    `bun:"update_at"`
	Password      *Password    `bun:"rel:has-one,join:id=user_id"`
}

func NewPatientUser(dto dto.UserDto) *User {
	user := &User{
		ID:        uuid.New(),
		Email:     dto.Email,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Rol:       PatientRole,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
	password := NewPassword(user.ID, dto.Password)
	_ = password.SetHash()
	user.Password = password

	return user
}

func (u *User) IsAdmin() bool {
	return u.Rol == AdminRole
}

func (u *User) UpdateRole(role string) error {
	switch role {
	case string(AdminRole):
		u.Rol = AdminRole
	case string(DoctorRole):
		u.Rol = DoctorRole
	case string(PatientRole):
		u.Rol = PatientRole
	case string(SecretaryRole):
		u.Rol = SecretaryRole
	default:
		return ErrInvalidRole
	}
	u.UpdateAt = time.Now()
	return nil
}
