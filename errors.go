package errors

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	NilErrorCode     = 0
	UnknownErrorCode = -1
	ek               = "common"
)

var (
	NotFound             = errors.New("Route not found")
	NoMethod             = errors.New("Method not allowed")
	ServerError          = errors.New("Internal server error")
	NoDataInRequestError = errors.New("No data in requests")
	NoRowsAffected       = errors.New("No rows affected")
	RecordNotFound       = errors.New("record not found")
)

func New(msg string) error {
	return errors.New(msg)
}

func HasErrors(err interface{}) bool {
	hasErrors := false
	switch e := err.(type) {
	case gorm.Errors:
		hasErrors = len(e) > 0
	case []error:
		hasErrors = len(e) > 0
	case map[string]error:
		hasErrors = len(e) > 0
	default:
		hasErrors = e != nil
	}

	return hasErrors
}
