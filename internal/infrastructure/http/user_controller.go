package http

import (
	"errors"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/internal/domain/service/token"
	"github.com/oaxacos/vitacare/internal/domain/service/user"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/middlewares"
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

const prefix = "/api/v0/users/auth"

func NewUserController(s *server.Server, userSvc *user.UserService, tokenSvc *token.TokenSvc) {
	userController := &UserController{
		c:            s.Mux,
		Config:       s.Config,
		userService:  userSvc,
		tokenService: tokenSvc,
	}
	userController.c.Route(prefix, func(r chi.Router) {
		r.Post("/register", userController.handleRegisterUser)
		r.Post("/login", userController.handleLogin)
		r.Post("/renew", userController.handleRenewToken)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(s.Config))
			r.Put("/logout", userController.handleLogout)
		})
	})
}

// @Router /api/v0/user/auth/register [post]
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags users
// @Success 200 {object} dto.UserLoggedInDto
// @Param user body dto.UserDto true "User data"
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
		RefreshToken: refreshToken,
		User: dto.User{
			ID:        newUser.ID,
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Email:     newUser.Email,
		},
	}
	response.SetRefreshTokenCookie(w, refreshToken)
	response.RenderJson(w, dataResponse, http.StatusCreated)
}

// @Router /api/v0/user/auth/login [post]
// @Summary login a user
// @Description login a user and set a cookie with the refresh token
// @Tags users
// @Success 200 {object} dto.UserLoggedInDto
// @Param user body dto.UserLoginDto true "User data"
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
		AccessToken: accessToken,
		User: dto.User{
			ID:        userWithCredentials.ID,
			FirstName: userWithCredentials.FirstName,
			LastName:  userWithCredentials.LastName,
			Email:     userWithCredentials.Email,
		},
	}
	response.SetRefreshTokenCookie(w, refreshToken)
	response.WriteJsonResponse(w, dataResponse, http.StatusOK)
}

// @Router /api/v0/user/auth/renew [post]
// @Summary renew access token
// @Description renew access token with refresh token
// @Tags users
// @Success 200 {object} dto.UserDto
// @Param user body dto.TokenRefreshRequest true "User data"
func (u *UserController) handleRenewToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetContextLogger(ctx)
	var refreshToken dto.TokenRefreshRequest
	err := utils.ReadFromRequest(r, &refreshToken)
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
		if errors.Is(err, token.ErrInvalidToken) {
			err := u.tokenService.DeleteRefreshToken(refreshToken.RefreshToken)
			log.Error(err)
			if err != nil {

				response.RenderFatalError(w, err)
				return
			}
			response.RenderUnauthorized(w)
			return
		}
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
		User: dto.User{
			ID:        userInDB.ID,
			FirstName: userInDB.FirstName,
			LastName:  userInDB.LastName,
			Email:     userInDB.Email,
		},
	}
	response.SetRefreshTokenCookie(w, refreshToken.RefreshToken)
	response.WriteJsonResponse(w, resp, http.StatusOK)
}

// @Router /api/v0/user/auth/logout [put]
// @Summary logout a user
// @Description logout a user and delete the refresh token
// @Tags users
// @Success 200 {object} string
func (u *UserController) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetContextLogger(ctx)
	claims := utils.GetClaimsFromContext(ctx)
	if claims == nil {
		response.RenderUnauthorized(w)
		return
	}
	err := u.tokenService.DeleteRefreshTokenByUser(claims.UserID)
	if err != nil {
		log.Error(err)
		response.RenderFatalError(w, err)
		return
	}
	response.DeleteRefreshTokenCookie(w)
	response.RenderJson(w, response.Envelop("message", "success"), http.StatusOK)
}
