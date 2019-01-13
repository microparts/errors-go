package errors

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gitlab.teamc.io/teamc.io/microservice/support/logs-go.git"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type ResponseObject struct {
	Error ErrorObject `json:"error"`
}

type ErrorObject struct {
	Message    string              `json:"message"`
	Code       int                 `json:"code,omitempty"`
	Validation map[string][]string `json:"validation,omitempty"`
	Debug      string              `json:"debug,omitempty"`
}

func Response(c *gin.Context, err interface{}) {
	var (
		jsonObj *ResponseObject
	)
	errCode := http.StatusBadRequest

	switch et := err.(type) {
	case gorm.Errors:
		if gorm.IsRecordNotFoundError(et) {
			errCode = http.StatusNotFound
		} else {
			errCode = http.StatusBadRequest
		}

		jsonObj = &ResponseObject{Error: ErrorObject{
			Message: et.Error(),
		}}

		logs.DBLogs.WithError(et).WithField("stage", "query").Error("Query error occurred")

	case validator.ValidationErrors:
		errCode = http.StatusUnprocessableEntity

		jsonObj = &ResponseObject{Error: ErrorObject{
			Message:    "Ошибка валидации",
			Validation: MakeErrorsSlice(et),
		}}

		logs.HttpLogs.
			WithError(et).
			WithFields(logrus.Fields{
				"validation": MakeErrorsSlice(et),
			}).
			Warn("validation error occurred")

	case error:
		errCode = getErrCode(et)

		jsonObj = &ResponseObject{Error: ErrorObject{Message: et.Error()}}

		logs.Log.
			WithError(et).
			WithFields(logrus.Fields{
				"message": et.Error(),
			}).
			Warn("error occurred")
	}

	c.AbortWithStatusJSON(errCode, jsonObj)
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
