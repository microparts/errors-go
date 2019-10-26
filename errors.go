package errors

import (
	"errors"
	"fmt"
)

//New returns new error with passed message
func New(msg string) error {
	return errors.New(msg)
}

//Newf returns new error with message sprintf'ed by format with passed params
func Newf(format string, params ...interface{}) error {
	return fmt.Errorf(format, params...)
}

//HasErrors checks if error occurs in passed err
func HasErrors(err interface{}) bool {
	hasErrors := false
	switch e := err.(type) {
	case []error:
		hasErrors = len(e) > 0
	case map[string]error:
		hasErrors = len(e) > 0
	default:
		hasErrors = e != nil
	}

	return hasErrors
}
