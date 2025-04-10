package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/response"
	"github.com/oaxacos/vitacare/pkg/utils"
)

func AuthMiddleware(config *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.GetContextLogger(r.Context())
			authorizationToken := utils.GetAuthorizationToken(r)
			if authorizationToken == "" {
				log.Debugf("authorization token is empty")
				response.RenderUnauthorized(w)
				return
			}
			claims, err := utils.VerifyAccessToken(authorizationToken, []byte(config.Token.PrivateKeyAccessToken))
			if err != nil {
				log.Errorf("error verifying access token: %s", err)
				response.RenderUnauthorized(w)
				return
			}
			newContext := context.WithValue(r.Context(), utils.AuthorizationKey, claims)
			r = r.WithContext(newContext)
			log.Debugf("user is authenticated")
			next.ServeHTTP(w, r)
		})
	}
}

func AdminMiddleware(config *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.GetContextLogger(r.Context())
			log.Debugf("admin middleware")
			claims := utils.GetClaimsFromContext(r.Context())
			if claims == nil {
				log.Debugf("claims is nil")
				response.RenderUnauthorized(w)
				return
			}

			userId := claims.UserID
			if userId == uuid.Nil {
				log.Debugf("user id is nil")
				response.RenderUnauthorized(w)
				return
			}

			if claims.Rol != model.AdminRole {
				log.Debugf("user is not admin, user role: %s", claims.Rol)
				response.RenderForbidden(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
