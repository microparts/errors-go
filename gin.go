package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gitlab.teamc.io/teamc.io/microservice/support/logs-go.git"
	"gopkg.in/go-playground/validator.v9"
)

type ResponseObject struct {
	Error ErrorObject `json:"error,omitempty"`
}

type ErrorObject struct {
	Message    interface{}         `json:"message"`
	Code       int                 `json:"code,omitempty"`
	Validation map[string][]string `json:"validation,omitempty"`
	Debug      string              `json:"debug,omitempty"`
}

func Response(c *gin.Context, err interface{}) {
	errCode, data := MakeResponse(err)
	resp := ResponseObject{Error: *data}
	c.AbortWithStatusJSON(errCode, resp)
}

func MakeResponse(err interface{}) (int, *ErrorObject) {
	errObj := &ErrorObject{}
	errCode := http.StatusBadRequest

	switch et := err.(type) {
	case gorm.Errors:
		if gorm.IsRecordNotFoundError(et) {
			errCode = http.StatusNotFound
		} else {
			errCode = http.StatusBadRequest
		}

		errObj.Message = et.Error()

		logs.DBLogs.WithError(et).WithField("stage", "query").Error("Query error occurred")

	case validator.ValidationErrors:
		errCode = http.StatusUnprocessableEntity

		errObj.Message = "Ошибка валидации"
		errObj.Validation = MakeErrorsSlice(et)

		logs.HttpLogs.
			WithError(et).
			WithFields(logrus.Fields{
				"validation": MakeErrorsSlice(et),
			}).
			Warn("validation error occurred")

	case error:
		errCode = getErrCode(et)

		errObj.Message = et.Error()

		logs.Log.
			WithError(et).
			WithFields(logrus.Fields{
				"message": et.Error(),
			}).
			Warn("error occurred")
	case map[string]error:
		msgs := make(map[string]string)
		for k, e := range et {
			msgs[k] = e.Error()
		}

		errObj.Message = msgs

		logs.Log.
			WithFields(logrus.Fields{
				"errors": msgs,
			}).
			Warn("error occurred")
	}

	return errCode, errObj
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
