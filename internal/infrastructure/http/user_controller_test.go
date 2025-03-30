package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/internal/domain/repository/password"
	tokenRepository "github.com/oaxacos/vitacare/internal/domain/repository/token"
	userRepository "github.com/oaxacos/vitacare/internal/domain/repository/user"
	"github.com/oaxacos/vitacare/internal/domain/service/token"
	"github.com/oaxacos/vitacare/internal/domain/service/user"
	"github.com/oaxacos/vitacare/internal/infrastructure/db"
	"github.com/oaxacos/vitacare/pkg/logger"
	"github.com/oaxacos/vitacare/pkg/server"
	"github.com/oaxacos/vitacare/pkg/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserController(t *testing.T) {

	configTest, err := config.NewConfig("test")
	if err != nil {
		t.Fatalf("error loading config %v", err)
	}

	repoDb, err := db.NewConnection(configTest)
	if err != nil {
		t.Fatalf("error creating db connection %v", err)
	}
	passRepo := password.NewPasswordRepository(repoDb)
	tokenRepo := tokenRepository.NewTokenRepository(repoDb)
	userRepo := userRepository.NewUserRepository(repoDb)

	userService := user.NewUserService(userRepo, passRepo)
	tokenSvc := token.NewTokenService(configTest, tokenRepo)

	s := server.NewServer(configTest)

	NewUserController(s, userService, tokenSvc)

	t.Run("login a user", func(t *testing.T) {
		data := map[string]interface{}{
			"first_name":            "Jose",
			"last_name":             "Ruiz",
			"email":                 "test@test.com",
			"password":              "supersecret",
			"password_confirmation": "supersecret",
		}

		log := logger.GetGlobalLogger()

		out, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/api/v0/users/auth/register", bytes.NewBuffer(out))
		if err != nil {
			t.Fatalf("error creating request %v", err)
		}

		response := executeRequest(req, s)

		assert.NotNil(t, response)
		assert.Equal(t, http.StatusCreated, response.Code)
		fmt.Printf("response: %v", response.Body.String())
		// get cookie
		cookie := response.Header().Get("refresh_token")
		assert.NotEmpty(t, cookie)

		// make again the request for the same user
		response = executeRequest(req, s)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.Code)

		utils.CleanUpDB(repoDb.DB, context.TODO(), []string{"users", "user_passwords", "tokens"})
	})

}

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *server.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	//s.Router.ServeHTTP(rr, req)
	s.Mux.ServeHTTP(rr, req)

	return rr
}
