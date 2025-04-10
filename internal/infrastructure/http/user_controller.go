package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/internal/domain/service/token"
	"github.com/oaxacos/vitacare/internal/domain/service/user"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/middlewares"
	"github.com/oaxacos/vitacare/pkg/server"

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
	validator    *validator.Validator
}

const prefix = "/api/v0/users"

func NewUserController(s *server.Server, userSvc *user.UserService, tokenSvc *token.TokenSvc, validator *validator.Validator) {
	userController := &UserController{
		c:            s.Mux,
		Config:       s.Config,
		userService:  userSvc,
		tokenService: tokenSvc,
		validator:    validator,
	}

	userController.c.Route(prefix, func(r chi.Router) {
		// authPath := fmt.Sprintf("%s/auth", prefix)
		// fmt.Printf("auth path: %s", authPath)
		r.Route("/auth", func(r chi.Router) {

			r.Post("/register", userController.handleRegisterUser)
			r.Post("/login", userController.handleLogin)
			r.Post("/renew", userController.handleRenewToken)
			r.Group(func(r chi.Router) {
				r.Use(middlewares.AuthMiddleware(s.Config))
				r.Put("/logout", userController.handleLogout)
			})

		})
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(s.Config), middlewares.AdminMiddleware(s.Config))
			r.Patch("/{id}/role", userController.handleUpdateUserRole)
			r.Patch("/", userController.handleUpdateUser)
		})
	})
}

// @Router /api/v0/users/auth/register [post]
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

	err = u.validator.ValidateStruct(userData)
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

// @Router /api/v0/users/auth/login [post]
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
	err = u.validator.ValidateStruct(loginData)
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

// @Router /api/v0/users/auth/renew [post]
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

	err = u.validator.ValidateStruct(refreshToken)
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

// @Router /api/v0/users/auth/logout [put]
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

// @Router /api/v0/users/{id}/role [patch]
// @Summary update user role
// @Security <YourTypeOfKey>
// @Description An admin can update the role of a user
// @Tags users
// @Security Token
// @Param id path string true "User ID"
// @Success 200 {object} string
func (u *UserController) handleUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// claims := utils.GetClaimsFromContext(ctx)
	userIDToUpdate := chi.URLParam(r, "id")

	var updateUserRole dto.UpdateUserRoleDto

	userIdParsed, err := uuid.Parse(userIDToUpdate)
	if err != nil {
		response.RenderError(w, http.StatusBadRequest, fmt.Sprintf("invalid user id: %s", userIDToUpdate))
		return
	}

	err = utils.ReadFromRequest(r, &updateUserRole)

	if err != nil {
		response.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = u.userService.UpdateUserRole(ctx, userIdParsed, updateUserRole.Role)
	if err != nil {
		response.RenderFatalError(w, err)
		return
	}
	resp := response.Envelop("message", "user role updated")

	response.RenderJson(w, resp, http.StatusOK)

}

// @Router /api/v0/users/ [patch]
// @Summary update user profile
// @Security <YourTypeOfKey>
// @Description Any user can update his profile, first name, last name, dni, phone and birthdate
// @Tags users
// @Security Token
// @Success 200 {object} string
// @Param user body dto.UpdateUserDto true "User data"
func (u *UserController) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaimsFromContext(r.Context())

	if claims == nil {
		response.RenderUnauthorized(w)
		return
	}

	userID := claims.UserID

	if userID == uuid.Nil {
		response.RenderUnauthorized(w)
		return
	}

	var updateUser dto.UpdateUserDto
	err := utils.ReadFromRequest(r, &updateUser)
	if err != nil {
		response.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = u.validator.ValidateStruct(updateUser)
	if err != nil {
		response.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	isEmpty := utils.IsEmptyStruct[dto.UpdateUserDto](updateUser)
	if isEmpty {
		response.RenderError(w, http.StatusBadRequest, "empty request body, at least one field is required")
		return
	}

	err = u.userService.UpdateUserInfo(r.Context(), userID, updateUser)

	if err != nil {
		response.RenderFatalError(w, err)
		return
	}

	resp := response.Envelop("user", updateUser)
	response.WriteJsonResponse(w, resp, http.StatusOK)
}
