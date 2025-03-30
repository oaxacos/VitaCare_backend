package dto

type ErrorResponse struct {
	Error ErrorDto `json:"error"`
}

type ErrorDto struct {
	Message string `json:"message"`
	Code    int    `json:"status_code"`
}

func MapResponseError(msg any, code int) ErrorResponse {
	message := "internal server error"
	if msg != nil {
		if err, ok := msg.(error); ok {
			message = err.Error()
		} else {
			message = msg.(string)
		}
	}
	return ErrorResponse{
		Error: ErrorDto{
			Message: message,
			Code:    code,
		},
	}
}
