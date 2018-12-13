package errors

import (
	"net/http"
)

type Response struct {
	Error ErrorObject `json:"error"`
}

type ErrorObject struct {
	Message    string              `json:"message"`
	Code       int                 `json:"code,omitempty"`
	Validation map[string][]string `json:"validation,omitempty"`
	Debug      string              `json:"debug,omitempty"`
}

func getErrCode(et error) (errCode int) {
	switch et.Error() {
	case NotFound.Error():
		errCode = http.StatusNotFound
	case NoMethod.Error():
		errCode = http.StatusMethodNotAllowed
	case ServerError.Error():
		errCode = http.StatusInternalServerError
	case RecordNotFound.Error():
		errCode = http.StatusNotFound
	default:
		errCode = http.StatusBadRequest
	}
	return
}
