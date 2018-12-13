package errors

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

type AppError struct {
	// embedded builtin types are not exported, so alias to Err
	Err  error
	Code int
}

type ErrorCoder interface {
	ErrorCode() int
}

func (e AppError) ErrorCode() int {
	return e.Code
}

func (e AppError) Cause() error {
	return errors.Cause(e.Err)
}

const (
	NilErrorCode     = 0
	UnknownErrorCode = -1
	ek               = "common"
)

func ErrorCode(err error) int {
	if err == nil {
		return NilErrorCode
	}
	if ec, ok := err.(ErrorCoder); ok {
		return ec.ErrorCode()
	}

	return UnknownErrorCode
}

var CommonValidationErrors = map[string]string{
	ek:         "Ошибка валидации для свойства `%s` с правилом `%s`",
	"required": "Свойство `%s` обязательно для заполнения",
	"gt":       "Свойство `%s` должно содержать более `%s` элементов",
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s; (Code:%d)", e.Err, e.Code)
}

// Implement MarshalJSON for writing to json logs, and WriteTo for writing to non-json logs...

var (
	NotFound             = errors.New("Route not found")
	NoMethod             = errors.New("Method not allowed")
	ServerError          = errors.New("Internal server error")
	NoDataInRequestError = errors.New("No data in requests")
	NoRowsAffected       = errors.New("No rows affected")
	RecordNotFound       = errors.New("record not found")
)

func HasErrors(err interface{}) bool {
	hasErrors := false
	switch e := err.(type) {
	case gorm.Errors:
		hasErrors = len(e) > 0
	case []error:
		hasErrors = len(e) > 0
	default:
		hasErrors = e != nil
	}

	return hasErrors
}

// Извлечение ошибок из gorm запросов (для update|delete)
func GetResultErrors(result *gorm.DB) []error {
	err := result.GetErrors()
	affected := result.RowsAffected
	if affected == 0 {
		err = append(err, NoRowsAffected)
	}
	return err
}

// validationErrors Формирование массива ошибок
func MakeErrorsSlice(err error) map[string][]string {
	ve := map[string][]string{}
	for _, e := range err.(validator.ValidationErrors) {
		field := getFieldName(e.Namespace(), e.Field())
		if _, ok := ve[field]; !ok {
			ve[field] = []string{}
		}
		ve[field] = append(
			ve[field],
			getErrMessage(e.ActualTag(), field, e.Param()),
		)
	}
	return ve
}
func getFieldName(namespace string, field string) string {
	namespace = strings.Replace(namespace, "]", "", -1)
	namespace = strings.Replace(namespace, "[", ".", -1)
	namespaceSlice := strings.Split(namespace, ".")
	fieldName := field

	if len(namespaceSlice) > 2 {
		fieldName = strings.Join([]string{strings.Join(namespaceSlice[1:len(namespaceSlice)-1], "."), field}, ".")
	}

	return fieldName
}

func getErrMessage(errorType string, field string, param string) string {
	errKey := errorType
	_, ok := CommonValidationErrors[errorType]
	if !ok {
		errKey = ek
	}

	if param != "" && errKey == ek {
		return fmt.Sprintf(CommonValidationErrors[errKey], field, errorType)
	}

	return fmt.Sprintf(CommonValidationErrors[errKey], field)
}
