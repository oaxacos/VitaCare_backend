package server

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer()
	var okMessage = map[string]string{
		"status": "ok",
	}

	req, err := http.NewRequest("GET", "/api/v0/healthcheck", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	s.ServeHTTP(recorder, req)
	recivedMessage := make(map[string]string)
	assert.Equal(t, http.StatusOK, recorder.Code)
	err = json.Unmarshal(recorder.Body.Bytes(), &recivedMessage)
	assert.NoError(t, err)
	assert.Equal(t, okMessage, recivedMessage)
}

func TestServerNotFound(t *testing.T) {
	s := NewServer()
	var notFound = map[string]string{
		"error": "page not found",
	}
	req, err := http.NewRequest("GET", "/api/v0/healthcheck/test", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	s.ServeHTTP(recorder, req)
	receivedMessage := make(map[string]string)
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	err = json.Unmarshal(recorder.Body.Bytes(), &receivedMessage)
	assert.NoError(t, err)
	assert.Equal(t, notFound, receivedMessage)
}
