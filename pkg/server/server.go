package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/response"
	"go.uber.org/zap"
	"net/http"
)

var defaultPort = ":8080"

type Server struct {
	*chi.Mux
}

func handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	response.RenderJson(w, map[string]string{
		"status": "ok",
	}, http.StatusOK)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	response.RenderNotFound(w)
}

func handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.RenderNotFound(w)
}

func NewServer() *Server {
	logs := logger.GetGlobalLogger()
	r := chi.NewRouter()
	r.Use(loggerMiddleware(logs))
	r.NotFound(handleNotFound)
	r.MethodNotAllowed(handleMethodNotAllowed)
	r.Get("/api/v0/healthcheck", handleHealthcheck)

	return &Server{
		r,
	}
}

func (s *Server) Start(conf *config.Config) error {
	logs := logger.GetGlobalLogger()
	port := defaultPort
	if conf.Server.Port != 0 {
		port = fmt.Sprintf(":%d", conf.Server.Port)
	}
	logs.Infof("start server on %s", port)
	return http.ListenAndServe(port, s)
}

func loggerMiddleware(log *zap.SugaredLogger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info(
				"method: ", r.Method,
				"	",
				"path: ", r.URL.Path,
			)
			ctx := logger.SetContextLogger(r.Context(), log)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
