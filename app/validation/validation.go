package validation

import (
	"errors"
	"net/http"
)

var (
	EmptyFieldError = errors.New("field should not be empty")
)

func AddFormFieldError(r *http.Request, field string, err error) {
}
