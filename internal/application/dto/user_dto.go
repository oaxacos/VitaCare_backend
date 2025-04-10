package dto

import (
	"github.com/google/uuid"
)

type UserDto struct {
	FirstName            string `json:"first_name" validate:"required,min=3"`
	LastName             string `json:"last_name" validate:"required,min=3"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password,min=6"`
}

type UserLoggedInDto struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type UserLoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type UpdateUserRoleDto struct {
	Role string `json:"role" validate:"required,oneof=admin doctor patient secretary"`
}

type UpdateUserDto struct {
	FirstName string `json:"first_name" validate:"omitempty,min=3"`
	LastName  string `json:"last_name" validate:"omitempty,min=3"`
	Dni       string `json:"dni" validate:"omitempty,min=3"`
	Phone     string `json:"phone" validate:"omitempty,min=3"`
	BirthDate string `json:"birth_date" validate:"omitempty,min=3"` // YYYY-MM-DD
}
