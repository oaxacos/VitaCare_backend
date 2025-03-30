package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/response"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

var defaultPort = ":8080"

type Server struct {
	*chi.Mux
	Config *config.Config
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

func NewServer(conf *config.Config) *Server {
	logs := logger.GetGlobalLogger()
	r := chi.NewRouter()
	r.Use(loggerMiddleware(logs))
	r.Use(enableCors(conf))

	r.Get("/api/v0/healthcheck", handleHealthcheck)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition
	))

	r.NotFound(handleNotFound)
	r.MethodNotAllowed(handleMethodNotAllowed)

	return &Server{
		r,
		conf,
	}
}

func (s *Server) Start() error {
	logs := logger.GetGlobalLogger()
	port := defaultPort
	portStr := strconv.Itoa(s.Config.Server.Port)
	if portStr != "" {
		port = ":" + portStr
	}
	logs.Infof("start server on %s", port)
	return http.ListenAndServe(port, s.Mux)
}

func (s *Server) StartWithTimeout(timeout time.Duration) error {
	logs := logger.GetGlobalLogger()
	port := defaultPort
	portStr := strconv.Itoa(s.Config.Server.Port)
	if portStr != "" {
		port = ":" + portStr
	}
	logs.Infof("start server on %s", port)
	logs.Infof("start timeout %d", timeout)

	s.Group(func(r chi.Router) {
		r.Use(middleware.Timeout(timeout))

	})

	return http.ListenAndServe(port, s.Mux)
}

func loggerMiddleware(log *zap.SugaredLogger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.SetContextLogger(r.Context(), log)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func enableCors(conf *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				for i := range conf.Cors.TrustedOrigins {
					if origin == conf.Cors.TrustedOrigins[i] {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
							w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, DELETE")
							w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
							w.WriteHeader(http.StatusOK)
						}
						break
					}
				}
			}
			next.ServeHTTP(w, r)
		})

	}
}
