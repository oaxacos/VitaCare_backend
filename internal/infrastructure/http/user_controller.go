package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/oaxacos/vitacare/internal/application/dto"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/internal/domain/service/user"
	"github.com/oaxacos/vitacare/pkg/response"
	"github.com/oaxacos/vitacare/pkg/utils"
	"github.com/oaxacos/vitacare/pkg/validator"
	"net/http"
)

type UserController struct {
	userService *user.UserService
	c           *chi.Mux
}

const prefix = "/api/v0"

func NewUserController(c *chi.Mux, userSvc *user.UserService) {
	userController := &UserController{
		c:           c,
		userService: userSvc,
	}
	c.Route(prefix, func(r chi.Router) {
		r.Post("/users", userController.handleCreateUser)
	})

}

func (u *UserController) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// create user
	var userData dto.UserDto
	err := utils.ReadFromRequest(r, &userData)

	validator.NewValidator()

	if err != nil {
		response.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

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
	
	err = response.WriteJsonResponse(w, userData, http.StatusOK)
	if err != nil {
		response.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
