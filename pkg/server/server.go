package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

var defaultPort = ":8080"

type Server struct {
	*chi.Mux
}

func handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func NewServer() *Server {
	logs := logger.GetGlobalLogger()
	r := chi.NewRouter()
	r.Use(loggerMiddleware(logs))
	r.Get("/v0/api/healthcheck", handleHealthcheck)

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
