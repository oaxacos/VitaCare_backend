package response

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/oaxacos/vitacare/pkg/logger"
	"net/http"
)

const (
	refreshTokenCookieName = "refresh_token"
)

func WriteJsonResponse(w http.ResponseWriter, data any, status int) error {
	d, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(d)
	return nil
}

func RenderJson(w http.ResponseWriter, data any, status int) {
	err := WriteJsonResponse(w, data, status)
	if err != nil {
		RenderServerError(w, err)
	}
}

func RenderError(w http.ResponseWriter, status int, message any) {
	res := map[string]any{
		"error": message,
	}
	log := logger.GetGlobalLogger()
	err := WriteJsonResponse(w, res, status)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Errorf("%s", message)
}

func RenderServerError(w http.ResponseWriter, err error) {
	logger.GetGlobalLogger().Error(err)
	message := "something went wrong"
	RenderError(w, http.StatusInternalServerError, message)
}

func RenderNotFound(w http.ResponseWriter) {
	message := "page not found"
	RenderError(w, http.StatusNotFound, message)
}

func RenderBadRequest(w http.ResponseWriter) {
	message := "bad request"
	RenderError(w, http.StatusBadRequest, message)
}

func RenderFatalError(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		RenderServerError(w, err)
	} else {
		RenderError(w, http.StatusInternalServerError, err.Error())
	}
}

func RenderUnauthorized(w http.ResponseWriter) {
	message := "unauthorized"
	RenderError(w, http.StatusUnauthorized, message)
}

type Cookie struct {
	Name     string
	Value    string
	MaxAge   int
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
}

func DeleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}

func SetRefreshTokenCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, cookie)
}

func DeleteRefreshTokenCookie(w http.ResponseWriter) {
	DeleteCookie(w, refreshTokenCookieName)
}

func Envelop(key string, data any) map[string]any {
	return map[string]any{key: data}
}
