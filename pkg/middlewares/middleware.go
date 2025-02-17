package middlewares

import (
	"context"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/response"
	"github.com/oaxacos/vitacare/pkg/utils"
	"net/http"
)

func AuthMiddleware(config *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do something
			log := logger.GetContextLogger(r.Context())
			authorizationToken := utils.GetAuthorizationToken(r)
			log.Debugf("authorization token: %s", authorizationToken)
			if authorizationToken == "" {
				log.Debugf("authorization token is empty")
				response.RenderUnauthorized(w)
				return
			}
			log.Infof("config %v", config.Token)
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
