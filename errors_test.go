package errors

import (
	"errors"
	"fmt"
	"testing"
)

type testCase struct {
	input    interface{}
	hasError bool
}

func TestHasErrors(t *testing.T) {
	tcs := map[string]testCase{
		"common error": {input: errors.New("some error"), hasError: true},
		"errors slice": {input: []error{errors.New("some error"), errors.New("some more error")}, hasError: true},
		"errors map":   {input: map[string]error{"some error": errors.New("some error"), "some more error": errors.New("some more error")}, hasError: true},
		//"gorm error":   {input: gorm.Errors{errors.New("some error"), errors.New("some more error")}, hasError: true},
	}

	for caseName, tc := range tcs {
		if HasErrors(tc.input) != tc.hasError {
			fmt.Printf("test fail in `%s`", caseName)
			t.FailNow()
		}
	}
}

func TestNew(t *testing.T) {
	errMsg := "err msg"
	err := New(errMsg)
	if err == nil {
		fmt.Print("no errors occurs")
		t.FailNow()
	}
	if errMsg != err.Error() {
		fmt.Print("error message is not equals to err.Error() string")
		t.FailNow()
	}
}

func TestNewf(t *testing.T) {
	errFormat := "err msg: %s | %s"
	errMsg1 := "alert 1!"
	errMsg2 := "alert 2!"
	err := Newf(errFormat, errMsg1, errMsg2)
	if err == nil {
		fmt.Print("no errors occurs")
		t.FailNow()
	}

	if fmt.Sprintf(errFormat, errMsg1, errMsg2) != err.Error() {
		fmt.Print("errors are not equals")
		t.FailNow()
	}
}
