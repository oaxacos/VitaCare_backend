package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
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

func IsEmptyStruct[T any](data any) bool {
	if data == nil {
		return true
	}
	if str, ok := data.(string); ok {
		return str == ""
	}
	if _, ok := data.(int); ok {
		return data == 0
	}

	if _, ok := data.(float64); ok {
		return data == 0.0
	}

	//ou yeah
	newT := new(T)
	return reflect.DeepEqual(data, *newT)
}
