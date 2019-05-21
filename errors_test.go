package errors

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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
		"gorm error":   {input: gorm.Errors{errors.New("some error"), errors.New("some more error")}, hasError: true},
	}

	for caseName, tc := range tcs {
		if !assert.Equal(t, HasErrors(tc.input), tc.hasError, caseName) {
			t.FailNow()
		}
	}
}

func TestNew(t *testing.T) {
	errMsg := "err msg"
	err := New(errMsg)
	if !assert.Error(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, errMsg, err.Error()) {
		t.FailNow()
	}
}

func TestNewf(t *testing.T) {
	errFormat := "err msg: %s | %s"
	errMsg1 := "alert 1!"
	errMsg2 := "alert 2!"
	err := Newf(errFormat, errMsg1, errMsg2)
	if !assert.Error(t, err) {
		t.FailNow()
	}

	if !assert.Equal(t, fmt.Sprintf(errFormat, errMsg1, errMsg2), err.Error()) {
		t.FailNow()
	}
}
