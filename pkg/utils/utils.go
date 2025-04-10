package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrorBodyReadAfterClose = errors.New("request body closed")
	ErrorEmptyRequestBody   = errors.New("empty request body")
	ErrorInvalidBodyRequest = errors.New("invalid request body")
)

func ReadFromRequest(r *http.Request, data any) error {
	// read from request
	err := json.NewDecoder(r.Body).Decode(&data)
	defer func() {
		_ = r.Body.Close()
	}()

	if err != nil {
		if errors.Is(err, http.ErrBodyReadAfterClose) {
			return ErrorBodyReadAfterClose
		}
		if _, ok := err.(*json.SyntaxError); ok {
			return ErrorInvalidBodyRequest
		}

		if err.Error() == "EOF" {
			return ErrorEmptyRequestBody
		}
	}
	return nil
}
