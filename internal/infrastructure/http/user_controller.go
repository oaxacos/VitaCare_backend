package http

import (
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/internal/domain/service/token"
	"github.com/oaxacos/vitacare/internal/domain/service/user"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/server"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oaxacos/vitacare/internal/application/dto"
	"github.com/oaxacos/vitacare/pkg/response"
	"github.com/oaxacos/vitacare/pkg/utils"
	"github.com/oaxacos/vitacare/pkg/validator"
)

type UserController struct {
	userService  *user.UserService
	tokenService *token.TokenSvc
	c            *chi.Mux
	Config       *config.Config
}

const prefix = "/api/v0"

func NewUserController(s *server.Server, userSvc *user.UserService, tokenSvc *token.TokenSvc) {
	userController := &UserController{
		c:            s.Mux,
		Config:       s.Config,
		userService:  userSvc,
		tokenService: tokenSvc,
	}
	userController.c.Route(prefix, func(r chi.Router) {
		r.Post("/users/auth/register", userController.handleRegisterUser)
		r.Post("/users/auth/login", userController.handleLogin)
		r.Group(func(r chi.Router) {
			r.Post("/users/auth/renew", userController.handleRenewToken)
		})
	})
}

func (u *UserController) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	// create user
	var userData dto.UserDto
	ctx := r.Context()
	err := utils.ReadFromRequest(r, &userData)
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
	err = u.userService.ExistUser(r.Context(), userData.Email)
	if err != nil {
		response.RenderFatalError(w, err)
		return
	}
	// create user
	newUser, err := u.userService.CreateUser(r.Context(), userData)
	if err != nil {
		response.RenderFatalError(w, err)
		return
	}
	// create refresh token and access token
	accessToken, refreshToken, err := u.tokenService.GenerateToken(ctx, newUser)
	if err != nil {
		response.RenderFatalError(w, err)
		return
	}
	dataResponse := dto.UserLoggedInDto{
		AccessToken:  accessToken,
		UserID:       newUser.ID.String(),
		RefreshToken: refreshToken,
	}

	response.RenderJson(w, dataResponse, http.StatusCreated)
}

func (u *UserController) handleLogin(w http.ResponseWriter, r *http.Request) {
	var loginData dto.UserLoginDto
	err := utils.ReadFromRequest(r, &loginData)
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
	userWithCredentials, err := u.userService.LoginUser(r.Context(), loginData)

	if err != nil {
		response.RenderFatalError(w, err)
		return
	}
	ctx := r.Context()
	// create refresh token and access token
	accessToken, refreshToken, err := u.tokenService.GenerateToken(ctx, userWithCredentials)
	if err != nil {
		response.RenderFatalError(w, err)
		return
	}
	dataResponse := dto.UserLoggedInDto{
		AccessToken:  accessToken,
		UserID:       userWithCredentials.ID.String(),
		RefreshToken: refreshToken,
	}

	response.WriteJsonResponse(w, dataResponse, http.StatusOK)
}

func (u *UserController) handleRenewToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetContextLogger(ctx)
	var refreshToken dto.TokenRefreshRequest
	err := utils.ReadFromRequest(r, &refreshToken)
	log.Debugf("refresh token: %s", refreshToken.RefreshToken)
	if err != nil {
		response.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	validator.NewValidator()
	err = validator.Validate(refreshToken)
	if err != nil {
		response.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}
	refreshTokenModel, err := u.tokenService.ValidateRefreshToken(ctx, refreshToken.RefreshToken)
	if err != nil {
		log.Error(err)
		response.RenderUnauthorized(w)
		return
	}
	userInDB, err := u.userService.GetByID(ctx, refreshTokenModel.UserID)
	if err != nil {
		log.Error(err)
		response.RenderFatalError(w, err)
		return
	}
	newAccessToken, err := u.tokenService.GenerateAccessToken(ctx, userInDB)

	resp := dto.TokenRefreshResponse{
		AccessToken: newAccessToken,
	}

	response.WriteJsonResponse(w, resp, http.StatusOK)
}
