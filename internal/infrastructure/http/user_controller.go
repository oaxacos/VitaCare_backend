package http

import (
	"net/http"

	"github.com/oaxacos/vitacare/internal/domain/service/token"
	"github.com/oaxacos/vitacare/internal/domain/service/user"

	"github.com/go-chi/chi/v5"
	"github.com/oaxacos/vitacare/internal/application/dto"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/response"
	"github.com/oaxacos/vitacare/pkg/utils"
	"github.com/oaxacos/vitacare/pkg/validator"
)

type UserController struct {
	userService  *user.UserService
	tokenService *token.TokenService
	c            *chi.Mux
}

const prefix = "/api/v0"

func NewUserController(c *chi.Mux, userSvc *user.UserService, tokenSvc *token.TokenService) {
	userController := &UserController{
		c:            c,
		userService:  userSvc,
		tokenService: tokenSvc,
	}
	c.Route(prefix, func(r chi.Router) {
		r.Post("/users/auth/register", userController.handleRegisterUser)
		r.Post("/users/auth/login", userController.handleLogin)
		r.Post("/users/auth/renew", userController.handleRenewToken)
	})

}

func (u *UserController) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	// create user
	var userData dto.UserDto
	err := utils.ReadFromRequest(r, &userData)
	log := logger.GetContextLogger(r.Context())

	if err != nil {
		response.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}
	validator.NewValidator()

	err = validator.Validate(userData)
	if err != nil {
		response.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}
	// create user
	newUser := model.NewPatientUser(userData)
	err = u.userService.CreateUser(r.Context(), newUser)
	if err != nil {
		response.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Debug("User created successfully")
	// create token
	refreshToken, err := u.tokenService.GenerateRefreshToken(r.Context(), newUser)
	if err != nil {
		response.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Debug("Refresh token created successfully")
	accessToken, err := u.tokenService.GenerateAccessToken(r.Context(), newUser)
	if err != nil {
		response.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dataResponse := dto.UserLoggedInDto{
		AccessToken:  accessToken,
		UserID:       newUser.ID.String(),
		RefreshToken: refreshToken,
	}

	err = response.WriteJsonResponse(w, dataResponse, http.StatusOK)

	if err != nil {
		response.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (u *UserController) handleLogin(w http.ResponseWriter, r *http.Request) {
	var loginData dto.UserLoginDto
	err := utils.ReadFromRequest(r, &loginData)
	log := logger.GetContextLogger(r.Context())
	if err != nil {
		response.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}
	validator.NewValidator()
	err = validator.Validate(loginData)
	if err != nil {
		response.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}
	userFind, err := u.userService.Login(r.Context(), loginData)
	log.Debugf("user: %v", userFind)
	if err != nil {
		response.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// create token
	accessToken, refreshToken, err := u.tokenService.RenewTokens(r.Context(), userFind)
	if err != nil {
		response.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dataResponse := dto.UserLoggedInDto{
		AccessToken:  accessToken,
		UserID:       userFind.ID.String(),
		RefreshToken: refreshToken,
	}

	response.WriteJsonResponse(w, dataResponse, http.StatusOK)
}

func (u *UserController) handleRenewToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenHeader := r.Header.Get("Authorization")
	log := logger.GetContextLogger(r.Context())
	if refreshTokenHeader == "" {
		log.Error("refresh token is required")
		response.RenderError(w, http.StatusBadRequest, "refresh token is required")
		return
	}
	tokenValidated, err := u.tokenService.VerifyAccessToken(r.Context(), refreshTokenHeader)
	if err != nil {
		log.Error(err)
		response.RenderError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	userfind, err := u.userService.GetByID(r.Context(), tokenValidated.UserID)
	if err != nil {
		log.Error(err)
		response.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := u.tokenService.RenewTokens(r.Context(), userfind)
	if err != nil {
		response.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dataResponse := dto.UserLoggedInDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.WriteJsonResponse(w, dataResponse, http.StatusOK)
}
