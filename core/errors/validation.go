package errors

import (
	"fmt"

	"github.com/rafaeljusto/druns/core/tr"
)

type Validation struct {
	Err  error
	Code tr.Code
}

func NewValidation(code tr.Code, err error) error {
	return Validation{
		Err:  err,
		Code: code,
	}
}

func (v Validation) Error() string {
	return fmt.Sprintf("%s", v.Code)
}
