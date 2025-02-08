package dto

type UserDto struct {
	FirstName            string `json:"first_name" validate:"required,min=3"`
	LastName             string `json:"last_name" validate:"required,min=3"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password,min=6"`
}

type UserLoggedInDto struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLoginDto struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}
